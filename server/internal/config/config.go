package config

import (
	"os"

	"go.uber.org/fx"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
	LLM       LLMConfig
}

type LLMConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	Enabled bool
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
	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &Config{
		Port:      port,
		JWTSecret: secret,
		DBPath:    dbPath,
		LLM: LLMConfig{
			APIKey:  apiKey,
			BaseURL: baseURL,
			Model:   model,
			Enabled: apiKey != "",
		},
	}
}

var Module = fx.Provide(NewConfig)
