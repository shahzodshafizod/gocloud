//go:generate mockgen -source=cache.go -package=mocks -destination=mocks/cache.go
package pkg

import (
	"context"
	"time"
)

type Cache interface {
	// Saves a string value in the cache with a specified key and expiration time.
	SaveString(ctx context.Context, key string, value string, expiration time.Duration) error
	// Retrieves a string value from the cache using the provided key.
	GetString(ctx context.Context, key string) (string, error)
	// Stores a struct or any serializable data in the cache with a key and expiration time.
	SaveStruct(ctx context.Context, key string, v any, expiration time.Duration) error
	// Retrieves and deserializes a struct or data from the cache into the provided variable.
	GetStruct(ctx context.Context, key string, v any) error
	// Deletes a cached value associated with the specified key.
	Del(ctx context.Context, key string) error
}
