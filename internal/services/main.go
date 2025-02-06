package services

type Service[T any] interface {
	execute(*T) error
}
