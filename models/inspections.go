package models

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/errors"
	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/go-pg/pg/types"
)

type TypeInspection string

const (
	TypeApprofondi TypeInspection = "approfondi"
	TypeCourant    TypeInspection = "courant"
	TypePonctuel   TypeInspection = "ponctuel"
)

type OrigineInspection string

const (
	OriginePlanControle     OrigineInspection = "plan_de_controle"
	OrigineCirconstancielle OrigineInspection = "circonstancielle"
)

type CirconstanceInspection string

const (
	CirconstanceIncident CirconstanceInspection = "incident"
	CirconstancePlainte  CirconstanceInspection = "plainte"
	CirconstanceAutre    CirconstanceInspection = "autre"
)

type EtatInspection string

const (
	EtatInconnu                  EtatInspection = "inconnu"
	EtatPreparation              EtatInspection = "preparation"
	EtatEnCours                  EtatInspection = "en_cours"
	EtatAttenteValidation        EtatInspection = "attente_validation"
	EtatTraitementNonConformites EtatInspection = "traitement_non_conformites"
	EtatClos                     EtatInspection = "close"
)

type Inspection struct {
	Id                   int64                  `json:"id"`
	Date                 util.DateString        `json:"date" sql:"type:date"`
	Type                 TypeInspection         `json:"type"`
	Annonce              bool                   `json:"annonce" sql:",notnull,default:false"`
	Origine              OrigineInspection      `json:"origine"`
	Circonstance         CirconstanceInspection `json:"circonstance"`
	DetailCirconstance   string                 `json:"detail_circonstance"`
	Contexte             string                 `json:"contexte"`
	Etat                 EtatInspection         `json:"etat"`
	EtablissementId      int64                  `json:"etablissement_id" sql:",notnull"`
	Themes               []string               `json:"themes" sql:",array"`
	SuiteId              int64                  `json:"-" sql:"on_delete:SET NULL"`
	DateValidation       types.NullTime         `json:"date_validation" sql:"type:timestamptz"`
	RapportId            int64                  `json:"-" sql:"on_delete:SET NULL"`
	ValidationRejetee    bool                   `json:"validation_rejetee,omitempty" sql:",notnull,default:false"`
	MotifRejetValidation string                 `json:"motif_rejet_validation,omitempty"`

	Commentaires     []Commentaire     `json:"commentaires,omitempty"`
	Etablissement    *Etablissement    `json:"etablissement,omitempty"`
	Inspecteurs      []User            `pg:"many2many:inspection_to_inspecteurs" json:"inspecteurs,omitempty"`
	PointsDeControle []PointDeControle `json:"points_de_controle,omitempty"`
	Suite            *Suite            `json:"suite,omitempty"`
	Evenements       []Evenement       `json:"evenements,omitempty"`
	Rapport          *Rapport          `json:"rapport,omitempty"`

	NbMessagesNonLus int `json:"nb_messages_non_lus" sql:"-"`
}

type InspectionToInspecteur struct {
	InspectionId int64 `sql:",pk"`
	UserId       int64 `sql:",pk"`

	Inspection *Inspection
	User       *User
}

type Commentaire struct {
	Id           int64     `json:"id"`
	Message      string    `json:"message"`
	Date         time.Time `json:"date"`
	AuteurId     int64     `json:"-" sql:",notnull"`
	InspectionId int64     `json:"-" sql:",notnull"`

	Auteur        *User         `json:"auteur,omitempty"`
	Inspection    *Inspection   `json:"-"`
	PiecesJointes []PieceJointe `json:"pieces_jointes,omitempty"`
}

type TypeSuite string

const (
	TypeSuiteAucune                   TypeSuite = "aucune"
	TypeSuiteObservation              TypeSuite = "observation"
	TypeSuitePropositionMiseEnDemeure TypeSuite = "proposition_mise_en_demeure"
	TypeSuitePropositionRenforcement  TypeSuite = "proposition_renforcement"
	TypeSuiteAutre                    TypeSuite = "autre"
)

type Suite struct {
	Id          int64     `json:"id"`
	Type        TypeSuite `json:"type"`
	Synthese    string    `json:"synthese"`
	PenalEngage bool      `json:"penal_engage" sql:",notnull,default:false"`
}

var (
	ErrSuiteManquante             = errors.NewErrBadInput("Pas de suite")
	ErrConstatManquant            = errors.NewErrBadInput("Un constat est manquant")
	ErrPointDeControleNonPublie   = errors.NewErrBadInput("Un point de contrôle n'est pas publié")
	ErrPresenceConstatNonConforme = errors.NewErrBadInput("Un constat n'est pas conforme")
	ErrAbsenceConstatNonConforme  = errors.NewErrBadInput("Aucun constat n'est non-conforme")
)

func (inspection *Inspection) CheckCoherenceSuiteConstats() error {
	if inspection.Suite == nil {
		return ErrSuiteManquante
	}
	if inspection.Suite.Type == TypeSuiteAucune {
		// si suite = aucune, tous les constats doivent être conformes
		for _, pointDeControle := range inspection.PointsDeControle {
			if !pointDeControle.Publie {
				return ErrPointDeControleNonPublie
			}
			if pointDeControle.Constat == nil {
				return ErrConstatManquant
			}
			if pointDeControle.Constat.Type != TypeConstatConforme {
				return ErrPresenceConstatNonConforme
			}
		}
	} else {
		// si suite <> aucune, au moins un constat doit être non-conforme
		auMoinsUnConstatNonConforme := false
		for _, pointDeControle := range inspection.PointsDeControle {
			if !pointDeControle.Publie {
				return ErrPointDeControleNonPublie
			}
			if pointDeControle.Constat == nil {
				return ErrConstatManquant
			}
			if pointDeControle.Constat.Type != TypeConstatConforme {
				auMoinsUnConstatNonConforme = true
				break
			}
		}
		if !auMoinsUnConstatNonConforme {
			return ErrAbsenceConstatNonConforme
		}
	}
	return nil
}
