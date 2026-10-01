package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/apperr"
	"github.com/cng1985/ai-learning-server/pkg/growth"
	"github.com/cng1985/ai-learning-server/pkg/llm"
)

// ProjectService 项目实践：Project → Task → Execution → Evidence。
type ProjectService struct {
	repo     *repository.ProjectRepo
	catalog  *CatalogService
	growth   *GrowthService
	settings *SettingsService
}

func NewProjectService(repo *repository.ProjectRepo, catalog *CatalogService, growth *GrowthService, settings *SettingsService) *ProjectService {
	return &ProjectService{repo: repo, catalog: catalog, growth: growth, settings: settings}
}

func (s *ProjectService) view(c *Catalog, p model.Project, levels map[string]int, best map[string]int) model.ProjectView {
	v := model.ProjectView{Project: p, SkillIDs: parseStrings(p.SkillIDs)}
	v.Project.Tasks = nil
	for _, t := range p.Tasks {
		reqs := parseRequirements(t.Requirements)
		tv := model.ProjectTaskView{ProjectTask: t, Requirements: reqs}
		if levels != nil {
			mv := matchView(growth.Fit(toEngineReqs(reqs), levels), c)
			tv.Match = &mv
		}
		if score, ok := best[t.ID]; ok {
			tv.Submitted, tv.BestScore = true, score
			if score >= 60 {
				v.Completed++
			}
		}
		v.Tasks = append(v.Tasks, tv)
	}
	v.Tasks = nonNil(v.Tasks)
	if len(p.Tasks) > 0 {
		v.Progress = v.Completed * 100 / len(p.Tasks)
	}
	return v
}

func (s *ProjectService) bestScores(userID string) map[string]int {
	best := map[string]int{}
	if userID == "" {
		return best
	}
	subs, _ := s.repo.ListSubmissions(userID, "", 0)
	for _, sub := range subs {
		if old, ok := best[sub.TaskID]; !ok || sub.Score > old {
			best[sub.TaskID] = sub.Score
		}
	}
	return best
}

// List 项目列表；userID 不为空时附带个人匹配度与完成进度。
func (s *ProjectService) List(userID, status string) ([]model.ProjectView, error) {
	projects, err := s.repo.List(status)
	if err != nil {
		return nil, err
	}
	c, _ := s.catalog.Load()
	var levels map[string]int
	if userID != "" {
		levels = s.growth.SkillLevels(userID)
	}
	best := s.bestScores(userID)
	out := make([]model.ProjectView, 0, len(projects))
	for _, p := range projects {
		p.Description = ""
		out = append(out, s.view(c, p, levels, best))
	}
	return out, nil
}

// ProjectDetail 项目详情与个人提交记录。
type ProjectDetail struct {
	model.ProjectView
	Submissions []model.SubmissionView `json:"submissions"`
}

func (s *ProjectService) Get(userID, id string) (*ProjectDetail, error) {
	p, err := s.repo.Find(id)
	if err != nil {
		return nil, apperr.NotFound("项目不存在")
	}
	c, _ := s.catalog.Load()
	d := &ProjectDetail{ProjectView: s.view(c, *p, s.growth.SkillLevels(userID), s.bestScores(userID))}
	subs, _ := s.repo.ListSubmissions(userID, id, 0)
	d.Submissions = s.submissionViews(subs)
	return d, nil
}

func (s *ProjectService) submissionViews(subs []model.TaskSubmission) []model.SubmissionView {
	out := make([]model.SubmissionView, 0, len(subs))
	projects := map[string]*model.Project{}
	for _, sub := range subs {
		p, ok := projects[sub.ProjectID]
		if !ok {
			p, _ = s.repo.Find(sub.ProjectID)
			projects[sub.ProjectID] = p
		}
		v := model.SubmissionView{TaskSubmission: sub}
		if p != nil {
			v.ProjectTitle = p.Title
			for _, t := range p.Tasks {
				if t.ID == sub.TaskID {
					v.TaskTitle = t.Title
				}
			}
		}
		out = append(out, v)
	}
	return out
}

func (s *ProjectService) Submissions(userID string, limit int) []model.SubmissionView {
	subs, _ := s.repo.ListSubmissions(userID, "", limit)
	return s.submissionViews(subs)
}

// Submit 提交任务成果，由 Review Agent 评审并生成技能证据。
func (s *ProjectService) Submit(ctx context.Context, userID, projectID, taskID string, req model.SubmissionRequest) (*model.SubmissionView, error) {
	req.Content = strings.TrimSpace(req.Content)
	if utf8.RuneCountInString(req.Content) < 20 {
		return nil, apperr.BadRequest("请至少用 20 个字描述你的实现方案与成果")
	}
	p, err := s.repo.Find(projectID)
	if err != nil {
		return nil, apperr.NotFound("项目不存在")
	}
	t, err := s.repo.FindTask(projectID, taskID)
	if err != nil {
		return nil, apperr.NotFound("任务不存在")
	}
	c, _ := s.catalog.Load()
	reqs := parseRequirements(t.Requirements)
	review := s.evaluate(ctx, c, p, t, reqs, req)
	now := time.Now().UnixMilli()
	sub := &model.TaskSubmission{
		ID: genID("sub"), UserID: userID, ProjectID: projectID, TaskID: taskID,
		Content: req.Content, RepoURL: strings.TrimSpace(req.RepoURL), Status: model.SubmissionReviewed,
		Score: review.Score, Feedback: review.Feedback, Evaluator: review.Evaluator, CreatedAt: now, ReviewedAt: now,
	}
	if err := s.repo.AddSubmission(sub); err != nil {
		return nil, err
	}
	skillIDs := make([]string, 0, len(reqs))
	for _, r := range reqs {
		skillIDs = append(skillIDs, r.SkillID)
	}
	ev := s.growth.AddEvidence(userID, skillIDs, model.EvidenceProject, sub.ID, p.Title+" · "+t.Title, review.Score, review.Summary)
	return &model.SubmissionView{TaskSubmission: *sub, TaskTitle: t.Title, ProjectTitle: p.Title, Evidence: ev}, nil
}

type reviewResult struct {
	Score     int
	Feedback  string
	Summary   string
	Evaluator string
}

// evaluate Review Agent：优先由大模型按评审标准打分；未配置大模型时使用可解释的规则评审。
func (s *ProjectService) evaluate(ctx context.Context, c *Catalog, p *model.Project, t *model.ProjectTask, reqs []model.SkillRequirement, sub model.SubmissionRequest) reviewResult {
	skills := make([]string, 0, len(reqs))
	for _, r := range reqs {
		skills = append(skills, fmt.Sprintf("%s L%d", c.SkillName(r.SkillID), r.Level))
	}
	if client := s.settings.LLMClient(); client != nil && client.Enabled() {
		prompt := fmt.Sprintf(`请作为资深技术评审，评审学习者提交的项目任务成果。
项目：%s
任务：%s
任务说明：%s
交付要求：%s
技能要求：%s
代码仓库：%s
提交内容：
<submission>
%s
</submission>

评分标准：方案完整性 30、技术正确性 30、工程质量（边界/异常/性能）25、表达与文档 15。
严格输出 JSON：{"score":0-100 的整数,"strengths":["优点"],"improvements":["改进建议"],"summary":"一句话评价"}`,
			p.Title, t.Title, t.Description, t.Deliverable, strings.Join(skills, "、"), sub.RepoURL, sub.Content)
		raw, err := client.Complete(ctx, []llm.Message{
			{Role: "system", Content: "你是 Review Agent，负责代码分析、项目评价与能力评估，只输出合法 JSON。"},
			{Role: "user", Content: prompt},
		})
		if err == nil {
			var out struct {
				Score        int      `json:"score"`
				Strengths    []string `json:"strengths"`
				Improvements []string `json:"improvements"`
				Summary      string   `json:"summary"`
			}
			if parseJSONResponse(raw, &out) == nil && out.Score > 0 {
				return reviewResult{
					Score: clampScore(out.Score), Summary: out.Summary, Evaluator: "ai",
					Feedback: formatFeedback(out.Summary, out.Strengths, out.Improvements),
				}
			}
		}
	}
	return ruleReview(c, t, reqs, sub)
}

var (
	reHeading = regexp.MustCompile(`(?m)^#{1,4}\s`)
	reList    = regexp.MustCompile(`(?m)^\s*([-*]|\d+\.)\s`)
	reCode    = regexp.MustCompile("```")
)

// ruleReview 规则评审：从篇幅、结构、技术覆盖与交付物四个维度给出可解释的评分。
func ruleReview(c *Catalog, t *model.ProjectTask, reqs []model.SkillRequirement, sub model.SubmissionRequest) reviewResult {
	content := sub.Content
	length := utf8.RuneCountInString(content)
	var strengths, improvements []string
	score := 35

	switch {
	case length >= 600:
		score += 20
		strengths = append(strengths, "方案描述充分")
	case length >= 250:
		score += 12
	default:
		score += 4
		improvements = append(improvements, "补充实现思路、关键设计决策与取舍")
	}

	structure := 0
	if reHeading.MatchString(content) {
		structure++
	}
	if len(reList.FindAllString(content, -1)) >= 3 {
		structure++
	}
	if reCode.MatchString(content) {
		structure++
		strengths = append(strengths, "包含关键代码或配置")
	} else {
		improvements = append(improvements, "附上核心代码片段、SQL 或配置示例")
	}
	score += structure * 5

	lower := strings.ToLower(content)
	covered := 0
	for _, r := range reqs {
		name := strings.ToLower(c.SkillName(r.SkillID))
		if hitsAny(lower, name, c.Skills[r.SkillID].Tags) {
			covered++
		}
	}
	if len(reqs) > 0 {
		score += covered * 15 / len(reqs)
		if covered < len(reqs) {
			improvements = append(improvements, "说明任务涉及的每项技能是如何应用的")
		} else {
			strengths = append(strengths, "覆盖了任务要求的全部技能点")
		}
	}

	quality := 0
	for _, kw := range []string{"测试", "test", "异常", "边界", "性能", "并发", "监控", "日志", "压测", "索引", "缓存", "回滚"} {
		if strings.Contains(lower, kw) {
			quality++
		}
	}
	if quality >= 3 {
		score += 10
		strengths = append(strengths, "考虑了测试、异常或性能等工程质量问题")
	} else {
		improvements = append(improvements, "补充测试用例、异常处理与性能验证")
	}
	if strings.HasPrefix(strings.TrimSpace(sub.RepoURL), "http") {
		score += 5
	}
	score = clampScore(score)
	summary := fmt.Sprintf("规则评审：%s完成度 %d 分（配置大模型后由 Review Agent 进行深度评审）", t.Title, score)
	return reviewResult{Score: score, Summary: summary, Evaluator: "rule", Feedback: formatFeedback(summary, strengths, improvements)}
}

func hitsAny(text, name string, tags []byte) bool {
	if name != "" && strings.Contains(text, name) {
		return true
	}
	for _, tag := range parseStrings(tags) {
		if tag != "" && strings.Contains(text, strings.ToLower(tag)) {
			return true
		}
	}
	return false
}

func clampScore(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func formatFeedback(summary string, strengths, improvements []string) string {
	var b strings.Builder
	b.WriteString(summary)
	if len(strengths) > 0 {
		b.WriteString("\n\n**亮点**\n")
		for _, it := range strengths {
			b.WriteString("- " + it + "\n")
		}
	}
	if len(improvements) > 0 {
		b.WriteString("\n**改进建议**\n")
		for _, it := range improvements {
			b.WriteString("- " + it + "\n")
		}
	}
	return strings.TrimSpace(b.String())
}

// ---------- 管理端 ----------

func (s *ProjectService) Save(id, authorID string, req model.ProjectRequest) (*model.Project, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, apperr.BadRequest("请填写项目名称")
	}
	p := &model.Project{}
	if id != "" {
		old, err := s.repo.Find(id)
		if err != nil {
			return nil, apperr.NotFound("项目不存在")
		}
		p = old
	} else {
		p.ID, p.AuthorID, p.CreatedAt = genID("proj"), authorID, time.Now().UnixMilli()
	}
	p.Title, p.Summary, p.Description, p.Domain = req.Title, req.Summary, req.Description, req.Domain
	p.Difficulty, p.Icon, p.Color = clampRange(req.Difficulty, 1, 5, 3), req.Icon, req.Color
	p.SkillIDs = toJSON(nonNil(req.SkillIDs))
	p.Status = req.Status
	if p.Status == "" {
		p.Status = "published"
	}
	oldTasks := p.Tasks
	p.Tasks = nil
	if err := s.repo.Save(p); err != nil {
		return nil, err
	}
	keep := map[string]bool{}
	for i, tr := range req.Tasks {
		if strings.TrimSpace(tr.Title) == "" {
			continue
		}
		tid := tr.ID
		if tid == "" {
			tid = genID("task")
		}
		keep[tid] = true
		task := model.ProjectTask{
			ID: tid, ProjectID: p.ID, Title: tr.Title, Description: tr.Description, Deliverable: tr.Deliverable,
			Requirements: toJSON(nonNil(tr.Requirements)), EstimatedHours: clampRange(tr.EstimatedHours, 1, 200, 4), Sort: i,
		}
		if err := s.repo.Save(&task); err != nil {
			return nil, err
		}
	}
	for _, t := range oldTasks {
		if !keep[t.ID] {
			_ = s.repo.DeleteTask(p.ID, t.ID)
		}
	}
	return s.repo.Find(p.ID)
}

func (s *ProjectService) Delete(id string) error { return s.repo.Delete(id) }

func (s *ProjectService) AllSubmissions(projectID string, limit int) []model.SubmissionView {
	subs, _ := s.repo.ListSubmissions("", projectID, limit)
	return s.submissionViews(subs)
}
