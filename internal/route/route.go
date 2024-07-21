package route

import (
	"awesomeProject12/internal/controller"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *controller.UserController
	AuthMiddleware fiber.Handler
	Config         *viper.Viper
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Static("/profile_pictures", c.Config.GetString("directories.profile_pictures"))
	c.App.Post("api/users/register", c.UserController.Register)
	c.App.Post("api/users/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use("api/users/current", c.AuthMiddleware)
	c.App.Get("api/users/current", c.UserController.Current)
	c.App.Patch("api/users/current", c.UserController.Update)
	c.App.Patch("api/users/current/password", c.UserController.PasswordUpdate)
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}
