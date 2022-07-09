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
	Rdb                  *redis.Client
	productsTableName    string
	productsQtyTableName string
}

type redisProduct struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func New(connStr string) (RedisDB, error) {
	const (
		productsTableName    = "product_list"
		productsQtyTableName = "product_qty"
	)
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
	return RedisDB{Rdb: rdb, productsTableName: productsTableName, productsQtyTableName: productsQtyTableName}, nil
}

func (db RedisDB) ListCreate(ctx context.Context, products []internal.Product) error {
	redisProducts, err := toRedisProductsList(products)
	if err != nil {
		return fmt.Errorf("unable to convert products list to redis list: %w ", err)
	}
	redisProductsQty, err := toRedisProductsQtyList(products)
	if err != nil {
		return fmt.Errorf("unable to convert qty list to redis list: %w ", err)
	}
	txPipe := db.Rdb.TxPipeline()
	txPipe.Del(ctx, db.productsTableName)
	txPipe.HSet(ctx, db.productsTableName, redisProducts)
	txPipe.Del(ctx, db.productsQtyTableName)
	txPipe.HSet(ctx, db.productsQtyTableName, redisProductsQty)
	_, err = txPipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipe error: %w", err)
	}
	return nil
}

func toRedisProductsList(products []internal.Product) (map[string]string, error) {
	redisProducts := make([]redisProduct, len(products))
	for i := range products {
		redisProducts[i].Id = products[i].Id
		redisProducts[i].Name = products[i].Name
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

func toRedisProductsQtyList(products []internal.Product) (map[string]string, error) {
	result := make(map[string]string)
	for _, val := range products {
		result[strconv.Itoa(val.Id)] = strconv.Itoa(val.Qty)
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

func (db RedisDB) Sale(ctx context.Context, id, cnt int) (int, error) {
	res, err := db.Rdb.HIncrBy(ctx, db.productsQtyTableName, strconv.Itoa(id), -int64(cnt)).Result()
	if err != nil {
		return 0, fmt.Errorf("cant update qty in redis:%w", err)
	}
	return int(res), nil
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
		products[i].Price = redisProducts[i].Price
	}
	return products, nil
}
