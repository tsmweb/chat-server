package db

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheDB interface with methods for manipulating data in the cache.
type CacheDB interface {
	// Set Redis `SET key value [expiration]` command.
	// Zero expiration means the key has no expiration time.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get Redis `GET key` command.
	Get(ctx context.Context, key string) (string, error)

	// Del Redis `DEL key [key ...]` command.
	Del(ctx context.Context, keys ...string) error

	// HSet accepts values in following formats:
	//   - HSet(ctx, "myhash", "key1", "value1", "key2", "value2")
	//   - HSet(ctx, "myhash", []string{"key1", "value1", "key2", "value2"})
	//   - HSet(ctx, "myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
	HSet(ctx context.Context, key string, values ...interface{}) error

	// HGet Redis `HGET key field` command.
	HGet(ctx context.Context, key, field string) (string, error)

	// HDel Redis `HDEL key field` command.
	HDel(ctx context.Context, key string, fields ...string) error

	// SAdd Redis `SADD key member` command.
	SAdd(ctx context.Context, key string, members ...interface{}) error

	// SMembers Redis `SMEMBERS key` command.
	SMembers(ctx context.Context, key string) ([]string, error)

	// Expire Redis `EXPIRE key expiration` command.
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// RedisCacheDB implements the CacheDB interface.
type RedisCacheDB struct {
	db *redis.Client
}

// NewRedisCacheDB returns an instance of CacheDB.
func NewRedisCacheDB(addr string, password string) CacheDB {
	db := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	return &RedisCacheDB{db: db}
}

func (c *RedisCacheDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	_, err := c.db.Set(ctx, key, value, expiration).Result()
	return err
}

func (c *RedisCacheDB) Get(ctx context.Context, key string) (string, error) {
	val, err := c.db.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", nil
		}
		return "", err
	}

	return val, nil
}

func (c *RedisCacheDB) Del(ctx context.Context, keys ...string) error {
	_, err := c.db.Del(ctx, keys...).Result()
	return err
}

func (c *RedisCacheDB) HSet(ctx context.Context, key string, values ...interface{}) error {
	_, err := c.db.HSet(ctx, key, values...).Result()
	return err
}

func (c *RedisCacheDB) HGet(ctx context.Context, key, field string) (string, error) {
	return c.db.HGet(ctx, key, field).Result()
}

func (c *RedisCacheDB) HDel(ctx context.Context, key string, fields ...string) error {
	_, err := c.db.HDel(ctx, key, fields...).Result()
	return err
}

func (c *RedisCacheDB) SAdd(ctx context.Context, key string, members ...interface{}) error {
	_, err := c.db.SAdd(ctx, key, members...).Result()
	return err
}

func (c *RedisCacheDB) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.db.SMembers(ctx, key).Result()
}

func (c *RedisCacheDB) Expire(ctx context.Context, key string, expiration time.Duration) error {
	_, err := c.db.Expire(ctx, key, expiration).Result()
	return err
}
