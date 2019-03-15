package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (repo *Repository) ListThemes() ([]models.Theme, error) {
	themes := []models.Theme{}
	err := repo.db.client.Model(&themes).Order("nom asc").Select()
	return themes, err
}

func (repo *Repository) CreateTheme(ctx *domain.UserContext, theme models.Theme) (int64, error) {
	theme.Id = 0
	err := repo.db.client.Insert(&theme)
	if err != nil {
		return 0, err
	}
	err = repo.eventsManager.DispatchUpdatedResources(ctx, "themes")
	return theme.Id, err
}

func (repo *Repository) DeleteTheme(ctx *domain.UserContext, idTheme int64) error {
	theme := models.Theme{
		Id: idTheme,
	}
	err := repo.db.client.Delete(&theme)
	if err != nil {
		return err
	}
	err = repo.eventsManager.DispatchUpdatedResources(ctx, "themes")
	return err
}
