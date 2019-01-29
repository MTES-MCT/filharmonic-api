package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *HttpServer) getUser(c *gin.Context) {
	ctx := server.retrieveUserContext(c)
	c.JSON(http.StatusOK, ctx.User)
}
