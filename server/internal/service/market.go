package service

import (
	"sort"
	"strings"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/growth"
	"github.com/cng1985/ai-learning-server/pkg/rbac"
)

// MarketService IT 任务平台：企业发布任务 → 技能模型匹配 → 人才画像 → 推荐人员 → 履约与收益。
type MarketService struct {
	repo    *repository.MarketRepo
	growth  *GrowthService
	catalog *CatalogService
	users   *repository.UserRepo
	gr      *repository.GrowthRepo
}

func NewMarketService(repo *repository.MarketRepo, growth *GrowthService, catalog *CatalogService, users *repository.UserRepo, gr *repository.GrowthRepo) *MarketService {
	return &MarketService{repo: repo, growth: growth, catalog: catalog, users: users, gr: gr}
}

func (s *MarketService) match(c *Catalog, reqs []model.SkillRequirement, levels map[string]int) *model.MatchView {
	mv := matchView(growth.Fit(toEngineReqs(reqs), levels), c)
	mv.Items = nonNil(mv.Items)
	return &mv
}

func (s *MarketService) view(c *Catalog, o model.Opportunity, userID string, levels map[string]int) model.OpportunityView {
	v := model.OpportunityView{Opportunity: o, Requirements: parseRequirements(o.Requirements)}
	apps, _ := s.repo.ListApplications(o.ID, "")
	v.ApplicantCount = len(apps)
	if levels != nil {
		v.Match = s.match(c, v.Requirements, levels)
	}
	for i := range apps {
		if apps[i].UserID == userID {
			v.MyApplication = &apps[i]
		}
	}
	return v
}

// List 任务市场列表；为人才计算个人匹配度，sortBy=match 时按匹配度排序。
func (s *MarketService) List(userID, role, status, sortBy string) ([]model.OpportunityView, error) {
	list, err := s.repo.List(status, "")
	if err != nil {
		return nil, err
	}
	c, _ := s.catalog.Load()
	var levels map[string]int
	if isTalent(role) {
		levels = s.growth.SkillLevels(userID)
	}
	out := make([]model.OpportunityView, 0, len(list))
	for _, o := range list {
		out = append(out, s.view(c, o, userID, levels))
	}
	if sortBy == "match" && levels != nil {
		sort.SliceStable(out, func(i, j int) bool { return out[i].Match.Score > out[j].Match.Score })
	}
	return out, nil
}

// TopMatches 为人才推荐匹配度最高的开放任务。
func (s *MarketService) TopMatches(userID, role string, limit int) []model.OpportunityView {
	if !isTalent(role) {
		return nil
	}
	list, _ := s.List(userID, role, model.OpportunityOpen, "match")
	if len(list) > limit {
		list = list[:limit]
	}
	return list
}

func isTalent(role string) bool {
	for _, r := range rbac.TalentRoles {
		if r == role {
			return true
		}
	}
	return false
}

// OpportunityDetail 任务详情；发布者可见申请列表。
type OpportunityDetail struct {
	model.OpportunityView
	Publisher    *model.UserBrief        `json:"publisher,omitempty"`
	IsOwner      bool                    `json:"isOwner"`
	Applications []model.ApplicationView `json:"applications"`
}

func (s *MarketService) Get(userID, role, id string, canManage bool) (*OpportunityDetail, error) {
	o, err := s.repo.Find(id)
	if err != nil {
		return nil, apperr.NotFound("任务不存在")
	}
	c, _ := s.catalog.Load()
	var levels map[string]int
	if isTalent(role) {
		levels = s.growth.SkillLevels(userID)
	}
	d := &OpportunityDetail{OpportunityView: s.view(c, *o, userID, levels), IsOwner: o.PublisherID == userID}
	if u, err := s.users.FindByID(o.PublisherID); err == nil {
		b := userBrief(*u)
		d.Publisher = &b
	}
	if d.IsOwner || canManage {
		d.Applications = s.applicationViews(c, o, "")
	}
	return d, nil
}

func (s *MarketService) applicationViews(c *Catalog, o *model.Opportunity, userID string) []model.ApplicationView {
	var apps []model.OpportunityApplication
	if o != nil {
		apps, _ = s.repo.ListApplications(o.ID, "")
	} else {
		apps, _ = s.repo.ListApplications("", userID)
	}
	out := make([]model.ApplicationView, 0, len(apps))
	for _, a := range apps {
		v := model.ApplicationView{OpportunityApplication: a}
		opp := o
		if opp == nil {
			opp, _ = s.repo.Find(a.OpportunityID)
		}
		if opp != nil {
			oc := *opp
			oc.Description = ""
			v.Opportunity = &oc
			if o != nil {
				v.Match = s.match(c, parseRequirements(opp.Requirements), s.growth.SkillLevels(a.UserID))
			}
		}
		if u, err := s.users.FindByID(a.UserID); err == nil {
			b := userBrief(*u)
			v.User = &b
		}
		out = append(out, v)
	}
	return out
}

func (s *MarketService) normalize(req *model.OpportunityRequest) error {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" || strings.TrimSpace(req.Description) == "" {
		return apperr.BadRequest("请填写任务标题与需求描述")
	}
	if len(req.Requirements) == 0 {
		return apperr.BadRequest("请至少设置一项技能要求")
	}
	c, err := s.catalog.Load()
	if err != nil {
		return err
	}
	for i := range req.Requirements {
		if _, ok := c.Skills[req.Requirements[i].SkillID]; !ok {
			return apperr.BadRequest("技能要求中存在未知技能")
		}
		req.Requirements[i].Level = clampRange(req.Requirements[i].Level, 1, 5, 2)
	}
	if req.Mode == "" {
		req.Mode = "remote"
	}
	return nil
}

// Publish 企业发布任务。
func (s *MarketService) Publish(publisherID string, req model.OpportunityRequest) (*model.Opportunity, error) {
	if err := s.normalize(&req); err != nil {
		return nil, err
	}
	if req.Company == "" {
		if u, err := s.users.FindByID(publisherID); err == nil {
			req.Company = u.Nickname
		}
	}
	now := time.Now().UnixMilli()
	o := &model.Opportunity{
		ID: genID("opp"), Title: req.Title, Company: req.Company, PublisherID: publisherID,
		Description: req.Description, Requirements: toJSON(req.Requirements), Budget: req.Budget,
		DurationDays: clampRange(req.DurationDays, 1, 365, 14), Mode: req.Mode, Status: model.OpportunityOpen,
		CreatedAt: now, UpdatedAt: now,
	}
	return o, s.repo.Save(o)
}

func (s *MarketService) Update(userID, id string, canManage bool, req model.OpportunityRequest) (*model.Opportunity, error) {
	o, err := s.owned(userID, id, canManage)
	if err != nil {
		return nil, err
	}
	if err := s.normalize(&req); err != nil {
		return nil, err
	}
	o.Title, o.Description, o.Requirements = req.Title, req.Description, toJSON(req.Requirements)
	o.Budget, o.DurationDays, o.Mode = req.Budget, clampRange(req.DurationDays, 1, 365, 14), req.Mode
	if req.Company != "" {
		o.Company = req.Company
	}
	switch req.Status {
	case model.OpportunityOpen, model.OpportunityInProgress, model.OpportunityClosed:
		o.Status = req.Status
	}
	o.UpdatedAt = time.Now().UnixMilli()
	return o, s.repo.Save(o)
}

func (s *MarketService) Delete(userID, id string, canManage bool) error {
	if _, err := s.owned(userID, id, canManage); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *MarketService) owned(userID, id string, canManage bool) (*model.Opportunity, error) {
	o, err := s.repo.Find(id)
	if err != nil {
		return nil, apperr.NotFound("任务不存在")
	}
	if o.PublisherID != userID && !canManage {
		return nil, apperr.Forbidden("只能管理自己发布的任务")
	}
	return o, nil
}

// Apply 人才申请任务，记录申请时的匹配度。
func (s *MarketService) Apply(userID, role, id string, req model.ApplyRequest) (*model.OpportunityApplication, error) {
	if !isTalent(role) {
		return nil, apperr.Forbidden("仅学习者与创作者可以申请任务")
	}
	o, err := s.repo.Find(id)
	if err != nil {
		return nil, apperr.NotFound("任务不存在")
	}
	if o.Status != model.OpportunityOpen {
		return nil, apperr.BadRequest("该任务已停止招募")
	}
	if _, err := s.repo.FindUserApplication(id, userID); err == nil {
		return nil, apperr.Conflict("你已申请过该任务")
	}
	fit := growth.Fit(toEngineReqs(parseRequirements(o.Requirements)), s.growth.SkillLevels(userID))
	now := time.Now().UnixMilli()
	a := &model.OpportunityApplication{
		ID: genID("app"), OpportunityID: id, UserID: userID, Message: strings.TrimSpace(req.Message),
		MatchScore: fit.Score, Status: model.ApplicationApplied, CreatedAt: now, UpdatedAt: now,
	}
	return a, s.repo.SaveApplication(a)
}

func (s *MarketService) MyApplications(userID string) []model.ApplicationView {
	c, _ := s.catalog.Load()
	return s.applicationViews(c, nil, userID)
}

// Published 发布者的任务及申请情况。
func (s *MarketService) Published(userID string) ([]OpportunityDetail, error) {
	list, err := s.repo.List("", userID)
	if err != nil {
		return nil, err
	}
	c, _ := s.catalog.Load()
	out := make([]OpportunityDetail, 0, len(list))
	for i := range list {
		o := list[i]
		out = append(out, OpportunityDetail{
			OpportunityView: s.view(c, o, "", nil), IsOwner: true, Applications: s.applicationViews(c, &o, ""),
		})
	}
	return out, nil
}

// AdminList 管理端：全部任务及申请。
func (s *MarketService) AdminList(status string) ([]OpportunityDetail, error) {
	list, err := s.repo.List(status, "")
	if err != nil {
		return nil, err
	}
	c, _ := s.catalog.Load()
	out := make([]OpportunityDetail, 0, len(list))
	for i := range list {
		o := list[i]
		d := OpportunityDetail{OpportunityView: s.view(c, o, "", nil), Applications: s.applicationViews(c, &o, "")}
		if u, err := s.users.FindByID(o.PublisherID); err == nil {
			b := userBrief(*u)
			d.Publisher = &b
		}
		out = append(out, d)
	}
	return out, nil
}

// Decide 发布者处理申请：accept 录用 / reject 婉拒 / complete 验收结算。
// 验收后生成收入记录与任务表现证据，驱动"真实任务 → 获得收益 → 能力提升"闭环。
func (s *MarketService) Decide(userID, appID string, canManage bool, d model.ApplicationDecision) (*model.OpportunityApplication, error) {
	a, err := s.repo.FindApplication(appID)
	if err != nil {
		return nil, apperr.NotFound("申请不存在")
	}
	o, err := s.owned(userID, a.OpportunityID, canManage)
	if err != nil {
		return nil, err
	}
	now := time.Now().UnixMilli()
	switch d.Action {
	case "accept":
		if a.Status != model.ApplicationApplied {
			return nil, apperr.BadRequest("只能录用待处理的申请")
		}
		a.Status = model.ApplicationAccepted
		o.Status = model.OpportunityInProgress
	case "reject":
		if a.Status != model.ApplicationApplied && a.Status != model.ApplicationAccepted {
			return nil, apperr.BadRequest("当前状态不能婉拒")
		}
		a.Status = model.ApplicationRejected
	case "complete":
		if a.Status != model.ApplicationAccepted {
			return nil, apperr.BadRequest("只能验收已录用的任务")
		}
		rating := clampRange(d.Rating, 1, 5, 4)
		payout := d.Payout
		if payout <= 0 {
			payout = o.Budget
		}
		a.Status, a.Rating, a.Review, a.Payout = model.ApplicationCompleted, rating, strings.TrimSpace(d.Review), payout
		o.Status = model.OpportunityClosed
		if payout > 0 {
			_ = s.gr.AddIncome(&model.IncomeRecord{
				ID: genID("inc"), UserID: a.UserID, SourceType: "opportunity", SourceID: o.ID,
				Title: o.Company + " · " + o.Title, Amount: payout, CreatedAt: now,
			})
		}
		skillIDs := []string{}
		for _, r := range parseRequirements(o.Requirements) {
			skillIDs = append(skillIDs, r.SkillID)
		}
		detail := a.Review
		if detail == "" {
			detail = "企业验收通过"
		}
		s.growth.AddEvidence(a.UserID, skillIDs, model.EvidenceOpportunity, o.ID, o.Company+" · "+o.Title, ratingScore(rating), detail)
	default:
		return nil, apperr.BadRequest("未知操作")
	}
	a.UpdatedAt, o.UpdatedAt = now, now
	if err := s.repo.SaveApplication(a); err != nil {
		return nil, err
	}
	_ = s.repo.Save(o)
	return a, nil
}

// ratingScore 将企业 1-5 星评价映射为任务表现分。
func ratingScore(rating int) int {
	return []int{0, 45, 60, 75, 88, 96}[clampRange(rating, 1, 5, 4)]
}

// Candidates 按技能模型为任务推荐人才。
func (s *MarketService) Candidates(userID, id string, canManage bool, limit int) ([]model.CandidateView, error) {
	o, err := s.owned(userID, id, canManage)
	if err != nil {
		return nil, err
	}
	c, err := s.catalog.Load()
	if err != nil {
		return nil, err
	}
	reqs := parseRequirements(o.Requirements)
	talents, err := s.users.ListByRoles(rbac.TalentRoles)
	if err != nil {
		return nil, err
	}
	all, _ := s.gr.ListAllSkillStates()
	byUser := map[string][]model.StudentSkillState{}
	for _, st := range all {
		byUser[st.UserID] = append(byUser[st.UserID], st)
	}
	out := []model.CandidateView{}
	for _, u := range talents {
		levels := map[string]int{}
		for _, st := range byUser[u.ID] {
			levels[st.SkillID] = st.Level
		}
		mv := s.match(c, reqs, levels)
		if mv.Score == 0 {
			continue
		}
		cv := model.CandidateView{User: userBrief(u), Match: *mv}
		if g := s.growth.Goal(u.ID); g != nil {
			if r, ok := c.Roles[g.RoleID]; ok {
				cv.Goal = r.Name
			}
		}
		for _, r := range reqs {
			for _, st := range byUser[u.ID] {
				if st.SkillID == r.SkillID {
					cv.Skills = append(cv.Skills, skillStateView(c.Skills[st.SkillID], st))
				}
			}
		}
		cv.Skills = nonNil(cv.Skills)
		out = append(out, cv)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Match.Score > out[j].Match.Score })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
