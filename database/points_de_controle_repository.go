package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (repo *Repository) CreatePointDeControle(ctx *domain.UserContext, idInspection int64, pointDeControle models.PointDeControle) (int64, error) {
	pointDeControle.Id = 0
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

func (repo *Repository) CheckUserAllowedPointDeControle(ctx *domain.UserContext, id int64) (bool, error) {
	if ctx.IsExploitant() {
		count, err := repo.db.client.Model(&models.PointDeControle{}).
			Join("JOIN inspections AS i").
			JoinOn("i.id = point_de_controle.inspection_id").
			Join("JOIN etablissements AS e").
			JoinOn("e.id = i.etablissement_id").
			Join("JOIN etablissement_to_exploitants AS ex").
			JoinOn("ex.etablissement_id = e.id").
			JoinOn("ex.user_id = ?", ctx.User.Id).
			Where("point_de_controle.id = ?", id).
			Count()
		return count == 1, err
	} else {
		count, err := repo.db.client.Model(&models.PointDeControle{}).
			Join("JOIN inspection_to_inspecteurs AS u").
			JoinOn("u.inspection_id = point_de_controle.inspection_id").
			JoinOn("u.user_id = ?", ctx.User.Id).
			Where("point_de_controle.id = ?", id).
			Count()
		return count == 1, err
	}
}
