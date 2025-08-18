package service

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already in use")
	ErrCGCAlreadyExists   = errors.New("cgc already in use")
	EntityNotFound        = errors.New("entity not found")
)
