package services

import (
	"auth-backend/internal/config"
	"auth-backend/internal/models"
	"auth-backend/internal/token"
	"fmt"
)

var EmailOrPasswordInvalid = fmt.Errorf("email or password invalid")
var refreshTokenExpired = fmt.Errorf("refresh token expired")

type AuthService struct {
	userService *UserService
	cfg         *config.AppConfig
}

type AuthResult struct {
	Token        string
	RefreshToken string
	User         *models.User
	Err          error
}

func NewAuthService(cfg *config.AppConfig, userService *UserService) *AuthService {
	return &AuthService{
		cfg:         cfg,
		userService: userService,
	}
}

func (s *AuthService) Login(email string, password string) *AuthResult {
	user, err := s.userService.GetUserByEmail(email)
	if err != nil {
		return &AuthResult{
			Err: err,
		}
	}

	valid := s.userService.CheckPasswordHash(password, user.Password)

	if !valid {
		return &AuthResult{
			Err: EmailOrPasswordInvalid,
		}
	}

	accessToken, refreshToken, err := s.generateTokens(user)

	if err != nil {
		return &AuthResult{
			Err: err,
		}
	}

	return &AuthResult{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
		Err:          nil,
	}
}

func (s *AuthService) Refresh(email string, refreshToken string) *AuthResult {
	user, err := s.userService.GetUserByEmail(email)
	if err != nil {
		return &AuthResult{
			Err: err,
		}
	}

	_, valid := token.VerifyToken(s.cfg.TokenSecret, refreshToken)
	if !valid {
		return &AuthResult{
			Err: refreshTokenExpired,
		}
	}
	if user.RefreshToken != refreshToken {
		return &AuthResult{
			Err: refreshTokenExpired,
		}
	}

	accessToken, refreshToken, err := s.generateTokens(user)

	if err != nil {
		return &AuthResult{
			Err: err,
		}
	}

	return &AuthResult{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
		Err:          nil,
	}

}

func (s *AuthService) generateTokens(user *models.User) (string, string, error) {
	accessToken, err := token.NewToken(s.cfg.TokenSecret, s.cfg.AccessTokenExpirationTimeHours, &token.TokenUserInfo{
		ID:    user.ID,
		Login: user.Login,
		Email: user.Email,
	})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := token.NewToken(s.cfg.TokenSecret, s.cfg.RefreshTokenExpirationTimeHours, nil)
	if err != nil {
		return "", "", err
	}
	err = s.userService.UpdateRefreshToken(user.Email, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
