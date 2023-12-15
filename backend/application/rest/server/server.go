package server

import (
	"database/sql"
	"log"

	"github.com/bersennaidoo/iathletestats/backend/application/rest/handlers"
	"github.com/bersennaidoo/iathletestats/backend/infrastructure/repositories/pgrepo"
	"github.com/bersennaidoo/iathletestats/backend/infrastructure/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config         *viper.Viper
	router         *gin.Engine
	runnersHandler *handlers.RunnersHandler
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	runnersRepo := pgrepo.NewRunnersRepo(dbHandler)
	runnersService := services.NewRunnersService(runnersRepo)
	runnersHandler := handlers.NewRunnersHandler(runnersService)

	router := gin.Default()

	router.POST("/runner", runnersHandler.CreateRunner)

	return HttpServer{
		config:         config,
		router:         router,
		runnersHandler: runnersHandler,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
