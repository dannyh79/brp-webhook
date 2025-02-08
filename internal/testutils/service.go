package testutils

import "fmt"

type MockServiceInterface[T any] interface {
	Execute(*T) error
	CalledTimes() int
}

type MockService[T any] struct {
	calledTimes int
	shouldFail  bool
	err         error
}

func (s *MockService[T]) Execute(dto *T) error {
	s.calledTimes++
	if !s.shouldFail {
		return nil
	}
	if s.err == nil {
		s.err = fmt.Errorf("failed to execute mocked service")
	}
	return s.err
}

func (s *MockService[T]) CalledTimes() int {
	return s.calledTimes
}

func NewMockService[T any](shouldFail bool, errs ...error) MockServiceInterface[T] {
	var err error
	if len(errs) > 0 {
		err = errs[0]
	}
	return &MockService[T]{
		shouldFail: shouldFail,
		err:        err,
	}
}
