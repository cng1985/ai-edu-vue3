package service

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/llm"
)

const settingDefaultVirtualModel = "llm.default_virtual_model"

// ModelRouter 虚拟模型 → 统一模型 → 厂商模型的路由解析。
type ModelRouter struct {
	repo     *repository.AiModelRepo
	settings *repository.SettingsRepo
	boot     *config.Config
	mu       sync.RWMutex
	clients  map[string]*llm.Client // key: resolved model identifier
}

func NewModelRouter(repo *repository.AiModelRepo, settings *repository.SettingsRepo, cfg *config.Config) *ModelRouter {
	return &ModelRouter{
		repo:     repo,
		settings: settings,
		boot:     cfg,
		clients:  make(map[string]*llm.Client),
	}
}

// DefaultVirtualModelCode 返回默认虚拟模型编码。
func (r *ModelRouter) DefaultVirtualModelCode() string {
	if val, err := r.settings.Get(settingDefaultVirtualModel); err == nil && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return "chat-default"
}

// SetDefaultVirtualModel 设置默认虚拟模型。
func (r *ModelRouter) SetDefaultVirtualModel(code string) error {
	now := timeNowMillis()
	return r.settings.Upsert(settingDefaultVirtualModel, strings.TrimSpace(code), now)
}

func timeNowMillis() int64 {
	return time.Now().UnixMilli()
}

// Resolve 解析虚拟模型到可用的 LLM 配置。
func (r *ModelRouter) Resolve(virtualModelCode string) (*model.ResolvedLLM, error) {
	if virtualModelCode == "" {
		virtualModelCode = r.DefaultVirtualModelCode()
	}

	vm, err := r.repo.FindVirtualModelByCode(virtualModelCode)
	if err != nil {
		// 回退到环境变量配置
		return r.fallbackResolved(virtualModelCode)
	}
	if vm.Status != 1 {
		return r.fallbackResolved(virtualModelCode)
	}

	mappings, err := r.repo.ListActiveMappingsByVirtual(vm.ID)
	if err != nil || len(mappings) == 0 {
		return r.fallbackResolved(virtualModelCode)
	}

	for _, mapping := range mappings {
		resolved, err := r.resolveCanonical(mapping.CanonicalModelID, virtualModelCode)
		if err == nil && resolved.Enabled {
			return resolved, nil
		}
	}
	return r.fallbackResolved(virtualModelCode)
}

func (r *ModelRouter) resolveCanonical(canonicalModelID, virtualModelCode string) (*model.ResolvedLLM, error) {
	canon, err := r.repo.FindCanonicalModel(canonicalModelID)
	if err != nil || canon.Status != 1 {
		return nil, errors.New("统一模型不可用")
	}

	providerModels, err := r.repo.ListActiveProviderModelsByCanonical(canonicalModelID)
	if err != nil || len(providerModels) == 0 {
		return nil, errors.New("无可用厂商模型")
	}

	// 按优先级分组，同优先级内按权重随机
	pm := pickProviderModel(providerModels)
	provider, err := r.repo.FindProvider(pm.ProviderID)
	if err != nil || provider.Status != 1 {
		return nil, errors.New("厂商不可用")
	}

	apiKey := provider.APIKey
	if apiKey == "" {
		apiKey = r.boot.LLM.APIKey
	}
	baseURL := provider.BaseURL
	if baseURL == "" {
		baseURL = r.boot.LLM.BaseURL
	}

	modelCode := pm.ModelCode
	if pm.DeploymentName != "" {
		modelCode = pm.DeploymentName
	}
	if modelCode == "" {
		modelCode = canon.Code
	}

	return &model.ResolvedLLM{
		VirtualModelCode:   virtualModelCode,
		CanonicalModelCode:   canon.Code,
		ProviderCode:       provider.Code,
		ModelCode:          modelCode,
		DeploymentName:     pm.DeploymentName,
		BaseURL:            baseURL,
		APIKey:             apiKey,
		Enabled:            apiKey != "",
	}, nil
}

func (r *ModelRouter) fallbackResolved(virtualModelCode string) (*model.ResolvedLLM, error) {
	cfg := r.boot.LLM
	return &model.ResolvedLLM{
		VirtualModelCode:   virtualModelCode,
		CanonicalModelCode: cfg.Model,
		ProviderCode:       "env",
		ModelCode:          cfg.Model,
		BaseURL:            cfg.BaseURL,
		APIKey:             cfg.APIKey,
		Enabled:            cfg.Enabled,
	}, nil
}

// ClientFor 获取解析后的 LLM 客户端。
func (r *ModelRouter) ClientFor(virtualModelCode string) (*llm.Client, *model.ResolvedLLM, error) {
	resolved, err := r.Resolve(virtualModelCode)
	if err != nil {
		return nil, nil, err
	}
	if !resolved.Enabled {
		return nil, resolved, errors.New("LLM 未配置")
	}
	key := resolved.ProviderCode + ":" + resolved.ModelCode
	r.mu.RLock()
	client, ok := r.clients[key]
	r.mu.RUnlock()
	if !ok {
		client = llm.NewClient(resolved.APIKey, resolved.BaseURL, resolved.ModelCode)
		r.mu.Lock()
		r.clients[key] = client
		r.mu.Unlock()
	}
	return client, resolved, nil
}

// InvalidateClients 清除客户端缓存。
func (r *ModelRouter) InvalidateClients() {
	r.mu.Lock()
	r.clients = make(map[string]*llm.Client)
	r.mu.Unlock()
}

func pickProviderModel(list []model.ProviderModel) model.ProviderModel {
	if len(list) == 1 {
		return list[0]
	}
	// 取最低优先级组
	minPriority := list[0].Priority
	for _, pm := range list {
		if pm.Priority < minPriority {
			minPriority = pm.Priority
		}
	}
	var candidates []model.ProviderModel
	totalWeight := 0
	for _, pm := range list {
		if pm.Priority == minPriority {
			candidates = append(candidates, pm)
			w := pm.Weight
			if w <= 0 {
				w = 100
			}
			totalWeight += w
		}
	}
	if len(candidates) == 1 {
		return candidates[0]
	}
	roll := rand.Intn(totalWeight)
	for _, pm := range candidates {
		w := pm.Weight
		if w <= 0 {
			w = 100
		}
		roll -= w
		if roll < 0 {
			return pm
		}
	}
	return candidates[0]
}
