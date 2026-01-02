package api

import (
	"fmt"
	"log"

	"github.com/farzadamr/event-manager-api/api/router"
	"github.com/farzadamr/event-manager-api/api/validation"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitServer(cfg *config.Config) {
	gin.SetMode(cfg.Server.RunMode)
	r := gin.New()
	RegisterValidatiors()

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
		// Event
		events := v1.Group("/events")
		router.Event(events, cfg)
	}
}

func RegisterValidatiors() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validation.IranianMobileNumberValidator, true)
		if err != nil {
			log.Fatalf("Unable to register validator -> %s", err.Error())
		}
		err = val.RegisterValidation("password", validation.PasswordValidator, true)
		if err != nil {
			log.Fatalf("Unable to register validator -> %s", err.Error())
		}
		err = val.RegisterValidation("date", validation.DateValidator, true)
		if err != nil {
			log.Fatalf("Unable to register validator -> %s", err.Error())
		}
	}
}
