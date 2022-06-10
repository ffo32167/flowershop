package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ffo32167/flowershop/internal"
	"github.com/ffo32167/flowershop/internal/postgres"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("hello, flower shop!")

	log, err := zap.NewProduction()
	if err != nil {
		fmt.Println(fmt.Errorf("cant start logger: %w", err))
	}
	defer func() {
		err = log.Sync()
		if err != nil {
			fmt.Println(fmt.Errorf("cant sync logger: %w", err))
		}
	}()
	var storage internal.Storage
	storage, err = postgres.New(context.Background(), os.Getenv("PG_CONN_STR"))
	if err != nil {
		log.Error("storage create error: ", zap.Error(err))
	}
	list, err := storage.List()
	if err != nil {
		log.Error("cant get list: ", zap.Error(err))
	}
	fmt.Println(list)
}
