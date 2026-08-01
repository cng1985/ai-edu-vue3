package service

import (
	"strings"
	"sync"
	"time"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/llm"
)

const (
	settingLLMAPIKey  = "llm.api_key"
	settingLLMBaseURL = "llm.base_url"
	settingLLMModel   = "llm.model"
)

type SettingsService struct {
	repo   *repository.SettingsRepo
	boot   *config.Config
	mu     sync.RWMutex
	llm    *llm.Client
	llmCfg config.LLMConfig
	meta   settingsMeta
}

type settingsMeta struct {
	updatedAt int64
	source    string
}

func NewSettingsService(repo *repository.SettingsRepo, cfg *config.Config) *SettingsService {
	s := &SettingsService{repo: repo, boot: cfg}
	s.reload()
	return s
}

func (s *SettingsService) reload() {
	rows, _ := s.repo.List()
	dbValues := map[string]string{}
	var updatedAt int64
	for _, row := range rows {
		dbValues[row.Key] = row.Value
		if row.UpdatedAt > updatedAt {
			updatedAt = row.UpdatedAt
		}
	}
	apiKey, baseURL, modelName, source := mergeLLMSettings(s.boot.LLM, dbValues)
	cfg := config.LLMConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   modelName,
		Enabled: apiKey != "",
	}
	client := llm.NewClient(apiKey, baseURL, modelName)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.llmCfg = cfg
	s.llm = client
	s.meta = settingsMeta{updatedAt: updatedAt, source: source}
}

func mergeLLMSettings(boot config.LLMConfig, db map[string]string) (apiKey, baseURL, modelName, source string) {
	apiKey = boot.APIKey
	baseURL = boot.BaseURL
	modelName = boot.Model
	source = "environment"
	if len(db) == 0 {
		return
	}
	if v, ok := db[settingLLMAPIKey]; ok {
		apiKey = v
		source = "database"
	}
	if v, ok := db[settingLLMBaseURL]; ok && strings.TrimSpace(v) != "" {
		baseURL = v
		source = "database"
	}
	if v, ok := db[settingLLMModel]; ok && strings.TrimSpace(v) != "" {
		modelName = v
		source = "database"
	}
	return
}

func (s *SettingsService) LLMClient() *llm.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.llm
}

func (s *SettingsService) LLMConfig() config.LLMConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.llmCfg
}

func (s *SettingsService) GetView() model.SettingsView {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return model.SettingsView{
		LLM: model.LLMSettingsView{
			APIKeyConfigured: s.llmCfg.APIKey != "",
			APIKeyMasked:     maskSecret(s.llmCfg.APIKey),
			BaseURL:          s.llmCfg.BaseURL,
			Model:            s.llmCfg.Model,
			Enabled:          s.llmCfg.Enabled,
			Source:           s.meta.source,
		},
		UpdatedAt: s.meta.updatedAt,
	}
}

func (s *SettingsService) Update(req model.SettingsUpdateRequest) (*model.SettingsView, error) {
	now := time.Now().UnixMilli()
	current, _ := s.repo.GetAll()

	apiKey := strings.TrimSpace(req.LLM.APIKey)
	if apiKey == "" {
		if v, ok := current[settingLLMAPIKey]; ok {
			apiKey = v
		} else {
			apiKey = s.boot.LLM.APIKey
		}
	}
	baseURL := strings.TrimSpace(req.LLM.BaseURL)
	if baseURL == "" {
		baseURL = s.boot.LLM.BaseURL
	}
	modelName := strings.TrimSpace(req.LLM.Model)
	if modelName == "" {
		modelName = s.boot.LLM.Model
	}

	if err := s.repo.Upsert(settingLLMAPIKey, apiKey, now); err != nil {
		return nil, err
	}
	if err := s.repo.Upsert(settingLLMBaseURL, baseURL, now); err != nil {
		return nil, err
	}
	if err := s.repo.Upsert(settingLLMModel, modelName, now); err != nil {
		return nil, err
	}

	s.reload()
	view := s.GetView()
	view.UpdatedAt = now
	return &view, nil
}

func maskSecret(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "****"
	}
	return value[:4] + "****" + value[len(value)-4:]
}
