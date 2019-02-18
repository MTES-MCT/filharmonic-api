package emails

import (
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"github.com/rs/zerolog/log"
)

type Config struct {
	APIPublicKey  string `default:""`
	APIPrivateKey string `default:""`
	FromEmail     string `default:"contact@filharmonic.beta.gouv.fr"`
	FromName      string `default:"Fil'Harmonic"`
}

type EmailService struct {
	config Config
	client *mailjet.Client
}

type Email struct {
	RecipientEmail string
	RecipientName  string
	Subject        string
	TextPart       string
	HTMLPart       string
}

type emailTemplate struct {
	Name      string
	MailJetID int
}

func New(config Config) *EmailService {
	return &EmailService{
		config: config,
		client: mailjet.NewMailjetClient(config.APIPublicKey, config.APIPrivateKey),
	}
}

// See https://dev.mailjet.com/guides/#send-api-v3-1
func (em *EmailService) Send(email Email) error {
	log.Debug().
		Str("recipient", email.RecipientEmail).
		Msg("send email")
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: em.config.FromEmail,
				Name:  em.config.FromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email.RecipientEmail,
					Name:  email.RecipientName,
				},
			},
			Subject:  email.Subject,
			TextPart: email.TextPart,
			HTMLPart: email.HTMLPart,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := em.client.SendMailV31(&messages)
	return err
}

// used in tests
func (em *EmailService) SetBaseURL(baseURL string) {
	em.client.SetBaseURL(baseURL)
}
