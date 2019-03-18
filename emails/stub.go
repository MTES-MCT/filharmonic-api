package emails

type StubEmailService struct {
}

func NewStub() *StubEmailService {
	return &StubEmailService{}
}

func (em *StubEmailService) Send(email Email) error {
	return nil
}

func (em *StubEmailService) Close() error {
	return nil
}
