package domain

type Service struct {
	repo    Repository
	storage Storage
}

func New(repo Repository, storage Storage, emails EmailService) *Service {
	return &Service{
		repo:    repo,
		storage: storage,
	}
}
