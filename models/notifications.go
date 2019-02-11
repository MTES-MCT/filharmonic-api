package models

import "time"

type TypeEvenement string

const (
	EvenementCreationMessage             TypeEvenement = "message"
	EvenementCreationCommentaire         TypeEvenement = "commentaire"
	EvenementLectureMessage              TypeEvenement = "lecture_message"
	EvenementCommentaireGeneral          TypeEvenement = "commentaire_general"
	EvenementCreationInspection          TypeEvenement = "creation_inspection"
	EvenementModificationInspection      TypeEvenement = "modification_inspection"
	EvenementPublicationInspection       TypeEvenement = "publication_inspection"
	EvenementDemandeValidationInspection TypeEvenement = "demande_validation_inspection"
	EvenementRejetValidationInspection   TypeEvenement = "rejet_validation_inspection"
	EvenementValidationInspection        TypeEvenement = "validation_inspection"
	EvenementCreationPointDeControle     TypeEvenement = "creation_point_de_controle"
	EvenementModificationPointDeControle TypeEvenement = "modification_point_de_controle"
	EvenementSuppressionPointDeControle  TypeEvenement = "suppression_point_de_controle"
	EvenementPublicationPointDeControle  TypeEvenement = "publication_point_de_controle"
	EvenementCreationConstat             TypeEvenement = "creation_constat"
	EvenementSuppressionConstat          TypeEvenement = "suppression_constat"
	EvenementCreationSuite               TypeEvenement = "creation_suite"
	EvenementModificationSuite           TypeEvenement = "modification_suite"
	EvenementSuppressionSuite            TypeEvenement = "suppression_suite"
)

type Notification struct {
	Id             int64 `json:"id"`
	Lue            bool  `json:"lue" sql:",notnull"`
	EvenementId    int64 `json:"evenement_id" sql:",notnull"`
	DestinataireId int64 `json:"destinataire_id" sql:",notnull"`

	Destinataire *User      `json:"-"`
	Evenement    *Evenement `json:"evenement,omitempty"`
}

type Evenement struct {
	Id           int64                  `json:"id"`
	Type         TypeEvenement          `json:"type"`
	CreatedAt    time.Time              `json:"created_at"`
	Data         map[string]interface{} `json:"data" sql:"type:jsonb"`
	AuteurId     int64                  `json:"auteur_id" sql:",notnull"`
	InspectionId int64                  `json:"inspection_id" sql:",notnull"`

	Auteur     *User       `json:"auteur,omitempty"`
	Inspection *Inspection `json:"-"`
}
