package domain

import (
	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
)

var (
	ErrCreationConstatImpossible    = errors.NewErrForbidden("Impossible de créer le constat")
	ErrSuppressionConstatImpossible = errors.NewErrForbidden("Impossible de supprimer le constat")
)

func (s *Service) CreateConstat(ctx *UserContext, idPointDeControle int64, constat models.Constat) (int64, error) {
	if !ctx.IsInspecteur() {
		return 0, ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}
	err = s.repo.CanCreateConstat(ctx, idPointDeControle)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateConstat(ctx, idPointDeControle, constat)
}

func (s *Service) DeleteConstat(ctx *UserContext, idPointDeControle int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}

	err := s.repo.CanDeleteConstat(ctx, idPointDeControle)
	if err != nil {
		return err
	}

	return s.repo.DeleteConstat(ctx, idPointDeControle)
}
