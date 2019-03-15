package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
)

func (repo *Repository) ListCanevas() ([]models.Canevas, error) {
	canevas := []models.Canevas{}
	err := repo.db.client.Model(&canevas).
		OrderExpr("lower(unaccent(nom)) asc").
		Select()
	return canevas, err
}

func (repo *Repository) GetCanevasByID(id int64) (*models.Canevas, error) {
	canevas := models.Canevas{
		Id: id,
	}
	err := repo.db.client.Model(&canevas).WherePK().Select()
	return &canevas, err
}

func (repo *Repository) ImportCanevas(ctx *domain.UserContext, inspectionId int64, canevas models.Canevas) error {
	pointsDeControle := make([]models.PointDeControle, 0)
	for _, dataPointDeControle := range canevas.Data.PointsDeControle {
		pointDeControle := models.PointDeControle{
			Sujet:                    dataPointDeControle.Sujet,
			ReferencesReglementaires: dataPointDeControle.ReferencesReglementaires,
			InspectionId:             inspectionId,
		}
		pointsDeControle = append(pointsDeControle, pointDeControle)
	}
	err := repo.db.client.Insert(&pointsDeControle)
	if err != nil {
		return err
	}
	messages := make([]models.Message, 0)
	for i, dataPointDeControle := range canevas.Data.PointsDeControle {
		if dataPointDeControle.Message != "" {
			message := models.Message{
				Message:           dataPointDeControle.Message,
				Date:              util.Now(),
				AuteurId:          ctx.User.Id,
				PointDeControleId: pointsDeControle[i].Id,
			}
			messages = append(messages, message)
		}
	}
	err = repo.db.client.Insert(&messages)
	return err
}

func (repo *Repository) CreateCanevas(ctx *domain.UserContext, idInspection int64, canevas models.Canevas) (int64, error) {
	canevas.Id = 0
	pointsDeControle := []models.CanevasPointDeControle{}
	_, err := repo.db.Query(&pointsDeControle, `select distinct(point_de_controle.sujet) as sujet,
		point_de_controle.references_reglementaires as references_reglementaires,
		(select m.message from messages as m
		where m.point_de_controle_id = point_de_controle.id
		order by m.date
		limit 1) as message
		from point_de_controles as point_de_controle
		left join messages as message
		on message.point_de_controle_id = point_de_controle.id
		where point_de_controle.inspection_id = ?
		order by point_de_controle.sujet asc`, idInspection)
	if err != nil {
		return 0, err
	}
	canevas.DataVersion = 1
	canevas.AuteurId = ctx.User.Id
	canevas.CreatedAt = util.Now()
	canevas.Data = models.CanevasData{
		PointsDeControle: pointsDeControle,
	}
	_, err = repo.db.Model(&canevas).
		OnConflict("(nom) DO UPDATE").
		Insert()
	if err != nil {
		return 0, err
	}
	err = repo.eventsManager.DispatchUpdatedResources(ctx, "canevas")
	return canevas.Id, err
}

func (repo *Repository) DeleteCanevas(ctx *domain.UserContext, id int64) error {
	canevas := models.Canevas{
		Id: id,
	}
	err := repo.db.client.Delete(&canevas)
	if err != nil {
		return err
	}
	err = repo.eventsManager.DispatchUpdatedResources(ctx, "canevas")
	return err
}
