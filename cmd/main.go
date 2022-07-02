package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ffo32167/flowershop/internal/cachedb"
	"github.com/ffo32167/flowershop/internal/cachedb/postgres"
	"github.com/ffo32167/flowershop/internal/cachedb/redis"
	"go.uber.org/zap"
)

func main() {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pg, err := postgres.New(context.Background(), os.Getenv("PG_CONN_STR"))
	if err != nil {
		log.Error("cant conn to postgres: ", zap.Error(err))
	}
	rd, err := redis.New(os.Getenv("REDIS_CONN_STR"))
	if err != nil {
		log.Error("cant conn to redis: ", zap.Error(err))
	}

	cdb, err := cachedb.New(ctx, pg, rd)
	if err != nil {
		log.Error("storage create error: ", zap.Error(err))
	}

	res, err := cdb.List()
	if err != nil {
		fmt.Println("list err:", err)
	}
	for _, val := range res {
		fmt.Println(val)
	}

}
