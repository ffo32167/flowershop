package storage_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ffo32167/flowershop/internal"
	"github.com/ffo32167/flowershop/internal/storage"
	"github.com/ffo32167/flowershop/internal/storage/mock"
)

//go:generate mockgen -source=./storage.go -destination=./mock/storage.go -package=mock
func TestProductRepository_List(t *testing.T) {
	tests := []struct {
		name            string
		ctx             context.Context
		expectedProduct []internal.Product
		expectedErrText string
		prepare         func(ctrl *gomock.Controller) (storage.SqlDB, storage.NoSqlDB)
	}{
		{
			name:            "01",
			ctx:             context.Background(),
			expectedProduct: []internal.Product{internal.Product{Id: 15, Name: "product-15", Qty: 1, Price: 150}},
			expectedErrText: "",
			prepare: func(ctrl *gomock.Controller) (storage.SqlDB, storage.NoSqlDB) {
				sqlDB := mock.NewMockSqlDB(ctrl)
				sqlDB.EXPECT().
					List(gomock.Any()).
					Return([]internal.Product{internal.Product{Id: 15, Name: "product-15", Qty: 1, Price: 150}}, nil)
				noSqlDB := mock.NewMockNoSqlDB(ctrl)
				noSqlDB.EXPECT().
					ListCreate(gomock.Any(), []internal.Product{internal.Product{Id: 15, Name: "product-15", Qty: 1, Price: 150}}).
					Return(nil)
				noSqlDB.EXPECT().
					List(gomock.Any()).
					Return([]internal.Product{internal.Product{Id: 15, Name: "product-15", Qty: 1, Price: 150}}, nil)
				return sqlDB, noSqlDB
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sql, noSql := test.prepare(ctrl)
			storage, err := storage.New(context.TODO(), sql, noSql)
			assert.NoError(t, err)
			product, err := storage.List(test.ctx)
			if test.expectedErrText != "" {
				assert.EqualError(t, err, test.expectedErrText)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedProduct, product)
			}
		})
	}
}

func TestProductRepository_RenewCache(t *testing.T) {
	tests := []struct {
		name            string
		ctx             context.Context
		expectedErrText string
		prepare         func(ctrl *gomock.Controller) (storage.SqlDB, storage.NoSqlDB)
	}{
		{
			name:            "01",
			ctx:             context.Background(),
			expectedErrText: "",
			prepare: func(ctrl *gomock.Controller) (storage.SqlDB, storage.NoSqlDB) {
				sqlDB := mock.NewMockSqlDB(ctrl)
				sqlDB.EXPECT().
					List(gomock.Any()).
					Return([]internal.Product{internal.Product{Id: 15, Name: "product-15", Qty: 1, Price: 150}}, nil)
				noSqlDB := mock.NewMockNoSqlDB(ctrl)
				noSqlDB.EXPECT().
					ListCreate(gomock.Any(), []internal.Product{internal.Product{Id: 15, Name: "product-15", Qty: 1, Price: 150}}).
					Return(nil)
				sqlDB.EXPECT().
					List(gomock.Any()).
					Return([]internal.Product{internal.Product{Id: 15, Name: "product-15", Qty: 1, Price: 150}}, nil)
				noSqlDB.EXPECT().
					ListCreate(gomock.Any(), []internal.Product{internal.Product{Id: 15, Name: "product-15", Qty: 1, Price: 150}}).
					Return(nil)
				return sqlDB, noSqlDB
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sql, noSql := test.prepare(ctrl)
			storage, err := storage.New(context.TODO(), sql, noSql)
			assert.NoError(t, err)
			err = storage.RenewCache(test.ctx)
			if test.expectedErrText != "" {
				assert.EqualError(t, err, test.expectedErrText)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
