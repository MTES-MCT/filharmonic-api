package models

type Etablissement struct {
	ID      int64
	S3IC    string
	Raison  string
	Adresse string
}
type User struct {
	ID       int64
	Email    string
	Password string
	Profile  string
}
