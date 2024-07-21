package config

import (
	"awesomeProject12/internal/adapter"
	"awesomeProject12/internal/controller"
	"awesomeProject12/internal/middleware"
	"awesomeProject12/internal/repository"
	"awesomeProject12/internal/route"
	"awesomeProject12/internal/service"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Validate *validator.Validate
	Config   *viper.Viper
	Modifier *mold.Transformer
}

func Bootstrap(config *BootstrapConfig) {

	storage := adapter.NewFileAdapter()

	userRepository := repository.NewUserRepository()

	userService := service.NewUserService(config.DB, config.Validate, storage, userRepository, config.Config)

	userController := controller.NewUserController(userService, config.Modifier)

	autMiddleware := middleware.NewAuth(userService)

	route := route.RouteConfig{
		config.App, userController, autMiddleware, config.Config,
	}
	route.Setup()

}
