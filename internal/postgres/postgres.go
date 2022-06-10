package postgres

import (
	"context"
	"fmt"

	"github.com/ffo32167/flowershop/internal"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgDb struct {
	pool *pgxpool.Pool
}

type pgProducts []pgProduct

type pgProduct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func New(ctx context.Context, connStr string) (PgDb, error) {
	db, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return PgDb{}, fmt.Errorf("unable to connect to database: %w", err)
	}
	return PgDb{pool: db}, nil
}

func (db PgDb) List() ([]internal.Product, error) {
	rows, err := db.pool.Query(context.Background(),
		`select id, name from "flowerShop".products p`)
	if err != nil {
		return nil, fmt.Errorf("unable to execute select query: %w ", err)
	}
	defer rows.Close()

	var pgProduct pgProduct
	var pgProducts pgProducts
	for rows.Next() {
		err = rows.Scan(&pgProduct.Id, &pgProduct.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to scan query: %w ", err)
		}
		pgProducts = append(pgProducts, pgProduct)
	}
	internalProducts, err := pgProducts.toDomain()
	if err != nil {
		return internalProducts, fmt.Errorf("unable to convert PG Rates to internal Rates: %w ", err)
	}
	return internalProducts, nil
}

func (p pgProducts) toDomain() ([]internal.Product, error) {
	result := make([]internal.Product, len(p), len(p))
	for i := range p {
		result[i].Id = p[i].Id
		result[i].Name = p[i].Name
	}
	return result, nil
}
