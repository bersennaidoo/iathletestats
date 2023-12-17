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
	usersHandler   *handlers.UsersHandler
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	runnersRepo := pgrepo.NewRunnersRepo(dbHandler)
	runnersService := services.NewRunnersService(runnersRepo)
	usersRepo := pgrepo.NewUsersRepo(dbHandler)
	usersService := services.NewUsersService(usersRepo)
	runnersHandler := handlers.NewRunnersHandler(runnersService, usersService)
	usersHandler := handlers.NewUsersHandler(usersService)

	router := gin.Default()

	router.POST("/runner", runnersHandler.CreateRunner)

	router.POST("/login", usersHandler.Login)
	router.POST("/logout", usersHandler.Logout)

	return HttpServer{
		config:         config,
		router:         router,
		runnersHandler: runnersHandler,
		usersHandler:   usersHandler,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
