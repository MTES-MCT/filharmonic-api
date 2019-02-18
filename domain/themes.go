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

func (s *Service) addMissingThemes(nomThemes []string) error {
	existingThemes, err := s.repo.ListThemes()
	if err != nil {
		return err
	}
	for _, nomTheme := range nomThemes {
		found := false
		for _, existingTheme := range existingThemes {
			if nomTheme == existingTheme.Nom {
				found = true
				break
			}
		}
		if !found {
			_, err = s.repo.CreateTheme(models.Theme{
				Nom: nomTheme,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
