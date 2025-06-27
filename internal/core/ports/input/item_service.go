package input

import (
	"context"

	"github.com/fesbarbosa/melivendas-api/internal/core/domain"
)

type ItemService interface {
	CreateItem(ctx context.Context, code, title, description string, price, stock int64) (*domain.Item, error)

	GetItem(ctx context.Context, id int64) (*domain.Item, error)

	UpdateItem(ctx context.Context, id int64, code, title, description string, price, stock int64) (*domain.Item, error)

	DeleteItem(ctx context.Context, id int64) error

	ListItems(ctx context.Context, status string, limit, page int) (*domain.PagedItems, error)
}
