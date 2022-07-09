package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ffo32167/flowershop/internal/http"
	"github.com/ffo32167/flowershop/internal/storage"
	"github.com/ffo32167/flowershop/internal/storage/postgres"
	"github.com/ffo32167/flowershop/internal/storage/redis"
	"go.uber.org/zap"
)

func main() {
	/*
		http://localhost:8080/list
		http://localhost:8080/sale/1/1

	*/
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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	pg, err := postgres.New(ctx, os.Getenv("PG_CONN_STR"))
	if err != nil {
		log.Error("cant conn to postgres: ", zap.Error(err))
	}
	rd, err := redis.New(os.Getenv("REDIS_CONN_STR"))
	if err != nil {
		log.Error("cant conn to redis: ", zap.Error(err))
	}

	storage, err := storage.New(ctx, pg, rd)
	if err != nil {
		log.Error("storage create error: ", zap.Error(err))
	}

	apiServer := http.New(storage, os.Getenv("HTTP_PORT"), log)

	err = apiServer.Run()
	if err != nil {
		log.Error("cant start api server:", zap.Error(err))
	}
}
