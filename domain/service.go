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

func (s *Service) CreateInspection(ctx *UserContext, inspection models.Inspection) (int64, error) {
	if ctx.IsExploitant() {
		return 0, ErrBesoinProfilInspecteur
	}
	inspecteursIds := make([]int64, 0)
	for _, inspecteur := range inspection.Inspecteurs {
		inspecteursIds = append(inspecteursIds, inspecteur.Id)
	}
	ok, err := s.repo.CheckUsersInspecteurs(inspecteursIds)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}
	return s.repo.CreateInspection(ctx, inspection)
}

func (s *Service) GetInspection(ctx *UserContext, id int64) (*models.Inspection, error) {
	return s.repo.GetInspectionByID(ctx, id)
}

func (s *Service) UpdateInspection(ctx *UserContext, inspection models.Inspection) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}
	inspecteursIds := make([]int64, 0)
	for _, inspecteur := range inspection.Inspecteurs {
		inspecteursIds = append(inspecteursIds, inspecteur.Id)
	}
	if len(inspecteursIds) == 0 {
		return ErrInvalidInput
	}
	ok, err := s.repo.CheckUsersInspecteurs(inspecteursIds)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	return s.repo.UpdateInspection(ctx, inspection)
}

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
