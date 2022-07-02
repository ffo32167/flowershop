package cachedb

import (
	"context"
	"fmt"

	"github.com/ffo32167/flowershop/internal"
)

type CacheProducts struct {
	sqlDB   SqlDB
	noSqlDB NoSqlDB
	ctx     context.Context
}

type SqlDB interface {
	List() ([]internal.Product, error)
	Sale(id, cnt int) (res int, err error)
}
type NoSqlDB interface {
	ListCreate(ctx context.Context, products []internal.Product) error
	List(ctx context.Context) ([]internal.Product, error)
	Sale(id, cnt int) (int, error)
}

// передавать внутрь готовый PG, который интерфейс, потому что DI
func New(ctx context.Context, sqlDB SqlDB, noSqlDB NoSqlDB) (CacheProducts, error) {

	c := CacheProducts{sqlDB: sqlDB, noSqlDB: noSqlDB, ctx: ctx}
	err := c.RenewCache()
	if err != nil {
		return CacheProducts{}, err
	}
	return c, nil
}

func (c CacheProducts) RenewCache() error {
	products, err := c.sqlDB.List()
	if err != nil {
		return err
	}
	err = c.noSqlDB.ListCreate(c.ctx, products)
	if err != nil {
		return fmt.Errorf("cant create redis list: %w", err)
	}
	return nil
}

func (c CacheProducts) List() ([]internal.Product, error) {
	return c.noSqlDB.List(c.ctx)
}

func (c CacheProducts) Sale(int, int) (int, error) {
	return 0, nil
}
