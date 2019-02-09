package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/go-pg/pg"
)

func (repo *Repository) ListNotifications(ctx *domain.UserContext, filter *domain.ListNotificationsFilter) ([]models.Notification, error) {
	notifications := []models.Notification{}
	query := repo.db.client.Model(&notifications).
		Relation("Lecteur").
		Relation("Evenement").
		Relation("Evenement.Auteur")
	if filter != nil {
		query.Where("lue is ?", filter.Lue)
	}
	err := query.OrderExpr("evenement.created_at DESC").Select()
	return notifications, err
}

func (repo *Repository) CreateNotification(ctx *domain.UserContext, notification models.Notification) (int64, error) {
	err := repo.db.client.Insert(&notification)
	return notification.Id, err
}

func (repo *Repository) UpdateNotifications(ctx *domain.UserContext, notification models.Notification, ids []int64) error {
	_, err := repo.db.client.Model(&notification).Where("id in (?)", pg.In(ids)).Column("lue", "lecteur_id").Update()
	return err
}

func (repo *Repository) CheckUserAllowedNotifications(ctx *domain.UserContext, ids []int64) (bool, error) {
	if ctx.IsExploitant() {
		query := repo.db.client.Model(&models.Notification{}).
			Join("JOIN evenements AS evenement").
			JoinOn("evenement.id = notification.evenement_id").
			Join("JOIN inspections AS inspection").
			JoinOn("inspection.id = evenement.inspection_id").
			Join("JOIN etablissements AS etablissement").
			JoinOn("etablissement.id = inspection.etablissement_id").
			Join("JOIN etablissement_to_exploitants AS exploitants").
			JoinOn("exploitants.etablissement_id = etablissement.id").
			JoinOn("exploitants.user_id = ?", ctx.User.Id)
		if len(ids) > 0 {
			query.Where("notification.id in (?)", pg.In(ids))
		}
		count, err := query.Count()
		return count > 0, err
	} else {
		query := repo.db.client.Model(&models.Notification{}).
			Join("JOIN evenements AS evenement").
			JoinOn("evenement.id = notification.evenement_id").
			Join("JOIN inspection_to_inspecteurs AS inspecteurs").
			JoinOn("inspecteurs.inspection_id = evenement.inspection_id").
			JoinOn("inspecteurs.user_id = ?", ctx.User.Id)
		if len(ids) > 0 {
			query.Where("notification.id in (?)", pg.In(ids))
		}
		count, err := query.Count()
		return count > 0, err
	}
}
