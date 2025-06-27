package services

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/fesbarbosa/melivendas-api/internal/core/domain"
	"github.com/fesbarbosa/melivendas-api/internal/core/ports/output"
)

var (
	ErrItemNotFound = errors.New("item não encontrado")

	ErrDuplicateCode = errors.New("um item com este código já existe")

	ErrInvalidData = errors.New("dados do item inválidos")
)

type ItemService struct {
	repo output.ItemRepository
}

func NewItemService(repo output.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

func (s *ItemService) CreateItem(ctx context.Context, code, title, description string, price, stock int64) (*domain.Item, error) {

	if code == "" || title == "" || description == "" {
		return nil, fmt.Errorf("%w: código, título e descrição são obrigatórios", ErrInvalidData)
	}

	if price <= 0 {
		return nil, fmt.Errorf("%w: preço deve ser maior que 0", ErrInvalidData)
	}

	if stock < 0 {
		return nil, fmt.Errorf("%w: estoque não pode ser negativo", ErrInvalidData)
	}

	exists, err := s.repo.ExistsByCode(ctx, code, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar unicidade do código: %w", err)
	}
	if exists {
		return nil, ErrDuplicateCode
	}

	item := domain.NewItem(code, title, description, price, stock)

	savedItem, err := s.repo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar item: %w", err)
	}

	return savedItem, nil
}

func (s *ItemService) GetItem(ctx context.Context, id int64) (*domain.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter item: %w", err)
	}

	if item == nil {
		return nil, ErrItemNotFound
	}

	return item, nil
}

func (s *ItemService) UpdateItem(ctx context.Context, id int64, code, title, description string, price, stock int64) (*domain.Item, error) {

	if code == "" || title == "" || description == "" {
		return nil, fmt.Errorf("%w: código, título e descrição são obrigatórios", ErrInvalidData)
	}

	if price <= 0 {
		return nil, fmt.Errorf("%w: preço deve ser maior que 0", ErrInvalidData)
	}

	if stock < 0 {
		return nil, fmt.Errorf("%w: estoque não pode ser negativo", ErrInvalidData)
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter item: %w", err)
	}

	if item == nil {
		return nil, ErrItemNotFound
	}

	if item.Code != code {
		exists, err := s.repo.ExistsByCode(ctx, code, id)
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar unicidade do código: %w", err)
		}
		if exists {
			return nil, ErrDuplicateCode
		}
	}

	item.UpdateItem(code, title, description, price, stock)

	err = s.repo.Update(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar item: %w", err)
	}

	return item, nil
}

func (s *ItemService) DeleteItem(ctx context.Context, id int64) error {

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting item: %w", err)
	}

	if item == nil {
		return ErrItemNotFound
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir item: %w", err)
	}

	return nil
}

func (s *ItemService) ListItems(ctx context.Context, status string, limit, page int) (*domain.PagedItems, error) {

	if limit <= 0 {
		limit = 10
	} else if limit > 20 {
		limit = 20
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	total, err := s.repo.Count(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("erro ao contar itens: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	items, err := s.repo.FindAll(ctx, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao recuperar itens: %w", err)
	}

	data := make([]domain.Item, 0, len(items))
	for _, item := range items {
		data = append(data, *item)
	}

	return &domain.PagedItems{
		TotalPaginas: totalPages,
		Dados:        data,
	}, nil
}
