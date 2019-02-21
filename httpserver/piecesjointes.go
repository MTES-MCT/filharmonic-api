package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) createPieceJointe(c *gin.Context) (int64, error) {
	formFile, err := c.FormFile("file")
	if err != nil {
		return badInputErrorI(err)
	}
	file, err := formFile.Open()
	if err != nil {
		return badInputErrorI(err)
	}
	pieceJointe := models.PieceJointeFile{
		Content: file,
		Taille:  formFile.Size,
		Nom:     formFile.Filename,
		Type:    formFile.Header.Get("Content-Type"),
	}
	return server.service.CreatePieceJointe(server.retrieveUserContext(c), pieceJointe)
}

func (server *HttpServer) getPieceJointe(c *gin.Context) (*models.PieceJointeFile, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorF(err)
	}
	return server.service.GetPieceJointe(server.retrieveUserContext(c), id)
}
