package domain

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

type ListEvenementsFilter struct {
}

func (s *Service) ListEvenements(ctx *UserContext, filter ListEvenementsFilter) ([]models.Evenement, error) {
	return s.repo.ListEvenements(ctx, filter)
}

func (s *Service) GetEvenement(ctx *UserContext, id int64) (*models.Evenement, error) {
	return s.repo.GetEvenementByID(ctx, id)
}
