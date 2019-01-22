package domain

import "errors"

var (
	ErrBesoinProfilApprobateur = errors.New("Il faut être approbateur")
	ErrBesoinProfilExploitant  = errors.New("Il faut être exploitant")
	ErrBesoinProfilInspecteur  = errors.New("Il faut être inspecteur")

	ErrInvalidInput = errors.New("Données invalides")
)
