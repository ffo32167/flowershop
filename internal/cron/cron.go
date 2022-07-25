package cron

import (
	"context"
	"time"

	"github.com/ffo32167/flowershop/internal/storage"
	"go.uber.org/zap"
)

type Cron struct {
	sp  storage.StorageProduct
	log *zap.Logger
}

func New(sp storage.StorageProduct, log *zap.Logger) Cron {
	return Cron{sp: sp, log: log}
}

func (c Cron) Action(ctx context.Context, renewInterval time.Duration) {
	for {
		go func() {
			err := c.sp.RenewCache(ctx)
			if err != nil {
				c.log.Error("cron call error: ", zap.Error(err))
			}
		}()
		time.Sleep(renewInterval)
	}
}
