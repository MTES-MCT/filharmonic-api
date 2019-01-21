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
	EtablissementId    int64                  `sql:",notnull" json:"etablissement_id"`

	Commentaires     []Commentaire     `json:"commentaires,omitempty"`
	Etablissement    *Etablissement    `json:"etablissement,omitempty"`
	Themes           []string          `json:"themes,omitempty" sql:",array"`
	Inspecteurs      []User            `pg:"many2many:inspection_to_inspecteurs" json:"inspecteurs,omitempty"`
	PointsDeControle []PointDeControle `json:"points_de_controle,omitempty"`
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

type PointDeControle struct {
	Id                       int64    `json:"id"`
	Sujet                    string   `json:"sujet"`
	ReferencesReglementaires []string `json:"references_reglementaires" sql:",array"`
	Publie                   bool     `json:"publie" sql:",notnull"`
	InspectionId             int64    `json:"-" sql:",notnull"`

	Inspection *Inspection `json:"-"`
	Messages   []Message   `json:"messages,omitempty"`
}

type Commentaire struct {
	Id           int64     `json:"id"`
	Message      string    `json:"message"`
	Date         time.Time `json:"date"`
	AuteurId     int64     `json:"-" sql:",notnull"`
	InspectionId int64     `json:"-" sql:",notnull"`

	Auteur     *User       `json:"auteur,omitempty"`
	Inspection *Inspection `json:"-"`
}

type Message struct {
	Id                int64     `json:"id"`
	Message           string    `json:"message"`
	Date              time.Time `json:"date"`
	Lu                bool      `json:"lu" sql:",notnull"`
	Interne           bool      `json:"interne" sql:",notnull"`
	AuteurId          int64     `json:"-" sql:",notnull"`
	PointDeControleId int64     `json:"-" sql:",notnull"`

	Auteur          *User            `json:"auteur,omitempty"`
	PointDeControle *PointDeControle `json:"-"`
}
