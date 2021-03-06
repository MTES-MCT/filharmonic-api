package models

import (
	"time"

	"github.com/MTES-MCT/filharmonic-api/util"
	"github.com/go-pg/pg/types"
)

type PointDeControle struct {
	Id                       int64     `json:"id"`
	Sujet                    string    `json:"sujet"`
	ReferencesReglementaires []string  `json:"references_reglementaires" sql:",array"`
	Publie                   bool      `json:"publie" sql:",notnull,default:false"`
	InspectionId             int64     `json:"-" sql:",notnull"`
	DeletedAt                time.Time `json:"-" pg:",soft_delete"`
	ConstatId                int64     `json:"-" sql:"on_delete:SET NULL"`
	Order                    int64     `json:"order"`

	Inspection *Inspection `json:"-"`
	Constat    *Constat    `json:"constat,omitempty"`
	Messages   []Message   `json:"messages,omitempty"`
}

type Message struct {
	Id                            int64     `json:"id"`
	Message                       string    `json:"message"`
	Date                          time.Time `json:"date"`
	Lu                            bool      `json:"lu" sql:",notnull,default:false"`
	Interne                       bool      `json:"interne" sql:",notnull,default:false"`
	AuteurId                      int64     `json:"-" sql:",notnull"`
	PointDeControleId             int64     `json:"-" sql:",notnull"`
	EtapeTraitementNonConformites bool      `json:"etape_traitement_non_conformites" sql:",notnull,default:false"`

	Auteur          *User            `json:"auteur,omitempty"`
	PointDeControle *PointDeControle `json:"-"`
	PiecesJointes   []PieceJointe    `json:"pieces_jointes,omitempty"`
}

type TypeConstat string

const (
	TypeConstatInconnu     TypeConstat = "inconnu"
	TypeConstatObservation TypeConstat = "observation"
	TypeConstatConforme    TypeConstat = "conforme"
	TypeConstatNonConforme TypeConstat = "non_conforme"
)

type Constat struct {
	Id                                 int64           `json:"id"`
	Type                               TypeConstat     `json:"type"`
	Remarques                          string          `json:"remarques"`
	DateResolution                     types.NullTime  `json:"date_resolution" sql:"type:timestamptz"`
	EcheanceResolution                 util.DateString `json:"echeance_resolution,omitempty" sql:"type:date"`
	DelaiNombre                        int32           `json:"delai_nombre"`
	DelaiUnite                         string          `json:"delai_unite"`
	NotificationRappelEcheanceEnvoyee  bool            `json:"_" sql:",notnull,default:false"`
	NotificationEcheanceExpireeEnvoyee bool            `json:"_" sql:",notnull,default:false"`
}
