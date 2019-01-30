package httpserver

import (
	"net/http"
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/rs/zerolog/log"
)

func (server *HttpServer) createPieceJointe(c *gin.Context) {
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	file, err := formFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	pieceJointe := models.PieceJointeFile{
		Content: file,
		Taille:  formFile.Size,
		Nom:     formFile.Filename,
		Type:    formFile.Header.Get("Content-Type"),
	}
	pieceJointeId, err := server.service.CreatePieceJointe(server.retrieveUserContext(c), pieceJointe)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": pieceJointeId,
	})
}

func (server *HttpServer) getPieceJointe(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	pieceJointeFile, err := server.service.GetPieceJointe(server.retrieveUserContext(c), id)
	if err != nil {
		log.Error().Err(err).Msg("Bad service response")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	headers := make(map[string]string)
	c.Render(http.StatusOK, render.Reader{
		Headers:       headers,
		ContentType:   pieceJointeFile.Type,
		ContentLength: pieceJointeFile.Taille,
		Reader:        pieceJointeFile.Content,
	})
}
