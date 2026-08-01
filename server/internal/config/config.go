package config

import (
	"os"

	"go.uber.org/fx"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

func NewConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "ai-learning-admin-secret-dev"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/ai-learning.db"
	}
	return &Config{Port: port, JWTSecret: secret, DBPath: dbPath}
}

var Module = fx.Provide(NewConfig)
