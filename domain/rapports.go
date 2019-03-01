package domain

import (
	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
)

var (
	ErrRapportNotFound = errors.NewErrBadInput("Rapport non trouvé")
)

func (s *Service) GetRapport(ctx *UserContext, idInspection int64) (*models.File, error) {
	rapport, err := s.repo.GetRapport(ctx, idInspection)
	if err != nil {
		return nil, err
	}
	if rapport == nil {
		return nil, ErrRapportNotFound
	}
	reader, err := s.storage.Get(rapport.StorageId)
	if err != nil {
		return nil, err
	}
	return &models.File{
		Nom:     rapport.Nom,
		Type:    rapport.Type,
		Taille:  rapport.Taille,
		Content: reader,
	}, nil
}
