package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenUserInfo struct {
	ID    uint   `json:"id"`
	Login string `json:"login" json:"login,omitempty"`
	Email string `json:"email" json:"email,omitempty"`
}

type JwtUserInfoClaims struct {
	jwt.RegisteredClaims
	User *TokenUserInfo `json:"user,omitempty"`
}

func NewToken(secret string, expirationTime int, userInfo *TokenUserInfo) (string, error) {
	claims := JwtUserInfoClaims{
		User: userInfo,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationTime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	return signedToken, err
}

func VerifyToken(secret string, token string) (*TokenUserInfo, bool) {
	t, err := jwt.ParseWithClaims(token, &JwtUserInfoClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, false
	}

	expTime, err := t.Claims.GetExpirationTime()
	if err != nil {
		return nil, false
	}

	if !t.Valid && expTime.Before(time.Now()) {
		return nil, false
	}

	if userInfo, ok := t.Claims.(*JwtUserInfoClaims); ok {
		return userInfo.User, true
	}
	return nil, false
}
