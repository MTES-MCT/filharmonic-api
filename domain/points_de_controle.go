package domain

import "github.com/MTES-MCT/filharmonic-api/models"

func (s *Service) CreatePointDeControle(ctx *UserContext, idInspection int64, pointDeControle models.PointDeControle) (int64, error) {
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

	return s.repo.CreatePointDeControle(ctx, idInspection, pointDeControle)
}

func (s *Service) UpdatePointDeControle(ctx *UserContext, idPointDeControle int64, pointDeControle models.PointDeControle) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}

	return s.repo.UpdatePointDeControle(ctx, idPointDeControle, pointDeControle)
}

func (s *Service) DeletePointDeControle(ctx *UserContext, idPointDeControle int64) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}

	return s.repo.DeletePointDeControle(ctx, idPointDeControle)
}

func (s *Service) PublishPointDeControle(ctx *UserContext, idPointDeControle int64) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}

	return s.repo.PublishPointDeControle(ctx, idPointDeControle)
}
