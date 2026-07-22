package model

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrDuplicate    = errors.New("resource already exists")
	ErrInternal     = errors.New("internal server error")
)
