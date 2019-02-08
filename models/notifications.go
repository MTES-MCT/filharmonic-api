package models

import "time"

type TypeEvenement string

const (
	CreationMessage             TypeEvenement = "message"
	CreationCommentaire         TypeEvenement = "commentaire"
	LectureMessage              TypeEvenement = "lecture_message"
	CommentaireGeneral          TypeEvenement = "commentaire_general"
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
	ReadAt      time.Time `json:"read_at"`
	Lue         bool      `json:"lue" sql:",notnull"`
	EvenementId int64     `json:"evenement_id" sql:",notnull"`
	LecteurId   int64     `json:"lecteur_id"`

	Lecteur   *User      `json:"lecteur,omitempty"`
	Evenement *Evenement `json:"evenement,omitempty"`
}

type Evenement struct {
	Id           int64         `json:"id"`
	Type         TypeEvenement `json:"type"`
	CreatedAt    time.Time     `json:"created_at"`
	Data         string        `json:"data" sql:"type:json"`
	AuteurId     int64         `json:"auteur_id" sql:",notnull"`
	InspectionId int64         `json:"inspection_id" sql:",notnull"`

	Auteur     *User       `json:"auteur,omitempty"`
	Inspection *Inspection `json:"-"`
}
