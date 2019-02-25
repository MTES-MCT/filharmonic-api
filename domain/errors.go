package domain

// when the input is bad (e.g. malformed query, invalid object state)
type ErrBadInput struct {
	message string
}

func NewErrBadInput(message string) *ErrBadInput {
	return &ErrBadInput{
		message: message,
	}
}

func (e *ErrBadInput) Error() string {
	return e.message
}

// when the user is not authorized (e.g. user is not authenticated)
type ErrUnauthorized struct {
	message string
}

func NewErrUnauthorized(message string) *ErrUnauthorized {
	return &ErrUnauthorized{
		message: message,
	}
}

func (e *ErrUnauthorized) Error() string {
	return e.message
}

// when the user is forbidden (e.g. user can't access a resource)
type ErrForbidden struct {
	message string
}

func NewErrForbidden(message string) *ErrForbidden {
	return &ErrForbidden{
		message: message,
	}
}

func (e *ErrForbidden) Error() string {
	return e.message
}

var (
	// common errors
	ErrBesoinProfilApprobateur = NewErrForbidden("Il faut être approbateur")
	ErrBesoinProfilExploitant  = NewErrForbidden("Il faut être exploitant")
	ErrBesoinProfilInspecteur  = NewErrForbidden("Il faut être inspecteur")
	ErrNonAssigneInspection    = NewErrForbidden("Il faut être assigné à l'inspection")
	ErrInvalidInput            = NewErrBadInput("Données invalides")
)
