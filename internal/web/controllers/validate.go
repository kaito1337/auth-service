package controllers

import (
	"errors"
	"net/mail"
)

/*
Custom errors
*/

var (
	emailRequired         = errors.New("Заполните поле электронной почты")
	passwordRequired      = errors.New("Заполните поле пароля")
	invalidPasswordFormat = errors.New("Длина пароля должна быть не менее 8 и не более 24 символов")
	invalidLoginFormat    = errors.New("Длина логина должна быть не менее 3 и не более 24 символов")
	passwordsDoNotMatch   = errors.New("Пароли не совпадают")
	invalidEmailFormat    = errors.New("Неверный формат почты")
)

func ValidateAuthData(r AuthRequest) (bool, error) {
	if len(r.Email) == 0 {
		return false, emailRequired
	}
	if len(r.Password) == 0 {
		return false, passwordRequired
	}
	return true, nil
}

func ValidateRegisterData(r RegisterRequest) (bool, error) {
	if len(r.Email) == 0 {
		return false, emailRequired
	}
	if len(r.Login) < 3 || len(r.Login) > 24 {
		return false, invalidLoginFormat
	}
	if len(r.Password) < 8 || len(r.Password) > 24 {
		return false, invalidPasswordFormat
	}
	if r.ConfirmPassword != r.Password {
		return false, passwordsDoNotMatch
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return false, invalidEmailFormat
	}
	return true, nil
}

func ValidateRefreshData(r RefreshAuthRequest) (bool, error) {
	if len(r.Email) == 0 {
		return false, emailRequired
	}
	if len(r.RefreshToken) == 0 {
		return false, invalidPasswordFormat
	}
	return true, nil
}
