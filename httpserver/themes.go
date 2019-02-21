package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) listThemes(c *gin.Context) (interface{}, error) {
	return server.service.ListThemes(server.retrieveUserContext(c))
}

func (server *HttpServer) createTheme(c *gin.Context) (int64, error) {
	var theme models.Theme
	if err := c.ShouldBindJSON(&theme); err != nil {
		return badInputErrorI(err)
	}
	return server.service.CreateTheme(server.retrieveUserContext(c), theme)
}

func (server *HttpServer) deleteTheme(c *gin.Context) error {
	idTheme, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.DeleteTheme(server.retrieveUserContext(c), idTheme)
}
