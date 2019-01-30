package models

import (
	"encoding/json"
	"time"
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
	CirconstanceAutres   CirconstanceInspection = "autres"
)

type EtatInspection string

const (
	EtatPreparation       EtatInspection = "preparation"
	EtatEnCours           EtatInspection = "en_cours"
	EtatAttenteValidation EtatInspection = "attente_validation"
	EtatValide            EtatInspection = "valide"
)

type Inspection struct {
	Id                 int64                  `json:"id"`
	Date               time.Time              `json:"date" sql:"type:date"`
	Type               TypeInspection         `json:"type"`
	Annonce            bool                   `json:"annonce" sql:",notnull"`
	Origine            OrigineInspection      `json:"origine"`
	Circonstance       CirconstanceInspection `json:"circonstance"`
	DetailCirconstance string                 `json:"detail_circonstance"`
	Contexte           string                 `json:"contexte"`
	Etat               EtatInspection         `json:"etat"`
	EtablissementId    int64                  `json:"etablissement_id" sql:",notnull"`
	Themes             []string               `json:"themes" sql:",array"`
	SuiteId            int64                  `json:"-" sql:"on_delete:SET NULL"`

	Commentaires     []Commentaire     `json:"commentaires,omitempty"`
	Etablissement    *Etablissement    `json:"etablissement,omitempty"`
	Inspecteurs      []User            `pg:"many2many:inspection_to_inspecteurs" json:"inspecteurs,omitempty"`
	PointsDeControle []PointDeControle `json:"points_de_controle,omitempty"`
	Suite            *Suite            `json:"suite,omitempty"`
}

func (i *Inspection) MarshalJSON() ([]byte, error) {
	type Alias Inspection
	return json.Marshal(&struct {
		Date string `json:"date"`
		*Alias
	}{
		Date:  i.Date.Format("2006-01-02"),
		Alias: (*Alias)(i),
	})
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
	Id       int64     `json:"id"`
	Type     TypeSuite `json:"type"`
	Synthese string    `json:"synthese"`
	Delai    int       `json:"delai"`
}
