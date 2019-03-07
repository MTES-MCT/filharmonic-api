package cron

import (
	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/robfig/cron"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Activity string `default:"0 0 1 * * 1-5"` // every day from monday to friday, at 1
}

type CronManager struct {
	config  Config
	cron    *cron.Cron
	service *domain.Service
}

func New(config Config, service *domain.Service) (*CronManager, error) {
	cronmanager := &CronManager{
		config:  config,
		cron:    cron.New(),
		service: service,
	}
	err := cronmanager.cron.AddFunc(config.Activity, func() {
		err := cronmanager.service.SendEmailsNouveauxMessages()
		if err != nil {
			log.Error().Err(err).Msg("error while sending emails")
		}
	})
	if err != nil {
		return nil, err
	}
	cronmanager.cron.Start()
	return cronmanager, nil
}
