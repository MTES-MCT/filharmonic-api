package cron

import (
	"html/template"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestRenderEmailTemplate(t *testing.T) {
	assert := require.New(t)
	tmpl, err := template.ParseFiles("templates/new-messages.tmpl")
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse template")
	}
	data := domain.NouveauxMessagesUser{
		Destinataire: models.User{
			Email: "test@localhost",
			Nom:   "Michel Exploitant1",
		},
		Messages: []domain.NouveauMessage{
			domain.NouveauMessage{
				DateInspection:       "2018-02-24",
				RaisonEtablissement:  "Etablissement 1",
				SujetPointDeControle: "Rejets Eau",
				Message:              "Il faut des photos",
				NomAuteur:            "Alain Champion",
			},
			domain.NouveauMessage{
				DateInspection:       "2018-02-26",
				RaisonEtablissement:  "Etablissement 2",
				SujetPointDeControle: "Rejets Air",
				Message:              "Il faut des documents",
				NomAuteur:            "Alain Champion",
			},
		},
	}

	htmlPart, err := renderEmailTemplate(tmpl, data)
	assert.NoError(err)
	assert.Contains(htmlPart, "Il faut des photos")
	assert.Contains(htmlPart, "Il faut des documents")
}
