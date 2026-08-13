package config

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/fx"
)

// DefaultJWTSecret 仅用于本地开发，生产环境必须通过 JWT_SECRET 覆盖。
const DefaultJWTSecret = "ai-learning-admin-secret-dev"

type Config struct {
	Port         string
	JWTSecret    string
	TokenTTL     time.Duration // 登录 token 有效期
	DBPath       string
	CORSOrigins  string // 允许跨域的来源，逗号分隔；* 表示不限制
	LLM          LLMConfig
	Embedding    EmbeddingConfig
	Vector       VectorConfig
}

type EmbeddingConfig struct {
	APIKey     string
	BaseURL    string
	Model      string
	Dimensions int
	Enabled    bool
}

type VectorConfig struct {
	TopK         int
	HybridAlpha  float64 // 向量分数权重，关键词权重为 1-alpha
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
		secret = DefaultJWTSecret
		fmt.Println("⚠️  JWT_SECRET 未配置，正在使用内置开发密钥；生产环境请务必设置 JWT_SECRET")
	}
	tokenTTL := 7 * 24 * time.Hour
	if v := os.Getenv("JWT_TTL_HOURS"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 {
			tokenTTL = time.Duration(n) * time.Hour
		}
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/ai-learning.db"
	}
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "*"
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
	embedKey := os.Getenv("EMBEDDING_API_KEY")
	if embedKey == "" {
		embedKey = apiKey
	}
	embedBaseURL := os.Getenv("EMBEDDING_BASE_URL")
	if embedBaseURL == "" {
		embedBaseURL = baseURL
	}
	embedModel := os.Getenv("EMBEDDING_MODEL")
	if embedModel == "" {
		embedModel = "text-embedding-3-small"
	}
	embedDim := 1536
	if v := os.Getenv("EMBEDDING_DIMENSIONS"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 {
			embedDim = n
		}
	}
	topK := 3
	if v := os.Getenv("VECTOR_TOP_K"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 {
			topK = n
		}
	}
	hybridAlpha := 0.7
	if v := os.Getenv("VECTOR_HYBRID_ALPHA"); v != "" {
		if f, err := parseFloat(v); err == nil {
			hybridAlpha = f
		}
	}
	return &Config{
		Port:        port,
		JWTSecret:   secret,
		TokenTTL:    tokenTTL,
		DBPath:      dbPath,
		CORSOrigins: corsOrigins,
		LLM: LLMConfig{
			APIKey:  apiKey,
			BaseURL: baseURL,
			Model:   model,
			Enabled: apiKey != "",
		},
		Embedding: EmbeddingConfig{
			APIKey:     embedKey,
			BaseURL:    embedBaseURL,
			Model:      embedModel,
			Dimensions: embedDim,
			Enabled:    embedKey != "",
		},
		Vector: VectorConfig{
			TopK:        topK,
			HybridAlpha: hybridAlpha,
		},
	}
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

var Module = fx.Provide(NewConfig)
