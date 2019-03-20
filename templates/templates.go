package templates

import (
	"bytes"
	html "html/template"
	text "text/template"

	"github.com/MTES-MCT/filharmonic-api/models"
)

type TemplateType string

const (
	// Emails
	TemplateEmailExpirationDelais TemplateType = "expiration-delais"
	TemplateEmailNouveauxMessages TemplateType = "nouveaux-messages"
	TemplateEmailRecapValidation  TemplateType = "recap-validation"
	TemplateEmailRappelEcheances  TemplateType = "rappel-echeances"

	// ODT
	TemplateODTLettreAnnonce TemplateType = "lettre-annonce"
	TemplateODTLettreSuite   TemplateType = "lettre-suite"
	TemplateODTRapport       TemplateType = "rapport"
)

var emailTemplates = []TemplateType{
	TemplateEmailExpirationDelais,
	TemplateEmailNouveauxMessages,
	TemplateEmailRecapValidation,
	TemplateEmailRappelEcheances,
}

var odtTemplates = []TemplateType{
	TemplateODTLettreAnnonce,
	TemplateODTLettreSuite,
	TemplateODTRapport,
}

type Template struct {
	HTML *html.Template
	Text *text.Template
}

func (t *Template) Render(data interface{}) (*models.RenderedTemplate, error) {
	result := &models.RenderedTemplate{}
	if t.HTML != nil {
		var tpl bytes.Buffer
		err := t.HTML.Execute(&tpl, data)
		if err != nil {
			return nil, err
		}
		result.HTML = tpl.String()
	}
	if t.Text != nil {
		var tpl bytes.Buffer
		err := t.Text.Execute(&tpl, data)
		if err != nil {
			return nil, err
		}
		result.Text = tpl.String()
	}
	return result, nil
}
