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

	var firstConfigured *model.ResolvedLLM
	for _, pm := range orderProviderModels(providerModels) {
		provider, err := r.repo.FindProvider(pm.ProviderID)
		if err != nil || provider.Status != 1 {
			continue
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
		resolved := &model.ResolvedLLM{
			VirtualModelCode:   virtualModelCode,
			CanonicalModelCode: canon.Code,
			ProviderCode:       provider.Code,
			ModelCode:          modelCode,
			DeploymentName:     pm.DeploymentName,
			BaseURL:            baseURL,
			APIKey:             apiKey,
			AuthType:           provider.AuthType,
			Enabled:            apiKey != "",
		}
		if resolved.Enabled {
			return resolved, nil
		}
		if firstConfigured == nil {
			firstConfigured = resolved
		}
	}
	if firstConfigured != nil {
		return firstConfigured, nil
	}
	return nil, errors.New("无可用厂商")
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
	key := resolved.ProviderCode + ":" + resolved.ModelCode + ":" + resolved.AuthType
	r.mu.RLock()
	client, ok := r.clients[key]
	r.mu.RUnlock()
	if !ok {
		client = llm.NewClientWithAuth(resolved.APIKey, resolved.BaseURL, resolved.ModelCode, resolved.AuthType)
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

// orderProviderModels 保留优先级顺序，并在同优先级内按权重生成无重复尝试顺序。
// 这样首选厂商不可用时仍可尝试同组及后续优先级的备选厂商。
func orderProviderModels(list []model.ProviderModel) []model.ProviderModel {
	remaining := append([]model.ProviderModel(nil), list...)
	ordered := make([]model.ProviderModel, 0, len(list))
	for len(remaining) > 0 {
		minPriority := remaining[0].Priority
		for _, pm := range remaining[1:] {
			if pm.Priority < minPriority {
				minPriority = pm.Priority
			}
		}
		var group, rest []model.ProviderModel
		for _, pm := range remaining {
			if pm.Priority == minPriority {
				group = append(group, pm)
			} else {
				rest = append(rest, pm)
			}
		}
		for len(group) > 0 {
			totalWeight := 0
			for _, pm := range group {
				totalWeight += providerModelWeight(pm)
			}
			roll := rand.Intn(totalWeight)
			selected := 0
			for i, pm := range group {
				roll -= providerModelWeight(pm)
				if roll < 0 {
					selected = i
					break
				}
			}
			ordered = append(ordered, group[selected])
			group = append(group[:selected], group[selected+1:]...)
		}
		remaining = rest
	}
	return ordered
}

func providerModelWeight(pm model.ProviderModel) int {
	if pm.Weight <= 0 {
		return 100
	}
	return pm.Weight
}
