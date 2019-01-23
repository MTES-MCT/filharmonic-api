package domain

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (s *Service) CreateCommentaire(ctx *UserContext, idInspection int64, commentaire models.Commentaire) (int64, error) {
	if ctx.IsExploitant() {
		return 0, ErrBesoinProfilInspecteur
	}
	ok, err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}

	return s.repo.CreateCommentaire(ctx, idInspection, commentaire)
}
