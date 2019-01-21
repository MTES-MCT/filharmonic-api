package models

type Profil string

const (
	ProfilInspecteur  Profil = "inspecteur"
	ProfilExploitant  Profil = "exploitant"
	ProfilApprobateur Profil = "approbateur"
)

type User struct {
	Id       int64  `json:"id"`
	Nom      string `json:"nom"`
	Prenom   string `json:"prenom"`
	Email    string `json:"email" sql:",unique"`
	Password string `json:"password,omitempty"`
	Profile  Profil `json:"profile"`
}
