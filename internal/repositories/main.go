package repositories

import (
	"errors"
	"net/http"
)

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Repository[T any] interface {
	Save(*T) (*T, error)
	Destroy(*T) error
}

// Record already exists.
var ErrorAlreadyExists = errors.New("Record already exists")

// Record not found.
var ErrorNotFound = errors.New("Record not found")
