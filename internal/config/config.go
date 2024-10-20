package config

import (
	"os"
)

// Структура для взаимодествия api, реализующая методы получения базовых настроек
type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h Handler) DBPath() string {
	dbPath, ok := os.LookupEnv(dbPathEnv)
	if !ok {
		dbPath = defaultDBPath
	}
	return dbPath
}

func (h Handler) Port() string {
	webDirPath, ok := os.LookupEnv(portEnv)
	if !ok {
		webDirPath = defaultPort
	}
	return webDirPath
}

func (h Handler) WebDirPath() string {
	webDirPath, ok := os.LookupEnv(webDirEnv)
	if !ok {
		webDirPath = defaultWebDir
	}
	return webDirPath
}

func (h Handler) Password() string {
	password, ok := os.LookupEnv(passwordEnv)
	if !ok {
		password = defaultPassword
	}
	return password
}
