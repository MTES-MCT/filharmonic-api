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

func (repo *Repository) UpdatePointDeControle(ctx *domain.UserContext, idPointDeControle int64, pointDeControle models.PointDeControle) error {
	pointDeControle.Id = idPointDeControle
	columns := []string{"sujet", "references_reglementaires"}
	_, err := repo.db.client.Model(&pointDeControle).Column(columns...).WherePK().Update()
	return err
}

func (repo *Repository) PublishPointDeControle(ctx *domain.UserContext, idPointDeControle int64) error {
	pointDeControle := models.PointDeControle{
		Id:     idPointDeControle,
		Publie: true,
	}
	columns := []string{"publie"}
	_, err := repo.db.client.Model(&pointDeControle).Column(columns...).WherePK().Update()
	return err
}

func (repo *Repository) DeletePointDeControle(ctx *domain.UserContext, idPointDeControle int64) error {
	pointDeControle := models.PointDeControle{
		Id: idPointDeControle,
	}
	return repo.db.client.Delete(&pointDeControle)
}

func (repo *Repository) CheckInspecteurAllowedPointDeControle(ctx *domain.UserContext, id int64) (bool, error) {
	count, err := repo.db.client.Model(&models.PointDeControle{}).
		Join("JOIN inspection_to_inspecteurs AS u").
		JoinOn("u.inspection_id = point_de_controle.inspection_id").
		JoinOn("u.user_id = ?", ctx.User.Id).
		Where("point_de_controle.id = ?", id).
		Count()
	return count == 1, err
}
