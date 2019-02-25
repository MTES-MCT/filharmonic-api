package domain

import (
	"strings"
	"time"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/MTES-MCT/filharmonic-api/util"
)

type Person struct {
	Nom       string
	Telephone string
	Email     string
}

type LettreAnnonce struct {
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

func NewLettreAnnonce(inspection *models.Inspection) LettreAnnonce {
	lettre := LettreAnnonce{}
	lettre.Inspection = *inspection
	lettre.DateLettre = util.FormatDate(time.Now())
	lettre.DateInspection = util.FormatDate(inspection.Date.Time)
	lettre.Prefet = "PRÉFET DU RHÔNE"
	lettre.DepartementUnite = "Rhône"
	lettre.HeureInspection = "9h30"
	lettre.NomDirection = "DREAL Auvergne-Rhône-Alpes"
	lettre.URLDirection = "www.auvergne.rhone-alpes.developpement-durable.gouv.fr"
	lettre.VilleUnite = "Lyon"

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

func (s *Service) GenererLettreAnnonce(ctx *UserContext, idInspection int64) (*models.PieceJointeFile, error) {
	if !ctx.IsInspecteur() {
		return nil, ErrBesoinProfilInspecteur
	}
	inspection, err := s.repo.GetInspectionByID(ctx, idInspection)
	if err != nil {
		return nil, err
	}
	if inspection == nil {
		return nil, ErrInvalidInput
	}

	contenuLettre, err := s.templateService.RenderLettreAnnonce(NewLettreAnnonce(inspection))
	if err != nil {
		return nil, err
	}
	return &models.PieceJointeFile{
		Nom:     "lettre-annonce.odt",
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
	rapport.Prefet = "PRÉFET DU RHÔNE"
	rapport.DepartementUnite = "Rhône"
	rapport.NomDirection = "DREAL Auvergne-Rhône-Alpes"
	rapport.URLDirection = "www.auvergne.rhone-alpes.developpement-durable.gouv.fr"
	rapport.VilleUnite = "Lyon"

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

func (s *Service) GenererRapport(ctx *UserContext, idInspection int64) (*models.PieceJointeFile, error) {
	if !ctx.IsInspecteur() {
		return nil, ErrBesoinProfilInspecteur
	}
	inspection, err := s.repo.GetInspectionByID(ctx, idInspection)
	if err != nil {
		return nil, err
	}
	if inspection == nil {
		return nil, ErrInvalidInput
	}

	contenuRapport, err := s.templateService.RenderRapport(NewRapport(inspection))
	if err != nil {
		return nil, err
	}
	return &models.PieceJointeFile{
		Nom:     "rapport.odt",
		Type:    "application/vnd.oasis.opendocument.text",
		Taille:  int64(len(contenuRapport)),
		Content: strings.NewReader(contenuRapport),
	}, nil
}
