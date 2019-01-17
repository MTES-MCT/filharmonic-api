package models

type Profil string

const (
	ProfilInspecteur  Profil = "inspecteur"
	ProfilExploitant  Profil = "exploitant"
	ProfilApprobateur Profil = "approbateur"
)

type Etablissement struct {
	Id          int64  `json:"id"`
	S3IC        string `json:"s3ic" sql:",unique"`
	Raison      string `json:"raison"`
	Adresse     string `json:"adresse"`
	Exploitants []User `pg:"many2many:etablissement_to_exploitants"`
}
type User struct {
	Id       int64  `json:"id"`
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
