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
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func New(connStr string) (RedisDB, error) {
	const productsTableName = "product_list"
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
		return RedisDB{productsTableName: productsTableName}, err
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
	_, err = db.Rdb.HSet(ctx, db.productsTableName, redisProducts).Result()
	if err != nil {
		return fmt.Errorf("zadd failed, err: %w", err)
	}
	return nil
}

func ToRedisProductsList(products []internal.Product) (map[string]string, error) {
	redisProducts := make([]redisProduct, len(products))
	for i := range products {
		redisProducts[i].Id = products[i].Id
		redisProducts[i].Name = products[i].Name
		redisProducts[i].Quantity = products[i].Qty
		redisProducts[i].Price = products[i].Price
	}
	result := make(map[string]string)
	for i, val := range redisProducts {
		encoded, err := json.Marshal(val)
		if err != nil {
			return make(map[string]string), fmt.Errorf("marshalling failed, err on id %v, error: %w", val.Id, err)
		}
		result[strconv.Itoa(redisProducts[i].Id)] = string(encoded[:])
	}
	return result, nil
}

func (db RedisDB) List(ctx context.Context) ([]internal.Product, error) {
	redisProducts, err := db.Rdb.HGetAll(ctx, db.productsTableName).Result()
	if err != nil {
		return []internal.Product{}, fmt.Errorf("cant hgetall:%w", err)
	}
	result, err := toDomain(redisProducts)
	if err != nil {
		return nil, fmt.Errorf("cant convert redis format to domain:%w", err)
	}
	return result, nil
}

func (db RedisDB) Sale(id, cnt int) (int, error) {
	return 0, nil
}

func toDomain(p map[string]string) ([]internal.Product, error) {
	redisProducts := make([]redisProduct, len(p))
	var i int
	var prod redisProduct
	for _, val := range p {
		err := json.Unmarshal([]byte(val), &prod)
		if err != nil {
			return []internal.Product{}, fmt.Errorf("cant unmarshal redis data:%w", err)
		}
		redisProducts[i] = prod
		i++
	}
	products := make([]internal.Product, len(redisProducts))
	for i := range redisProducts {
		products[i].Id = redisProducts[i].Id
		products[i].Name = redisProducts[i].Name
		products[i].Qty = redisProducts[i].Quantity
		products[i].Price = redisProducts[i].Price
	}
	return products, nil
}
