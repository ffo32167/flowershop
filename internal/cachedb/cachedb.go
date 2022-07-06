package cachedb

import (
	"context"
	"fmt"

	"github.com/ffo32167/flowershop/internal"
)

// перенести в internal
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
	ListCreate(products []internal.Product) error
	List() ([]internal.Product, error)
	Sale(id, cnt int) (int, error)
}

func New(ctx context.Context, sqlDB SqlDB, noSqlDB NoSqlDB) (CacheProducts, error) {
	c := CacheProducts{sqlDB: sqlDB, noSqlDB: noSqlDB, ctx: ctx}
	err := c.RenewCache()
	if err != nil {
		return CacheProducts{}, fmt.Errorf("cant renew cache: %w", err)
	}
	return c, nil
}

func (c CacheProducts) RenewCache() error {
	products, err := c.sqlDB.List()
	if err != nil {
		return fmt.Errorf("cant get product list: %w", err)
	}
	err = c.noSqlDB.ListCreate(products)
	if err != nil {
		return fmt.Errorf("cant create redis list: %w", err)
	}
	return nil
}

func (c CacheProducts) List() ([]internal.Product, error) {
	return c.noSqlDB.List()
}

func (c CacheProducts) Sale(id int, cnt int) (int, error) {
	_, err := c.sqlDB.Sale(id, cnt)
	if err != nil {
		return 0, fmt.Errorf("cant sale in sql: %w", err)
	}
	_, err = c.noSqlDB.Sale(id, cnt)
	if err != nil {
		return 0, fmt.Errorf("cant sale in nosql: %w", err)
	}
	return 0, nil
}
