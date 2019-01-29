package domain

import "github.com/MTES-MCT/filharmonic-api/models"

func (s *Service) CreateSuite(ctx *UserContext, idInspection int64, suite models.Suite) (int64, error) {
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
	ok, err = s.repo.CheckCanCreateSuite(ctx, idInspection)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}
	return s.repo.CreateSuite(ctx, idInspection, suite)
}

func (s *Service) UpdateSuite(ctx *UserContext, idInspection int64, suite models.Suite) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	ok, err = s.repo.CheckCanDeleteSuite(ctx, idInspection)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	return s.repo.UpdateSuite(ctx, idInspection, suite)
}

func (s *Service) DeleteSuite(ctx *UserContext, idInspection int64) error {
	if ctx.IsExploitant() {
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
