package domain

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/models"
)

type ListNotificationsFilter struct {
	Lue bool `form:"lue"`
}

func (s *Service) ListNotifications(ctx *UserContext, filter ListNotificationsFilter) ([]models.Notification, error) {
	return s.repo.ListNotifications(filter)
}

func (s *Service) CreateNotification(ctx *UserContext, notification models.Notification) (int64, error) {
	notification.CreatedAt = time.Now()
	notification.Lue = false
	notification.LecteurId = ctx.User.Id
	return s.repo.CreateNotification(notification)
}
