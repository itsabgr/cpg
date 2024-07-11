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

func (cpg *CPG) EnqueueCheckWallet(ctx context.Context, wallet string, asset string) error {
	if !cpg.assets.Has(asset) {
		return errors.New("unknown asset")
	}
	_, err := model.Transaction(ctx, cpg.db, &sql.TxOptions{Isolation: sql.LevelSerializable}, func(ctx context.Context, tx model.Tx) (struct{}, error) {
		isWatched, err := cpg.IsWatchedWallet(ctx, tx, wallet)
		if err != nil {
			return struct{}{}, err
		}
		if !isWatched {
			return struct{}{}, errors.New("wallet not watched")
		}
		return struct{}{}, cpg.enqueueCheckWallet(ctx, cpg.db, wallet, asset)
	})
	return err
}

func (cpg *CPG) enqueueCheckWallet(ctx context.Context, tx model.Tx, wallet string, asset string) error {
	_, err := cpg.riverClient.InsertTx(ctx, tx, CheckWalletArgs{Wallet: wallet}, &river.InsertOpts{
		MaxAttempts: cpg.config.MaxCheckWalletAttempts,
		Metadata:    nil,
		Pending:     false,
		Priority:    0,
		Queue:       asset,
		ScheduledAt: time.Time{},
		Tags:        nil,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	return err
}

func (cpg *CPG) EnqueueCheckout(ctx context.Context, invoiceId uuid.UUID, asset string) error {
	if !cpg.assets.Has(asset) {
		return errors.New("unknown asset")
	}
	_, err := model.Transaction(ctx, cpg.db, &sql.TxOptions{Isolation: sql.LevelSerializable}, func(ctx context.Context, tx model.Tx) (struct{}, error) {
		invoice, err := model.SelectInvoice(ctx, cpg.db, invoiceId)
		if err != nil {
			return struct{}{}, err
		}
		if invoice == nil {
			return struct{}{}, errors.New("invoice not found")
		}
		if invoice.Status != model.InvoiceStatusDone {
			return struct{}{}, errors.New("invoice not done")
		}
		return struct{}{}, cpg.enqueueCheckout(ctx, cpg.db, invoiceId, asset)
	})
	return err
}

func (cpg *CPG) enqueueCheckout(ctx context.Context, tx model.Tx, invoiceId uuid.UUID, asset string) error {
	_, err := cpg.riverClient.InsertTx(ctx, tx, CheckoutArgs{
		InvoiceId: invoiceId,
	}, &river.InsertOpts{
		MaxAttempts: cpg.config.MaxCheckWalletAttempts,
		Metadata:    nil,
		Pending:     false,
		Priority:    0,
		Queue:       asset,
		ScheduledAt: time.Time{},
		Tags:        nil,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	return err
}
