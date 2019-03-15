package database

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) CreateEvenementTx(tx *pg.Tx, ctx *domain.UserContext, typeEvenement models.TypeEvenement, idInspection int64, data map[string]interface{}) error {
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
	if err != nil {
		return err
	}
	err = repo.eventsManager.DispatchUpdatedResource(ctx, "inspection", idInspection)
	if err != nil {
		return err
	}
	userIds, err := repo.ListUsersAssignedToInspection(idInspection)
	if err != nil {
		return err
	}
	err = repo.eventsManager.DispatchUpdatedResourcesToUsers("inspections", userIds)
	return err
}

func (repo *Repository) CreateEvenement(ctx *domain.UserContext, typeEvenement models.TypeEvenement, idInspection int64, data map[string]interface{}) error {
	return repo.db.client.RunInTransaction(func(tx *pg.Tx) error {
		return repo.CreateEvenementTx(tx, ctx, typeEvenement, idInspection, data)
	})
}
