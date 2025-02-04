package onprem

import (
	"bytes"
	"context"
	"encoding/gob"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type cache struct {
	client *redis.Client
}

func NewCache() (pkg.Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, errors.Wrap(err, "rdb.Ping")
	}

	return &cache{client: rdb}, nil
}

func (c *cache) SaveString(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := c.client.Set(ctx, key, value, expiration).Err()
	if c.notNil(err) {
		return errors.Wrap(err, "c.client.Set")
	}
	return nil
}

func (c *cache) GetString(ctx context.Context, key string) (string, error) {
	var value string
	err := c.client.Get(ctx, key).Scan(&value)
	if c.notNil(err) {
		return "", errors.Wrap(err, "c.client.Get.Scan")
	}
	return value, nil
}

func (c *cache) SaveStruct(ctx context.Context, key string, value any, expiration time.Duration) error {
	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(value)
	if err != nil {
		return errors.Wrap(err, "gob.NewEncoder.Encode")
	}
	err = c.client.Set(ctx, key, buffer.Bytes(), expiration).Err()
	if c.notNil(err) {
		return errors.Wrap(err, "c.client.Set")
	}
	return nil
}

func (c *cache) GetStruct(ctx context.Context, key string, value any) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if c.notNil(err) {
		return errors.Wrap(err, "c.client.Get.Bytes")
	}
	return gob.NewDecoder(bytes.NewReader(data)).Decode(value)
}

func (c *cache) Del(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if c.notNil(err) {
		return errors.Wrap(err, "c.client.Del")
	}
	return nil
}

func (c *cache) notNil(err error) bool {
	return err != nil && !errors.Is(err, redis.Nil)
}
