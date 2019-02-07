package models

import "time"

type TypeEvenement string

const (
	CreationMessage             TypeEvenement = "message"
	LectureMessage              TypeEvenement = "lecture_message"
	CommentaireGeneral          TypeEvenement = "commentaire_general"
	CreationCommentaire         TypeEvenement = "commentaire"
	CreationInspection          TypeEvenement = "creation_inspection"
	ModificationInspection      TypeEvenement = "modification_inspection"
	PublicationInspection       TypeEvenement = "publication_inspection"
	DemandeValidationInspection TypeEvenement = "demande_validation_inspection"
	RejetValidationInspection   TypeEvenement = "rejet_validation_inspection"
	ValidationInspection        TypeEvenement = "validation_inspection"
	CreationPointDeControle     TypeEvenement = "creation_point_de_controle"
	ModificationPointDeControle TypeEvenement = "modification_point_de_controle"
	SuppressionPointDeControle  TypeEvenement = "suppression_point_de_controle"
	PublicationPointDeControle  TypeEvenement = "publication_point_de_controle"
	CreationConstat             TypeEvenement = "creation_constat"
	SuppressionConstat          TypeEvenement = "suppression_constat"
	CreationSuite               TypeEvenement = "creation_suite"
	ModificationSuite           TypeEvenement = "modification_suite"
	SuppressionSuite            TypeEvenement = "suppression_suite"
)

type Notification struct {
	Id          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Lue         bool      `json:"lue" sql:",notnull"`
	EvenementId int64     `json:"-" sql:",notnull"`
	LecteurId   int64     `json:"-" sql:",notnull"`

	Lecteur   *User      `json:"lecteur,omitempty"`
	Evenement *Evenement `json:"evenement,omitempty"`
}

type Evenement struct {
	Id           int64         `json:"id"`
	Type         TypeEvenement `json:"type"`
	CreatedAt    time.Time     `json:"created_at"`
	Data         string        `json:"data" sql:"type:json"`
	AuteurId     int64         `json:"-" sql:",notnull"`
	InspectionId int64         `json:"-" sql:",notnull"`

	Auteur     *User       `json:"-"`
	Inspection *Inspection `json:"-"`
}
