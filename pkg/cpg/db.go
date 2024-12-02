package cpg

import (
	"context"
	"cpg/pkg/ent/database"
	"cpg/pkg/ent/database/invoice"
	"github.com/google/uuid"
	"github.com/itsabgr/ge"
	"github.com/itsabgr/retry"
	"sync"
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

func (db *DB) tryAcquireInvoice(ctx context.Context, id string, lockHolder uuid.UUID) (bool, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	affected, err := db.client.Invoice.Update().Where(invoice.ID(id), invoice.LockExpireAtLT(time.Now())).SetLockHolder(lockHolder).SetLockExpireAt(time.Now().Add(db.lockTTL)).Save(timeout)
	return affected > 0, err
}

func (db *DB) invoiceExists(ctx context.Context, id string) (bool, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	return db.client.Invoice.Query().Where(invoice.ID(id)).Exist(timeout)
}

func (db *DB) pingLock(ctx context.Context, lockHolder uuid.UUID) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	affected, err := db.client.Invoice.Update().Where(invoice.LockHolder(lockHolder)).SetLockExpireAt(time.Now().Add(db.lockTTL)).Save(timeout)
	if err != nil {
		return err
	}
	if affected <= 0 {
		return ge.New("no lock affected")
	}
	return nil
}

func (db *DB) unlock(ctx context.Context, lockHolder uuid.UUID) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	return db.client.Invoice.Update().Where(invoice.LockHolder(lockHolder)).SetLockExpireAt(time.Now()).Exec(timeout)
}

func (db *DB) acquireInvoice(ctx context.Context, id string, lockHolder uuid.UUID) (err error) {
	acquired := false
	for range retry.Retry(ctx, 1, 0, 5, time.Second) {
		acquired, err = db.tryAcquireInvoice(ctx, id, lockHolder)
		if err != nil {
			return ge.Wrap(ge.New("failed to try acquire invoice"), err)
		}
		if acquired {
			break
		}
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	if !acquired {
		return ge.New("failed to acquire invoice")
	}
	return nil
}

func (db *DB) LockInvoice(ctx context.Context, id string, fn func(ctx context.Context)) error {

	if exists, err := db.invoiceExists(ctx, id); err != nil {
		return ge.Wrap(ge.New("failed to check invoice existence"), err)
	} else if !exists {
		return ge.New("invoice not exists")
	}

	lockHolder := uuid.New()
	if err := db.acquireInvoice(ctx, id, lockHolder); err != nil {
		return err
	}
	defer func() { _ = db.unlock(context.Background(), lockHolder) }()

	ctx, cancel := db.pinger(ctx, lockHolder)
	defer cancel()

	fn(ctx)

	return nil
}

func (db *DB) pinger(ctx context.Context, lockHolder uuid.UUID) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer cancel(ge.New("pinger returned"))
		for range retry.Retry(ctx, db.lockTTL/2) {
			if pingErr := db.pingLock(ctx, lockHolder); pingErr != nil {
				cancel(ge.Wrap(ge.New("pinger failed"), pingErr))
				return
			}
		}
	}()

	return ctx, func() {
		defer wg.Wait()
		cancel(nil)
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
