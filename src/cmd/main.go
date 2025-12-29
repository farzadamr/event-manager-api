package main

import (
	"log"

	"github.com/farzadamr/event-manager-api/api"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/infra/migration"
)

func main() {
	cfg := config.GetConfig()

	err := database.InitDb(cfg)
	defer database.CloseDb()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	migration.Up_1()

	api.InitServer(cfg)

}
