package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) addCommentaire(c *gin.Context) (int64, error) {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorI(err)
	}
	var commentaire models.Commentaire
	if err = c.ShouldBindJSON(&commentaire); err != nil {
		return badInputErrorI(err)
	}
	return server.service.CreateCommentaire(server.retrieveUserContext(c), idInspection, commentaire)
}
