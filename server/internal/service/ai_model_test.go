package service

import (
	"testing"

	"github.com/cng1985/ai-learning-server/internal/config"
	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newAiModelTestService(t *testing.T) (*AiModelService, *ModelRouter, *repository.AiModelRepo) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(
		&model.SystemSetting{},
		&model.CanonicalModel{},
		&model.Capability{},
		&model.CapabilityModel{},
		&model.Provider{},
		&model.ProviderModel{},
		&model.VirtualModel{},
		&model.VirtualModelMapping{},
	); err != nil {
		t.Fatal(err)
	}
	repo := repository.NewAiModelRepo(db)
	router := NewModelRouter(repo, repository.NewSettingsRepo(db), &config.Config{})
	return NewAiModelService(repo, router), router, repo
}

func TestUpdateProviderSupportsZeroStatusAndNeverReturnsSecret(t *testing.T) {
	svc, _, repo := newAiModelTestService(t)
	created, err := svc.CreateProvider(model.Provider{
		Code: "openai", Name: "OpenAI", APIKey: "sk-super-secret", Status: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.APIKey != "" || created.APIKeyMasked == "" {
		t.Fatalf("创建响应泄露或未掩码 API Key: %#v", created)
	}

	disabled := 0
	clear := true
	updated, err := svc.UpdateProvider(created.ID, model.ProviderUpdateRequest{
		Status: &disabled, ClearAPIKey: clear,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != 0 || updated.APIKey != "" || updated.APIKeyMasked != "" {
		t.Fatalf("禁用或清除密钥未生效: %#v", updated)
	}
	stored, err := repo.FindProvider(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status != 0 || stored.APIKey != "" {
		t.Fatalf("数据库状态不正确: %#v", stored)
	}
}

func TestModelRouterFallsBackToAnotherProviderModel(t *testing.T) {
	_, router, repo := newAiModelTestService(t)
	now := int64(1)
	records := []any{
		&model.CanonicalModel{ID: "cm", Code: "chat", Name: "Chat", Status: 1, CreatedAt: now, UpdatedAt: now},
		&model.Provider{ID: "disabled", Code: "disabled", Name: "Disabled", Status: 1, APIKey: "bad", CreatedAt: now, UpdatedAt: now},
		&model.Provider{ID: "ready", Code: "ready", Name: "Ready", Status: 1, BaseURL: "https://example.com/v1", APIKey: "good", CreatedAt: now, UpdatedAt: now},
		&model.ProviderModel{ID: "pm1", ProviderID: "disabled", CanonicalModelID: "cm", ModelCode: "bad-model", Priority: 10, Weight: 100, Status: 1, CreatedAt: now, UpdatedAt: now},
		&model.ProviderModel{ID: "pm2", ProviderID: "ready", CanonicalModelID: "cm", ModelCode: "good-model", Priority: 20, Weight: 100, Status: 1, CreatedAt: now, UpdatedAt: now},
		&model.VirtualModel{ID: "vm", Code: "chat-default", Name: "Default", Status: 1, CreatedAt: now, UpdatedAt: now},
		&model.VirtualModelMapping{ID: "mapping", VirtualModelID: "vm", CanonicalModelID: "cm", Priority: 10, CreatedAt: now},
	}
	for _, record := range records {
		switch v := record.(type) {
		case *model.CanonicalModel:
			if err := repo.CreateCanonicalModel(v); err != nil {
				t.Fatal(err)
			}
		case *model.Provider:
			if err := repo.CreateProvider(v); err != nil {
				t.Fatal(err)
			}
		case *model.ProviderModel:
			if err := repo.CreateProviderModel(v); err != nil {
				t.Fatal(err)
			}
		case *model.VirtualModel:
			if err := repo.CreateVirtualModel(v); err != nil {
				t.Fatal(err)
			}
		case *model.VirtualModelMapping:
			if err := repo.CreateVirtualModelMapping(v); err != nil {
				t.Fatal(err)
			}
		}
	}
	disabledProvider, err := repo.FindProvider("disabled")
	if err != nil {
		t.Fatal(err)
	}
	disabledProvider.Status = 0
	if err := repo.UpdateProvider(disabledProvider); err != nil {
		t.Fatal(err)
	}

	resolved, err := router.Resolve("chat-default")
	if err != nil {
		t.Fatal(err)
	}
	if resolved.ProviderCode != "ready" || resolved.ModelCode != "good-model" || !resolved.Enabled {
		t.Fatalf("未路由到可用备选厂商: %#v", resolved)
	}
}

func TestQuickSetupIsIdempotentAndUpdatesExistingChain(t *testing.T) {
	svc, _, repo := newAiModelTestService(t)
	req := model.AiModelQuickSetupRequest{
		ProviderCode: "custom", ProviderName: "Custom", BaseURL: "https://example.com/v1",
		APIKey: "first-key", CanonicalCode: "chat", CanonicalName: "Chat",
		ModelCode: "chat-v1", VirtualCode: "chat-default", VirtualName: "Default",
		ContextWindow: 32000,
	}
	if _, err := svc.QuickSetup(req); err != nil {
		t.Fatal(err)
	}
	req.APIKey = "second-key"
	req.CanonicalName = "Chat Updated"
	req.ContextWindow = 64000
	if _, err := svc.QuickSetup(req); err != nil {
		t.Fatal(err)
	}

	provider, err := repo.FindProviderByCode("custom")
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := repo.FindCanonicalModelByCode("chat")
	if err != nil {
		t.Fatal(err)
	}
	providerModels, err := repo.ListProviderModels(provider.ID, canonical.ID)
	if err != nil {
		t.Fatal(err)
	}
	mappings, err := repo.ListVirtualModelMappings("")
	if err != nil {
		t.Fatal(err)
	}
	if provider.APIKey != "second-key" || canonical.Name != "Chat Updated" || canonical.ContextWindow != 64000 {
		t.Fatalf("已有配置未更新: provider=%#v canonical=%#v", provider, canonical)
	}
	if len(providerModels) != 1 || len(mappings) != 1 {
		t.Fatalf("重复初始化创建了重复记录: providerModels=%d mappings=%d", len(providerModels), len(mappings))
	}
}
