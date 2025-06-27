package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError representa um erro padronizado da API
type APIError struct {
	Status  int    `json:"status"`
	Codigo  string `json:"codigo"`
	Mensagem string `json:"mensagem"`
}

var (
	// ErrNotFound representa um erro de recurso não encontrado
	ErrNotFound = errors.New("recurso não encontrado")

	// ErrBadRequest representa um erro de requisição inválida
	ErrBadRequest = errors.New("requisição inválida")

	// ErrConflict representa um conflito com dados existentes
	ErrConflict = errors.New("conflito com dados existentes")

	// ErrInternalServer representa um erro inesperado do servidor
	ErrInternalServer = errors.New("erro interno do servidor")
)

// NewAPIError creates a new API error from an error
func NewAPIError(err error) *APIError {
	if err == nil {
		return nil
	}

	var status int
	var code string

	switch {
	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
		code = "NAO_ENCONTRADO"
	case errors.Is(err, ErrBadRequest):
		status = http.StatusBadRequest
		code = "REQUISICAO_INVALIDA"
	case errors.Is(err, ErrConflict):
		status = http.StatusConflict
		code = "CONFLITO"
	default:
		status = http.StatusInternalServerError
		code = "ERRO_INTERNO_SERVIDOR"
	}

	return &APIError{
		Status:   status,
		Codigo:   code,
		Mensagem: err.Error(),
	}
}

// Error retorna a mensagem de erro
func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Codigo, e.Mensagem)
}
