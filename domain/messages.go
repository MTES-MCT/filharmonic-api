package domain

import (
	"github.com/MTES-MCT/filharmonic-api/models"
)

func (s *Service) CreateMessage(ctx *UserContext, idPointDeControle int64, message models.Message) (int64, error) {
	if ctx.IsApprobateur() {
		return 0, ErrInvalidInput
	}
	ok, err := s.repo.CheckUserAllowedPointDeControle(ctx, idPointDeControle)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}

	if ctx.IsExploitant() {
		message.Interne = false
	}
	etatInspection, err := s.repo.GetEtatInspectionByPointDeControleID(idPointDeControle)
	if err != nil {
		return 0, err
	}
	message.EtapeTraitementNonConformites = etatInspection == models.EtatTraitementNonConformites
	return s.repo.CreateMessage(ctx, idPointDeControle, message)
}

func (s *Service) LireMessage(ctx *UserContext, idMessage int64) error {
	ok, err := s.repo.CheckUserAllowedMessage(ctx, idMessage)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	ok, err = s.repo.CheckUserIsRecipient(ctx, idMessage)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	return s.repo.LireMessage(ctx, idMessage)
}

// func (s *Service) ListNouveauxMessages() ([]NouveauxMessagesUser, error) {
// 	return s.repo.ListNouveauxMessages()
// }
