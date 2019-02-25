package models

type RegimeEtablissement string

const (
	RegimeAucun          RegimeEtablissement = "aucun"
	RegimeAutorisation   RegimeEtablissement = "autorisation"
	RegimeDeclaration    RegimeEtablissement = "declaration"
	RegimeEnregistrement RegimeEtablissement = "enregistrement"
	RegimeInconnu        RegimeEtablissement = "inconnu"
)

func RegimeFromString(regime string) RegimeEtablissement {
	switch regime {
	case "NULL":
		return RegimeAucun
	case "Autorisation":
		return RegimeAutorisation
	case "Déclaration":
		return RegimeDeclaration
	case "Enregistrement":
		return RegimeEnregistrement
	default:
		return RegimeInconnu
	}
}

type Etablissement struct {
	Id         int64               `json:"id"`
	S3IC       string              `json:"s3ic" sql:",unique"`
	Nom        string              `json:"nom"`
	Raison     string              `json:"raison"`
	Seveso     string              `json:"seveso"`
	Activite   string              `json:"activite"`
	Iedmtd     bool                `json:"iedmtd" sql:",notnull"`
	Adresse1   string              `json:"adresse1"`
	Adresse2   string              `json:"adresse2"`
	CodePostal string              `json:"code_postal"`
	Commune    string              `json:"commune"`
	Regime     RegimeEtablissement `json:"regime"`

	Exploitants []User       `pg:"many2many:etablissement_to_exploitants" json:"exploitants,omitempty"`
	Inspections []Inspection `json:"inspections,omitempty"`
}

type EtablissementToExploitant struct {
	EtablissementId int64 `sql:",pk"`
	UserId          int64 `sql:",pk"`

	Etablissement *Etablissement
	User          *User
}

type FindEtablissementResults struct {
	Etablissements []Etablissement `json:"etablissements"`
	Total          int             `json:"total"`
}
