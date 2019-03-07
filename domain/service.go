package domain

type Config struct {
	SeuilRappelEcheances float32 `default:"0.2"`
}

type Service struct {
	config          Config
	repo            Repository
	storage         Storage
	templateService TemplateService
	emailService    EmailService
}

func New(config Config, repo Repository, storage Storage, templateService TemplateService, emailService EmailService) *Service {
	return &Service{
		config:          config,
		repo:            repo,
		storage:         storage,
		templateService: templateService,
		emailService:    emailService,
	}
}
