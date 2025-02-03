//go:generate mockgen -source=nosql.go -package=mocks -destination=mocks/nosql.go
package pkg

import (
	"context"
)

type NoSQL interface {
	// Inserts a new item into the specified table and returns the ID of the inserted item or an error.
	Insert(ctx context.Context, table string, item Map) (string, error)
	// Retrieves a single item from the specified table using the provided keys.
	GetItem(ctx context.Context, table string, keys Map) (Map, error)
	// Fetches multiple items from the specified table based on a filter condition.
	GetItems(ctx context.Context, table string, filter Map) ([]Map, error)
	// Updates items in the specified table that match the filter, applying the provided update, and returns the updated item or an error.
	Update(ctx context.Context, table string, filter Map, update Map) (Map, error)
}

type Map map[string]any
