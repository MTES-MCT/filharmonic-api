package models

type Profil string

const (
	ProfilInspecteur  Profil = "inspecteur"
	ProfilExploitant  Profil = "exploitant"
	ProfilApprobateur Profil = "approbateur"
)

type User struct {
	Id      int64        `json:"id"`
	Nom     string       `json:"nom"`
	Prenom  string       `json:"prenom"`
	Email   string       `json:"email" sql:",unique"`
	Profile Profil       `json:"profile"`
	Favoris []Inspection `pg:"many2many:user_to_favoris" json:"favoris,omitempty"`
}

type UserToFavori struct {
	InspectionId int64 `sql:",pk"`
	UserId       int64 `sql:",pk"`

	Inspection *Inspection
	User       *User
}
