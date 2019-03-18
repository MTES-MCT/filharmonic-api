package domain

import (
	"github.com/MTES-MCT/filharmonic-api/emails"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/rs/zerolog/log"
)

type NouveauxMessagesUser struct {
	Destinataire models.User
	Messages     []NouveauMessage
}

type NouveauMessage struct {
	DateInspection       string `json:"date_inspection"`
	NomAuteur            string `json:"nom_auteur"`
	RaisonEtablissement  string `json:"raison_etablissement"`
	SujetPointDeControle string `json:"sujet_point_de_controle"`
	Message              string `json:"message"`
	DateMessage          string `json:"date_message"`
	InspectionId         int64  `json:"inspection_id"`
	PointDeControleId    int64  `json:"point_de_controle_id"`
	MessageId            int64  `json:"message_id"`
}

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
			To:       nouveauxMessagesUser.Destinataire.AsRecipient(),
			Subject:  "Fil'Harmonic : Nouveaux messages",
			TextPart: "",
			HTMLPart: htmlPart,
		})
		if err != nil {
			log.Error().Err(err).Msg("error while sending email")
		}
	}
	return nil
}

type RecapValidationInspection struct {
	Destinataire         models.User
	InspectionId         int64  `json:"inspection_id"`
	DateInspection       string `json:"date_inspection"`
	RaisonEtablissement  string `json:"raison_etablissement"`
	AdresseEtablissement string `json:"adresse_etablissement"`
	NonConformites       bool   `json:"non_conformites"`
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
			To:       recapValidation.Destinataire.AsRecipient(),
			Subject:  "Fil'Harmonic : Inspection validée",
			TextPart: "",
			HTMLPart: htmlPart,
		})
		if err != nil {
			log.Error().Err(err).Msg("error while sending email")
		}
	}
	return nil
}

type InspectionExpirationDelais struct {
	Destinataire         models.User
	InspectionId         int64  `json:"inspection_id"`
	ConstatId            int64  `json:"constat_id"`
	DateInspection       string `json:"date_inspection"`
	RaisonEtablissement  string `json:"raison_etablissement"`
	AdresseEtablissement string `json:"adresse_etablissement"`
}

func (s *Service) SendEmailsExpirationDelais() error {
	inspectionsDelaisExpires, err := s.repo.ListInspectionsExpirationDelais()
	if err != nil {
		return err
	}
	constatsIds := []int64{}
	previousInspectionId := int64(0)
	previousRecipientEmail := ""
	for _, inspectionDelaisExpires := range inspectionsDelaisExpires {
		htmlPart, err := s.templateService.RenderHTMLEmailExpirationDelais(inspectionDelaisExpires)
		if err != nil {
			return err
		}
		if inspectionDelaisExpires.InspectionId != previousInspectionId || inspectionDelaisExpires.Destinataire.Email != previousRecipientEmail {
			err = s.emailService.Send(emails.Email{
				To:       inspectionDelaisExpires.Destinataire.AsRecipient(),
				Subject:  "Fil'Harmonic : Expiration des délais",
				TextPart: "",
				HTMLPart: htmlPart,
			})
			if err != nil {
				log.Error().Err(err).Msg("error while sending email")
			}
			previousInspectionId = inspectionDelaisExpires.InspectionId
			previousRecipientEmail = inspectionDelaisExpires.Destinataire.Email
		}
		constatsIds = append(constatsIds, inspectionDelaisExpires.ConstatId)
	}
	err = s.repo.UpdateNotificationEcheanceExpireeEnvoyee(constatsIds)
	return err
}

type InspectionEcheancesProches struct {
	Destinataire         models.User
	InspectionId         int64  `json:"inspection_id"`
	ConstatId            int64  `json:"constat_id"`
	DateInspection       string `json:"date_inspection"`
	RaisonEtablissement  string `json:"raison_etablissement"`
	AdresseEtablissement string `json:"adresse_etablissement"`
}

func (s *Service) SendEmailsRappelEcheances() error {
	inspectionsEcheancesProches, err := s.repo.ListInspectionsEcheancesProches(s.config.SeuilRappelEcheances)
	if err != nil {
		return err
	}
	constatsIds := []int64{}
	previousInspectionId := int64(0)
	for _, inspectionEcheancesProches := range inspectionsEcheancesProches {
		var htmlPart string
		htmlPart, err = s.templateService.RenderHTMLEmailRappelEcheances(inspectionEcheancesProches)
		if err != nil {
			return err
		}
		if inspectionEcheancesProches.InspectionId != previousInspectionId {
			err = s.emailService.Send(emails.Email{
				To:       inspectionEcheancesProches.Destinataire.AsRecipient(),
				Subject:  "Fil'Harmonic : Rappel des échéances",
				TextPart: "",
				HTMLPart: htmlPart,
			})
			if err != nil {
				log.Error().Err(err).Msg("error while sending email")
			}
			previousInspectionId = inspectionEcheancesProches.InspectionId
		}
		constatsIds = append(constatsIds, inspectionEcheancesProches.ConstatId)
	}
	err = s.repo.UpdateNotificationRappelEcheanceEnvoyee(constatsIds)
	return err
}
