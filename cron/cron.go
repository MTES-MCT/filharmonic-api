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
	Activity string `default:"0 0 1 * * *"`
}

type CronManager struct {
	config       Config
	cron         *cron.Cron
	service      *domain.Service
	emailService *emails.EmailService
}

var nouveauxMessagesTemplate *template.Template

func init() {
	var err error
	nouveauxMessagesTemplate, err = template.ParseFiles("templates/new-messages.tmpl")
	// nouveauxMessagesTemplate, err = template.ParseFiles("cron/templates/new-messages.tmpl")
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse template")
	}
}

func New(config Config, service *domain.Service, emailService *emails.EmailService) (*CronManager, error) {
	cronmanager := &CronManager{
		config:       config,
		cron:         cron.New(),
		service:      service,
		emailService: emailService,
	}
	err := cronmanager.cron.AddFunc(config.Activity, cronmanager.sendEmailsNouveauxMessages)
	cronmanager.cron.Start()
	return cronmanager, err
}

func (cron *CronManager) sendEmailsNouveauxMessages() {
	log.Info().Msg("sending cron emails nouveaux messages")
	nouveauxMessagesUsers, err := cron.service.ListNouveauxMessages()
	if err != nil {
		log.Error().Err(err).Msg("error while fetching data from database to be sent by emails")
	}
	for _, nouveauxMessagesUser := range nouveauxMessagesUsers {
		htmlPart, err := renderEmailTemplate(nouveauxMessagesTemplate, nouveauxMessagesUser)
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
