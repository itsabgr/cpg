package cpg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/itsabgr/cpg/pkg/model"
	"github.com/riverqueue/river"
	"time"
)

type CheckWalletArgs struct {
	Wallet string
}

func (CheckWalletArgs) Kind() string { return "CheckWallet" }

func (cpg *CPG) WorkCheckWallet(ctx context.Context, job river.Job[CheckWalletArgs]) error {
	assetName := job.Queue
	asset, err := cpg.assets.Get(assetName)
	if err != nil {
		return err
	}
	wallet, err := model.SelectWallet(ctx, cpg.db, job.Args.Wallet)
	if err != nil {
		return err
	}
	if wallet == nil {
		return river.JobCancel(errors.New("wallet not found"))
	}
	balanceOfWallet, err := asset.GetBalanceOfWallet(ctx, wallet)
	if err != nil {
		return err
	}
	_, err = model.Transaction(ctx, cpg.db, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}, func(ctx context.Context, tx model.Tx) (struct{}, error) {
		invoice, err := model.SelectInvoiceByWallet(ctx, cpg.db, wallet.Address)
		if err != nil {
			return struct{}{}, err
		}
		if invoice == nil {
			_ = model.DeleteWallet(ctx, cpg.db, wallet.Address)
			return struct{}{}, river.JobCancel(errors.New("invoice not found"))
		}
		if invoice.Status == model.InvoiceStatusDone {
			return struct{}{}, river.JobCancel(errors.New("invoice already done"))
		}
		if invoice.Amount.Cmp(balanceOfWallet) > 0 {
			return struct{}{}, river.JobCancel(errors.New("insufficient balance"))
		}
		err = model.SetInvoiceStatus(ctx, tx, invoice.Id, model.InvoiceStatusDone)
		if err != nil {
			return struct{}{}, err
		}
		err = cpg.enqueueCheckout(ctx, tx, invoice.Id, assetName)
		return struct{}{}, err
	})
	return err
}

type CheckoutArgs struct {
	InvoiceId uuid.UUID
}

func (CheckoutArgs) Kind() string { return "Checkout" }

func (cpg *CPG) WorkCheckout(ctx context.Context, job river.Job[CheckoutArgs]) error {
	asset, err := cpg.assets.Get(job.Queue)
	if err != nil {
		return err
	}
	models, err := model.Transaction(ctx, cpg.db, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	}, func(ctx context.Context, tx model.Tx) (*struct {
		Invoice *model.Invoice
		Wallet  *model.Wallet
	}, error) {
		invoice, err := model.SelectInvoice(ctx, tx, job.Args.InvoiceId)
		if err != nil {
			return nil, err
		}
		if invoice == nil {
			return nil, river.JobCancel(errors.New("invoice not found"))
		}
		wallet, err := model.SelectWallet(ctx, tx, invoice.Wallet)
		if err != nil {
			return nil, err
		}
		if wallet == nil {
			return nil, river.JobCancel(errors.New("wallet not found"))
		}
		return &struct {
			Invoice *model.Invoice
			Wallet  *model.Wallet
		}{invoice, wallet}, nil
	})
	if err != nil {
		return err
	}
	err = asset.Checkout(ctx, models.Wallet, models.Invoice)
	if err != nil {
		return river.JobSnooze(max(asset.Info().Delay/2, time.Second*5))
	}
	return nil
}
