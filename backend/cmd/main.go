package main

import (
	"log"
	"os"

	"github.com/bersennaidoo/iathletestats/backend/application/rest/server"
	"github.com/bersennaidoo/iathletestats/backend/physical/config"
	"github.com/bersennaidoo/iathletestats/backend/physical/dbconn"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting IAthletestats App")

	log.Println("Initializing configuration")
	config := config.InitConfig(getConfigFileName())

	log.Println("Initializing database")
	dbHandler := dbconn.InitDatabase(config)

	log.Println("Initializig HTTP sever")
	httpServer := server.InitHttpServer(config, dbHandler)

	httpServer.Start()
}

func getConfigFileName() string {
	env := os.Getenv("ENV")

	if env != "" {
		return "runners-" + env
	}

	return "runners"
}
