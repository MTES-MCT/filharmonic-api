package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (repo *Repository) CreatePointDeControle(ctx *domain.UserContext, idInspection int64, pointDeControle models.PointDeControle) (int64, error) {
	pointDeControle.InspectionId = idInspection
	err := repo.db.client.Insert(&pointDeControle)
	if err != nil {
		return 0, err
	}
	return pointDeControle.Id, nil
}
