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
		WithQuery("lue", "false").
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	notifications.Length().Equal(3)
	firstNotification := notifications.First().Object()
	firstNotification.ValueEqual("id", 1)
	firstNotification.ValueEqual("lue", false)
	evenement := firstNotification.Value("evenement").Object()
	evenement.ValueEqual("id", 1)
	auteur := evenement.Value("auteur").Object()
	auteur.ValueEqual("id", 3)
}

func TestCreateNotification(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	notificationInput := models.Notification{
		EvenementId: 4,
	}

	tests.AuthInspecteur(e.POST("/notifications")).WithJSON(notificationInput).
		Expect().
		Status(http.StatusOK)

	notifications := tests.AuthInspecteur(e.GET("/notifications")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	notifications.Length().Equal(4)
	lastNotification := notifications.Last().Object()
	lastNotification.ValueEqual("id", 4)
	lastNotification.ValueEqual("lue", false)
	evenement := lastNotification.Value("evenement").Object()
	evenement.ValueEqual("id", 4)
	auteur := evenement.Value("auteur").Object()
	auteur.ValueEqual("id", 3)
}

func TestLireNotificationsAllowedInspecteur(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	ids := []interface{}{1, 2}

	tests.AuthInspecteur(e.POST("/notifications/lire")).WithJSON(ids).
		Expect().
		Status(http.StatusOK)

	notifications := tests.AuthInspecteur(e.GET("/notifications")).
		WithQuery("lue", "true").
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	notifications.Length().Equal(2)
	lastNotification := notifications.Last().Object()
	lastNotification.ValueEqual("id", 2)
	lecteur := lastNotification.Value("lecteur").Object()
	lecteur.ValueEqual("id", 3)
	evenement := lastNotification.Value("evenement").Object()
	evenement.ValueEqual("id", 2)
	auteur := evenement.Value("auteur").Object()
	auteur.ValueEqual("id", 3)
}

func TestLireNotificationsDisallowedExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	ids := []interface{}{1, 2}

	tests.AuthExploitant(e.POST("/notifications/lire")).WithJSON(ids).
		Expect().
		Status(http.StatusBadRequest)
}
