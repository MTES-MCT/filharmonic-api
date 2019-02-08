package domain

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

type ListNotificationsFilter struct {
	Lue bool `form:"lue" default:"false"`
}

func (s *Service) ListNotifications(ctx *UserContext, filter *ListNotificationsFilter) ([]models.Notification, error) {
	ok, err := s.repo.CheckUserAllowedNotifications(ctx, []int64{})
	if err != nil {
		return []models.Notification{}, err
	}
	if !ok {
		return []models.Notification{}, ErrInvalidInput
	}
	return s.repo.ListNotifications(ctx, filter)
}

func (s *Service) CreateNotification(ctx *UserContext, notification models.Notification) (int64, error) {
	ok, err := s.repo.CheckUserAllowedNotifications(ctx, []int64{})
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}
	return s.repo.CreateNotification(ctx, notification)
}

func (s *Service) LireNotifications(ctx *UserContext, ids []int64) error {
	ok, err := s.repo.CheckUserAllowedNotifications(ctx, ids)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	return s.repo.UpdateNotifications(ctx, ids)
}
