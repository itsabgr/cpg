package cpg

import (
	"context"
	"cpg/pkg/crypto"
	"github.com/google/uuid"
	"github.com/itsabgr/ge"
	"math/big"
	"sync"
	"time"
)

type CPG struct {
	_             sync.Mutex
	assets        *Assets
	db            *DB
	backupKeyring *crypto.KeyRing
	saltKeyring   *crypto.KeyRing
}

func NewCPG(assets *Assets, db *DB, backupKeyring, saltKeyring *crypto.KeyRing) *CPG {
	return &CPG{
		assets:        assets,
		db:            db,
		backupKeyring: backupKeyring,
		saltKeyring:   saltKeyring,
	}
}

type RecoverInvoiceParams struct {
	InvoiceID     string
	InvoiceBackup []byte
}

func constantTimeOr(conds ...bool) bool {
	n := 0
	for _, c := range conds {
		if c {
			n += 3
		} else {
			n += 2
		}
	}
	if len(conds)*2 < n {
		return false
	}
	return true
}
func (cpg *CPG) RecoverInvoice(ctx context.Context, params RecoverInvoiceParams) (err error) {

	inv, err := DecryptInvoice(cpg.backupKeyring, params.InvoiceBackup)

	//inv is always non-nil

	if !constantTimeOr(
		err != nil,
		inv.ID != params.InvoiceID,
		inv.MinAmount.Cmp(big.NewInt(0)) <= 0,
		len(inv.Metadata) >= 256,
		!inv.Deadline.After(inv.CreateAt),
		inv.Recipient == inv.Beneficiary,
	) {
		err = ge.Wrap(ge.New("failed to recover backup"), err)
		return
	}

	assetProvider := cpg.assets.Get(inv.Asset)
	if assetProvider == nil {
		err = ge.New("asset is not supported no more")
		return
	}

	if err = assetProvider.PrepareInvoice(ctx, inv); err != nil {
		err = ge.Wrap(ge.New("failed to prepare recovered invoice"), err)
		return
	}
	if err = cpg.db.InsertInvoice(ctx, inv, true); err != nil {
		err = ge.Wrap(ge.New("failed to insert invoice into db"), err)
		return
	}
	return
}

type CreateInvoiceParams struct {
	AssetName   string
	Metadata    string
	Recipient   string
	Beneficiary string
	MinAmount   *big.Int
	Deadline    time.Time
}

type CreateInvoiceResult struct {
	InvoiceID     string
	InvoiceBackup []byte
}

func (cpg *CPG) CreateInvoice(ctx context.Context, params CreateInvoiceParams) (result CreateInvoiceResult, err error) {
	if params.Beneficiary == params.Recipient {
		err = ge.New("same beneficiary and recipient")
		return
	}
	if params.MinAmount.Cmp(big.NewInt(0)) <= 0 {
		err = ge.New("non positive min amount")
		return
	}
	if len(params.Metadata) >= 256 {
		err = ge.New("too big metadata")
		return
	}
	if false == !params.Deadline.After(time.Now()) {
		err = ge.New("past deadline")
		return
	}
	assetProvider := cpg.assets.Get(params.AssetName)
	if assetProvider == nil {
		err = ge.New("asset not found")
		return
	}

	assetInfo := assetProvider.Info()

	inv := &Invoice{saltKeyring: cpg.saltKeyring}
	inv.Metadata = params.Metadata
	inv.Recipient = params.Recipient
	inv.Beneficiary = params.Beneficiary
	inv.Asset = params.AssetName
	inv.Deadline = params.Deadline
	inv.ID = uuid.NewString()
	inv.CreateAt = time.Now()
	inv.MinAmount.Set(params.MinAmount)
	inv.EncryptedSalt = randomEncryptedSalt(cpg.saltKeyring, assetInfo.SaltLength)

	inv.WalletAddress = ""
	if err = assetProvider.PrepareInvoice(ctx, inv); err != nil {
		err = ge.Wrap(ge.New("failed to prepare invoice"), err)
		return
	}
	ge.Assert(inv.WalletAddress != "")

	result.InvoiceID = inv.ID
	result.InvoiceBackup = inv.Encrypt(cpg.backupKeyring)

	if err = cpg.db.InsertInvoice(ctx, inv, false); err != nil {
		err = ge.Wrap(ge.New("failed to insert invoice into db"), err)
		return
	}

	return
}

func randomEncryptedSalt(saltKeyring *crypto.KeyRing, saltLength int) []byte {
	return saltKeyring.Encrypt(crypto.ReadN(nil, saltLength), (*[24]byte)(crypto.ReadN(nil, 24)))
}

type CancelInvoiceParams struct {
	InvoiceID string
}

func (cpg *CPG) CancelInvoice(ctx context.Context, params CancelInvoiceParams) (err error) {
	var inv *Invoice
	inv, err = cpg.db.GetInvoice(ctx, params.InvoiceID, false)
	if err != nil {
		err = ge.Wrap(ge.New("failed to get invoice"), err)
		return
	}
	if inv == nil {
		err = ge.New("invoice not found")
		return
	}

	invoiceStatus := inv.Status()
	if invoiceStatus != InvoiceStatusPending {
		err = ge.Detail(ge.New("invoice status is not pending"), ge.D{"invoiceStatus": invoiceStatus})
		return
	}
	now := time.Now()
	if err = cpg.db.SetInvoiceCancelAt(ctx, params.InvoiceID, now); err != nil {
		err = ge.Wrap(ge.New("failed to update invoice cancelAt"), err)
		return
	}
	return
}

type GetInvoiceParams struct {
	InvoiceID string
}

type GetInvoiceResult struct {
	MinAmount     big.Int
	Recipient     string
	Beneficiary   string
	Asset         string
	Metadata      string
	CreateAt      time.Time
	Deadline      time.Time
	FillAt        *time.Time
	CancelAt      *time.Time
	WalletAddress string
	Status        InvoiceStatus
}

func (cpg *CPG) GetInvoice(ctx context.Context, params GetInvoiceParams) (result GetInvoiceResult, err error) {
	inv, err := cpg.db.GetInvoice(ctx, params.InvoiceID, false)
	if err != nil {
		err = ge.Wrap(ge.New("failed to get invoice"), err)
		return
	}
	if inv == nil {
		err = ge.New("invoice not found")
		return
	}

	result = GetInvoiceResult{
		MinAmount:     inv.MinAmount,
		Recipient:     inv.Recipient,
		Beneficiary:   inv.Beneficiary,
		Asset:         inv.Asset,
		Metadata:      inv.Metadata,
		CreateAt:      inv.CreateAt,
		Deadline:      inv.Deadline,
		FillAt:        inv.FillAt,
		CancelAt:      inv.CancelAt,
		WalletAddress: inv.WalletAddress,
		Status:        inv.Status(),
	}

	return
}

type CheckInvoiceParams struct {
	InvoiceID string
}

type CheckInvoiceResult struct {
	InvoiceStatus InvoiceStatus
}

func (cpg *CPG) CheckInvoice(ctx context.Context, params CheckInvoiceParams) (result CheckInvoiceResult, err error) {
	inv, _, err := cpg.checkInvoice(ctx, params.InvoiceID, true)

	if err != nil {
		return result, err
	}

	result.InvoiceStatus = inv.Status()

	return
}

type TryCheckoutInvoiceParams struct {
	InvoiceID    string
	CheckBalance bool
}

func (cpg *CPG) TryCheckoutInvoice(ctx context.Context, params TryCheckoutInvoiceParams) (err error) {

	inv, asset, err := cpg.checkInvoice(ctx, params.InvoiceID, params.CheckBalance)

	if err != nil {
		return err
	}

	switch inv.Status() {
	case InvoiceStatusExpired, InvoiceStatusCanceled, InvoiceStatusFilled:
	case InvoiceStatusPending:
		return ge.New("invoice is pending to fill")
	default:
		panic(ge.UNREACHABLE)

	}

	err = asset.TryFlush(ctx, inv)

	if err != nil {
		return ge.Wrap(ge.New("failed to flush invoice"), err)
	}

	return
}

func (cpg *CPG) checkInvoice(ctx context.Context, id string, getBalance bool) (inv *Invoice, assetProvider Asset, err error) {
	inv, err = cpg.db.GetInvoice(ctx, id, true)
	if err != nil {
		err = ge.Wrap(ge.New("failed to get invoice"), err)
		return nil, nil, err
	}
	if inv == nil {
		err = ge.New("invoice not found")
		return nil, nil, err
	}

	assetProvider = cpg.assets.Get(inv.Asset)
	if assetProvider == nil {
		err = ge.New("asset is not supported no more")
		return nil, nil, err
	}

	switch inv.Status() {
	case InvoiceStatusExpired, InvoiceStatusCanceled, InvoiceStatusFilled:
	case InvoiceStatusPending:
		if getBalance {
			var invoiceBalance *big.Int
			invoiceBalance, err = assetProvider.GetBalance(ctx, inv)
			if err != nil {
				err = ge.Wrap(ge.New("failed to get invoice balance"), err)
				return
			}
			if invoiceBalance.Cmp(&inv.MinAmount) < 0 {
				err = ge.Detail(ge.New("insufficient wallet balance"), ge.D{"balance": invoiceBalance.String()})
				return
			}
			now := time.Now()
			if err = cpg.db.SetInvoiceFillAt(ctx, inv.ID, now); err != nil {
				err = ge.Wrap(ge.New("failed to update invoice fill_at"), err)
				return
			}
			inv.FillAt = &now
		} else {
			err = ge.New("invoice is pending to fill")
			return nil, nil, err
		}
	default:
		panic(ge.UNREACHABLE)
	}
	return inv, assetProvider, nil
}
