package httpserver

import (
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/MTES-MCT/filharmonic-api/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/rs/zerolog/log"
)

// call a service function and return a success message
func returnOk(serviceFunc func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := serviceFunc(c)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "operation succeeded",
		})
	}
}

// call a service function returning an id and return it inside a json object
func returnId(serviceFunc func(*gin.Context) (int64, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := serviceFunc(c)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}

// call a service function and return its result
func returnResult(serviceFunc func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := serviceFunc(c)
		if err != nil {
			handleError(c, err)
			return
		}
		if result == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not_found",
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// call a service function returning a file and transfer it
func returnFile(serviceFunc func(*gin.Context) (*models.PieceJointeFile, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := serviceFunc(c)
		if err != nil {
			handleError(c, err)
			return
		}
		c.Render(http.StatusOK, render.Reader{
			Headers:       map[string]string{},
			ContentType:   file.Type,
			ContentLength: file.Taille,
			Reader:        file.Content,
		})
	}
}

// log an error and return the appropriate HTTP response
func handleError(c *gin.Context, err error) {
	statusCode := 0
	logEvent := log.Error().Err(err)
	switch err.(type) {
	case *domain.ErrBadInput:
		logEvent.Msg("service returned a bad input error")
		statusCode = http.StatusBadRequest
	case *domain.ErrUnauthorized:
		logEvent.Msg("service returned an unauthorized error")
		statusCode = http.StatusUnauthorized
	case *domain.ErrForbidden:
		logEvent.Msg("service returned an forbidden error")
		statusCode = http.StatusForbidden
	default:
		logEvent.Msg("service returned an unknown error")
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{
		"message": err.Error(),
	})
}

// shorthand to return a bad input error
func badInputError(err error) error {
	return domain.NewErrBadInput(err.Error())
}

// shorthand to return 0 and a bad input error
func badInputErrorI(err error) (int64, error) {
	return 0, domain.NewErrBadInput(err.Error())
}

// shorthand to return a nil result and a bad input error
func badInputErrorN(err error) (interface{}, error) {
	return nil, domain.NewErrBadInput(err.Error())
}

// shorthand to return a nil file and a bad input error
func badInputErrorF(err error) (*models.PieceJointeFile, error) {
	return nil, domain.NewErrBadInput(err.Error())
}
