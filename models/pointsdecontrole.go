package models

import "time"

type PointDeControle struct {
	Id                       int64     `json:"id"`
	Sujet                    string    `json:"sujet"`
	ReferencesReglementaires []string  `json:"references_reglementaires" sql:",array"`
	Publie                   bool      `json:"publie" sql:",notnull"`
	InspectionId             int64     `json:"-" sql:",notnull"`
	DeletedAt                time.Time `json:"-" pg:",soft_delete"`
	ConstatId                int64     `json:"-" sql:"on_delete:SET NULL"`

	Inspection *Inspection `json:"-"`
	Constat    *Constat    `json:"constat,omitempty"`
	Messages   []Message   `json:"messages,omitempty"`
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

type TypeConstat string

const (
	TypeConstatObservation TypeConstat = "observation"
	TypeConstatConforme    TypeConstat = "conforme"
	TypeConstatNonConforme TypeConstat = "non_conforme"
)

type Constat struct {
	Id        int64       `json:"id"`
	Type      TypeConstat `json:"type"`
	Remarques string      `json:"remarques"`
}