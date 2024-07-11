package model

import (
	"context"
	"github.com/google/uuid"
	"math/big"
	"time"
)

type Invoice struct {
	Id        uuid.UUID
	CreatedAt time.Time
	Wallet    string
	Asset     string
	Amount    *big.Int
	Recipient string
	Status    InvoiceStatus
}

func InsertInvoice(ctx context.Context, tx Tx, iv *Invoice) error {
	return Affected(tx.ExecContext(ctx, `insert into invoice(id, "createdAt", wallet, asset, amount, recipient, status) values (?,?,?,?,?,?,?);`, iv.Id, iv.CreatedAt, iv.Wallet, iv.Asset, iv.Amount, iv.Recipient, iv.Status)).ExactAffect(1)
}

func DeleteInvoice(ctx context.Context, tx Tx, id uuid.UUID) error {
	return Affected(tx.ExecContext(ctx, `delete from invoice where id = ?;`, id)).ExactAffect(1)
}
func SetInvoiceStatus(ctx context.Context, tx Tx, id uuid.UUID, status InvoiceStatus) error {
	return Affected(tx.ExecContext(ctx, `update invoice set status = ? where id = ?;`, status.String(), id)).ExactAffect(1)
}
func SelectInvoice(ctx context.Context, tx Tx, id uuid.UUID) (*Invoice, error) {
	invoice := &Invoice{Id: id}
	err := tx.QueryRowContext(ctx, `select "createdAt", wallet, asset, amount, recipient, status from invoice where id = ?`, id).Scan(&invoice.CreatedAt, &invoice.Wallet, &invoice.Asset, &invoice.Amount, &invoice.Recipient, &invoice.Status)
	if IsNotFound(err) {
		return nil, nil
	}
	return invoice, err
}
func SelectInvoiceByWallet(ctx context.Context, tx Tx, wallet string) (*Invoice, error) {
	invoice := &Invoice{Wallet: wallet}
	err := tx.QueryRowContext(ctx, `select id, "createdAt", asset, amount, recipient, status from invoice where wallet = ?`, wallet).Scan(&invoice.Id, &invoice.CreatedAt, &invoice.Asset, &invoice.Amount, &invoice.Recipient, &invoice.Status)
	if IsNotFound(err) {
		return nil, nil
	}
	return invoice, err
}
func SelectInvoiceStatus(ctx context.Context, tx Tx, id uuid.UUID) (InvoiceStatus, error) {
	var status string
	err := tx.QueryRowContext(ctx, `select status from invoice where id = ?`, id).Scan(&status)
	if IsNotFound(err) {
		return "", nil
	}
	return InvoiceStatus(status), nil
}
