package database

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (repo *Repository) ListThemes() ([]models.Theme, error) {
	themes := []models.Theme{}
	err := repo.db.client.Model(&themes).Order("nom asc").Select()
	return themes, err
}

func (repo *Repository) CreateTheme(theme models.Theme) (int64, error) {
	theme.Id = 0
	err := repo.db.client.Insert(&theme)
	return theme.Id, err
}

func (repo *Repository) DeleteTheme(idTheme int64) error {
	theme := models.Theme{
		Id: idTheme,
	}
	return repo.db.client.Delete(&theme)
}
