package domain

import "github.com/MTES-MCT/filharmonic-api/models"

func (s *Service) ListInspections(ctx *UserContext) ([]models.Inspection, error) {
	return s.repo.ListInspections(ctx)
}

func (s *Service) CreateInspection(ctx *UserContext, inspection models.Inspection) (int64, error) {
	if ctx.IsExploitant() {
		return 0, ErrBesoinProfilInspecteur
	}
	inspecteursIds := make([]int64, 0)
	for _, inspecteur := range inspection.Inspecteurs {
		inspecteursIds = append(inspecteursIds, inspecteur.Id)
	}
	ok, err := s.repo.CheckUsersInspecteurs(inspecteursIds)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrInvalidInput
	}
	return s.repo.CreateInspection(ctx, inspection)
}

func (s *Service) GetInspection(ctx *UserContext, id int64) (*models.Inspection, error) {
	return s.repo.GetInspectionByID(ctx, id)
}

func (s *Service) UpdateInspection(ctx *UserContext, inspection models.Inspection) error {
	if ctx.IsExploitant() {
		return ErrBesoinProfilInspecteur
	}
	inspecteursIds := make([]int64, 0)
	for _, inspecteur := range inspection.Inspecteurs {
		inspecteursIds = append(inspecteursIds, inspecteur.Id)
	}
	if len(inspecteursIds) == 0 {
		return ErrInvalidInput
	}
	ok, err := s.repo.CheckUsersInspecteurs(inspecteursIds)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	return s.repo.UpdateInspection(ctx, inspection)
}

func (s *Service) PublishInspection(ctx *UserContext, idInspection int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}
	return s.changeEtatInspection(ctx, idInspection, models.EtatPreparation, models.EtatEnCours)
}

func (s *Service) AskValidateInspection(ctx *UserContext, idInspection int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}
	return s.changeEtatInspection(ctx, idInspection, models.EtatEnCours, models.EtatAttenteValidation)
}

func (s *Service) ValidateInspection(ctx *UserContext, idInspection int64) error {
	if !ctx.IsApprobateur() {
		return ErrBesoinProfilApprobateur
	}
	return s.changeEtatInspection(ctx, idInspection, models.EtatAttenteValidation, models.EtatValide)
}

func (s *Service) RejectInspection(ctx *UserContext, idInspection int64) error {
	if !ctx.IsApprobateur() {
		return ErrBesoinProfilApprobateur
	}
	return s.changeEtatInspection(ctx, idInspection, models.EtatAttenteValidation, models.EtatEnCours)
}

func (s *Service) changeEtatInspection(ctx *UserContext, idInspection int64, fromEtat models.EtatInspection, toEtat models.EtatInspection) error {
	ok, err := s.repo.CheckEtatInspection(idInspection, []models.EtatInspection{fromEtat})
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}

	return s.repo.UpdateEtatInspection(ctx, idInspection, toEtat)
}
