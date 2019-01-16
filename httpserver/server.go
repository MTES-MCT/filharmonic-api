package httpserver

import (
	"log"
	"net/http"

	"github.com/MTES-MCT/filharmonic-api/domain"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	service *domain.Service
}

func New(service *domain.Service) *HttpServer {
	return &HttpServer{
		service: service,
	}
}

func (s *HttpServer) Start() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/etablissements", s.listEtablissements)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return server
}
