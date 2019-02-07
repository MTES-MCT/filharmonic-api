package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestListNotifications(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	notifications := tests.AuthInspecteur(e.GET("/notifications")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	notifications.Length().Equal(2)
	firstNotification := notifications.First().Object()
	firstNotification.ValueEqual("id", 1)
	firstNotification.ValueEqual("lue", false)
	lecteur := firstNotification.Value("lecteur").Object()
	lecteur.ValueEqual("id", 3)
	evenement := firstNotification.Value("evenement").Object()
	evenement.ValueEqual("id", 1)
}

func TestCreateNotification(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	notificationInput := models.Notification{
		EvenementId: 3,
	}

	tests.AuthInspecteur(e.POST("/notifications")).WithJSON(notificationInput).
		Expect().
		Status(http.StatusOK)

	notifications := tests.AuthInspecteur(e.GET("/notifications")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	notifications.Length().Equal(3)
	lastNotification := notifications.Last().Object()
	lastNotification.ValueEqual("id", 3)
	lastNotification.ValueEqual("lue", false)
	lecteur := lastNotification.Value("lecteur").Object()
	lecteur.ValueEqual("id", 1)
	evenement := lastNotification.Value("evenement").Object()
	evenement.ValueEqual("id", 3)
}
