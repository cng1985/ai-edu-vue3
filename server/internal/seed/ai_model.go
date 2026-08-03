package seed

import (
	"fmt"
	"os"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
)

// SeedAiModels 初始化 AI 大模型分层配置（每次启动检查，缺失时补齐）。
func SeedAiModels(repo *repository.AiModelRepo, settings *repository.SettingsRepo) error {
	n, err := repo.CountProviders()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}

	fmt.Println("📦 初始化 AI 大模型配置...")
	now := time.Now().UnixMilli()
	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	modelName := os.Getenv("LLM_MODEL")
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}

	// 能力标签
	caps := []model.Capability{
		{ID: "cap_chat", Code: "CHAT", Name: "文本对话", Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "cap_tool", Code: "TOOL_CALLING", Name: "工具调用", Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "cap_vision", Code: "VISION", Name: "视觉理解", Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "cap_embed", Code: "EMBEDDING", Name: "向量嵌入", Status: 1, CreatedAt: now, UpdatedAt: now},
	}
	for _, c := range caps {
		if err := repo.CreateCapability(&c); err != nil {
			return err
		}
	}

	// 统一模型
	canonicals := []model.CanonicalModel{
		{ID: "cm_gpt4o_mini", Code: "gpt-4o-mini", Name: "GPT-4o Mini", Status: 1, ContextWindow: 128000, CreatedAt: now, UpdatedAt: now},
		{ID: "cm_gpt4o", Code: "gpt-4o", Name: "GPT-4o", Status: 1, ContextWindow: 128000, CreatedAt: now, UpdatedAt: now},
		{ID: "cm_deepseek", Code: "deepseek-chat", Name: "DeepSeek Chat", Status: 1, ContextWindow: 64000, CreatedAt: now, UpdatedAt: now},
	}
	for _, cm := range canonicals {
		if err := repo.CreateCanonicalModel(&cm); err != nil {
			return err
		}
	}

	// 能力关联
	capLinks := []model.CapabilityModel{
		{ID: "cml_1", CanonicalModelID: "cm_gpt4o_mini", CapabilityID: "cap_chat", CreatedAt: now},
		{ID: "cml_2", CanonicalModelID: "cm_gpt4o_mini", CapabilityID: "cap_tool", CreatedAt: now},
		{ID: "cml_3", CanonicalModelID: "cm_gpt4o", CapabilityID: "cap_chat", CreatedAt: now},
		{ID: "cml_4", CanonicalModelID: "cm_gpt4o", CapabilityID: "cap_tool", CreatedAt: now},
		{ID: "cml_5", CanonicalModelID: "cm_gpt4o", CapabilityID: "cap_vision", CreatedAt: now},
		{ID: "cml_6", CanonicalModelID: "cm_deepseek", CapabilityID: "cap_chat", CreatedAt: now},
	}
	for _, l := range capLinks {
		if err := repo.CreateCapabilityModel(&l); err != nil {
			return err
		}
	}

	// 厂商
	providers := []model.Provider{
		{ID: "pv_openai", Code: "openai", Name: "OpenAI", Status: 1, BaseURL: "https://api.openai.com/v1", AuthType: "Bearer", APIKey: apiKey, CreatedAt: now, UpdatedAt: now},
		{ID: "pv_deepseek", Code: "deepseek", Name: "DeepSeek", Status: 1, BaseURL: "https://api.deepseek.com/v1", AuthType: "Bearer", APIKey: "", CreatedAt: now, UpdatedAt: now},
	}
	for _, p := range providers {
		if err := repo.CreateProvider(&p); err != nil {
			return err
		}
	}

	// 厂商模型
	providerModels := []model.ProviderModel{
		{ID: "pm_openai_mini", ProviderID: "pv_openai", CanonicalModelID: "cm_gpt4o_mini", ModelCode: "gpt-4o-mini", Priority: 10, Weight: 100, Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "pm_openai_4o", ProviderID: "pv_openai", CanonicalModelID: "cm_gpt4o", ModelCode: "gpt-4o", Priority: 10, Weight: 100, Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "pm_deepseek_chat", ProviderID: "pv_deepseek", CanonicalModelID: "cm_deepseek", ModelCode: "deepseek-chat", Priority: 20, Weight: 100, Status: 1, CreatedAt: now, UpdatedAt: now},
	}
	for _, pm := range providerModels {
		if err := repo.CreateProviderModel(&pm); err != nil {
			return err
		}
	}

	// 虚拟模型
	virtuals := []model.VirtualModel{
		{ID: "vm_chat_default", Code: "chat-default", Name: "默认对话模型", Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "vm_chat_smart", Code: "chat-smart", Name: "高级对话模型", Status: 1, CreatedAt: now, UpdatedAt: now},
	}
	for _, v := range virtuals {
		if err := repo.CreateVirtualModel(&v); err != nil {
			return err
		}
	}

	// 虚拟模型映射
	mappings := []model.VirtualModelMapping{
		{ID: "vmm_1", VirtualModelID: "vm_chat_default", CanonicalModelID: "cm_gpt4o_mini", Priority: 10, CreatedAt: now},
		{ID: "vmm_2", VirtualModelID: "vm_chat_default", CanonicalModelID: "cm_deepseek", Priority: 20, CreatedAt: now},
		{ID: "vmm_3", VirtualModelID: "vm_chat_smart", CanonicalModelID: "cm_gpt4o", Priority: 10, CreatedAt: now},
		{ID: "vmm_4", VirtualModelID: "vm_chat_smart", CanonicalModelID: "cm_gpt4o_mini", Priority: 20, CreatedAt: now},
	}
	for _, m := range mappings {
		if err := repo.CreateVirtualModelMapping(&m); err != nil {
			return err
		}
	}

	// 默认虚拟模型
	_ = settings.Upsert("llm.default_virtual_model", "chat-default", now)

	fmt.Println("✅ AI 大模型配置已初始化")
	fmt.Printf("   默认虚拟模型: chat-default → %s\n", modelName)
	return nil
}
