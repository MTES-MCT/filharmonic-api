package domain

import "github.com/MTES-MCT/filharmonic-api/models"

func (s *Service) ListCanevas(ctx *UserContext) ([]models.Canevas, error) {
	if ctx.IsExploitant() {
		return nil, ErrBesoinProfilInspecteur
	}
	return s.repo.ListCanevas()
}

func (s *Service) CreateCanevas(ctx *UserContext, idInspection int64, canevas models.Canevas) (int64, error) {
	if ctx.IsExploitant() {
		return 0, ErrBesoinProfilInspecteur
	}
	return s.repo.CreateCanevas(ctx, idInspection, canevas)
}

func (s *Service) DeleteCanevas(ctx *UserContext, idCanevas int64) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}
	return s.repo.DeleteCanevas(ctx, idCanevas)
}
