package httpserver

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/authentication"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type LoginHTTPRequest struct {
	Ticket string `json:"ticket" binding:"required"`
}

func (server *HttpServer) login(c *gin.Context) {
	var request LoginHTTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := server.authenticationService.Login(request.Ticket)
	if err != nil {
		if err != authentication.ErrMissingUser {
			log.Error().Err(err).Msg("Bad service response")
		}
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, result)
}

type AuthenticateHTTPRequest struct {
	Token string `json:"token" binding:"required"`
}

func (server *HttpServer) authenticate(c *gin.Context) {
	var request AuthenticateHTTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := server.authenticationService.ValidateToken(request.Token)
	if err != nil {
		if err != authentication.ErrMissingUser {
			log.Error().Err(err).Msg("Bad service response")
		}
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, result.User)
}

func (server *HttpServer) logout(c *gin.Context) {
	err := server.authenticationService.Logout(c.Request.Header.Get(AuthorizationHeader))
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "disconnected"})
}
