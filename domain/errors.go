package domain

import "github.com/MTES-MCT/filharmonic-api/errors"

var (
	// common errors
	ErrBesoinProfilApprobateur = errors.NewErrForbidden("Il faut être approbateur")
	ErrBesoinProfilExploitant  = errors.NewErrForbidden("Il faut être exploitant")
	ErrBesoinProfilInspecteur  = errors.NewErrForbidden("Il faut être inspecteur")
	ErrNonAssigneInspection    = errors.NewErrForbidden("Il faut être assigné à l'inspection")
	ErrInvalidInput            = errors.NewErrBadInput("Données invalides")
)
