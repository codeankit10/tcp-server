package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// redis interface
type RedisCache struct {
	ctx    context.Context
	client *redis.Client
}

// create new instance of redis

func NewRedisCache(ctx context.Context, host string, port int) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})

	// Test the connection to Redis.
	err := rdb.Set(ctx, "key", "value", 0).Err()

	return &RedisCache{
		ctx:    ctx,
		client: rdb,
	}, err
}

func (c *RedisCache) Add(key int, expiration int64) error {
	return c.client.Set(c.ctx, strconv.Itoa(key), "value", time.Duration(expiration*1e9)*time.Second).Err()
}

func (c *RedisCache) Get(key int) (bool, error) {
	val, err := c.client.Get(c.ctx, strconv.Itoa(key)).Result()
	return val != "", err
}

// Delete removes an item from the cache.
func (c *RedisCache) Delete(key int) {
	c.client.Del(c.ctx, strconv.Itoa(key))
}
