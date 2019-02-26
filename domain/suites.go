package domain

import (
	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
)

var (
	ErrInspecteurNonAffecte    = errors.NewErrForbidden("Utilisateur non affecté à cette inspection")
	ErrCreationSuiteImpossible = errors.NewErrForbidden("Impossible de créer la suite")
)

func (s *Service) CreateSuite(ctx *UserContext, idInspection int64, suite models.Suite) (int64, error) {
	if !ctx.IsInspecteur() {
		return 0, ErrBesoinProfilInspecteur
	}

	err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return 0, err
	}
	err = s.repo.CheckCanCreateSuite(ctx, idInspection)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateSuite(ctx, idInspection, suite)
}

func (s *Service) UpdateSuite(ctx *UserContext, idInspection int64, suite models.Suite) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}

	err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return err
	}
	ok, err := s.repo.CheckCanDeleteSuite(ctx, idInspection)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	return s.repo.UpdateSuite(ctx, idInspection, suite)
}

func (s *Service) DeleteSuite(ctx *UserContext, idInspection int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckCanDeleteSuite(ctx, idInspection)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}

	return s.repo.DeleteSuite(ctx, idInspection)
}
