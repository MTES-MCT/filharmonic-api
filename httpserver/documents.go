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

func (server *HttpServer) genererLettreSuite(c *gin.Context) (*models.PieceJointeFile, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorF(err)
	}
	return server.service.GenererLettreSuite(server.retrieveUserContext(c), id)
}

func (server *HttpServer) genererRapport(c *gin.Context) (*models.PieceJointeFile, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorF(err)
	}
	return server.service.GenererRapport(server.retrieveUserContext(c), id)
}
