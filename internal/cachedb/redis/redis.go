package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ffo32167/flowershop/internal"
	"github.com/go-redis/redis/v9"
)

type RedisDB struct {
	Rdb               *redis.Client
	productsTableName string
}

type redisProduct struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity string  `json:"quantity"`
	Price    float64 `json:"price"`
}

const productsTableName = "product_list"

func New(connStr string) (RedisDB, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     connStr,
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 10,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return RedisDB{}, err
	}
	return RedisDB{Rdb: rdb, productsTableName: productsTableName}, nil
}

func (db RedisDB) ListCreate(ctx context.Context, products []internal.Product) error {
	redisProducts, err := ToRedisProductsList(products)
	if err != nil {
		return fmt.Errorf("unable to convert products list to redis list: %w ", err)
	}
	_, err = db.Rdb.Del(ctx, db.productsTableName).Result()
	if err != nil {
		return fmt.Errorf("cant delete set: %w", err)
	}
	num, err := db.Rdb.HSet(ctx, db.productsTableName, redisProducts).Result()
	if err != nil {
		return fmt.Errorf("zadd failed, err: %w", err)
	}
	fmt.Printf("zadd %d succ.\n", num)
	return nil
}

func ToRedisProductsList(products []internal.Product) (map[string]string, error) {
	redisProducts := make([]redisProduct, len(products))
	for i := range products {
		redisProducts[i].Id = products[i].Id
		redisProducts[i].Name = products[i].Name
		redisProducts[i].Quantity = internal.TranslateQtyToStr(products[i].Qty)
		redisProducts[i].Price = products[i].Price
	}
	result := make(map[string]string)
	for i, val := range redisProducts {
		encoded, err := json.Marshal(val)
		if err != nil {
			return make(map[string]string), fmt.Errorf("marshalling failed, err: %w", err)
		}
		result[strconv.Itoa(redisProducts[i].Id)] = string(encoded[:])
	}
	return result, nil
}

func (db RedisDB) List(ctx context.Context) ([]internal.Product, error) {
	//	iter, err := db.Rdb.HGetAll(ctx, listName).Result()
	redisList := make([]string, 0)
	/*	for iter.Next(ctx) {
			redisList = append(redisList, iter.Val())
		}
	*/result, err := toDomain(redisList)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db RedisDB) Sale(id, cnt int) (int, error) {
	return 0, nil
}

func toDomain(p []string) ([]internal.Product, error) {
	result := make([]internal.Product, len(p)/2)
	var j int
	var err error
	for i := 0; i < len(p)/2; i++ {
		result[i].Id, err = strconv.Atoi(p[j])
		if err != nil {
			return nil, err
		}
		j++
		result[i].Name = p[j]
		j++
	}
	return result, nil
}
