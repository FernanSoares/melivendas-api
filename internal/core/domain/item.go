package domain

import (
	"time"
)

type ItemStatus string

const (
	ItemStatusActive ItemStatus = "ACTIVE"

	ItemStatusInactive ItemStatus = "INACTIVE"
)

type Item struct {
	ID          int64      `json:"id" db:"id"`
	Code        string     `json:"code" db:"code"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Price       int64      `json:"price" db:"price"`
	Stock       int64      `json:"stock" db:"stock"`
	Status      ItemStatus `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

func NewItem(code, title, description string, price, stock int64) *Item {
	status := ItemStatusActive
	if stock == 0 {
		status = ItemStatusInactive
	}

	now := time.Now()

	return &Item{
		Code:        code,
		Title:       title,
		Description: description,
		Price:       price,
		Stock:       stock,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (i *Item) UpdateStock(stock int64) {
	i.Stock = stock
	if stock == 0 {
		i.Status = ItemStatusInactive
	} else {
		i.Status = ItemStatusActive
	}
	i.UpdatedAt = time.Now()
}

func (i *Item) UpdateItem(code, title, description string, price, stock int64) {
	i.Code = code
	i.Title = title
	i.Description = description
	i.Price = price
	i.UpdateStock(stock)
}

type PagedItems struct {
	TotalPaginas int    `json:"totalPaginas"`
	Dados        []Item `json:"dados"`
}
