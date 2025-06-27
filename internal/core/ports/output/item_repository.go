package output

import (
	"context"

	"github.com/fesbarbosa/melivendas-api/internal/core/domain"
)

type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) (*domain.Item, error)

	GetByID(ctx context.Context, id int64) (*domain.Item, error)

	Update(ctx context.Context, item *domain.Item) error

	Delete(ctx context.Context, id int64) error

	FindAll(ctx context.Context, status string, limit, offset int) ([]*domain.Item, error)

	Count(ctx context.Context, status string) (int, error)

	ExistsByCode(ctx context.Context, code string, excludeID int64) (bool, error)
}
