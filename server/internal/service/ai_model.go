package service

import (
	"errors"
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
)

type AiModelService struct {
	repo   *repository.AiModelRepo
	router *ModelRouter
}

func NewAiModelService(repo *repository.AiModelRepo, router *ModelRouter) *AiModelService {
	return &AiModelService{repo: repo, router: router}
}

func (s *AiModelService) Overview() (*model.AiModelOverview, error) {
	providers, _ := s.repo.CountProviders()
	canonical, _ := s.repo.CountCanonicalModels()
	caps, _ := s.repo.CountCapabilities()
	virtual, _ := s.repo.CountVirtualModels()
	pm, _ := s.repo.CountProviderModels()
	return &model.AiModelOverview{
		ProviderCount:       providers,
		CanonicalModelCount: canonical,
		CapabilityCount:     caps,
		VirtualModelCount:   virtual,
		ProviderModelCount:  pm,
		DefaultVirtualModel: s.router.DefaultVirtualModelCode(),
	}, nil
}

func (s *AiModelService) SetDefaultVirtualModel(code string) error {
	if _, err := s.repo.FindVirtualModelByCode(code); err != nil {
		return errors.New("虚拟模型不存在")
	}
	if err := s.router.SetDefaultVirtualModel(code); err != nil {
		return err
	}
	s.router.InvalidateClients()
	return nil
}

func (s *AiModelService) ResolveTest(virtualModelCode string) (*model.ResolvedLLM, error) {
	resolved, err := s.router.Resolve(virtualModelCode)
	if err != nil {
		return nil, err
	}
	// 不返回 API Key
	return resolved, nil
}

// --- CanonicalModel ---

func (s *AiModelService) ListCanonicalModels(keyword string, page, pageSize int) (*model.PageResult[model.CanonicalModel], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	list, total, err := s.repo.ListCanonicalModels(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.CanonicalModel{}
	}
	return &model.PageResult[model.CanonicalModel]{List: list, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *AiModelService) CreateCanonicalModel(m model.CanonicalModel) (*model.CanonicalModel, error) {
	if m.Code == "" || m.Name == "" {
		return nil, errors.New("编码和名称必填")
	}
	if _, err := s.repo.FindCanonicalModelByCode(m.Code); err == nil {
		return nil, errors.New("编码已存在")
	}
	now := time.Now().UnixMilli()
	m.ID = genID("cm")
	if m.Status == 0 {
		m.Status = 1
	}
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := s.repo.CreateCanonicalModel(&m); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	return &m, nil
}

func (s *AiModelService) UpdateCanonicalModel(id string, req model.CanonicalModel) (*model.CanonicalModel, error) {
	m, err := s.repo.FindCanonicalModel(id)
	if err != nil {
		return nil, errors.New("统一模型不存在")
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.ContextWindow > 0 {
		m.ContextWindow = req.ContextWindow
	}
	if req.Status != 0 {
		m.Status = req.Status
	}
	m.UpdatedAt = time.Now().UnixMilli()
	if err := s.repo.UpdateCanonicalModel(m); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	return m, nil
}

func (s *AiModelService) DeleteCanonicalModel(id string) error {
	if _, err := s.repo.FindCanonicalModel(id); err != nil {
		return errors.New("统一模型不存在")
	}
	s.router.InvalidateClients()
	return s.repo.DeleteCanonicalModel(id)
}

// --- Capability ---

func (s *AiModelService) ListCapabilities(keyword string, page, pageSize int) (*model.PageResult[model.Capability], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	list, total, err := s.repo.ListCapabilities(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.Capability{}
	}
	return &model.PageResult[model.Capability]{List: list, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *AiModelService) CreateCapability(m model.Capability) (*model.Capability, error) {
	if m.Code == "" || m.Name == "" {
		return nil, errors.New("编码和名称必填")
	}
	now := time.Now().UnixMilli()
	m.ID = genID("cap")
	if m.Status == 0 {
		m.Status = 1
	}
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := s.repo.CreateCapability(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *AiModelService) UpdateCapability(id string, req model.Capability) (*model.Capability, error) {
	m, err := s.repo.FindCapability(id)
	if err != nil {
		return nil, errors.New("能力标签不存在")
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Status != 0 {
		m.Status = req.Status
	}
	m.UpdatedAt = time.Now().UnixMilli()
	if err := s.repo.UpdateCapability(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AiModelService) DeleteCapability(id string) error {
	if _, err := s.repo.FindCapability(id); err != nil {
		return errors.New("能力标签不存在")
	}
	return s.repo.DeleteCapability(id)
}

// --- CapabilityModel ---

func (s *AiModelService) ListCapabilityModels(canonicalModelID string) ([]model.CapabilityModel, error) {
	list, err := s.repo.ListCapabilityModels(canonicalModelID)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if cm, err := s.repo.FindCanonicalModel(list[i].CanonicalModelID); err == nil {
			list[i].CanonicalModelCode = cm.Code
			list[i].CanonicalModelName = cm.Name
		}
		if cap, err := s.repo.FindCapability(list[i].CapabilityID); err == nil {
			list[i].CapabilityCode = cap.Code
			list[i].CapabilityName = cap.Name
		}
	}
	return list, nil
}

func (s *AiModelService) CreateCapabilityModel(m model.CapabilityModel) (*model.CapabilityModel, error) {
	if m.CanonicalModelID == "" || m.CapabilityID == "" {
		return nil, errors.New("统一模型和能力标签必填")
	}
	if s.repo.ExistsCapabilityModel(m.CanonicalModelID, m.CapabilityID) {
		return nil, errors.New("关联已存在")
	}
	m.ID = genID("cml")
	m.CreatedAt = time.Now().UnixMilli()
	if err := s.repo.CreateCapabilityModel(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *AiModelService) DeleteCapabilityModel(id string) error {
	if _, err := s.repo.FindCapabilityModel(id); err != nil {
		return errors.New("关联不存在")
	}
	return s.repo.DeleteCapabilityModel(id)
}

// --- Provider ---

func (s *AiModelService) ListProviders(keyword string, page, pageSize int) (*model.PageResult[model.Provider], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	list, total, err := s.repo.ListProviders(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.Provider{}
	}
	for i := range list {
		list[i].APIKeyMasked = maskSecret(list[i].APIKey)
		list[i].APIKey = ""
	}
	return &model.PageResult[model.Provider]{List: list, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *AiModelService) GetProvider(id string) (*model.Provider, error) {
	p, err := s.repo.FindProvider(id)
	if err != nil {
		return nil, errors.New("厂商不存在")
	}
	p.APIKeyMasked = maskSecret(p.APIKey)
	return p, nil
}

func (s *AiModelService) CreateProvider(m model.Provider) (*model.Provider, error) {
	if m.Code == "" || m.Name == "" {
		return nil, errors.New("编码和名称必填")
	}
	now := time.Now().UnixMilli()
	m.ID = genID("pv")
	if m.Status == 0 {
		m.Status = 1
	}
	if m.AuthType == "" {
		m.AuthType = "Bearer"
	}
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := s.repo.CreateProvider(&m); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	m.APIKeyMasked = maskSecret(m.APIKey)
	return &m, nil
}

func (s *AiModelService) UpdateProvider(id string, req model.Provider) (*model.Provider, error) {
	p, err := s.repo.FindProvider(id)
	if err != nil {
		return nil, errors.New("厂商不存在")
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.BaseURL != "" {
		p.BaseURL = req.BaseURL
	}
	if req.AuthType != "" {
		p.AuthType = req.AuthType
	}
	if strings.TrimSpace(req.APIKey) != "" {
		p.APIKey = strings.TrimSpace(req.APIKey)
	}
	if req.Status != 0 {
		p.Status = req.Status
	}
	p.UpdatedAt = time.Now().UnixMilli()
	if err := s.repo.UpdateProvider(p); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	p.APIKeyMasked = maskSecret(p.APIKey)
	return p, nil
}

func (s *AiModelService) DeleteProvider(id string) error {
	if _, err := s.repo.FindProvider(id); err != nil {
		return errors.New("厂商不存在")
	}
	s.router.InvalidateClients()
	return s.repo.DeleteProvider(id)
}

// --- ProviderModel ---

func (s *AiModelService) ListProviderModels(providerID, canonicalModelID string) ([]model.ProviderModel, error) {
	list, err := s.repo.ListProviderModels(providerID, canonicalModelID)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if p, err := s.repo.FindProvider(list[i].ProviderID); err == nil {
			list[i].ProviderCode = p.Code
			list[i].ProviderName = p.Name
		}
		if cm, err := s.repo.FindCanonicalModel(list[i].CanonicalModelID); err == nil {
			list[i].CanonicalModelCode = cm.Code
			list[i].CanonicalModelName = cm.Name
		}
	}
	return list, nil
}

func (s *AiModelService) CreateProviderModel(m model.ProviderModel) (*model.ProviderModel, error) {
	if m.ProviderID == "" || m.CanonicalModelID == "" || m.ModelCode == "" {
		return nil, errors.New("厂商、统一模型和模型标识必填")
	}
	now := time.Now().UnixMilli()
	m.ID = genID("pm")
	if m.Status == 0 {
		m.Status = 1
	}
	if m.Priority == 0 {
		m.Priority = 10
	}
	if m.Weight == 0 {
		m.Weight = 100
	}
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := s.repo.CreateProviderModel(&m); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	return &m, nil
}

func (s *AiModelService) UpdateProviderModel(id string, req model.ProviderModel) (*model.ProviderModel, error) {
	m, err := s.repo.FindProviderModel(id)
	if err != nil {
		return nil, errors.New("厂商模型不存在")
	}
	if req.ModelCode != "" {
		m.ModelCode = req.ModelCode
	}
	m.DeploymentName = req.DeploymentName
	if req.Priority > 0 {
		m.Priority = req.Priority
	}
	if req.Weight > 0 {
		m.Weight = req.Weight
	}
	m.ReasoningSupported = req.ReasoningSupported
	if req.Status != 0 {
		m.Status = req.Status
	}
	m.UpdatedAt = time.Now().UnixMilli()
	if err := s.repo.UpdateProviderModel(m); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	return m, nil
}

func (s *AiModelService) DeleteProviderModel(id string) error {
	if _, err := s.repo.FindProviderModel(id); err != nil {
		return errors.New("厂商模型不存在")
	}
	s.router.InvalidateClients()
	return s.repo.DeleteProviderModel(id)
}

// --- VirtualModel ---

func (s *AiModelService) ListVirtualModels(keyword string, page, pageSize int) (*model.PageResult[model.VirtualModel], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	list, total, err := s.repo.ListVirtualModels(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []model.VirtualModel{}
	}
	return &model.PageResult[model.VirtualModel]{List: list, Total: int(total), Page: page, PageSize: pageSize}, nil
}

func (s *AiModelService) CreateVirtualModel(m model.VirtualModel) (*model.VirtualModel, error) {
	if m.Code == "" || m.Name == "" {
		return nil, errors.New("编码和名称必填")
	}
	now := time.Now().UnixMilli()
	m.ID = genID("vm")
	if m.Status == 0 {
		m.Status = 1
	}
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := s.repo.CreateVirtualModel(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *AiModelService) UpdateVirtualModel(id string, req model.VirtualModel) (*model.VirtualModel, error) {
	m, err := s.repo.FindVirtualModel(id)
	if err != nil {
		return nil, errors.New("虚拟模型不存在")
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Status != 0 {
		m.Status = req.Status
	}
	m.UpdatedAt = time.Now().UnixMilli()
	if err := s.repo.UpdateVirtualModel(m); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	return m, nil
}

func (s *AiModelService) DeleteVirtualModel(id string) error {
	if _, err := s.repo.FindVirtualModel(id); err != nil {
		return errors.New("虚拟模型不存在")
	}
	s.router.InvalidateClients()
	return s.repo.DeleteVirtualModel(id)
}

// --- VirtualModelMapping ---

func (s *AiModelService) ListVirtualModelMappings(virtualModelID string) ([]model.VirtualModelMapping, error) {
	list, err := s.repo.ListVirtualModelMappings(virtualModelID)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if vm, err := s.repo.FindVirtualModel(list[i].VirtualModelID); err == nil {
			list[i].VirtualModelCode = vm.Code
			list[i].VirtualModelName = vm.Name
		}
		if cm, err := s.repo.FindCanonicalModel(list[i].CanonicalModelID); err == nil {
			list[i].CanonicalModelCode = cm.Code
			list[i].CanonicalModelName = cm.Name
		}
	}
	return list, nil
}

func (s *AiModelService) CreateVirtualModelMapping(m model.VirtualModelMapping) (*model.VirtualModelMapping, error) {
	if m.VirtualModelID == "" || m.CanonicalModelID == "" {
		return nil, errors.New("虚拟模型和统一模型必填")
	}
	m.ID = genID("vmm")
	if m.Priority == 0 {
		m.Priority = 10
	}
	m.CreatedAt = time.Now().UnixMilli()
	if err := s.repo.CreateVirtualModelMapping(&m); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	return &m, nil
}

func (s *AiModelService) UpdateVirtualModelMapping(id string, req model.VirtualModelMapping) (*model.VirtualModelMapping, error) {
	m, err := s.repo.FindVirtualModelMapping(id)
	if err != nil {
		return nil, errors.New("映射不存在")
	}
	if req.Priority > 0 {
		m.Priority = req.Priority
	}
	if err := s.repo.UpdateVirtualModelMapping(m); err != nil {
		return nil, err
	}
	s.router.InvalidateClients()
	return m, nil
}

func (s *AiModelService) ListVirtualModelOptions() ([]model.VirtualModel, error) {
	res, err := s.ListVirtualModels("", 1, 200)
	if err != nil {
		return nil, err
	}
	return res.List, nil
}

func (s *AiModelService) DeleteVirtualModelMapping(id string) error {
	if _, err := s.repo.FindVirtualModelMapping(id); err != nil {
		return errors.New("映射不存在")
	}
	s.router.InvalidateClients()
	return s.repo.DeleteVirtualModelMapping(id)
}

// QuickSetup 一键创建厂商 → 统一模型 → 厂商模型 → 虚拟模型 → 映射的完整链路。
func (s *AiModelService) QuickSetup(req model.AiModelQuickSetupRequest) (*model.AiModelOverview, error) {
	req.ProviderCode = strings.TrimSpace(req.ProviderCode)
	req.ProviderName = strings.TrimSpace(req.ProviderName)
	req.BaseURL = strings.TrimSpace(req.BaseURL)
	req.APIKey = strings.TrimSpace(req.APIKey)
	req.CanonicalCode = strings.TrimSpace(req.CanonicalCode)
	req.CanonicalName = strings.TrimSpace(req.CanonicalName)
	req.ModelCode = strings.TrimSpace(req.ModelCode)
	req.VirtualCode = strings.TrimSpace(req.VirtualCode)
	req.VirtualName = strings.TrimSpace(req.VirtualName)

	if req.ProviderCode == "" || req.ProviderName == "" {
		return nil, errors.New("厂商编码和名称必填")
	}
	if req.APIKey == "" {
		return nil, errors.New("API Key 必填")
	}
	if req.CanonicalCode == "" || req.CanonicalName == "" {
		return nil, errors.New("统一模型编码和名称必填")
	}
	if req.ModelCode == "" {
		req.ModelCode = req.CanonicalCode
	}
	if req.VirtualCode == "" {
		req.VirtualCode = "chat-default"
	}
	if req.VirtualName == "" {
		req.VirtualName = "默认对话模型"
	}
	if req.BaseURL == "" {
		req.BaseURL = "https://api.openai.com/v1"
	}
	if req.ContextWindow == 0 {
		req.ContextWindow = 128000
	}

	now := time.Now().UnixMilli()

	// 1. 厂商
	provider, err := s.repo.FindProviderByCode(req.ProviderCode)
	if err != nil {
		provider = &model.Provider{
			ID: genID("pv"), Code: req.ProviderCode, Name: req.ProviderName,
			BaseURL: req.BaseURL, AuthType: "Bearer", APIKey: req.APIKey,
			Status: 1, CreatedAt: now, UpdatedAt: now,
		}
		if err := s.repo.CreateProvider(provider); err != nil {
			return nil, err
		}
	} else {
		provider.Name = req.ProviderName
		provider.BaseURL = req.BaseURL
		provider.APIKey = req.APIKey
		provider.Status = 1
		provider.UpdatedAt = now
		if err := s.repo.UpdateProvider(provider); err != nil {
			return nil, err
		}
	}

	// 2. 统一模型
	canonical, err := s.repo.FindCanonicalModelByCode(req.CanonicalCode)
	if err != nil {
		canonical = &model.CanonicalModel{
			ID: genID("cm"), Code: req.CanonicalCode, Name: req.CanonicalName,
			Status: 1, ContextWindow: req.ContextWindow, CreatedAt: now, UpdatedAt: now,
		}
		if err := s.repo.CreateCanonicalModel(canonical); err != nil {
			return nil, err
		}
	}

	// 3. 厂商模型
	pms, _ := s.repo.ListProviderModels(provider.ID, canonical.ID)
	if len(pms) == 0 {
		pm := model.ProviderModel{
			ID: genID("pm"), ProviderID: provider.ID, CanonicalModelID: canonical.ID,
			ModelCode: req.ModelCode, Priority: 10, Weight: 100, Status: 1,
			CreatedAt: now, UpdatedAt: now,
		}
		if err := s.repo.CreateProviderModel(&pm); err != nil {
			return nil, err
		}
	}

	// 4. 虚拟模型
	virtual, err := s.repo.FindVirtualModelByCode(req.VirtualCode)
	if err != nil {
		virtual = &model.VirtualModel{
			ID: genID("vm"), Code: req.VirtualCode, Name: req.VirtualName,
			Status: 1, CreatedAt: now, UpdatedAt: now,
		}
		if err := s.repo.CreateVirtualModel(virtual); err != nil {
			return nil, err
		}
	}

	// 5. 虚拟映射
	mappings, _ := s.repo.ListVirtualModelMappings(virtual.ID)
	hasMapping := false
	for _, m := range mappings {
		if m.CanonicalModelID == canonical.ID {
			hasMapping = true
			break
		}
	}
	if !hasMapping {
		mapping := model.VirtualModelMapping{
			ID: genID("vmm"), VirtualModelID: virtual.ID, CanonicalModelID: canonical.ID,
			Priority: 10, CreatedAt: now,
		}
		if err := s.repo.CreateVirtualModelMapping(&mapping); err != nil {
			return nil, err
		}
	}

	// 6. 设为默认
	if err := s.SetDefaultVirtualModel(req.VirtualCode); err != nil {
		return nil, err
	}

	s.router.InvalidateClients()
	return s.Overview()
}
