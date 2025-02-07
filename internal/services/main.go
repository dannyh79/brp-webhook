package services

type Service[T any] interface {
	execute(*T) error
}

type ServiceContext struct {
	RegistrationService
}

func NewServiceContext(rs *RegistrationService) *ServiceContext {
	return &ServiceContext{
		RegistrationService: *rs,
	}
}
