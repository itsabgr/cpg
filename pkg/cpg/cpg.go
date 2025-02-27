package cpg

import (
	"context"
	"cpg/pkg/crypto"
	"github.com/google/uuid"
	"github.com/itsabgr/ge"
	"log/slog"
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

	if err = cpg.tryAutoCheckout(ctx, inv.ID); err != nil {
		return err
	}

	return
}

type CreateInvoiceParams struct {
	AssetName    string
	Metadata     string
	Recipient    string
	Beneficiary  string
	AutoCheckout bool
	MinAmount    *big.Int
	Deadline     time.Time
}

type CreateInvoiceResult struct {
	InvoiceID     string
	InvoiceBackup []byte
}

func (cpg *CPG) CreateInvoice(ctx context.Context, params CreateInvoiceParams) (result CreateInvoiceResult, err error) {
	if params.Beneficiary == params.Recipient {
		err = ge.New("same beneficiary and recipient")
		return result, err
	}
	if params.MinAmount.Cmp(big.NewInt(0)) <= 0 {
		err = ge.New("non positive min amount")
		return result, err
	}
	if len(params.Metadata) >= 256 {
		err = ge.New("too big metadata")
		return result, err
	}
	if !params.Deadline.After(time.Now()) {
		err = ge.New("past deadline")
		return result, err
	}
	assetProvider := cpg.assets.Get(params.AssetName)
	if assetProvider == nil {
		err = ge.New("asset is not supported")
		return result, err
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
	inv.AuthCheckout = params.AutoCheckout
	inv.MinAmount.Set(params.MinAmount)
	inv.EncryptedSalt = randomEncryptedSalt(cpg.saltKeyring, assetInfo.SaltLength)

	inv.WalletAddress = ""
	if err = assetProvider.PrepareInvoice(ctx, inv); err != nil {
		err = ge.Wrap(ge.New("failed to prepare invoice"), err)
		return result, err
	}
	ge.Assert(inv.WalletAddress != "")

	result.InvoiceID = inv.ID
	result.InvoiceBackup = inv.Encrypt(cpg.backupKeyring)

	if err = cpg.db.InsertInvoice(ctx, inv, false); err != nil {
		err = ge.Wrap(ge.New("failed to insert invoice into db"), err)
		return result, err
	}

	return
}

func randomEncryptedSalt(saltKeyring *crypto.KeyRing, saltLength int) []byte {
	return saltKeyring.Encrypt(crypto.ReadN(nil, saltLength), (*[24]byte)(crypto.ReadN(nil, 24)))
}

type RequestCheckoutParams struct {
	InvoiceID string
}

func (cpg *CPG) RequestCheckout(ctx context.Context, params RequestCheckoutParams) (err error) {

	if params.InvoiceID == "" {
		return ge.New("invoice id is empty")
	}

	var inv *Invoice
	inv, err = cpg.db.GetInvoice(ctx, params.InvoiceID, "", false)
	if err != nil {
		return ge.Wrap(ge.New("failed to get invoice"), err)
	}
	if inv == nil {
		return ge.New("invoice not found")
	}

	if inv.CheckoutRequestAt != nil {
		return ge.New("already requested to checkout")
	}

	invoiceStatus := inv.Status()

	switch invoiceStatus {
	case InvoiceStatusExpired, InvoiceStatusCanceled, InvoiceStatusFilled, InvoiceStatusCheckout:
	case InvoiceStatusPending:
		return ge.New("invoice status is pending")
	default:
		return ErrInvalidInvoiceStatus
	}

	if err = cpg.db.SetInvoiceCheckoutRequestAt(ctx, params.InvoiceID); err != nil {
		err = ge.Wrap(ge.New("failed to update invoice checkout request"), err)
		return err
	}

	return err
}

type CancelInvoiceParams struct {
	InvoiceID     string
	WalletAddress string // to make sure
}

func (cpg *CPG) CancelInvoice(ctx context.Context, params CancelInvoiceParams) (err error) {

	if params.InvoiceID == "" {
		return ge.New("invoice id is empty")
	}

	if params.WalletAddress == "" {
		return ge.New("wallet address is empty")
	}

	var inv *Invoice
	inv, err = cpg.db.GetInvoice(ctx, params.InvoiceID, params.WalletAddress, false)
	if err != nil {
		err = ge.Wrap(ge.New("failed to get invoice"), err)
		return err
	}
	if inv == nil {
		err = ge.New("invoice not found")
		return err
	}

	invoiceStatus := inv.Status()

	switch invoiceStatus {
	case InvoiceStatusPending:
	case InvoiceStatusExpired, InvoiceStatusCanceled, InvoiceStatusFilled, InvoiceStatusCheckout:
		err = ge.Detail(ge.New("invoice status is not pending"), ge.D{"invoiceStatus": invoiceStatus})
		return err
	default:
		return ErrInvalidInvoiceStatus
	}

	if err = cpg.db.SetInvoiceCancelAt(ctx, params.InvoiceID); err != nil {
		err = ge.Wrap(ge.New("failed to update invoice cancelAt"), err)
		return err
	}

	if err = cpg.tryAutoCheckout(ctx, inv.ID); err != nil {
		return err
	}

	return err
}

type GetInvoiceParams struct {
	InvoiceID string
}

type GetInvoiceResult struct {
	MinAmount         big.Int
	Recipient         string
	Beneficiary       string
	Asset             string
	Metadata          string
	CreateAt          time.Time
	Deadline          time.Time
	FillAt            *time.Time
	CancelAt          *time.Time
	LastCheckoutAt    *time.Time
	CheckoutRequestAt *time.Time
	AutoCheckout      bool
	WalletAddress     string
	Status            InvoiceStatus
}

func (cpg *CPG) GetInvoice(ctx context.Context, params GetInvoiceParams) (result GetInvoiceResult, err error) {

	if params.InvoiceID == "" {
		return result, ge.New("invoice id is empty")
	}

	inv, err := cpg.db.GetInvoice(ctx, params.InvoiceID, "", false)
	if err != nil {
		err = ge.Wrap(ge.New("failed to get invoice"), err)
		return
	}
	if inv == nil {
		err = ge.New("invoice not found")
		return
	}

	result = GetInvoiceResult{
		MinAmount:         inv.MinAmount,
		Recipient:         inv.Recipient,
		Beneficiary:       inv.Beneficiary,
		Asset:             inv.Asset,
		Metadata:          inv.Metadata,
		CreateAt:          inv.CreateAt,
		Deadline:          inv.Deadline,
		FillAt:            inv.FillAt,
		CancelAt:          inv.CancelAt,
		LastCheckoutAt:    inv.LastCheckoutAt,
		CheckoutRequestAt: inv.CheckoutRequestAt,
		AutoCheckout:      inv.AuthCheckout,
		WalletAddress:     inv.WalletAddress,
		Status:            inv.Status(),
	}

	return
}

type CheckInvoiceParams struct {
	InvoiceID     string
	WalletAddress string
}

type CheckInvoiceResult struct {
	InvoiceStatus InvoiceStatus
}

func (cpg *CPG) CheckInvoice(ctx context.Context, params CheckInvoiceParams) (result CheckInvoiceResult, err error) {
	inv, err := cpg.db.GetInvoice(ctx, params.InvoiceID, params.WalletAddress, true)
	if err != nil {
		return result, ge.Wrap(ge.New("failed to get invoice"), err)
	}
	if inv == nil {
		return result, ge.New("invoice not found")
	}

	asset := cpg.assets.Get(inv.Asset)
	if asset == nil {
		return result, ge.New("asset is not supported")
	}

	switch result.InvoiceStatus = inv.Status(); result.InvoiceStatus {

	case InvoiceStatusExpired, InvoiceStatusCanceled, InvoiceStatusFilled, InvoiceStatusCheckout:

		break

	case InvoiceStatusPending:

		inv.saltKeyring = cpg.saltKeyring

		invoiceBalance, err := asset.GetBalance(ctx, inv)
		if err != nil {
			return result, ge.Wrap(ge.New("failed to get invoice balance"), err)
		}

		if invoiceBalance.Cmp(&inv.MinAmount) < 0 {
			return result, nil
		}

		now := time.Now()

		if err = cpg.db.SetInvoiceFillAt(ctx, inv.ID); err != nil {
			return result, ge.Wrap(ge.New("failed to update invoice fill_at"), err)
		}

		inv.FillAt = &now

	default:
		return result, ErrInvalidInvoiceStatus
	}

	if err = cpg.tryAutoCheckout(ctx, inv.ID); err != nil {
		return result, err
	}

	return result, nil
}

type TryCheckoutInvoiceParams struct {
	InvoiceID string
}

func (cpg *CPG) TryCheckoutInvoice(ctx context.Context, params TryCheckoutInvoiceParams) (err error) {

	inv, err := cpg.db.GetInvoice(ctx, params.InvoiceID, "", true)

	if err != nil {
		return ge.Wrap(ge.New("failed to get invoice"), err)
	}

	if inv == nil {
		return ge.New("invoice not found")
	}

	if inv.CheckoutRequestAt == nil {
		return ge.New("invoice not requested to checkout")
	}

	if false == inv.CheckoutRequestAt.Before(time.Now()) {
		return ge.New("invoice is already checking out")
	}

	asset := cpg.assets.Get(inv.Asset)
	if asset == nil {
		return ge.New("asset is not supported")
	}

	switch invoiceStatus := inv.Status(); invoiceStatus {

	case InvoiceStatusPending:

		return ge.New("invoice status is pending")

	case InvoiceStatusExpired, InvoiceStatusCanceled, InvoiceStatusFilled, InvoiceStatusCheckout:

		inv.saltKeyring = cpg.saltKeyring

		if err = asset.TryFlush(ctx, inv); err != nil {
			return ge.Wrap(ge.New("failed to flush invoice"), err)
		}

		if err = cpg.db.SetInvoiceLastCheckoutAt(ctx, inv.ID); err != nil {
			slog.Warn("failed to update invoice last_checkout_at", slog.String("invoice", inv.ID), slog.String("error", err.Error()))
		}

		return nil

	default:
		return ErrInvalidInvoiceStatus
	}
}

func (cpg *CPG) tryAutoCheckout(ctx context.Context, invoiceId string) error {
	if err := cpg.db.TrySetAutoCheckout(ctx, invoiceId); err != nil {
		return ge.Wrap(ge.Detail(ge.New("failed to update checkout request of auto-checkout invoice"), ge.D{"invoice": invoiceId}), err)
	}
	return nil
}
