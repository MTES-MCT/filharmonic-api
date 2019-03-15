package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/go-pg/pg"
)

func (repo *Repository) ListNotifications(ctx *domain.UserContext, filter *domain.ListNotificationsFilter) ([]models.Notification, error) {
	notifications := []models.Notification{}
	query := repo.db.client.Model(&notifications).
		Relation("Destinataire").
		Relation("Evenement").
		Relation("Evenement.Auteur").
		Where("destinataire_id = ?", ctx.User.Id)
	if filter != nil {
		query.Where("lue is ?", filter.Lue)
	}
	err := query.OrderExpr("evenement.created_at DESC").Select()
	return notifications, err
}

//nolint: gocyclo
func (repo *Repository) createNotifications(tx *pg.Tx, ctx *domain.UserContext, evenement models.Evenement) error {
	notifications := make([]models.Notification, 0)

	addNotification := func(userId int64) {
		notifications = append(notifications, models.Notification{
			DestinataireId: userId,
			EvenementId:    evenement.Id,
		})
	}

	groupesDestinataires := domain.NotificationsEvenements[evenement.Type](ctx)
	if util.ContainsString(groupesDestinataires, "inspecteurs") {
		inspecteurs := []models.User{}
		err := tx.Model(&inspecteurs).
			Join("JOIN inspection_to_inspecteurs AS inspecteurs").
			JoinOn(`inspecteurs.user_id = "user".id`).
			Where("inspecteurs.inspection_id = ?", evenement.InspectionId).
			Where(`"user".id <> ?`, ctx.User.Id).
			Column("user.id").
			Select()
		if err != nil {
			return err
		}
		for _, inspecteur := range inspecteurs {
			addNotification(inspecteur.Id)
		}
	}
	if util.ContainsString(groupesDestinataires, "exploitants") {
		exploitants := []models.User{}
		err := tx.Model(&exploitants).
			Join("JOIN etablissement_to_exploitants AS exploitants").
			JoinOn(`exploitants.user_id = "user".id`).
			Join("JOIN inspections").
			JoinOn("inspections.etablissement_id = exploitants.etablissement_id").
			Where("inspections.id = ?", evenement.InspectionId).
			Where(`"user".id <> ?`, ctx.User.Id).
			Column("user.id").
			Select()
		if err != nil {
			return err
		}
		for _, exploitant := range exploitants {
			addNotification(exploitant.Id)
		}
	}
	if util.ContainsString(groupesDestinataires, "approbateurs") {
		approbateurs := []models.User{}
		err := tx.Model(&approbateurs).
			Where("profile = ?", models.ProfilApprobateur).
			Where("id <> ?", ctx.User.Id).
			Column("id").
			Select()
		if err != nil {
			return err
		}
		for _, approbateur := range approbateurs {
			addNotification(approbateur.Id)
		}
	}
	if len(notifications) == 0 {
		return nil
	}
	err := tx.Insert(&notifications)
	if err != nil {
		return err
	}

	destinatairesId := []int64{}
	for _, notification := range notifications {
		destinatairesId = append(destinatairesId, notification.DestinataireId)
	}
	err = repo.eventsManager.DispatchUpdatedResourcesToUsers("notifications", destinatairesId)
	return err
}

func (repo *Repository) LireNotifications(ctx *domain.UserContext, ids []int64) error {
	notification := models.Notification{
		Lue: true,
	}
	_, err := repo.db.client.Model(&notification).
		Where("id in (?)", pg.In(ids)).
		Where("destinataire_id = ?", ctx.User.Id).
		Column("lue").
		Update()
	return err
}
