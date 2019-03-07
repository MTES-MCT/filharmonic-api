package domain

type Service struct {
	repo            Repository
	storage         Storage
	templateService TemplateService
	emailService    EmailService
}

func New(repo Repository, storage Storage, templateService TemplateService, emailService EmailService) *Service {
	return &Service{
		repo:            repo,
		storage:         storage,
		templateService: templateService,
		emailService:    emailService,
	}
}
