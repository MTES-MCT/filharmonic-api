package domain

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

var NotificationsEvenements = make(map[models.TypeEvenement](func(*UserContext) []string))

func init() {
	NotificationsEvenements[models.EvenementCreationMessage] = func(ctx *UserContext) []string {
		if ctx.IsExploitant() {
			return []string{"inspecteurs"}
		}
		return []string{"exploitants"}
	}
	NotificationsEvenements[models.EvenementCreationCommentaire] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementLectureMessage] = func(ctx *UserContext) []string {
		if ctx.IsExploitant() {
			return []string{"inspecteurs"}
		}
		return []string{"exploitants"}
	}
	NotificationsEvenements[models.EvenementCommentaireGeneral] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementCreationInspection] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementModificationInspection] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementPublicationInspection] = func(ctx *UserContext) []string {
		return []string{"inspecteurs", "exploitants"}
	}
	NotificationsEvenements[models.EvenementDemandeValidationInspection] = func(ctx *UserContext) []string {
		return []string{"inspecteurs", "approbateurs"}
	}
	NotificationsEvenements[models.EvenementRejetValidationInspection] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementValidationInspection] = func(ctx *UserContext) []string {
		return []string{"inspecteurs", "exploitants"}
	}
	NotificationsEvenements[models.EvenementClotureInspection] = func(ctx *UserContext) []string {
		return []string{"inspecteurs", "exploitants"}
	}
	NotificationsEvenements[models.EvenementCreationPointDeControle] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementModificationPointDeControle] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementSuppressionPointDeControle] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementPublicationPointDeControle] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementCreationConstat] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementModificationConstat] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementSuppressionConstat] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementResolutionConstat] = func(ctx *UserContext) []string {
		return []string{"inspecteurs", "exploitants"}
	}
	NotificationsEvenements[models.EvenementCreationSuite] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementModificationSuite] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
	NotificationsEvenements[models.EvenementSuppressionSuite] = func(ctx *UserContext) []string {
		return []string{"inspecteurs"}
	}
}

type ListNotificationsFilter struct {
	Lue bool `form:"lue" default:"false"`
}

func (s *Service) ListNotifications(ctx *UserContext, filter *ListNotificationsFilter) ([]models.Notification, error) {
	return s.repo.ListNotifications(ctx, filter)
}

func (s *Service) LireNotifications(ctx *UserContext, ids []int64) error {
	return s.repo.LireNotifications(ctx, ids)
}
