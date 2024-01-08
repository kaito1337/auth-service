package controllers

import (
	"auth-backend/internal/models"
	"auth-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *services.AuthService
}
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Ok           bool         `json:"ok"`
	Token        string       `json:"token,omitempty"`
	RefreshToken string       `json:"refreshToken,omitempty"`
	User         *models.User `json:"user,omitempty"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

type RefreshAuthRequest struct {
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) GetGroup() string {
	return "/auth"
}

func (c *AuthController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method:  "POST",
			Path:    "/signIn",
			Handler: c.authHandler(),
		},
		&Handler{
			Method:  "POST",
			Path:    "/refresh",
			Handler: c.refreshHandler(),
		},
	}
}

func (c *AuthController) authHandler() func(*fiber.Ctx) error {
	return func(fc *fiber.Ctx) error {
		fc.Accepts("application/json")
		var request AuthRequest
		if err := fc.BodyParser(&request); err != nil {
			return err
		}
		if valid, err := ValidateAuthData(request); !valid {
			return fc.Status(400).JSON(AuthResponse{
				Ok:           false,
				ErrorMessage: err.Error(),
			})
		}

		authResult := c.authService.Login(request.Email, request.Password)
		if authResult.Err != nil {
			return fc.Status(400).JSON(AuthResponse{
				Ok:           false,
				ErrorMessage: authResult.Err.Error(),
			})
		}

		return fc.Status(200).JSON(AuthResponse{
			Ok:           true,
			Token:        authResult.Token,
			RefreshToken: authResult.RefreshToken,
			User:         authResult.User,
		})
	}
}

func (c *AuthController) refreshHandler() func(*fiber.Ctx) error {
	return func(fc *fiber.Ctx) error {
		fc.Accepts("application/json")
		var request RefreshAuthRequest
		if err := fc.BodyParser(&request); err != nil {
			return err
		}
		if valid, err := ValidateRefreshData(request); !valid {
			return fc.Status(400).JSON(AuthResponse{
				Ok:           false,
				ErrorMessage: err.Error(),
			})
		}

		authResult := c.authService.Refresh(request.Email, request.RefreshToken)
		if authResult.Err != nil {
			return fc.Status(400).JSON(AuthResponse{
				Ok:           false,
				ErrorMessage: authResult.Err.Error(),
			})
		}

		return fc.Status(200).JSON(AuthResponse{
			Ok:           true,
			Token:        authResult.Token,
			RefreshToken: authResult.RefreshToken,
			User:         authResult.User,
		})
	}
}
