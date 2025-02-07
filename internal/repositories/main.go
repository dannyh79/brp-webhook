package repositories

import "errors"

type Repository[T any] interface {
	Save(*T) (*T, error)
}

// Record already exists.
var ErrorAlreadyExists = errors.New("Record already exists")
