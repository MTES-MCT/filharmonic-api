package database

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) ListEvenements(ctx *domain.UserContext, filter domain.ListEvenementsFilter) ([]models.Evenement, error) {
	evenements := []models.Evenement{}
	err := repo.db.client.Model(&evenements).
		Relation("Auteur").
		Where("auteur_id = ?", ctx.User.Id).
		Select()
	return evenements, err
}

func (repo *Repository) CreateEvenement(tx *pg.Tx, ctx *domain.UserContext, typeEvenement models.TypeEvenement, idInspection int64, data map[string]interface{}) error {
	evenement := models.Evenement{
		AuteurId:     ctx.User.Id,
		CreatedAt:    time.Now(),
		Type:         typeEvenement,
		InspectionId: idInspection,
		Data:         data,
	}
	err := tx.Insert(&evenement)
	if err != nil {
		return err
	}

	err = repo.createNotifications(tx, ctx, evenement)
	return err
}
