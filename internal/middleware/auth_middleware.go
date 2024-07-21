package middleware

import (
	"awesomeProject12/internal/model"
	"awesomeProject12/internal/service"
	"github.com/gofiber/fiber/v3"
	"log"
)

func NewAuth(userUserCase *service.UserService) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		log.Printf("Authorization : %s", request.Token)

		auth, err := userUserCase.Verify(ctx.UserContext(), request)
		if err != nil {
			log.Printf(err.Error())
			return fiber.ErrUnauthorized
		}

		log.Printf("user : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}
func GetUser(ctx fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
