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

	tests.Auth(e.GET("/etablissements")).WithQuery("s3ic", "23").
		Expect().
		Status(http.StatusOK).
		JSON().Array().
		Element(0).Object().ValueEqual("s3ic", "1234")
}

func TestFindEtablissementsOwnedByExploitant(t *testing.T) {
	e, close := tests.Init(t, initTestEtablissementsDB)
	defer close()

	results := tests.Auth(e.GET("/etablissements")).
		Expect().
		Status(http.StatusOK).
		JSON().Array()
	results.Length().Equal(1)
	results.First().Object().ValueEqual("s3ic", "1234")
}

func initTestEtablissementsDB(db *database.Database, assert *require.Assertions) {
	etablissements := []interface{}{
		&models.Etablissement{
			Id:      1,
			S3IC:    "1234",
			Raison:  "Raison sociale",
			Adresse: "1 rue des fleurs 75000 Paris",
		},
		&models.Etablissement{
			Id:      2,
			S3IC:    "4567",
			Raison:  "Raison sociale 2",
			Adresse: "1 rue des plantes 44000 Nantes",
		},
	}

	assert.NoError(db.Insert(etablissements...))
	assert.NoError(db.Insert(&models.EtablissementToExploitant{
		EtablissementId: 1,
		UserId:          1,
	}))
	assert.NoError(db.Insert(&models.EtablissementToExploitant{
		EtablissementId: 2,
		UserId:          2,
	}))
}
