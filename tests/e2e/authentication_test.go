package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAuthentication(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	tests.AuthInspecteur(e.GET("/ping")).
		Expect().
		Status(http.StatusOK)
}
