package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fesbarbosa/melivendas-api/internal/core/domain"
	"github.com/fesbarbosa/melivendas-api/internal/core/ports/input"
	"github.com/fesbarbosa/melivendas-api/internal/core/services"
	apiErrors "github.com/fesbarbosa/melivendas-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	itemService input.ItemService
}

func NewItemHandler(itemService input.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

type ItemRequest struct {
	Code        string `json:"code" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	Stock       int64  `json:"stock" binding:"required,gte=0"`
}

type ItemResponse struct {
	Sucesso  bool        `json:"sucesso"`
	Mensagem string      `json:"mensagem,omitempty"`
	Dados    interface{} `json:"dados,omitempty"`
}

type PagedItems struct {
	TotalPaginas int           `json:"totalPaginas"`
	Dados        []domain.Item `json:"dados"`
}

func (h *ItemHandler) Create(c *gin.Context) {
	var req ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiErrors.NewAPIError(
			errors.Join(apiErrors.ErrBadRequest, err),
		))
		return
	}

	item, err := h.itemService.CreateItem(
		c.Request.Context(),
		req.Code,
		req.Title,
		req.Description,
		req.Price,
		req.Stock,
	)

	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, services.ErrDuplicateCode):
			statusCode = http.StatusConflict
		case errors.Is(err, services.ErrInvalidData):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"sucesso": false, "erro": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ItemResponse{
		Sucesso:  true,
		Mensagem: "Item criado com sucesso",
		Dados:    item,
	})
}

func (h *ItemHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"sucesso": false, "erro": "ID de item inválido"})
		return
	}

	item, err := h.itemService.GetItem(c.Request.Context(), id)
	if err != nil {
		var statusCode int
		if errors.Is(err, services.ErrItemNotFound) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"sucesso": false, "erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ItemResponse{
		Sucesso: true,
		Dados:   item,
	})
}

func (h *ItemHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"sucesso": false, "erro": "ID de item inválido"})
		return
	}

	var req ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiErrors.NewAPIError(
			errors.Join(apiErrors.ErrBadRequest, err),
		))
		return
	}

	item, err := h.itemService.UpdateItem(
		c.Request.Context(),
		id,
		req.Code,
		req.Title,
		req.Description,
		req.Price,
		req.Stock,
	)

	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, services.ErrItemNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, services.ErrDuplicateCode):
			statusCode = http.StatusConflict
		case errors.Is(err, services.ErrInvalidData):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"sucesso": false, "erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ItemResponse{
		Sucesso:  true,
		Mensagem: "Item atualizado com sucesso",
		Dados:    item,
	})
}

func (h *ItemHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"sucesso": false, "erro": "ID de item inválido"})
		return
	}

	err = h.itemService.DeleteItem(c.Request.Context(), id)
	if err != nil {
		var statusCode int
		if errors.Is(err, services.ErrItemNotFound) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"sucesso": false, "erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ItemResponse{
		Sucesso:  true,
		Mensagem: "Item excluído com sucesso",
	})
}

func (h *ItemHandler) List(c *gin.Context) {

	status := c.Query("status")

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 20 {
				limit = 20
			}
		}
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	result, err := h.itemService.ListItems(c.Request.Context(), status, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"sucesso": false, "erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
