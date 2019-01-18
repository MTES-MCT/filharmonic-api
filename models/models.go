package models

type Profil string

const (
	ProfilInspecteur  Profil = "inspecteur"
	ProfilExploitant  Profil = "exploitant"
	ProfilApprobateur Profil = "approbateur"
)

type Etablissement struct {
	Id       int64  `json:"id"`
	S3IC     string `json:"s3ic" sql:",unique"`
	Nom      string `json:"nom"`
	Raison   string `json:"raison"`
	Adresse  string `json:"adresse"`
	Seveso   string `json:"seveso"`
	Activite string `json:"activite"`
	Iedmtd   bool   `json:"iedmtd"`

	Exploitants []User `pg:"many2many:etablissement_to_exploitants"`
}
type User struct {
	Id       int64  `json:"id"`
	Nom      string `json:"nom"`
	Prenom   string `json:"prenom"`
	Email    string `json:"email" sql:",unique"`
	Password string `json:"password"`
	Profile  Profil `json:"profile"`
}

type EtablissementToExploitant struct {
	EtablissementId int64 `sql:",pk"`
	UserId          int64 `sql:",pk"`

	Etablissement *Etablissement
	User          *User
}
