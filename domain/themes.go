package domain

import "github.com/MTES-MCT/filharmonic-api/models"

func (s *Service) ListThemes(ctx *UserContext) ([]models.Theme, error) {
	return s.repo.ListThemes()
}

func (s *Service) CreateTheme(ctx *UserContext, theme models.Theme) (int64, error) {
	if ctx.IsExploitant() {
		return 0, ErrBesoinProfilInspecteur
	}
	return s.repo.CreateTheme(theme)
}

func (s *Service) DeleteTheme(ctx *UserContext, idTheme int64) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}
	return s.repo.DeleteTheme(idTheme)
}
