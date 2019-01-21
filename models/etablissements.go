package models

type Etablissement struct {
	Id       int64  `json:"id"`
	S3IC     string `json:"s3ic" sql:",unique"`
	Nom      string `json:"nom"`
	Raison   string `json:"raison"`
	Adresse  string `json:"adresse"`
	Seveso   string `json:"seveso"`
	Activite string `json:"activite"`
	Iedmtd   bool   `json:"iedmtd" sql:",notnull"`

	Exploitants []User `pg:"many2many:etablissement_to_exploitants" json:"exploitants,omitempty"`
}

type EtablissementToExploitant struct {
	EtablissementId int64 `sql:",pk"`
	UserId          int64 `sql:",pk"`

	Etablissement *Etablissement
	User          *User
}
