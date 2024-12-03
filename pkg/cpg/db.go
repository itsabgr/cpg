package cpg

import (
	"context"
	"cpg/pkg/ent/database"
	"cpg/pkg/ent/database/invoice"
	"time"
)

type DB struct {
	client  *database.Client
	lockTTL time.Duration
}

func NewDB(client *database.Client) *DB {
	return &DB{
		client:  client,
		lockTTL: time.Second * 10,
	}
}

func (db *DB) SetInvoiceCancelAt(ctx context.Context, id string, cancelAt time.Time) error {
	return db.client.Invoice.UpdateOneID(id).SetCancelAt(cancelAt).Exec(ctx)
}

func (db *DB) SetInvoiceFillAt(ctx context.Context, id string, fillAt time.Time) error {
	return db.client.Invoice.UpdateOneID(id).SetFillAt(fillAt).Exec(ctx)
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

func (db *DB) GetInvoice(ctx context.Context, id string, withSalt bool) (*Invoice, error) {
	fields := []string{
		invoice.FieldMinAmount,
		invoice.FieldRecipient,
		invoice.FieldBeneficiary,
		invoice.FieldAsset,
		invoice.FieldMetadata,
		invoice.FieldCreateAt,
		invoice.FieldDeadline,
		invoice.FieldFillAt,
		invoice.FieldCancelAt,
		invoice.FieldWalletAddress,
	}
	if withSalt {
		fields = append(fields, invoice.FieldEncryptedSalt)
	}

	found, err := db.client.Invoice.Query().Where(invoice.ID(id)).Select(fields...).Only(ctx)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	inv := &Invoice{
		ID:            id,
		MinAmount:     *found.MinAmount,
		Recipient:     found.Recipient,
		Beneficiary:   found.Beneficiary,
		Asset:         found.Asset,
		Metadata:      found.Metadata,
		CreateAt:      found.CreateAt,
		Deadline:      found.Deadline,
		FillAt:        found.FillAt,
		CancelAt:      found.CancelAt,
		WalletAddress: found.WalletAddress,
		EncryptedSalt: found.EncryptedSalt,
	}

	return inv, nil
}
