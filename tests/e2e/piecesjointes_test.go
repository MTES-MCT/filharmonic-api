package e2e

import (
	"net/http"
	"testing"

	"github.com/MTES-MCT/filharmonic-api/tests"
)

func TestAddGetPieceJointe(t *testing.T) {
	e, close := tests.Init(t)
	defer close()

	pieceJointeId := tests.AuthExploitant(e.POST("/piecesjointes")).
		WithMultipart().WithFile("file", "../testdata/pdf-sample.pdf").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().Value("id").Raw()

	tests.AuthExploitant(e.GET("/piecesjointes/{id}")).WithPath("id", pieceJointeId).
		Expect().
		Status(http.StatusOK).
		ContentType("application/octet-stream")
}
