package model

// CanonicalModel 统一模型定义，抽象不同厂商的同类模型。
type CanonicalModel struct {
	ID            string `gorm:"primaryKey;size:64" json:"id"`
	Code          string `gorm:"uniqueIndex;size:50" json:"code"`
	Name          string `gorm:"size:50" json:"name"`
	Status        int    `gorm:"default:1" json:"status"`
	ContextWindow int    `json:"contextWindow"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

// Capability 模型能力标签定义（如 CHAT、TOOL_CALLING、VISION）。
type Capability struct {
	ID        string `gorm:"primaryKey;size:64" json:"id"`
	Code      string `gorm:"uniqueIndex;size:50" json:"code"`
	Name      string `gorm:"size:50" json:"name"`
	Status    int    `gorm:"default:1" json:"status"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

// CapabilityModel 统一模型与能力标签的关联关系。
type CapabilityModel struct {
	ID               string `gorm:"primaryKey;size:64" json:"id"`
	CanonicalModelID string `gorm:"index;size:64" json:"canonicalModelId"`
	CapabilityID     string `gorm:"index;size:64" json:"capabilityId"`
	CreatedAt        int64  `json:"createdAt"`
	// 关联展示
	CanonicalModelCode string `gorm:"-" json:"canonicalModelCode,omitempty"`
	CanonicalModelName string `gorm:"-" json:"canonicalModelName,omitempty"`
	CapabilityCode     string `gorm:"-" json:"capabilityCode,omitempty"`
	CapabilityName     string `gorm:"-" json:"capabilityName,omitempty"`
}

// Provider AI 厂商定义。
type Provider struct {
	ID        string `gorm:"primaryKey;size:64" json:"id"`
	Code      string `gorm:"uniqueIndex;size:50" json:"code"`
	Name      string `gorm:"size:50" json:"name"`
	Status    int    `gorm:"default:1" json:"status"`
	BaseURL   string `gorm:"size:500" json:"baseUrl"`
	AuthType  string `gorm:"size:50" json:"authType"`
	APIKey    string `gorm:"size:200" json:"apiKey,omitempty"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	// 展示字段
	APIKeyMasked string `gorm:"-" json:"apiKeyMasked,omitempty"`
	ModelCount   int    `gorm:"-" json:"modelCount,omitempty"`
}

// ProviderModel 厂商模型实现定义，绑定厂商与统一模型。
type ProviderModel struct {
	ID               string `gorm:"primaryKey;size:64" json:"id"`
	ProviderID       string `gorm:"index;size:64" json:"providerId"`
	CanonicalModelID string `gorm:"index;size:64" json:"canonicalModelId"`
	ModelCode        string `gorm:"size:100" json:"modelCode"`
	DeploymentName   string `gorm:"size:100" json:"deploymentName"`
	Priority         int    `gorm:"default:10" json:"priority"`
	Weight           int    `gorm:"default:100" json:"weight"`
	ReasoningSupported bool `gorm:"default:false" json:"reasoningSupported"`
	Status           int    `gorm:"default:1" json:"status"`
	CreatedAt        int64  `json:"createdAt"`
	UpdatedAt        int64  `json:"updatedAt"`
	// 关联展示
	ProviderCode       string `gorm:"-" json:"providerCode,omitempty"`
	ProviderName       string `gorm:"-" json:"providerName,omitempty"`
	CanonicalModelCode string `gorm:"-" json:"canonicalModelCode,omitempty"`
	CanonicalModelName string `gorm:"-" json:"canonicalModelName,omitempty"`
}

// VirtualModel 虚拟模型定义，对外暴露稳定模型名。
type VirtualModel struct {
	ID        string `gorm:"primaryKey;size:64" json:"id"`
	Code      string `gorm:"uniqueIndex;size:50" json:"code"`
	Name      string `gorm:"size:50" json:"name"`
	Status    int    `gorm:"default:1" json:"status"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	MappingCount int `gorm:"-" json:"mappingCount,omitempty"`
}

// VirtualModelMapping 虚拟模型与统一模型的映射关系。
type VirtualModelMapping struct {
	ID               string `gorm:"primaryKey;size:64" json:"id"`
	VirtualModelID   string `gorm:"index;size:64" json:"virtualModelId"`
	CanonicalModelID string `gorm:"index;size:64" json:"canonicalModelId"`
	Priority         int    `gorm:"default:10" json:"priority"`
	CreatedAt        int64  `json:"createdAt"`
	// 关联展示
	VirtualModelCode   string `gorm:"-" json:"virtualModelCode,omitempty"`
	VirtualModelName string `gorm:"-" json:"virtualModelName,omitempty"`
	CanonicalModelCode string `gorm:"-" json:"canonicalModelCode,omitempty"`
	CanonicalModelName string `gorm:"-" json:"canonicalModelName,omitempty"`
}

// ResolvedLLM 模型路由解析结果。
type ResolvedLLM struct {
	VirtualModelCode   string `json:"virtualModelCode,omitempty"`
	CanonicalModelCode string `json:"canonicalModelCode"`
	ProviderCode       string `json:"providerCode"`
	ModelCode          string `json:"modelCode"`
	DeploymentName     string `json:"deploymentName,omitempty"`
	BaseURL            string `json:"baseUrl"`
	APIKey             string `json:"-"`
	Enabled            bool   `json:"enabled"`
}

// AiModelOverview AI 模型配置概览。
type AiModelOverview struct {
	ProviderCount       int64 `json:"providerCount"`
	CanonicalModelCount int64 `json:"canonicalModelCount"`
	CapabilityCount     int64 `json:"capabilityCount"`
	VirtualModelCount   int64 `json:"virtualModelCount"`
	ProviderModelCount  int64 `json:"providerModelCount"`
	DefaultVirtualModel string `json:"defaultVirtualModel"`
}
