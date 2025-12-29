package api

import (
	"fmt"
	"log"

	"github.com/farzadamr/event-manager-api/api/router"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
)

func InitServer(cfg *config.Config) {
	gin.SetMode(cfg.Server.RunMode)
	r := gin.New()

	RegisterRoutes(r, cfg)
	err := r.Run(fmt.Sprintf(":%s", cfg.Server.InternalPort))
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		// User
		users := v1.Group("/users")
		router.User(users, cfg)

	}

}
