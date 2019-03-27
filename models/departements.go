package models

type Departement struct {
	Id              int64  `json:"id"`
	CodeInsee       string `json:"code_insee" sql:",unique,notnull"`
	Nom             string `json:"nom"`
	Charniere       string `json:"charniere"`
	Region          string `json:"region"`
	CharniereRegion string `json:"charniere_region"`
}

func (d Departement) AvecCharniere() string {
	return d.Charniere + d.Nom
}

func (d Departement) RegionAvecCharniere() string {
	return d.CharniereRegion + d.Region
}
