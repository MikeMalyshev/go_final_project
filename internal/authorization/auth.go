package authorization

import (
	"fmt"

	"go_final_project/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// Структура для взаимодествия api, реализующая методы авторизации
type Handler struct {
	hashedPassword [32]byte
	token          jwt.Token
	cfg            config.Handler
}

func Create() *Handler {
	return &Handler{}
}

func (auth *Handler) VerifyPassword(password string) bool {
	return password == auth.cfg.Password()
}

func (auth *Handler) PasswordSetted() bool {
	return len(auth.cfg.Password()) > 0
}

func (auth *Handler) CreateToken() (string, error) {
	if !auth.PasswordSetted() {
		return "", fmt.Errorf("authorization.Handler.CreateToken: password not setted ")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	ss, err := token.SignedString([]byte(auth.cfg.Password()))
	if err != nil {
		return "", fmt.Errorf("authorization.Handler.CreateTocken: %v", err)
	}
	return ss, nil
}

func (auth *Handler) VerifyTocken(tokenString string) bool {
	token, err := jwt.Parse(tokenString,
		func(t *jwt.Token) (interface{}, error) {
			return []byte([]byte(auth.cfg.Password())), nil
		})

	if err != nil {
		return false
	}
	return token.Valid
}
