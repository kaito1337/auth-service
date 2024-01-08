package controllers

import (
	"auth-backend/internal/models"
	"auth-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type RegisterController struct {
	userService *services.UserService
}

type RegisterRequest struct {
	Email           string `json:"email"`
	Login           string `json:"login"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type RegisterResponse struct {
	Ok           bool         `json:"ok"`
	User         *models.User `json:"user,omitempty"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

func NewRegisterController(userService *services.UserService) *RegisterController {
	return &RegisterController{
		userService: userService,
	}
}

func (c *RegisterController) GetGroup() string {
	return "/auth"
}

func (c *RegisterController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method:  "POST",
			Path:    "/signUp",
			Handler: c.registerHandler(),
		},
	}
}

func (c *RegisterController) registerHandler() func(*fiber.Ctx) error {
	return func(fc *fiber.Ctx) error {
		fc.Accepts("application/json")
		var request RegisterRequest
		if err := fc.BodyParser(&request); err != nil {
			return err
		}
		if valid, err := ValidateRegisterData(request); !valid {
			return fc.Status(400).JSON(RegisterResponse{
				Ok:           false,
				ErrorMessage: err.Error(),
			})
		}

		user, err := c.userService.CreateUser(request.Email, request.Login, request.Password)
		if err != nil {
			return fc.Status(400).JSON(RegisterResponse{
				Ok:           false,
				ErrorMessage: err.Error(),
			})
		}
		return fc.Status(201).JSON(RegisterResponse{
			Ok:   true,
			User: user,
		})
	}
}
