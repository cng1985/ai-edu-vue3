package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"errors"

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
	repo      *repository.SettingsRepo
	boot      *config.Config
	router    *ModelRouter
	aiModel   *AiModelService
	knowledge *KnowledgeService
	mu        sync.RWMutex
	llm       *llm.Client
	llmCfg    config.LLMConfig
	meta      settingsMeta
}

type settingsMeta struct {
	updatedAt int64
	source    string
}

func NewSettingsService(
	repo *repository.SettingsRepo,
	cfg *config.Config,
	router *ModelRouter,
	aiModel *AiModelService,
	knowledge *KnowledgeService,
) *SettingsService {
	s := &SettingsService{repo: repo, boot: cfg, router: router, aiModel: aiModel, knowledge: knowledge}
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
	if s.router != nil {
		s.router.InvalidateClients()
	}
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
	if s.router != nil {
		client, _, err := s.router.ClientFor("")
		if err == nil && client != nil && client.Enabled() {
			return client
		}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.llm
}

func (s *SettingsService) LLMConfig() config.LLMConfig {
	if s.router != nil {
		resolved, err := s.router.Resolve("")
		if err == nil && resolved != nil {
			return config.LLMConfig{
				APIKey:  resolved.APIKey,
				BaseURL: resolved.BaseURL,
				Model:   resolved.ModelCode,
				Enabled: resolved.Enabled,
			}
		}
	}
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

// GetSystemView 聚合大模型分层配置与知识库状态。
func (s *SettingsService) GetSystemView() (*model.SystemSettingsView, error) {
	overview, err := s.aiModel.Overview()
	if err != nil {
		return nil, err
	}
	providerRes, err := s.aiModel.ListProviders("", 1, 100)
	if err != nil {
		return nil, err
	}
	virtualModels, err := s.aiModel.ListVirtualModelOptions()
	if err != nil {
		return nil, err
	}
	resolved, _ := s.ResolvedLLM()
	kbStatus, _ := s.knowledge.Status()

	providers := providerRes.List
	if providers == nil {
		providers = []model.Provider{}
	}
	if virtualModels == nil {
		virtualModels = []model.VirtualModel{}
	}

	return &model.SystemSettingsView{
		AiModel:       *overview,
		Resolved:      resolved,
		Providers:     providers,
		VirtualModels: virtualModels,
		Knowledge:     kbStatus,
	}, nil
}

func (s *SettingsService) ResolveVirtualModel(code string) (*model.ResolvedLLM, error) {
	if s.router == nil {
		return nil, errors.New("模型路由未初始化")
	}
	return s.router.Resolve(code)
}

func (s *SettingsService) SetDefaultVirtualModel(code string) error {
	return s.aiModel.SetDefaultVirtualModel(code)
}

func (s *SettingsService) SaveProvider(p model.Provider) (*model.Provider, error) {
	return s.aiModel.CreateProvider(p)
}

func (s *SettingsService) UpdateProvider(id string, req model.ProviderUpdateRequest) (*model.Provider, error) {
	return s.aiModel.UpdateProvider(id, req)
}

func (s *SettingsService) QuickSetup(req model.AiModelQuickSetupRequest) (*model.SystemSettingsView, error) {
	if _, err := s.aiModel.QuickSetup(req); err != nil {
		return nil, err
	}
	return s.GetSystemView()
}

func (s *SettingsService) ReindexKnowledge(ctx context.Context) (*model.KnowledgeStatus, error) {
	return s.knowledge.RebuildIndex(ctx)
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

func (s *SettingsService) DefaultVirtualModel() string {
	if s.router != nil {
		return s.router.DefaultVirtualModelCode()
	}
	return ""
}

func (s *SettingsService) ResolvedLLM() (*model.ResolvedLLM, error) {
	if s.router != nil {
		return s.router.Resolve("")
	}
	return nil, errors.New("模型路由未初始化")
}

func (s *SettingsService) LLMClientFor(virtualModelCode string) (*llm.Client, *model.ResolvedLLM, error) {
	if s.router != nil {
		return s.router.ClientFor(virtualModelCode)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.llm != nil && s.llm.Enabled() {
		return s.llm, &model.ResolvedLLM{
			ModelCode: s.llmCfg.Model,
			BaseURL:   s.llmCfg.BaseURL,
			Enabled:   true,
		}, nil
	}
	return nil, nil, errors.New("LLM 未配置")
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
