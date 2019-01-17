package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAuthentication(t *testing.T) {
	e, close := tests.Init(t, nil)
	defer close()

	tests.Auth(e.GET("/ping")).
		Expect().
		Status(http.StatusOK)
}
