package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/database"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
	"github.com/stretchr/testify/require"
)

func TestFindEtablissementsByS3IC(t *testing.T) {
	e, close := tests.Init(t, initTestEtablissementsDB)
	defer close()

	e.GET("/etablissements").WithQuery("s3ic", "23").
		Expect().
		Status(http.StatusOK).
		JSON().Array().
		Element(0).Object().ValueEqual("S3IC", "1234")
}

func initTestEtablissementsDB(db *database.Database, assert *require.Assertions) {
	etablissement := &models.Etablissement{
		S3IC:    "1234",
		Raison:  "Raison sociale",
		Adresse: "1 rue des fleurs 75000 Paris",
	}
	err := db.Insert(etablissement)
	assert.NoError(err)
}
