package domain

type Service struct {
	repo            Repository
	storage         Storage
	templateService TemplateService
}

func New(repo Repository, storage Storage, templateService TemplateService) *Service {
	return &Service{
		repo:            repo,
		storage:         storage,
		templateService: templateService,
	}
}
