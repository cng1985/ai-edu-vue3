package service

import (
	"sort"
	"strings"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
)

// TalentService 人才画像与人才库。
type TalentService struct {
	growth    *GrowthService
	projects  *ProjectService
	market    *MarketService
	community *CommunityService
	catalog   *CatalogService
	users     *repository.UserRepo
	gr        *repository.GrowthRepo
	repos     ecoCounters
}

type ecoCounters struct {
	catalog   *repository.CatalogRepo
	projects  *repository.ProjectRepo
	market    *repository.MarketRepo
	community *repository.CommunityRepo
	kernel    *repository.KernelRepo
}

func NewTalentService(growth *GrowthService, projects *ProjectService, market *MarketService, community *CommunityService,
	catalog *CatalogService, users *repository.UserRepo, gr *repository.GrowthRepo,
	catalogRepo *repository.CatalogRepo, projectRepo *repository.ProjectRepo, marketRepo *repository.MarketRepo,
	communityRepo *repository.CommunityRepo, kernelRepo *repository.KernelRepo) *TalentService {
	return &TalentService{
		growth: growth, projects: projects, market: market, community: community, catalog: catalog, users: users, gr: gr,
		repos: ecoCounters{catalog: catalogRepo, projects: projectRepo, market: marketRepo, community: communityRepo, kernel: kernelRepo},
	}
}

func contributionCount(c model.ContributionStats) int {
	return int(c.Posts + c.Answers + c.Resources)
}

// Profile 人才画像：职业目标、技能水平、知识掌握、项目经历、任务记录、社区贡献、收入记录。
func (s *TalentService) Profile(userID string) (*model.TalentProfile, error) {
	u, err := s.users.FindByID(userID)
	if err != nil {
		return nil, apperr.NotFound("用户不存在")
	}
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	p := &model.TalentProfile{User: userBrief(*u), JoinedAt: u.JoinedAt, Goal: s.growth.Goal(userID)}
	if p.Goal != nil {
		if career, ok := c.Careers[p.Goal.CareerID]; ok {
			p.GoalCareer = &career
		}
		if role, ok := c.Roles[p.Goal.RoleID]; ok {
			p.GoalRole = &role
		}
		if gap, err := s.growth.Gap(userID, p.Goal.RoleID); err == nil && gap != nil {
			p.Readiness = gap.Readiness
		}
	}
	p.Skills = s.growth.skillStatesWith(c, userID)
	p.Knowledge = s.knowledgeSummary(userID)
	p.Projects = s.projects.Submissions(userID, 20)
	p.Tasks = s.market.MyApplications(userID)
	ev, _ := s.gr.ListEvidence(userID)
	p.Evidence = nonNil(ev)
	p.Contributions = s.community.Contributions(userID)
	income, _ := s.gr.ListIncome(userID)
	p.Income = nonNil(income)
	for _, i := range income {
		p.TotalIncome += i.Amount
	}
	p.Flywheel = s.growth.Flywheel(userID, contributionCount(p.Contributions))
	p.Activity = s.growth.Activity(userID, 84)
	return p, nil
}

func (s *TalentService) knowledgeSummary(userID string) model.KnowledgeSummary {
	views, _ := s.growth.KnowledgeStates(userID)
	sum := model.KnowledgeSummary{Total: len(views)}
	domains := map[string]*model.DomainMastery{}
	var total float64
	for _, v := range views {
		switch v.Status {
		case model.KnowledgeStatusMastered:
			sum.Mastered++
		case model.KnowledgeStatusLearning:
			sum.Learning++
		}
		if v.Due {
			sum.Due++
		}
		total += v.EffectiveMastery
		d, ok := domains[v.Knowledge.Domain]
		if !ok {
			d = &model.DomainMastery{Domain: v.Knowledge.Domain}
			domains[v.Knowledge.Domain] = d
		}
		d.Count++
		d.Mastery += v.EffectiveMastery
	}
	if len(views) > 0 {
		sum.AvgMastery = round2(total / float64(len(views)))
	}
	for _, d := range domains {
		d.Mastery = round2(d.Mastery / float64(d.Count))
		sum.ByDomain = append(sum.ByDomain, *d)
	}
	sort.SliceStable(sum.ByDomain, func(i, j int) bool { return sum.ByDomain[i].Mastery > sum.ByDomain[j].Mastery })
	sum.ByDomain = nonNil(sum.ByDomain)
	return sum
}

func round2(v float64) float64 { return float64(int(v*100+0.5)) / 100 }

// List 人才库：按平均技能等级排序。
func (s *TalentService) List(keyword string) ([]model.TalentBrief, error) {
	users, err := s.users.ListByRoles(rbac.TalentRoles)
	if err != nil {
		return nil, err
	}
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	all, _ := s.gr.ListAllSkillStates()
	byUser := map[string][]model.StudentSkillState{}
	for _, st := range all {
		byUser[st.UserID] = append(byUser[st.UserID], st)
	}
	out := []model.TalentBrief{}
	for _, u := range users {
		if keyword != "" && !containsFold(u.Nickname, keyword) && !containsFold(u.Username, keyword) {
			continue
		}
		tb := model.TalentBrief{User: userBrief(u)}
		states := byUser[u.ID]
		sort.SliceStable(states, func(i, j int) bool {
			if states[i].Level != states[j].Level {
				return states[i].Level > states[j].Level
			}
			return states[i].Score > states[j].Score
		})
		total := 0
		for i, st := range states {
			total += st.Level
			tb.EvidenceCount += st.EvidenceCount
			if i < 4 {
				tb.TopSkills = append(tb.TopSkills, skillStateView(c.Skills[st.SkillID], st))
			}
		}
		if len(states) > 0 {
			tb.AvgLevel = round2(float64(total) / float64(len(states)))
		}
		tb.TopSkills = nonNil(tb.TopSkills)
		if g := s.growth.Goal(u.ID); g != nil {
			if r, ok := c.Roles[g.RoleID]; ok {
				tb.Goal = r.Name
			}
			if gap, err := s.growth.Gap(u.ID, g.RoleID); err == nil && gap != nil {
				tb.Readiness = gap.Readiness
			}
		}
		income, _ := s.gr.ListIncome(u.ID)
		for _, i := range income {
			tb.TotalIncome += i.Amount
		}
		out = append(out, tb)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].AvgLevel > out[j].AvgLevel })
	return out, nil
}

// Dashboard 管理端生态运营看板。
func (s *TalentService) Dashboard() model.EcoDashboard {
	r := s.repos
	d := model.EcoDashboard{
		Careers: r.catalog.Count(&model.Career{}), Roles: r.catalog.Count(&model.JobRole{}),
		Skills: r.catalog.Count(&model.Skill{}), Knowledge: r.catalog.Count(&model.KnowledgePoint{}),
		Relations: r.catalog.Count(&model.KnowledgeRelation{}),
		Projects:  r.projects.Count(), Submissions: r.projects.CountSubmissions(),
		Opportunities: r.market.CountByStatus(""), OpenOpportunities: r.market.CountByStatus(model.OpportunityOpen),
		Applications: r.market.CountApplications(""), CompletedTasks: r.market.CountApplications(model.ApplicationCompleted),
		Posts: r.community.CountPosts(), Resources: r.community.CountResourcesByType(),
		Goals: s.gr.CountGoals(), LearningEvents: s.gr.CountAllEvents(), Evidence: s.gr.CountEvidence(),
		TotalIncome: s.gr.SumIncome(), PipelineRuns: r.kernel.Count(),
		LevelDistribution: make([]int, 6),
	}
	all, _ := s.gr.ListAllSkillStates()
	for _, st := range all {
		d.LevelDistribution[clampLevel(st.Level)]++
	}
	return d
}

func containsFold(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}
