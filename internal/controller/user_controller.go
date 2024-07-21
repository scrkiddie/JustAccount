package controller

import (
	"awesomeProject12/internal/middleware"
	"awesomeProject12/internal/model"
	"awesomeProject12/internal/service"
	"github.com/go-playground/mold/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
	"log"
)

type UserController struct {
	UserService *service.UserService
	Modifier    *mold.Transformer
}

func NewUserController(userService *service.UserService, modifier *mold.Transformer) *UserController {
	return &UserController{UserService: userService, Modifier: modifier}
}

func (c *UserController) Register(ctx fiber.Ctx) error {
	request := new(model.RegisterUserRequest)

	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	if err := c.UserService.Create(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": "User created successfully"})
}

func (c *UserController) Login(ctx fiber.Ctx) error {
	request := new(model.LoginUserRequest)

	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	response, err := c.UserService.Login(ctx.UserContext(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": response.Token})
}

func (c *UserController) Current(ctx fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UserService.Current(ctx.UserContext(), request)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": response})
}

func (c *UserController) Update(ctx fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserRequest)
	if err := ctx.Bind().Body(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if err := c.Modifier.Struct(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return fiber.ErrInternalServerError
	}

	request.ID = auth.ID

	file := new(model.File)
	profilePicture, err := ctx.FormFile("profilePicture")
	if err != nil && err != fasthttp.ErrMissingFile {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	if profilePicture != nil {
		file.FileHeader = profilePicture
	}

	if err := c.UserService.Update(ctx.UserContext(), request, file); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "User updated successfully"})
}

func (c *UserController) PasswordUpdate(ctx fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdatePasswordRequest)
	if err := ctx.Bind().JSON(request); err != nil {
		log.Println(err.Error())
		return fiber.ErrBadRequest
	}

	request.ID = auth.ID
	if err := c.UserService.UpdatePassword(ctx.UserContext(), request); err != nil {
		log.Println(err.Error())
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": "Password updated successfully"})
}
