package domain

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListEtablissements(ctx *UserContext, s3ic string) ([]models.Etablissement, error) {
	return s.repo.FindEtablissementsByS3IC(ctx, s3ic)
}

func (s *Service) GetEtablissement(ctx *UserContext, id int64) (*models.Etablissement, error) {
	return s.repo.GetEtablissementByID(ctx, id)
}

func (s *Service) ListInspections(ctx *UserContext) ([]models.Inspection, error) {
	return s.repo.ListInspections(ctx)
}

func (s *Service) GetInspection(ctx *UserContext, id int64) (*models.Inspection, error) {
	return s.repo.GetInspectionByID(ctx, id)
}
