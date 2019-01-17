package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const authUserKey = "authUserKey"
const AuthorizationHeader = "Authorization"

func (server *HttpServer) authRequired(c *gin.Context) {
	token := c.Request.Header.Get(AuthorizationHeader)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	}
	userId, err := server.sso.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	}
	c.Set(authUserKey, userId)
	c.Next()
}
