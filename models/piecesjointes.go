package models

import "io"

type PieceJointe struct {
	Id            int64  `json:"id"`
	Nom           string `json:"nom"`
	Type          string `json:"type"`
	Taille        int64  `json:"taille"`
	StorageId     string `json:"-" sql:",unique"`
	MessageId     int64  `json:"-"`
	CommentaireId int64  `json:"-"`
	AuteurId      int64  `json:"-" sql:",notnull"`

	Message     *Message     `json:"-"`
	Commentaire *Commentaire `json:"-"`
	Auteur      *User        `json:"-"`
}

type File struct {
	Content io.Reader
	Taille  int64
	Nom     string
	Type    string
}

type Rapport struct {
	Id        int64  `json:"id"`
	Nom       string `json:"nom"`
	Type      string `json:"type"`
	Taille    int64  `json:"taille"`
	StorageId string `json:"-" sql:",unique"`
	AuteurId  int64  `json:"-" sql:",notnull"`

	Auteur *User `json:"-"`
}
