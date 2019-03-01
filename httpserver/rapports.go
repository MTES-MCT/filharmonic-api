package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) getRapport(c *gin.Context) (*models.File, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorF(err)
	}
	return server.service.GetRapport(server.retrieveUserContext(c), id)
}
