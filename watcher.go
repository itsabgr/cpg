package cpg

import (
	"cirello.io/pglock"
	"context"
	"database/sql"
	"errors"
	"github.com/itsabgr/cpg/pkg/model"
	"time"
)

func (cpg *CPG) processAsset(ctx context.Context, asset AssetImpl) (snooze time.Duration, err error) {
	lock, err := cpg.pgLock.AcquireContext(ctx, asset.Info().Name, nil)
	if err != nil {
		if errors.Is(err, pglock.ErrNotAcquired) {
			return asset.Info().Delay, nil
		}
		return 0, err
	}
	defer lock.Close()
	return cpg.processNextBlock(ctx, asset)
}
func (cpg *CPG) processNextBlock(ctx context.Context, asset AssetImpl) (snooze time.Duration, err error) {
	lastBlock, err := model.GetLastBlock(ctx, cpg.db, asset.Info().Name)
	if err != nil {
		return 0, err
	}
	return cpg.processBlock(ctx, asset, lastBlock+1)
}
func (cpg *CPG) processBlock(ctx context.Context, asset AssetImpl, n int64) (snooze time.Duration, err error) {
	assetInfo := asset.Info()
	wallets, err := asset.GetBlockAddresses(ctx, n)
	if err != nil {
		return 0, err
	}
	if wallets == nil {
		return assetInfo.Delay, nil
	}
	if len(wallets) == 0 {
		return 0, nil
	}
	return model.Transaction(ctx, cpg.db, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}, func(ctx context.Context, tx model.Tx) (time.Duration, error) {
		for _, wallet := range wallets {
			isWatched, err := cpg.IsWatchedWallet(ctx, tx, wallet)
			if err != nil {
				return 0, err
			}
			if !isWatched {
				continue
			}
			if err = cpg.enqueueCheckWallet(ctx, tx, wallet, assetInfo.Name); err != nil {
				return 0, err
			}
		}
		err := model.SetLastBlock(ctx, tx, assetInfo.Name, n)
		return 0, err
	})
}
