package models

import "time"

type Canevas struct {
	Id          int64       `json:"id"`
	Nom         string      `json:"nom" sql:",unique,notnull"`
	AuteurId    int64       `json:"-" sql:",notnull"`
	DataVersion int         `json:"data_version"`
	Data        CanevasData `json:"data"`
	CreatedAt   time.Time   `json:"created_at"`

	Auteur *User `json:"auteur,omitempty"`
}

type CanevasData struct {
	PointsDeControle []CanevasPointDeControle `json:"points_de_controle"`
}

type CanevasPointDeControle struct {
	Sujet                    string   `json:"sujet"`
	ReferencesReglementaires []string `json:"references_reglementaires" sql:",array"`
	Message                  string   `json:"message,omitempty"`
}
