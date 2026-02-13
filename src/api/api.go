package api

import (
	"fmt"

	"github.com/farzadamr/event-manager-api/api/middleware"
	"github.com/farzadamr/event-manager-api/api/router"
	"github.com/farzadamr/event-manager-api/api/validation"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var logger = logging.NewLogger(config.GetConfig())

func InitServer(cfg *config.Config) {
	gin.SetMode(cfg.Server.RunMode)
	r := gin.New()
	RegisterValidators()

	r.Use(middleware.DefaultStructuredLogger(cfg))
	r.Use(gin.Logger())

	RegisterRoutes(r, cfg)
	logger.Info(logging.General, logging.Startup, "Started", nil)
	err := r.Run(fmt.Sprintf(":%s", cfg.Server.InternalPort))
	if err != nil {
		logger.Fatal(logging.General, logging.Startup, err.Error(), nil)
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
		//Registrations
		registrations := v1.Group("/registrations")
		router.Registration(registrations, cfg)
		//Certificate
		certificates := v1.Group("/certificates")
		router.Certificate(certificates, cfg)
		//Role
		roles := v1.Group("/roles")
		router.Role(roles, cfg)
	}
}

func RegisterValidators() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validation.IranianMobileNumberValidator, true)
		if err != nil {
			logger.Error(logging.Validation, logging.Startup, err.Error(), nil)
		}
		err = val.RegisterValidation("password", validation.PasswordValidator, true)
		if err != nil {
			logger.Error(logging.Validation, logging.Startup, err.Error(), nil)
		}
		err = val.RegisterValidation("date", validation.DateValidator, true)
		if err != nil {
			logger.Error(logging.Validation, logging.Startup, err.Error(), nil)
		}
	}
}
