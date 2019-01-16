package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/app"
	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/stretchr/testify/require"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func TestFindEtablissementsByS3IC(t *testing.T) {
	assert := require.New(t)

	config := app.LoadConfig()
	config.Database.InitSchema = true
	db, server := app.Bootstrap(config)
	defer server.Close()

	initTestDB(db, assert)

	e := httpexpect.New(t, "http://localhost:8080")

	e.GET("/etablissements").WithQuery("s3ic", "23").
		Expect().
		Status(http.StatusOK).
		JSON().Array().
		Element(0).Object().ValueEqual("s3ic", "1234")
}

func initTestDB(db *database.Database, assert *require.Assertions) {
	etablissement := &models.Etablissement{
		S3IC:    "1234",
		Raison:  "Raison sociale",
		Adresse: "1 rue des fleurs 75000 Paris",
	}
	err := db.Insert(etablissement)
	assert.NoError(err)
}
