package models

import "time"

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
	Annonce            bool                   `json:"annonce"`
	Origine            OrigineInspection      `json:"origine"`
	Circonstance       CirconstanceInspection `json:"circonstance"`
	DetailCirconstance string                 `json:"detail_circonstance"`
	Contexte           string                 `json:"contexte"`
	Etat               EtatInspection         `json:"etat"`
	EtablissementId    int64                  `sql:",notnull" json:"etablissement_id"`

	Etablissement *Etablissement    `json:"etablissement,omitempty"`
	Themes        []ThemeInspection `json:"themes,omitempty"`
	Inspecteurs   []User            `pg:"many2many:inspection_to_inspecteurs" json:"inspecteurs,omitempty"`
}

type InspectionToInspecteur struct {
	InspectionId int64 `sql:",pk"`
	UserId       int64 `sql:",pk"`

	Inspection *Inspection
	User       *User
}

type ThemeInspection struct {
	Id           int64  `json:"id"`
	Nom          string `json:"nom"`
	InspectionId int64  `json:"-"`

	Inspection *Inspection
}
