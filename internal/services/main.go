package services

type Service[T any] interface {
	execute(*T) error
}

type ServiceContext struct {
	RegistrationService
	ReplyService
}

func NewServiceContext(regS *RegistrationService, replyS *ReplyService) *ServiceContext {
	return &ServiceContext{
		RegistrationService: *regS,
		ReplyService:        *replyS,
	}
}
