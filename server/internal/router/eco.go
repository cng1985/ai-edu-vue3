package router

import (
	"github.com/cng1985/ai-learning-server/internal/handler"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
	"github.com/gin-gonic/gin"
)

type permFunc func(string) gin.HandlerFunc

// registerEco 注册成长生态接口（学习者、创作者、企业共用）。
func registerEco(authed *gin.RouterGroup, h *handler.Handlers, perm permFunc) {
	write := perm(rbac.PermGrowthWrite)
	post := perm(rbac.PermCommunityPost)

	eco := authed.Group("/eco", perm(rbac.PermEcoRead))
	{
		eco.GET("/careers", h.Growth.Careers)
		eco.GET("/roles/:id", h.Growth.Role)
		eco.GET("/skills", h.Growth.Skills)
		eco.GET("/knowledge/graph", h.Growth.Graph)
		eco.GET("/knowledge/:id", h.Growth.Knowledge)

		eco.GET("/projects", h.Practice.Projects)
		eco.GET("/projects/:id", h.Practice.Project)
		eco.POST("/projects/:id/tasks/:taskId/submissions", write, h.Practice.Submit)

		eco.GET("/opportunities", h.Practice.Opportunities)
		eco.GET("/opportunities/:id", h.Practice.Opportunity)
		eco.POST("/opportunities/:id/apply", write, h.Practice.Apply)

		eco.GET("/community/posts", h.Community.Posts)
		eco.GET("/community/posts/:id", h.Community.Post)
		eco.POST("/community/posts", post, h.Community.CreatePost)
		eco.POST("/community/posts/:id/answers", post, h.Community.Answer)
		eco.POST("/community/posts/:id/like", post, h.Community.LikePost)
		eco.POST("/community/posts/:id/summarize", post, h.Community.Summarize)
		eco.POST("/community/answers/:id/like", post, h.Community.LikeAnswer)
		eco.POST("/community/answers/:id/accept", post, h.Community.AcceptAnswer)

		eco.GET("/resources", h.Community.Resources)
		eco.GET("/resources/:id", h.Community.Resource)
		eco.POST("/resources", perm(rbac.PermResourcePublish), h.Community.CreateResource)
	}

	me := authed.Group("/me", perm(rbac.PermEcoRead))
	{
		me.GET("/overview", h.Growth.Overview)
		me.GET("/goal", h.Growth.Goal)
		me.PUT("/goal", write, h.Growth.SetGoal)
		me.GET("/gap", h.Growth.Gap)
		me.GET("/skills", h.Growth.SkillStates)
		me.GET("/knowledge-states", h.Growth.KnowledgeStates)
		me.GET("/recommendations", h.Growth.Recommendations)
		me.GET("/reviews", h.Growth.DueReviews)
		me.GET("/profile", h.Growth.Profile)
		me.GET("/submissions", h.Practice.MySubmissions)
		me.GET("/applications", h.Practice.MyApplications)
		me.POST("/events", write, h.Growth.RecordEvent)
		me.POST("/knowledge/:id/practice", write, h.Growth.Practice)
		me.POST("/chapters/complete", write, h.Growth.CompleteChapter)
	}

	enterprise := authed.Group("/enterprise", perm(rbac.PermOpportunityPublish))
	{
		enterprise.GET("/opportunities", h.Practice.Published)
		enterprise.POST("/opportunities", h.Practice.Publish)
		enterprise.PUT("/opportunities/:id", h.Practice.UpdateOpportunity)
		enterprise.DELETE("/opportunities/:id", h.Practice.DeleteOpportunity)
		enterprise.GET("/opportunities/:id/candidates", h.Practice.Candidates)
		enterprise.POST("/applications/:id/decision", h.Practice.Decide)
	}

	talents := authed.Group("/talents", perm(rbac.PermTalentRead))
	{
		talents.GET("", h.Growth.Talents)
		talents.GET("/:id", h.Growth.TalentProfile)
	}
}

// registerEcoAdmin 注册管理端生态治理接口（已挂载后台角色校验）。
func registerEcoAdmin(admin *gin.RouterGroup, h *handler.Handlers, perm permFunc) {
	m := admin.Group("/manage")
	m.GET("/dashboard", perm(rbac.PermDashboard), h.CatalogAdmin.Dashboard)

	catalog := m.Group("", perm(rbac.PermEcoManage))
	{
		catalog.GET("/careers", h.Growth.Careers)
		catalog.POST("/careers", h.CatalogAdmin.SaveCareer)
		catalog.PUT("/careers/:id", h.CatalogAdmin.SaveCareer)
		catalog.DELETE("/careers/:id", h.CatalogAdmin.DeleteCareer)
		catalog.GET("/roles/:id", h.CatalogAdmin.Role)
		catalog.POST("/roles", h.CatalogAdmin.SaveRole)
		catalog.PUT("/roles/:id", h.CatalogAdmin.SaveRole)
		catalog.DELETE("/roles/:id", h.CatalogAdmin.DeleteRole)
		catalog.POST("/capabilities", h.CatalogAdmin.SaveCapability)
		catalog.PUT("/capabilities/:id", h.CatalogAdmin.SaveCapability)
		catalog.DELETE("/capabilities/:id", h.CatalogAdmin.DeleteCapability)
		catalog.GET("/skills", h.Growth.Skills)
		catalog.POST("/skills", h.CatalogAdmin.SaveSkill)
		catalog.PUT("/skills/:id", h.CatalogAdmin.SaveSkill)
		catalog.DELETE("/skills/:id", h.CatalogAdmin.DeleteSkill)
		catalog.GET("/knowledge", h.CatalogAdmin.Knowledge)
		catalog.POST("/knowledge", h.CatalogAdmin.SaveKnowledge)
		catalog.PUT("/knowledge/:id", h.CatalogAdmin.SaveKnowledge)
		catalog.DELETE("/knowledge/:id", h.CatalogAdmin.DeleteKnowledge)
		catalog.GET("/relations", h.CatalogAdmin.Relations)
		catalog.POST("/relations", h.CatalogAdmin.AddRelation)
		catalog.DELETE("/relations/:id", h.CatalogAdmin.DeleteRelation)
		catalog.GET("/chapter-knowledge", h.CatalogAdmin.ChapterMappings)
		catalog.PUT("/chapter-knowledge/:courseId/:chapterId", h.CatalogAdmin.SetChapterKnowledge)

		catalog.GET("/projects", h.Practice.AdminProjects)
		catalog.GET("/projects/:id", h.Practice.AdminProject)
		catalog.POST("/projects", h.Practice.SaveProject)
		catalog.PUT("/projects/:id", h.Practice.SaveProject)
		catalog.DELETE("/projects/:id", h.Practice.DeleteProject)
		catalog.GET("/submissions", h.Practice.AdminSubmissions)
	}

	market := m.Group("", perm(rbac.PermOpportunityManage))
	{
		market.GET("/opportunities", h.Practice.AdminOpportunities)
		market.PUT("/opportunities/:id", h.Practice.UpdateOpportunity)
		market.DELETE("/opportunities/:id", h.Practice.DeleteOpportunity)
		market.GET("/opportunities/:id/candidates", h.Practice.Candidates)
		market.POST("/applications/:id/decision", h.Practice.Decide)
	}

	m.GET("/talents", perm(rbac.PermTalentRead), h.Growth.Talents)
	m.GET("/talents/:id", perm(rbac.PermTalentRead), h.Growth.TalentProfile)

	community := m.Group("", perm(rbac.PermCommunityManage))
	{
		community.GET("/posts", h.Community.Posts)
		community.DELETE("/posts/:id", h.Community.DeletePost)
		community.GET("/resources", h.Community.Resources)
		community.DELETE("/resources/:id", h.Community.DeleteResource)
	}

	m.GET("/kernel/runs", perm(rbac.PermKernelRead), h.Kernel.AllRuns)
	m.GET("/kernel/stages", perm(rbac.PermKernelRead), h.Kernel.Stages)
	m.GET("/agents", perm(rbac.PermKernelRead), h.Kernel.Agents)
}
