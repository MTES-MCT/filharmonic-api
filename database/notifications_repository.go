package database

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (repo *Repository) ListNotifications(filter domain.ListNotificationsFilter) ([]models.Notification, error) {
	notifications := []models.Notification{}
	err := repo.db.client.Model(&notifications).
		Relation("Lecteur").
		Relation("Evenement").
		Where("lue is ?", filter.Lue).
		Select()
	return notifications, err
}

func (repo *Repository) CreateNotification(notification models.Notification) (int64, error) {
	notification.Id = 0
	err := repo.db.client.Insert(&notification)
	return notification.Id, err
}
