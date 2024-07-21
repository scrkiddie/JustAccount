package config

import "github.com/gofiber/fiber/v3"

func NewFiber() *fiber.App {
	return fiber.New(fiber.Config{ErrorHandler: NewErrorHandler()})
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		if code == fiber.StatusNotFound {
			return ctx.Status(code).JSON(fiber.Map{
				"error": "Not found",
			})
		}
		return ctx.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
