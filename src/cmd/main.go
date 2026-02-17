package main

import (
	"log"

	"github.com/farzadamr/event-manager-api/api"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/infra/database"
	"github.com/farzadamr/event-manager-api/infra/migration"
	"github.com/farzadamr/event-manager-api/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	err := database.InitDb(cfg)
	defer database.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	migration.Up_1()
	api.InitServer(cfg)
	log.Printf("api is running on %s:%s", cfg.Server.Domain, cfg.Server.ExternalPort)
}
