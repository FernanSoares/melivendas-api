package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fesbarbosa/melivendas-api/internal/core/domain"
	"github.com/jmoiron/sqlx"
)

type ItemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	query := `
		INSERT INTO items (code, title, description, price, stock, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		item.Code,
		item.Title,
		item.Description,
		item.Price,
		item.Stock,
		item.Status,
		item.CreatedAt,
		item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	item.ID = id
	return item, nil
}

func (r *ItemRepository) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	query := "SELECT * FROM items WHERE id = ?"

	var item domain.Item
	err := r.db.GetContext(ctx, &item, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepository) Update(ctx context.Context, item *domain.Item) error {
	query := `
		UPDATE items
		SET code = ?, title = ?, description = ?, price = ?, stock = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		item.Code,
		item.Title,
		item.Description,
		item.Price,
		item.Stock,
		item.Status,
		item.UpdatedAt,
		item.ID,
	)

	return err
}

func (r *ItemRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM items WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ItemRepository) FindAll(ctx context.Context, status string, limit, offset int) ([]*domain.Item, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = "SELECT * FROM items WHERE status = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?"
		args = []interface{}{status, limit, offset}
	} else {
		query = "SELECT * FROM items ORDER BY updated_at DESC LIMIT ? OFFSET ?"
		args = []interface{}{limit, offset}
	}

	items := []*domain.Item{}
	err := r.db.SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) Count(ctx context.Context, status string) (int, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = "SELECT COUNT(*) FROM items WHERE status = ?"
		args = []interface{}{status}
	} else {
		query = "SELECT COUNT(*) FROM items"
		args = []interface{}{}
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ItemRepository) ExistsByCode(ctx context.Context, code string, excludeID int64) (bool, error) {
	var query string
	var args []interface{}

	if excludeID > 0 {
		query = "SELECT COUNT(*) FROM items WHERE code = ? AND id != ?"
		args = []interface{}{code, excludeID}
	} else {
		query = "SELECT COUNT(*) FROM items WHERE code = ?"
		args = []interface{}{code}
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
