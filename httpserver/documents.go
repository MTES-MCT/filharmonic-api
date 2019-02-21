package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) genererLettreAnnonce(c *gin.Context) (*models.PieceJointeFile, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorF(err)
	}
	return server.service.GenererLettreAnnonce(server.retrieveUserContext(c), id)
}
