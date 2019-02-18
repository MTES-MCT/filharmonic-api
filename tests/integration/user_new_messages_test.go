package integration

import (
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestListNouveauxMessage(t *testing.T) {
	assert, application := tests.InitDB(t)
	nouveauxMessagesUsers, err := application.Repo.ListNouveauxMessages()

	assert.NoError(err)
	assert.Len(nouveauxMessagesUsers, 3)
	nouveauxMessagesExploitant := nouveauxMessagesUsers[0]
	assert.Equal("exploitant1@filharmonic.com", nouveauxMessagesExploitant.Destinataire.Email)
	assert.Len(nouveauxMessagesExploitant.Messages, 1)
	message := nouveauxMessagesExploitant.Messages[0]
	assert.Equal("Alain Champion", message.NomAuteur)
	assert.Equal("Raison sociale", message.RaisonEtablissement)
	assert.Equal("Il manque un document.", message.Message)
	assert.Equal("Mesure des émissions atmosphériques canalisées par un organisme extérieur", message.SujetPointDeControle)
	assert.Equal("2018-09-01", message.DateInspection)
}
