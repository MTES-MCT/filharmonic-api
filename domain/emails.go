package domain

import (
	"github.com/MTES-MCT/filharmonic-api/emails"
	"github.com/rs/zerolog/log"
)

func (s *Service) SendEmailsNouveauxMessages() error {
	nouveauxMessagesUsers, err := s.repo.ListNouveauxMessages()
	if err != nil {
		return err
	}
	for _, nouveauxMessagesUser := range nouveauxMessagesUsers {
		htmlPart, err := s.templateService.RenderHTMLEmailNouveauxMessages(nouveauxMessagesUser)
		if err != nil {
			return err
		}

		err = s.emailService.Send(emails.Email{
			Subject:        "Fil'Harmonic : Nouveaux messages",
			RecipientEmail: nouveauxMessagesUser.Destinataire.Email,
			RecipientName:  nouveauxMessagesUser.Destinataire.Nom,
			TextPart:       "",
			HTMLPart:       htmlPart,
		})
		if err != nil {
			log.Error().Err(err).Msg("error while sending email")
		}
	}
	return nil
}

func (s *Service) SendEmailsRecapValidation(idInspection int64) error {
	recaps, err := s.repo.GetRecapsValidation(idInspection)
	if err != nil {
		return err
	}

	for _, recapValidation := range recaps {
		htmlPart, err := s.templateService.RenderHTMLEmailRecapValidation(recapValidation)
		if err != nil {
			return err
		}

		err = s.emailService.Send(emails.Email{
			Subject:        "Fil'Harmonic : Inspection validée",
			RecipientEmail: recapValidation.Destinataire.Email,
			RecipientName:  recapValidation.Destinataire.Nom,
			TextPart:       "",
			HTMLPart:       htmlPart,
		})
		if err != nil {
			log.Error().Err(err).Msg("error while sending email")
		}
	}
	return nil
}
