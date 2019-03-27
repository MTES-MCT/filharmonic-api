package domain

import (
	"strings"
	"time"

	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
)

type Person struct {
	Nom       string
	Telephone string
	Email     string
}

type Lettre struct {
	Inspection       models.Inspection
	DateInspection   string
	DateLettre       string
	VilleUnite       string
	DepartementUnite string
	NomDirection     string
	URLDirection     string
	Prefet           string
	HeureInspection  string
	Auteur           Person
	AutresAuteurs    []Person
	Exploitant       Person
}

func NewLettre(inspection *models.Inspection) Lettre {
	lettre := Lettre{}
	lettre.Inspection = *inspection
	lettre.DateLettre = util.FormatDate(time.Now())
	lettre.DateInspection = util.FormatDate(inspection.Date.Time)
	if inspection.Etablissement.Departement != nil {
		lettre.Prefet = "Préfet " + inspection.Etablissement.Departement.AvecCharniere()
		lettre.DepartementUnite = inspection.Etablissement.Departement.AvecCharniere()
		lettre.NomDirection = "Direction régionale de l'environnement, de l'aménagement et du logement " + inspection.Etablissement.Departement.Region
	}
	lettre.HeureInspection = "9h30"
	lettre.URLDirection = "www.developpement-durable.gouv.fr"
	lettre.VilleUnite = inspection.Etablissement.Commune

	if len(inspection.Inspecteurs) > 0 {
		lettre.Auteur = Person{
			Nom:       inspection.Inspecteurs[0].Prenom + " " + inspection.Inspecteurs[0].Nom,
			Telephone: "04 01 02 03 04",
			Email:     inspection.Inspecteurs[0].Email,
		}
	}
	if len(inspection.Etablissement.Exploitants) > 0 {
		lettre.Exploitant = Person{
			Nom: inspection.Etablissement.Exploitants[0].Prenom + " " + inspection.Etablissement.Exploitants[0].Nom,
		}
	}

	lettre.AutresAuteurs = []Person{}
	if len(inspection.Inspecteurs) > 1 {
		for _, inspecteur := range inspection.Inspecteurs[1:] {
			lettre.AutresAuteurs = append(lettre.AutresAuteurs, Person{
				Nom: inspecteur.Prenom + " " + inspecteur.Nom,
			})
		}
	}
	return lettre
}

func (s *Service) GenererLettreAnnonce(ctx *UserContext, idInspection int64) (*models.File, error) {
	if !ctx.IsInspecteur() {
		return nil, ErrBesoinProfilInspecteur
	}
	err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return nil, err
	}
	filter := InspectionFilter{
		OmitPointsDeControleNonPublies: true,
	}
	inspection, err := s.repo.GetInspectionByID(ctx, idInspection, filter)
	if err != nil {
		return nil, err
	}
	if inspection == nil {
		return nil, ErrInvalidInput
	}
	if inspection.Etat != models.EtatEnCours {
		return nil, errors.NewErrForbidden("L'inspection doit être à l'état en cours.")
	}

	result, err := s.templateService.RenderODTLettreAnnonce(NewLettre(inspection))
	if err != nil {
		return nil, err
	}
	contenuLettre := result.Text
	return &models.File{
		Nom:     "lettre-annonce.odt",
		Type:    "application/vnd.oasis.opendocument.text",
		Taille:  int64(len(contenuLettre)),
		Content: strings.NewReader(contenuLettre),
	}, nil
}

func (s *Service) GenererLettreSuite(ctx *UserContext, idInspection int64) (*models.File, error) {
	if !ctx.IsInspecteur() {
		return nil, ErrBesoinProfilInspecteur
	}
	err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return nil, err
	}
	inspection, err := s.repo.GetInspectionByID(ctx, idInspection, InspectionFilter{})
	if err != nil {
		return nil, err
	}
	if inspection == nil {
		return nil, ErrInvalidInput
	}
	if inspection.Etat != models.EtatTraitementNonConformites && inspection.Etat != models.EtatClos {
		return nil, errors.NewErrForbidden("L'inspection ne doit pas être à l'état de traitement des non conformités ni clos.")
	}

	if inspection.Suite == nil {
		return nil, ErrInvalidInput
	}

	result, err := s.templateService.RenderODTLettreSuite(NewLettre(inspection))
	if err != nil {
		return nil, err
	}
	contenuLettre := result.Text
	return &models.File{
		Nom:     "lettre-suite.odt",
		Type:    "application/vnd.oasis.opendocument.text",
		Taille:  int64(len(contenuLettre)),
		Content: strings.NewReader(contenuLettre),
	}, nil
}

type Rapport struct {
	Inspection       models.Inspection
	DateInspection   string
	DateRapport      string
	VilleUnite       string
	DepartementUnite string
	NomDirection     string
	URLDirection     string
	Prefet           string
	Auteur           Person
	AutresAuteurs    []Person
	Exploitant       Person
}

func NewRapport(inspection *models.Inspection) Rapport {
	rapport := Rapport{}
	rapport.Inspection = *inspection
	rapport.DateRapport = util.FormatDate(time.Now())
	rapport.DateInspection = util.FormatDate(inspection.Date.Time)
	if inspection.Etablissement.Departement != nil {
		rapport.Prefet = "Préfet " + inspection.Etablissement.Departement.AvecCharniere()
		rapport.DepartementUnite = inspection.Etablissement.Departement.AvecCharniere()
		rapport.NomDirection = "Direction régionale de l'environnement, de l'aménagement et du logement " + inspection.Etablissement.Departement.Region
	}
	rapport.URLDirection = "www.developpement-durable.gouv.fr"
	rapport.VilleUnite = inspection.Etablissement.Commune

	if len(inspection.Inspecteurs) > 0 {
		rapport.Auteur = Person{
			Nom:       inspection.Inspecteurs[0].Prenom + " " + inspection.Inspecteurs[0].Nom,
			Telephone: "04 01 02 03 04",
			Email:     inspection.Inspecteurs[0].Email,
		}
	}
	if len(inspection.Etablissement.Exploitants) > 0 {
		rapport.Exploitant = Person{
			Nom: inspection.Etablissement.Exploitants[0].Prenom + " " + inspection.Etablissement.Exploitants[0].Nom,
		}
	}

	rapport.AutresAuteurs = []Person{}
	if len(inspection.Inspecteurs) > 1 {
		for _, inspecteur := range inspection.Inspecteurs[1:] {
			rapport.AutresAuteurs = append(rapport.AutresAuteurs, Person{
				Nom: inspecteur.Prenom + " " + inspecteur.Nom,
			})
		}
	}
	return rapport
}

func (s *Service) GenererRapport(ctx *UserContext, idInspection int64) (*models.File, error) {
	if !ctx.IsInspecteur() {
		return nil, ErrBesoinProfilInspecteur
	}
	err := s.repo.CheckInspecteurAllowedInspection(ctx, idInspection)
	if err != nil {
		return nil, err
	}
	filter := InspectionFilter{
		OmitPointsDeControleNonPublies: true,
	}
	inspection, err := s.repo.GetInspectionByID(ctx, idInspection, filter)
	if err != nil {
		return nil, err
	}
	if inspection == nil {
		return nil, ErrInvalidInput
	}
	if inspection.Etat != models.EtatEnCours {
		return nil, errors.NewErrForbidden("L'inspection doit être à l'état en cours.")
	}

	result, err := s.templateService.RenderODTRapport(NewRapport(inspection))
	if err != nil {
		return nil, err
	}
	contenuRapport := result.Text
	return &models.File{
		Nom:     "rapport.odt",
		Type:    "application/vnd.oasis.opendocument.text",
		Taille:  int64(len(contenuRapport)),
		Content: strings.NewReader(contenuRapport),
	}, nil
}
