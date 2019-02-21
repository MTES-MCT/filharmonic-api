package httpserver

import (
	"github.com/gin-gonic/gin"
)

type LoginHTTPRequest struct {
	Ticket string `json:"ticket" binding:"required"`
}

func (server *HttpServer) login(c *gin.Context) (interface{}, error) {
	var request LoginHTTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return badInputErrorN(err)
	}
	return server.authenticationService.Login(request.Ticket)
}

type AuthenticateHTTPRequest struct {
	Token string `json:"token" binding:"required"`
}

func (server *HttpServer) authenticate(c *gin.Context) (interface{}, error) {
	var request AuthenticateHTTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return badInputErrorN(err)
	}
	result, err := server.authenticationService.ValidateToken(request.Token)
	if err != nil {
		return nil, err
	}
	return result.User, nil
}

func (server *HttpServer) logout(c *gin.Context) error {
	return server.authenticationService.Logout(c.Request.Header.Get(AuthorizationHeader))
}
