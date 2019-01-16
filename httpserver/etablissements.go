package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *HttpServer) listEtablissements(c *gin.Context) {
	s3ic := c.Query("s3ic")
	etablissements, err := server.service.ListEtablissements(s3ic)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, etablissements)
}
