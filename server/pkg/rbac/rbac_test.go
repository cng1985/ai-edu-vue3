package rbac

import "testing"

func TestIsAdminRole(t *testing.T) {
	for _, role := range []string{RoleAdmin, RoleReviewer, RoleOperator} {
		if !IsAdminRole(role) {
			t.Errorf("%s 应为后台角色", role)
		}
	}
	for _, role := range []string{RoleLearner, RoleGuest, "", "hacker"} {
		if IsAdminRole(role) {
			t.Errorf("%s 不应为后台角色", role)
		}
	}
}

func TestIsValidRole(t *testing.T) {
	for _, role := range Roles() {
		if !IsValidRole(role) {
			t.Errorf("%s 应为合法角色", role)
		}
	}
	if IsValidRole("superuser") {
		t.Error("未定义角色不应合法")
	}
}

func TestResolverDefaults(t *testing.T) {
	r := NewResolver()
	if !r.Has(RoleAdmin, PermUserDelete) {
		t.Error("admin 默认应有 user:delete")
	}
	if r.Has(RoleLearner, PermUserDelete) {
		t.Error("learner 默认不应有 user:delete")
	}
	if r.Has("unknown", PermAIChat) {
		t.Error("未知角色不应有任何权限")
	}
}

func TestResolverCustomOverride(t *testing.T) {
	r := NewResolver()
	r.SetCustom(map[string][]string{
		RoleReviewer: {PermCourseRead},
	})
	if !r.Has(RoleReviewer, PermCourseRead) {
		t.Error("自定义后 reviewer 应有 course:read")
	}
	if r.Has(RoleReviewer, PermReviewApprove) {
		t.Error("自定义覆盖后 reviewer 不应再有默认的 review:approve")
	}
	// 未自定义的角色仍回退默认
	if !r.Has(RoleAdmin, PermRoleManage) {
		t.Error("admin 未被自定义，应回退默认权限")
	}
}

func TestResolverPermissionsReturnsCopy(t *testing.T) {
	r := NewResolver()
	perms := r.Permissions(RoleGuest)
	if len(perms) == 0 {
		t.Fatal("guest 应有默认权限")
	}
	perms[0] = "tampered"
	if r.Permissions(RoleGuest)[0] == "tampered" {
		t.Error("Permissions 应返回副本，外部修改不应影响内部数据")
	}
}
