package domain

import (
	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
)

var (
	ErrCreationPointDeControleImpossible     = errors.NewErrForbidden("Impossible de créer le point de contrôle")
	ErrModificationPointDeControleImpossible = errors.NewErrForbidden("Impossible de modifier le point de contrôle")
)

func (s *Service) canCreatePointDeControle(ctx *UserContext, idInspection int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}

	err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return err
	}
	return s.repo.CanCreatePointDeControle(ctx, idInspection)
}

func (s *Service) canUpdatePointDeControle(ctx *UserContext, idPointDeControle int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}

	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	return s.repo.CanUpdatePointDeControle(ctx, idPointDeControle)
}

func (s *Service) CreatePointDeControle(ctx *UserContext, idInspection int64, pointDeControle models.PointDeControle) (int64, error) {
	err := s.canCreatePointDeControle(ctx, idInspection)
	if err != nil {
		return 0, err
	}
	return s.repo.CreatePointDeControle(ctx, idInspection, pointDeControle)
}

func (s *Service) UpdatePointDeControle(ctx *UserContext, idPointDeControle int64, pointDeControle models.PointDeControle) error {
	err := s.canUpdatePointDeControle(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	return s.repo.UpdatePointDeControle(ctx, idPointDeControle, pointDeControle)
}

func (s *Service) DeletePointDeControle(ctx *UserContext, idPointDeControle int64) error {
	err := s.canUpdatePointDeControle(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	return s.repo.DeletePointDeControle(ctx, idPointDeControle)
}

func (s *Service) PublishPointDeControle(ctx *UserContext, idPointDeControle int64) error {
	err := s.canUpdatePointDeControle(ctx, idPointDeControle)
	if err != nil {
		return err
	}
	return s.repo.PublishPointDeControle(ctx, idPointDeControle)
}

func (s *Service) OrderPointsDeControle(ctx *UserContext, idInspection int64, pointsDeControleIds []int64) error {
	return s.repo.OrderPointsDeControle(ctx, idInspection, pointsDeControleIds)
}
