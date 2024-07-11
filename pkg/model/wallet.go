package model

import (
	"context"
)

type Wallet struct {
	Address string
	Key     []byte
}

func InsertWallet(ctx context.Context, tx Tx, wallet *Wallet) error {
	return Affected(tx.ExecContext(ctx, `insert into wallet ("address", "key") values (?,?);`, wallet.Address, wallet.Key)).ExactAffect(1)
}
func DeleteWallet(ctx context.Context, tx Tx, address string) error {
	return Affected(tx.ExecContext(ctx, `delete from wallet where address = ?;`, address)).ExactAffect(1)
}
func SelectWallet(ctx context.Context, tx Tx, address string) (*Wallet, error) {
	wallet := &Wallet{Address: address}
	err := tx.QueryRowContext(ctx, `select key from wallet where address = ?;`, address).Scan(&wallet.Key)
	if IsNotFound(err) {
		return nil, nil
	}
	return wallet, err
}
