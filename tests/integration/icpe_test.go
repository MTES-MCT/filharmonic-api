package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/database/importcsv"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestLoadEtablissementsFromCSV(t *testing.T) {
	assert, application := tests.InitDB(t)
	assert.NoError(application.DB.Insert(&[]interface{}{
		&models.Departement{
			CodeInsee: "06",
		},
		&models.Departement{
			CodeInsee: "13",
		},
		&models.Departement{
			CodeInsee: "14",
		},
		&models.Departement{
			CodeInsee: "21",
		},
		&models.Departement{
			CodeInsee: "2A",
		},
		&models.Departement{
			CodeInsee: "2B",
		},
		&models.Departement{
			CodeInsee: "50",
		},
		&models.Departement{
			CodeInsee: "61",
		},
		&models.Departement{
			CodeInsee: "67",
		},
		&models.Departement{
			CodeInsee: "68",
		},
		&models.Departement{
			CodeInsee: "71",
		},
		&models.Departement{
			CodeInsee: "83",
		},
		&models.Departement{
			CodeInsee: "972",
		},
	}))
	assert.NoError(application.DB.Insert(&models.Etablissement{
		S3IC:       "0521.00217",
		Nom:        "FROMAGERIE OLD",
		Raison:     "FROMAGERIE OLD",
		Activite:   "Fromage",
		Seveso:     "Haut",
		Iedmtd:     true,
		Adresse1:   "PLACE DE LA LAITERIE",
		Adresse2:   "",
		CodePostal: "56000",
		Commune:    "Vannes",
		Regime:     models.RegimeEnregistrement,
	}))

	err := importcsv.LoadEtablissementsCSV("s3ic_ic_gen_fabnum.mini.csv", application.DB)
	assert.NoError(err)
	etablissements, err := application.Repo.ListEtablissements()
	assert.NoError(err)
	assert.NotEmpty(etablissements)
	assert.Len(etablissements, 54)
	etablissement := etablissements[len(etablissements)-1-3]
	assert.Equal("0521.00217", etablissement.S3IC)
	assert.Equal("FROMAGERIE BERTHAUT", etablissement.Raison)
	assert.Equal("FROMAGERIE BERTHAUT", etablissement.Nom)
	assert.Equal("Fabrication de fromage", etablissement.Activite)
	assert.Equal("PLACE DU CHAMP DE FOIRE", etablissement.Adresse1)
	assert.Equal("BP 5", etablissement.Adresse2)
	assert.Equal("21460", etablissement.CodePostal)
	assert.Equal("EPOISSES", etablissement.Commune)
	assert.Equal("Non Seveso", etablissement.Seveso)
	assert.Equal(models.RegimeDeclaration, etablissement.Regime)
	assert.False(etablissement.Iedmtd)
}
