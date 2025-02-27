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
	client *database.Client
}

func NewDB(client *database.Client) *DB {
	return &DB{
		client: client,
	}
}

func (db *DB) SetInvoiceCancelAt(ctx context.Context, id string) error {
	at := time.Now()
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.DeadlineGT(at),
		invoice.FillAtIsNil(),
		invoice.LastCheckoutAtIsNil(),
		invoice.CancelAtIsNil(),
	).SetCancelAt(at).Save(ctx)

	if inv == nil || (err != nil && database.IsNotFound(err)) {
		return ge.New("invoice not found or can not cancel")
	}

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) SetInvoiceFillAt(ctx context.Context, id string) error {
	at := time.Now()
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.DeadlineGT(at),
		invoice.FillAtIsNil(),
		invoice.LastCheckoutAtIsNil(),
		invoice.CancelAtIsNil(),
	).SetFillAt(at).Save(ctx)

	if inv == nil || (err != nil && database.IsNotFound(err)) {
		return ge.New("invoice not found or can not fill")
	}

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) SetInvoiceLastCheckoutAt(ctx context.Context, id string) error {
	at := time.Now()
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.Or(
			invoice.DeadlineLT(at),
			invoice.FillAtNotNil(),
			invoice.CancelAtNotNil(),
		),
	).SetLastCheckoutAt(at).Save(ctx)

	if inv == nil || (err != nil && database.IsNotFound(err)) {
		return ge.New("invoice not found or can not checkout")
	}

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) SetInvoiceCheckoutRequestAt(ctx context.Context, id string) error {
	at := time.Now()
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.And(
			invoice.CheckoutRequestAtIsNil(),
			invoice.Or(
				invoice.DeadlineLT(at),
				invoice.FillAtNotNil(),
				invoice.CancelAtNotNil(),
			),
		),
	).SetCheckoutRequestAt(at).Save(ctx)

	if inv == nil || (err != nil && database.IsNotFound(err)) {
		return ge.New("invoice not found or can not checkout or already requested to checkout")
	}

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) TrySetAutoCheckout(ctx context.Context, id string) error {
	at := time.Now()
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.And(
			invoice.AutoCheckout(true),
			invoice.CheckoutRequestAtIsNil(),
			invoice.Or(
				invoice.DeadlineLT(at),
				invoice.FillAtNotNil(),
				invoice.CancelAtNotNil(),
			),
		),
	).SetCheckoutRequestAt(at).Save(ctx)

	if inv == nil || (err != nil && database.IsNotFound(err)) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DoneInvoiceCheckoutRequestAt(ctx context.Context, id string) error {
	inv, err := db.client.Invoice.UpdateOneID(id).Where(
		invoice.And(
			invoice.Or(
				invoice.CheckoutRequestAtNotNil(),
				invoice.CheckoutRequestAtLT(time.Now()),
			),
			invoice.Or(
				invoice.DeadlineLT(time.Now()),
				invoice.FillAtNotNil(),
				invoice.CancelAtNotNil(),
			),
		),
	).ClearCheckoutRequestAt().Save(ctx)

	if inv == nil || (err != nil && database.IsNotFound(err)) {
		return ge.New("invoice not found or can not checkout or already not requested to checkout")
	}

	if err != nil {
		return err
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
		SetAutoCheckout(inv.AuthCheckout).
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
		invoice.FieldCheckoutRequestAt,
		invoice.FieldCancelAt,
		invoice.FieldWalletAddress,
		invoice.FieldAutoCheckout,
	}
	if withSalt {
		fields = append(fields, invoice.FieldEncryptedSalt)
	}

	where := make([]predicate.Invoice, 0, 2)

	if id != "" {
		where = append(where, invoice.ID(id))
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
		ID:                id,
		MinAmount:         *found.MinAmount,
		Recipient:         found.Recipient,
		Beneficiary:       found.Beneficiary,
		Asset:             found.Asset,
		Metadata:          found.Metadata,
		CreateAt:          found.CreateAt,
		Deadline:          found.Deadline,
		FillAt:            found.FillAt,
		LastCheckoutAt:    found.LastCheckoutAt,
		CheckoutRequestAt: found.CheckoutRequestAt,
		CancelAt:          found.CancelAt,
		AuthCheckout:      found.AutoCheckout,
		WalletAddress:     found.WalletAddress,
		EncryptedSalt:     found.EncryptedSalt,
	}

	return inv, nil
}
