package httpserver

import (
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) getUser(c *gin.Context) (interface{}, error) {
	ctx := server.retrieveUserContext(c)
	return ctx.User, nil
}
