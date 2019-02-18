package cron

import (
	"bytes"
	"html/template"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/emails"
	"github.com/robfig/cron"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Activity     string `default:"0 0 1 * * *"`
	TemplatesDir string `default:"cron/templates"`
}

type CronManager struct {
	config       Config
	cron         *cron.Cron
	service      *domain.Service
	emailService *emails.EmailService

	nouveauxMessagesTemplate *template.Template
}

func New(config Config, service *domain.Service, emailService *emails.EmailService) (*CronManager, error) {
	cronmanager := &CronManager{
		config:       config,
		cron:         cron.New(),
		service:      service,
		emailService: emailService,
	}
	err := cronmanager.cron.AddFunc(config.Activity, cronmanager.sendEmailsNouveauxMessages)
	if err != nil {
		return nil, err
	}
	cronmanager.nouveauxMessagesTemplate, err = template.ParseFiles(config.TemplatesDir + "/new-messages.tmpl")
	if err != nil {
		return nil, err
	}
	cronmanager.cron.Start()
	return cronmanager, nil
}

func (cron *CronManager) sendEmailsNouveauxMessages() {
	log.Info().Msg("sending cron emails nouveaux messages")
	nouveauxMessagesUsers, err := cron.service.ListNouveauxMessages()
	if err != nil {
		log.Error().Err(err).Msg("error while fetching data from database to be sent by emails")
	}
	for _, nouveauxMessagesUser := range nouveauxMessagesUsers {
		htmlPart, err := renderEmailTemplate(cron.nouveauxMessagesTemplate, nouveauxMessagesUser)
		if err != nil {
			log.Error().Err(err).Msg("error while rendering email")
		}

		err = cron.emailService.Send(emails.Email{
			Subject:        "Fil'Harmonic : Nouveaux messages",
			RecipientEmail: nouveauxMessagesUser.Destinataire.Email,
			RecipientName:  nouveauxMessagesUser.Destinataire.Nom,
			TextPart:       "TEMPLATE TEXT",
			HTMLPart:       htmlPart,
		})
		if err != nil {
			log.Error().Err(err).Msg("error while sending activity by emails")
		}
	}
	log.Info().Int("emails", len(nouveauxMessagesUsers)).Msg("emails sent")
}

func renderEmailTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
