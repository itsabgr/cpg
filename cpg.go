package cpg

import (
	"cirello.io/pglock"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/itsabgr/cpg/pkg/model"
	"github.com/itsabgr/cpg/pkg/registry"
	"github.com/riverqueue/river"
	"math/big"
	"time"
)

type AssetInfo struct {
	Name  string
	Delay time.Duration
}

type CreateInvoiceInput struct {
	Asset     string
	Recipient string
	Amount    *big.Int
}

type AssetImpl interface {
	Info() AssetInfo
	ValidateAddress(address string) bool
	WalletForInvoice(*CreateInvoiceInput) (*model.Wallet, error)
	GetBalanceOfWallet(ctx context.Context, wallet *model.Wallet) (*big.Int, error)
	Checkout(ctx context.Context, wallet *model.Wallet, invoice *model.Invoice) error
	GetBlockAddresses(ctx context.Context, block int64) ([]string, error)
}
type Config struct {
	MaxCheckWalletAttempts int
	MaxCheckoutAttempts    int
}
type CPG struct {
	assets      registry.Registry[string, AssetImpl]
	db          model.DB
	riverClient river.Client[any]
	config      Config
	pgLock      *pglock.Client
}

func (cpg *CPG) Assets() map[string]AssetInfo {
	res := make(map[string]AssetInfo)
	for _, asset := range cpg.assets {
		info := asset.Info()
		res[info.Name] = info
	}
	return res
}
func (cpg *CPG) CreateInvoice(ctx context.Context, input *CreateInvoiceInput) (uuid.UUID, error) {
	assetImpl, err := cpg.assets.Get(input.Asset)
	if err != nil {
		return uuid.UUID{}, err
	}
	if !assetImpl.ValidateAddress(input.Recipient) {
		return uuid.UUID{}, errors.New("invalid address")
	}
	wallet, err := assetImpl.WalletForInvoice(input)
	if err != nil {
		return uuid.UUID{}, err
	}
	invoice := &model.Invoice{
		Id:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		Wallet:    wallet.Address,
		Asset:     input.Asset,
		Amount:    input.Amount,
		Recipient: input.Recipient,
		Status:    model.InvoiceStatusPend,
	}
	_, err = model.Transaction(ctx, cpg.db, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}, func(ctx context.Context, tx model.Tx) (struct{}, error) {
		err := model.InsertInvoice(ctx, tx, invoice)
		if err != nil {
			return struct{}{}, err
		}
		err = model.InsertWallet(ctx, tx, wallet)
		if err != nil {
			return struct{}{}, err
		}
		return struct{}{}, nil
	})
	return invoice.Id, err
}

func (cpg *CPG) GetInvoice(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	return model.SelectInvoice(ctx, cpg.db, id)
}

func (cpg *CPG) IsWatchedWallet(ctx context.Context, tx model.Tx, wallet string) (bool, error) {
	invoice, err := model.SelectInvoiceByWallet(ctx, tx, wallet)
	if err != nil {
		return false, err
	}
	return invoice != nil && invoice.Status == model.InvoiceStatusPend, nil
}
