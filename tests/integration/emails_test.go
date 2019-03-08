package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestSendEmailsRecapValidation(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	assert.NoError(application.Service.SendEmailsRecapValidation(int64(5)))
}

func TestSendEmailsExpirationDelais(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	util.SetTime(util.Date("2019-03-07").Time)
	assert.NoError(application.Service.SendEmailsExpirationDelais())
	constats := []models.Constat{}
	_, err := application.DB.Query(&constats, "select * from constats where notification_echeance_expiree_envoyee is true")
	assert.NoError(err)
	assert.Len(constats, 1)
}

func TestSendEmailsRappelEcheances(t *testing.T) {
	assert, application, close := tests.InitService(t)
	defer close()

	util.SetTime(util.Date("2019-03-05").Time)
	assert.NoError(application.Service.SendEmailsRappelEcheances())
	constats := []models.Constat{}
	_, err := application.DB.Query(&constats, "select * from constats where notification_rappel_echeance_envoyee is true")
	assert.NoError(err)
	assert.Len(constats, 1)
}
