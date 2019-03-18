package emails

import (
	"crypto/tls"
	"errors"
	"net/smtp"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Config struct {
	SmtpHost string `default:"mail.filharmonic.beta.gouv.fr"`
	SmtpPort int    `default:"25"`
	SmtpUser string `default:"noreply@filharmonic.beta.gouv.fr"`
	SmtpPass string `default:""`
	FromName string `default:"Fil'Harmonic"`
}

func (config *Config) SmtpAddr() string {
	return config.SmtpHost + ":" + strconv.Itoa(config.SmtpPort)
}

type EmailService struct {
	config Config
	client *smtp.Client
}

func New(config Config) (*EmailService, error) {
	log.Info().Msgf("connecting to SMTP on %s", config.SmtpAddr())
	client, err := smtp.Dial(config.SmtpAddr())
	if err != nil {
		return nil, err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: config.SmtpHost}
		if err = client.StartTLS(config); err != nil {
			return nil, err
		}
	}

	auth := smtp.PlainAuth("", config.SmtpUser, config.SmtpPass, config.SmtpHost)
	if ok, _ := client.Extension("AUTH"); !ok {
		return nil, errors.New("smtp: server doesn't support AUTH")
	}
	if err = client.Auth(auth); err != nil {
		return nil, err
	}
	log.Info().Msg("connected to SMTP")

	return &EmailService{
		config: config,
		client: client,
	}, nil
}

func (em *EmailService) Send(email Email) error {
	log.Info().
		Str("recipient", email.To).
		Msg("send email")

	err := em.client.Mail(em.config.SmtpUser)
	if err != nil {
		return err
	}
	err = em.client.Rcpt(email.To)
	if err != nil {
		return err
	}
	w, err := em.client.Data()
	if err != nil {
		return err
	}
	email.From = em.config.FromName + " <" + em.config.SmtpUser + ">"
	content, err := email.ToBytes()
	if err != nil {
		return err
	}
	_, err = w.Write(content)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return em.client.Quit()
}

func (em *EmailService) Close() error {
	return em.client.Close()
}
