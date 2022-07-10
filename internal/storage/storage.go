package storage

import (
	"context"
	"fmt"

	"github.com/ffo32167/flowershop/internal"
)

type StorageProducts struct {
	sqlDB   SqlDB
	noSqlDB NoSqlDB
}

type SqlDB interface {
	List(ctx context.Context) ([]internal.Product, error)
	Sale(ctx context.Context, id, cnt int) (int, error)
}

type NoSqlDB interface {
	ListCreate(ctx context.Context, products []internal.Product) error
	List(ctx context.Context) ([]internal.Product, error)
	Sale(ctx context.Context, id, cnt int) error
}

func New(ctx context.Context, sqlDB SqlDB, noSqlDB NoSqlDB) (StorageProducts, error) {
	c := StorageProducts{sqlDB: sqlDB, noSqlDB: noSqlDB}
	err := c.RenewCache(ctx)
	if err != nil {
		return StorageProducts{}, fmt.Errorf("cant renew cache: %w", err)
	}
	return c, nil
}

func (c StorageProducts) RenewCache(ctx context.Context) error {
	products, err := c.sqlDB.List(ctx)
	if err != nil {
		return fmt.Errorf("cant get product list: %w", err)
	}
	err = c.noSqlDB.ListCreate(ctx, products)
	if err != nil {
		return fmt.Errorf("cant create redis list: %w", err)
	}
	return nil
}

func (c StorageProducts) List(ctx context.Context) ([]internal.Product, error) {
	return c.noSqlDB.List(ctx)
}

func (c StorageProducts) Sale(ctx context.Context, id int, cnt int) error {
	_, err := c.sqlDB.Sale(ctx, id, cnt)
	if err != nil {
		return fmt.Errorf("cant sale in sql: %w", err)
	}
	err = c.noSqlDB.Sale(ctx, id, cnt)
	if err != nil {
		return fmt.Errorf("cant sale in nosql: %w", err)
	}
	return nil
}
