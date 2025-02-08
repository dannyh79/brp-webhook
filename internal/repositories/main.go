package repositories

import "errors"

type Repository[T any] interface {
	Save(*T) (*T, error)
	Destroy(*T) error
}

// Record already exists.
var ErrorAlreadyExists = errors.New("Record already exists")

// Record not found.
var ErrorNotFound = errors.New("Record not found")
