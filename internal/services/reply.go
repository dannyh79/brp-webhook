package services

type ReplyService struct{}

func (s *ReplyService) Execute(to string) error {
	return nil
}

func NewReplyService() *ReplyService {
	return &ReplyService{}
}
