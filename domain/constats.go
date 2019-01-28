package domain

import "github.com/MTES-MCT/filharmonic-api/models"

func (s *Service) CreateConstat(ctx *UserContext, idPointDeControle int64, constat models.Constat) (int64, error) {
	if ctx.IsExploitant() {
		return 0, ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}
	ok, err = s.repo.CheckCanCreateConstat(ctx, idPointDeControle)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}
	return s.repo.CreateConstat(ctx, idPointDeControle, constat)
}

func (s *Service) DeleteConstat(ctx *UserContext, idPointDeControle int64) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckCanDeleteConstat(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}

	return s.repo.DeleteConstat(ctx, idPointDeControle)
}
