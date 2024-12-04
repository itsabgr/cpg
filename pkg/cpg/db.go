package cpg

import (
	"context"
	"cpg/pkg/ent/database"
	"cpg/pkg/ent/database/invoice"
	"cpg/pkg/ent/database/predicate"
	"github.com/itsabgr/ge"
	"time"
)

type DB struct {
	client  *database.Client
	lockTTL time.Duration
}

var ErrInvoiceNotUpdate = ge.New("invoice did not update")

func NewDB(client *database.Client) *DB {
	return &DB{
		client:  client,
		lockTTL: time.Second * 10,
	}
}

func (db *DB) SetInvoiceCancelAt(ctx context.Context, id string, at time.Time) error {
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.DeadlineGT(time.Now()),
		invoice.FillAtIsNil(),
		invoice.LastCheckoutAtIsNil(),
		invoice.CancelAtIsNil(),
	).SetCancelAt(at).Save(ctx)
	if err != nil {
		if database.IsNotFound(err) {
			return ErrInvoiceNotUpdate
		}
		return err
	}
	if inv == nil {
		return ErrInvoiceNotUpdate
	}
	return nil
}

func (db *DB) SetInvoiceFillAt(ctx context.Context, id string, at time.Time) error {
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.DeadlineGT(time.Now()),
		invoice.FillAtIsNil(),
		invoice.LastCheckoutAtIsNil(),
		invoice.CancelAtIsNil(),
	).SetFillAt(at).Save(ctx)
	if err != nil {
		if database.IsNotFound(err) {
			return ErrInvoiceNotUpdate
		}
		return err
	}
	if inv == nil {
		return ErrInvoiceNotUpdate
	}
	return nil
}

func (db *DB) SetInvoiceLastCheckoutAt(ctx context.Context, id string, at time.Time) error {
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.Or(
			invoice.DeadlineLT(time.Now()),
			invoice.FillAtNotNil(),
			invoice.CancelAtNotNil(),
		),
	).SetLastCheckoutAt(at).Save(ctx)
	if err != nil {
		if database.IsNotFound(err) {
			return ErrInvoiceNotUpdate
		}
		return err
	}
	if inv == nil {
		return ErrInvoiceNotUpdate
	}
	return nil
}

func (db *DB) InsertInvoice(ctx context.Context, inv *Invoice, recovered bool) error {
	return db.client.Invoice.Create().
		SetID(inv.ID).
		SetMinAmount(&inv.MinAmount).
		SetRecipient(inv.Recipient).
		SetBeneficiary(inv.Beneficiary).
		SetAsset(inv.Asset).
		SetMetadata(inv.Metadata).
		SetCreateAt(inv.CreateAt).
		SetDeadline(inv.Deadline).
		SetWalletAddress(inv.WalletAddress).
		SetEncryptedSalt(inv.EncryptedSalt).
		Exec(ctx)
}

func (db *DB) GetInvoice(ctx context.Context, id, walletAddress string, withSalt bool) (*Invoice, error) {
	fields := []string{
		invoice.FieldMinAmount,
		invoice.FieldRecipient,
		invoice.FieldBeneficiary,
		invoice.FieldAsset,
		invoice.FieldMetadata,
		invoice.FieldCreateAt,
		invoice.FieldDeadline,
		invoice.FieldFillAt,
		invoice.FieldLastCheckoutAt,
		invoice.FieldCancelAt,
		invoice.FieldWalletAddress,
	}
	if withSalt {
		fields = append(fields, invoice.FieldEncryptedSalt)
	}

	where := make([]predicate.Invoice, 2)

	if id != "" {
		where = append(where, invoice.WalletAddress(id))
	}

	if walletAddress != "" {
		where = append(where, invoice.WalletAddress(walletAddress))
	}

	if len(where) <= 0 {
		return nil, ge.New("no invoice id or wallet address")
	}

	found, err := db.client.Invoice.Query().Where(where...).Select(fields...).Only(ctx)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	inv := &Invoice{
		ID:             id,
		MinAmount:      *found.MinAmount,
		Recipient:      found.Recipient,
		Beneficiary:    found.Beneficiary,
		Asset:          found.Asset,
		Metadata:       found.Metadata,
		CreateAt:       found.CreateAt,
		Deadline:       found.Deadline,
		FillAt:         found.FillAt,
		LastCheckoutAt: found.LastCheckoutAt,
		CancelAt:       found.CancelAt,
		WalletAddress:  found.WalletAddress,
		EncryptedSalt:  found.EncryptedSalt,
	}

	return inv, nil
}
