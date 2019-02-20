package templates

import (
	"bytes"
	html "html/template"
	text "text/template"
)

type Config struct {
	Dir string `default:"templates/templates/"`
}

type TemplateService struct {
	config Config

	emailNouveauxMessagesTemplate *html.Template
	lettreAnnonceTemplate         *text.Template
}

func New(config Config) (*TemplateService, error) {
	service := &TemplateService{
		config: config,
	}
	var err error
	service.emailNouveauxMessagesTemplate, err = html.ParseFiles(config.Dir + "email-new-messages.html")
	if err != nil {
		return nil, err
	}

	service.lettreAnnonceTemplate, err = text.ParseFiles(config.Dir + "modele_de_lettre_annonce_inspection.fodt")
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *TemplateService) RenderHTMLEmailNouveauxMessages(data interface{}) (string, error) {
	return s.renderHTMLTemplate(s.emailNouveauxMessagesTemplate, data)
}

func (s *TemplateService) RenderLettreAnnonce(data interface{}) (string, error) {
	return s.renderTextTemplate(s.lettreAnnonceTemplate, data)
}

func (s *TemplateService) renderHTMLTemplate(tmpl *html.Template, data interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (s *TemplateService) renderTextTemplate(tmpl *text.Template, data interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
