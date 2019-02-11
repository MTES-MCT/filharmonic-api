package e2e

import (
	"net/http"
	"testing"

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
	firstNotification.ValueEqual("id", 3)
	firstNotification.ValueEqual("lue", false)
	evenement := firstNotification.Value("evenement").Object()
	evenement.ValueEqual("id", 3)
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
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	notifications.Length().Equal(1)
}

func TestLireNotificationsDisallowedExploitant(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	ids := []interface{}{1}

	tests.AuthUser(e.POST("/notifications/lire"), 2).WithJSON(ids).
		Expect().
		Status(http.StatusOK)
}
