package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Status   int    `json:"status"`
	Codigo   string `json:"codigo"`
	Mensagem string `json:"mensagem"`
}

var (
	ErrNotFound       = errors.New("recurso não encontrado")
	ErrBadRequest     = errors.New("requisição inválida")
	ErrConflict       = errors.New("conflito com dados existentes")
	ErrInternalServer = errors.New("erro interno do servidor")
)

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

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Codigo, e.Mensagem)
}
