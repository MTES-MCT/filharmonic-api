package httpserver

import (
	"strconv"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) listInspections(c *gin.Context) (interface{}, error) {
	filter := domain.ListInspectionsFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		return badInputErrorN(err)
	}
	return server.service.ListInspections(server.retrieveUserContext(c), filter)
}

func (server *HttpServer) getInspection(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputErrorN(err)
	}
	return server.service.GetInspection(server.retrieveUserContext(c), id)
}

func (server *HttpServer) createInspection(c *gin.Context) (int64, error) {
	var inspection models.Inspection
	if err := c.ShouldBindJSON(&inspection); err != nil {
		return badInputErrorI(err)
	}
	return server.service.CreateInspection(server.retrieveUserContext(c), inspection)
}

func (server *HttpServer) updateInspection(c *gin.Context) error {
	var inspection models.Inspection
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	if err = c.ShouldBindJSON(&inspection); err != nil {
		return badInputError(err)
	}
	inspection.Id = id
	return server.service.UpdateInspection(server.retrieveUserContext(c), inspection)
}

func (server *HttpServer) publishInspection(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.PublishInspection(server.retrieveUserContext(c), idInspection)
}

func (server *HttpServer) askValidateInspection(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.AskValidateInspection(server.retrieveUserContext(c), idInspection)
}

func (server *HttpServer) rejectInspection(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.RejectInspection(server.retrieveUserContext(c), idInspection)
}

func (server *HttpServer) validateInspection(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	formFile, err := c.FormFile("file")
	if err != nil {
		return badInputError(err)
	}
	file, err := formFile.Open()
	if err != nil {
		return badInputError(err)
	}
	rapportFile := models.File{
		Content: file,
		Taille:  formFile.Size,
		Nom:     formFile.Filename,
		Type:    formFile.Header.Get("Content-Type"),
	}
	return server.service.ValidateInspection(server.retrieveUserContext(c), idInspection, rapportFile)
}

func (server *HttpServer) cloreInspection(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.CloreInspection(server.retrieveUserContext(c), idInspection)
}

func (server *HttpServer) addFavoriToInspection(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.AddFavoriToInspection(server.retrieveUserContext(c), idInspection)
}

func (server *HttpServer) removeFavoriToInspection(c *gin.Context) error {
	idInspection, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return badInputError(err)
	}
	return server.service.RemoveFavoriToInspection(server.retrieveUserContext(c), idInspection)
}

func (server *HttpServer) listInspectionsFavorites(c *gin.Context) (interface{}, error) {
	return server.service.ListInspectionsFavorites(server.retrieveUserContext(c))
}
