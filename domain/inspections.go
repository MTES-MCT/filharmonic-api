package domain

import (
	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
)

var (
	ErrInspectionNotFound          = errors.NewErrForbidden("Inspection non trouvée")
	ErrClotureInspectionImpossible = errors.NewErrForbidden("Impossible de clore l'inspection")
)

type ListInspectionsFilter struct {
	Assigned bool `form:"assigned"`
}

func (s *Service) ListInspections(ctx *UserContext, filter ListInspectionsFilter) ([]models.Inspection, error) {
	return s.repo.ListInspections(ctx, filter)
}

func (s *Service) ListInspectionsFavorites(ctx *UserContext) ([]models.Inspection, error) {
	return s.repo.ListInspectionsFavorites(ctx)
}

func (s *Service) CreateInspection(ctx *UserContext, inspection models.Inspection) (int64, error) {
	if !ctx.IsInspecteur() {
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
	inspectionId, err := s.repo.CreateInspection(ctx, inspection)
	if err != nil {
		return 0, err
	}
	return inspectionId, s.addMissingThemes(inspection.Themes)
}

type InspectionFilter struct {
	OmitPointsDeControleNonPublies bool
}

func (s *Service) GetInspection(ctx *UserContext, id int64) (*models.Inspection, error) {
	inspection, err := s.repo.GetInspectionByID(ctx, id, InspectionFilter{})
	if err != nil {
		return nil, err
	}
	if inspection == nil {
		return nil, ErrInspectionNotFound
	}
	return inspection, nil
}

func (s *Service) UpdateInspection(ctx *UserContext, inspection models.Inspection) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}
	inspecteursIds := make([]int64, 0)
	for _, inspecteur := range inspection.Inspecteurs {
		inspecteursIds = append(inspecteursIds, inspecteur.Id)
	}
	if len(inspecteursIds) == 0 {
		return ErrInvalidInput
	}
	err := s.repo.CheckInspecteurAllowedInspection(ctx, inspection.Id)
	if err != nil {
		return err
	}
	ok, err := s.repo.CheckUsersInspecteurs(inspecteursIds)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}
	err = s.repo.UpdateInspection(ctx, inspection)
	if err != nil {
		return err
	}
	return s.addMissingThemes(inspection.Themes)
}

func (s *Service) PublishInspection(ctx *UserContext, idInspection int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}
	err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return err
	}
	err = s.changeEtatInspection(ctx, idInspection, models.EtatPreparation, models.EtatEnCours)
	if err != nil {
		return err
	}
	err = s.repo.CreateEvenement(ctx, models.EvenementPublicationInspection, idInspection, nil)
	return err
}

func (s *Service) AskValidateInspection(ctx *UserContext, idInspection int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}
	err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return err
	}

	inspection, err := s.repo.GetInspectionTypesConstatsSuiteByID(idInspection)
	if err != nil {
		return err
	}
	err = inspection.CheckCoherenceSuiteConstats()
	if err != nil {
		return err
	}
	err = s.changeEtatInspection(ctx, idInspection, models.EtatEnCours, models.EtatAttenteValidation)
	if err != nil {
		return err
	}
	err = s.repo.CreateEvenement(ctx, models.EvenementDemandeValidationInspection, idInspection, nil)
	return err
}

func (s *Service) ValidateInspection(ctx *UserContext, idInspection int64, rapportFile models.File) error {
	if !ctx.IsApprobateur() {
		return ErrBesoinProfilApprobateur
	}

	ok, err := s.repo.CheckEtatInspection(idInspection, []models.EtatInspection{models.EtatAttenteValidation})
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidInput
	}

	hasNonConformites, err := s.repo.CheckInspectionHasNonConformites(idInspection)
	if err != nil {
		return err
	}
	etatCible := models.EtatClos
	if hasNonConformites {
		etatCible = models.EtatTraitementNonConformites
	}
	storageId, err := s.storage.Put(rapportFile)
	if err != nil {
		return err
	}
	rapport := models.Rapport{
		Nom:       rapportFile.Nom,
		Type:      rapportFile.Type,
		Taille:    rapportFile.Taille,
		StorageId: storageId,
		AuteurId:  ctx.User.Id,
	}
	err = s.repo.CreateRapport(idInspection, rapport)
	if err != nil {
		return err
	}

	err = s.repo.ValidateInspection(idInspection, etatCible)
	if err != nil {
		return err
	}
	err = s.repo.CreateEvenement(ctx, models.EvenementValidationInspection, idInspection, nil)
	return err
}

func (s *Service) RejectInspection(ctx *UserContext, idInspection int64) error {
	if !ctx.IsApprobateur() {
		return ErrBesoinProfilApprobateur
	}
	err := s.changeEtatInspection(ctx, idInspection, models.EtatAttenteValidation, models.EtatEnCours)
	if err != nil {
		return err
	}
	err = s.repo.CreateEvenement(ctx, models.EvenementRejetValidationInspection, idInspection, nil)
	return err
}

func (s *Service) CloreInspection(ctx *UserContext, idInspection int64) error {
	if !ctx.IsInspecteur() {
		return ErrBesoinProfilInspecteur
	}
	ok, err := s.repo.CheckEtatInspection(idInspection, []models.EtatInspection{models.EtatTraitementNonConformites})
	if err != nil {
		return err
	}
	if !ok {
		return ErrBesoinEtatTraitementNonConformites
	}
	err = s.repo.CanCloreInspection(ctx, idInspection)
	if err != nil {
		return err
	}

	err = s.repo.UpdateEtatInspection(ctx, idInspection, models.EtatClos)
	if err != nil {
		return err
	}
	return s.repo.CreateEvenement(ctx, models.EvenementClotureInspection, idInspection, nil)
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

func (s *Service) AddFavoriToInspection(ctx *UserContext, idInspection int64) error {
	return s.repo.AddFavoriToInspection(ctx, idInspection)
}

func (s *Service) RemoveFavoriToInspection(ctx *UserContext, idInspection int64) error {
	return s.repo.RemoveFavoriToInspection(ctx, idInspection)
}
