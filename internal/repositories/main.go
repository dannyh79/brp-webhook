package repositories

type Repository[T any] interface {
	Save(*T) (*T, error)
}
