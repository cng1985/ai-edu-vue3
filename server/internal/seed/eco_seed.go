package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/internal/service"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// EcoDeps 生态种子所需依赖。
type EcoDeps struct {
	DB        *gorm.DB
	Users     *repository.UserRepo
	Catalog   *repository.CatalogRepo
	Projects  *repository.ProjectRepo
	Market    *repository.MarketRepo
	Community *repository.CommunityRepo
	Growth    *service.GrowthService
	Practice  *service.ProjectService
	MarketSvc *service.MarketService
}

func js(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}

// SeedEcosystem 首次启动时写入成长生态数据；已存在职业数据时跳过，兼容旧数据库升级。
func SeedEcosystem(db *gorm.DB, users *repository.UserRepo, catalog *repository.CatalogRepo, projects *repository.ProjectRepo,
	market *repository.MarketRepo, community *repository.CommunityRepo, growth *service.GrowthService,
	practice *service.ProjectService, marketSvc *service.MarketService) error {
	if catalog.CountCareers() > 0 {
		return nil
	}
	fmt.Println("🌱 正在初始化 AI 成长生态数据（职业体系 / 知识图谱 / 项目 / 任务市场 / 社区）...")
	d := EcoDeps{DB: db, Users: users, Catalog: catalog, Projects: projects, Market: market, Community: community,
		Growth: growth, Practice: practice, MarketSvc: marketSvc}
	steps := []func(EcoDeps) error{seedCatalog, seedEcoUsers, seedProjects, seedOpportunities, seedGrowthHistory, seedCommunity}
	for _, step := range steps {
		if err := step(d); err != nil {
			return err
		}
	}
	fmt.Println("✅ 成长生态数据初始化完成")
	return nil
}

func seedCatalog(d EcoDeps) error {
	now := time.Now().UnixMilli()
	return d.DB.Transaction(func(tx *gorm.DB) error {
		for _, s := range skillDefs {
			if err := tx.Create(&model.Skill{ID: s.ID, Name: s.Name, Category: s.Category, Description: s.Desc, Tags: js(s.Tags), CreatedAt: now}).Error; err != nil {
				return err
			}
		}
		n := 0
		for _, k := range kpDefs {
			if err := tx.Create(&model.KnowledgePoint{
				ID: k.ID, Name: k.Name, Domain: k.Domain, Summary: k.Summary, Content: k.Content,
				Difficulty: k.Difficulty, EstimatedMinutes: k.Minutes, Questions: js(k.Questions), Tags: js([]string{}),
				CreatedAt: now, UpdatedAt: now,
			}).Error; err != nil {
				return err
			}
			for _, sid := range k.Skills {
				n++
				if err := tx.Create(&model.SkillKnowledge{ID: fmt.Sprintf("sk_%s_%s", sid, k.ID), SkillID: sid, KnowledgeID: k.ID, Weight: 1}).Error; err != nil {
					return err
				}
			}
		}
		for i, r := range relationDefs {
			if err := tx.Create(&model.KnowledgeRelation{ID: fmt.Sprintf("rel_%03d", i+1), FromID: r.From, ToID: r.To, Type: r.Type}).Error; err != nil {
				return err
			}
		}
		for ci, c := range careerDefs {
			if err := tx.Create(&model.Career{ID: c.ID, Name: c.Name, Category: c.Category, Description: c.Desc, Icon: c.Icon,
				Color: c.Color, Demand: c.Demand, SalaryRange: c.Salary, Sort: ci, CreatedAt: now}).Error; err != nil {
				return err
			}
			for ri, r := range c.Roles {
				if err := tx.Create(&model.JobRole{ID: r.ID, CareerID: c.ID, Name: r.Name, Level: r.Level, Description: r.Desc, Sort: ri}).Error; err != nil {
					return err
				}
				for capIdx, cp := range r.Caps {
					if err := tx.Create(&model.RoleCapability{ID: cp.ID, RoleID: r.ID, Name: cp.Name, Description: cp.Desc, Weight: cp.Weight, Sort: capIdx}).Error; err != nil {
						return err
					}
					for _, sk := range cp.Skills {
						sid, lvl := sk[0].(string), sk[1].(int)
						if err := tx.Create(&model.CapabilitySkill{ID: fmt.Sprintf("cs_%s_%s", cp.ID, sid), CapabilityID: cp.ID, SkillID: sid, RequiredLevel: lvl, Weight: 1}).Error; err != nil {
							return err
						}
					}
				}
			}
		}
		for courseID, chapters := range chapterKnowledgeDefs {
			for chapterID, kids := range chapters {
				for _, kid := range kids {
					if err := tx.Create(&model.ChapterKnowledge{ID: fmt.Sprintf("ck_%s_%s_%s", courseID, chapterID, kid), CourseID: courseID, ChapterID: chapterID, KnowledgeID: kid}).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

type ecoUser struct {
	ID, Username, Password, Nickname, Role, Avatar, Color string
	DaysAgo                                               int
}

var ecoUsers = []ecoUser{
	{"company_001", "company", "company123", "星云科技", "enterprise", "星", "#0ea5e9", 60},
	{"creator_001", "expert", "expert123", "李架构", "creator", "李", "#8b5cf6", 120},
	{"talent_bob", "bob", "bob123", "周航", "learner", "周", "#f97316", 90},
	{"talent_alice", "alice", "alice123", "林晓", "learner", "林", "#14b8a6", 75},
	{"talent_chen", "chen", "chen123", "陈一", "learner", "陈", "#ef4444", 40},
}

func seedEcoUsers(d EcoDeps) error {
	for _, u := range ecoUsers {
		if _, err := d.Users.FindByUsername(u.Username); err == nil {
			continue
		}
		hash, _ := authutil.HashPassword(u.Password)
		joined := time.Now().AddDate(0, 0, -u.DaysAgo).UnixMilli()
		if err := d.Users.Create(&model.User{ID: u.ID, Username: u.Username, Nickname: u.Nickname, PasswordHash: hash,
			Role: u.Role, Status: model.UserStatusActive, Avatar: u.Avatar, AvatarColor: u.Color, JoinedAt: joined}); err != nil {
			return err
		}
	}
	return nil
}

type taskDef struct {
	ID, Title, Desc, Deliverable string
	Hours                        int
	Reqs                         []model.SkillRequirement
}

func r(skill string, level int) model.SkillRequirement { return model.SkillRequirement{SkillID: skill, Level: level} }

var projectDefs = []struct {
	ID, Title, Summary, Desc, Domain, Icon, Color string
	Difficulty                                    int
	Skills                                        []string
	Tasks                                         []taskDef
}{
	{"ecommerce-order", "电商订单系统", "从表设计到高并发库存扣减，完整实践一个电商核心交易链路。",
		"## 项目背景\n\n某电商平台需要重构订单中心，支撑大促期间 **5000 QPS** 的下单峰值。\n\n## 你将实践\n\n- 订单表设计与索引\n- 下单接口与幂等\n- 高并发库存扣减（防超卖）\n- 多级缓存设计\n- 压测与性能优化\n\n每个任务完成后提交方案说明与关键代码，由 Review Agent 评审并生成技能证据。",
		"后端开发", "code", "#6b5cff", 4, []string{"java-core", "spring-boot", "mysql", "redis", "java-concurrency", "performance"},
		[]taskDef{
			{"order-schema", "设计订单表", "设计订单主表、订单明细表与支付流水表，说明主键、索引与快照字段的设计理由。", "建表 SQL + 索引设计说明", 4, []model.SkillRequirement{r("mysql", 2)}},
			{"order-api", "开发订单接口", "基于 Spring Boot 实现创建订单、查询订单、取消订单接口，保证下单幂等。", "接口设计文档 + 核心代码 + 幂等方案", 8, []model.SkillRequirement{r("spring-boot", 2), r("java-core", 2)}},
			{"stock-deduct", "库存扣减", "在高并发下实现库存扣减且不超卖，对比数据库乐观锁、Redis Lua、分布式锁三种方案。", "方案对比 + 实现代码 + 并发测试结果", 10, []model.SkillRequirement{r("java-concurrency", 3), r("redis", 2), r("mysql", 3)}},
			{"order-cache", "缓存设计", "为订单查询与商品详情设计缓存，处理穿透、击穿、雪崩与一致性问题。", "缓存架构图 + 一致性方案 + 代码", 6, []model.SkillRequirement{r("redis", 3)}},
			{"order-perf", "性能优化", "对下单链路进行压测，定位瓶颈并优化，给出优化前后对比数据。", "压测报告 + 优化措施 + 前后指标", 8, []model.SkillRequirement{r("performance", 3), r("mysql", 3)}},
		}},
	{"rag-kb", "企业知识库问答系统", "构建低幻觉、可溯源的企业知识库问答，掌握 RAG 全链路。",
		"## 项目背景\n\n企业内部文档分散，员工查找制度与技术资料效率低。你需要构建一个 **可溯源、低幻觉** 的知识库问答系统。\n\n## 你将实践\n\n- 文档切分与向量化\n- 混合检索与重排序\n- 提示词与回答生成\n- Agent 工具调用\n- 评估与持续优化",
		"AI 应用", "sparkles", "#0fb981", 4, []string{"rag", "prompt-eng", "ai-app-dev", "python"},
		[]taskDef{
			{"rag-ingest", "文档切分与向量化", "实现 Markdown/PDF 文档的结构感知切分与向量化入库，说明 chunk 大小与重叠策略。", "切分策略说明 + 入库代码", 6, []model.SkillRequirement{r("rag", 2), r("python", 2)}},
			{"rag-retrieval", "混合检索与重排序", "实现向量 + BM25 混合检索与 Rerank，并用测试问题对比召回效果。", "检索实现 + 召回率对比", 8, []model.SkillRequirement{r("rag", 3)}},
			{"rag-prompt", "提示词与回答生成", "设计带引用溯源的回答提示词，处理知识库无答案的情况。", "提示词模板 + 回答样例", 4, []model.SkillRequirement{r("prompt-eng", 3)}},
			{"rag-agent", "Agent 工具调用", "让助手可调用\"查询工单\"\"创建工单\"工具，并防止死循环与越权。", "工具定义 + 调用流程 + 防护措施", 8, []model.SkillRequirement{r("ai-app-dev", 3)}},
			{"rag-eval", "评估与优化", "构建 30 条黄金问答集，评估准确率与忠实度并迭代优化。", "评估集 + 指标报告 + 优化记录", 6, []model.SkillRequirement{r("ai-app-dev", 3), r("rag", 3)}},
		}},
	{"learning-dashboard", "学习数据看板前端", "用 Vue 构建一个学习数据可视化看板，练习组件化与异步数据处理。",
		"## 项目背景\n\n为学习平台开发一个数据看板，展示学习时长、掌握度与成长趋势。\n\n## 你将实践\n\n- 响应式页面布局\n- 组件拆分与复用\n- 异步数据加载与错误处理\n- 状态管理与渲染性能",
		"前端开发", "chart", "#2a8cf4", 2, []string{"html-css", "javascript", "vue"},
		[]taskDef{
			{"dash-layout", "页面布局", "使用语义化 HTML 与 Flex/Grid 完成看板响应式布局。", "页面代码 + 多端截图说明", 3, []model.SkillRequirement{r("html-css", 2)}},
			{"dash-components", "组件设计", "拆分统计卡片、趋势图、排行列表等组件，定义清晰的 props/emits。", "组件结构说明 + 代码", 5, []model.SkillRequirement{r("vue", 2)}},
			{"dash-async", "异步数据加载", "并发加载多个接口数据，处理加载态、错误态与重试。", "数据层代码 + 异常处理说明", 4, []model.SkillRequirement{r("javascript", 3)}},
			{"dash-state", "状态管理与性能", "用 Pinia 管理筛选状态，优化大列表渲染性能。", "状态设计 + 性能优化前后对比", 5, []model.SkillRequirement{r("vue", 3)}},
		}},
}

func seedProjects(d EcoDeps) error {
	now := time.Now().UnixMilli()
	for pi, p := range projectDefs {
		if err := d.Projects.Save(&model.Project{ID: p.ID, Title: p.Title, Summary: p.Summary, Description: p.Desc, Domain: p.Domain,
			Difficulty: p.Difficulty, Icon: p.Icon, Color: p.Color, AuthorID: "creator_001", SkillIDs: js(p.Skills),
			Status: "published", CreatedAt: now + int64(pi)}); err != nil {
			return err
		}
		for ti, t := range p.Tasks {
			if err := d.Projects.Save(&model.ProjectTask{ID: t.ID, ProjectID: p.ID, Title: t.Title, Description: t.Desc,
				Deliverable: t.Deliverable, Requirements: js(t.Reqs), EstimatedHours: t.Hours, Sort: ti}); err != nil {
				return err
			}
		}
	}
	return nil
}

var opportunityDefs = []struct {
	ID, Title, Desc, Mode string
	Budget                float64
	Days                  int
	Reqs                  []model.SkillRequirement
	DaysAgo               int
}{
	{"opp-stock-query", "开发库存查询模块", "为仓储系统开发库存查询服务：支持按 SKU、仓库、批次多维查询，P99 < 100ms，需要设计缓存与索引方案。", "remote", 6000, 14,
		[]model.SkillRequirement{r("spring-boot", 3), r("mysql", 3), r("redis", 2)}, 2},
	{"opp-rag-assistant", "企业内部知识库 RAG 助手", "基于公司 2000+ 份制度与技术文档搭建问答助手，要求回答可溯源，并接入企业 IM。", "remote", 12000, 21,
		[]model.SkillRequirement{r("rag", 3), r("prompt-eng", 2), r("ai-app-dev", 3)}, 3},
	{"opp-admin-refactor", "管理后台前端页面重构", "将 jQuery 老后台迁移到 Vue3，涉及 12 个页面与通用表格、表单组件封装。", "remote", 5000, 10,
		[]model.SkillRequirement{r("vue", 3), r("javascript", 3), r("html-css", 2)}, 5},
	{"opp-order-perf", "订单服务性能优化", "大促前对订单服务进行全链路压测与优化，目标下单接口 QPS 从 800 提升至 3000。", "onsite", 15000, 20,
		[]model.SkillRequirement{r("performance", 4), r("mysql", 4), r("java-concurrency", 3)}, 1},
	{"opp-prompt-tuning", "客服话术 Prompt 优化", "优化电商客服机器人的回复提示词，降低答非所问率，提供 Few-shot 示例与评估结果。", "remote", 1500, 5,
		[]model.SkillRequirement{r("prompt-eng", 2)}, 4},
}

func seedOpportunities(d EcoDeps) error {
	for _, o := range opportunityDefs {
		at := time.Now().AddDate(0, 0, -o.DaysAgo).UnixMilli()
		if err := d.Market.Save(&model.Opportunity{ID: o.ID, Title: o.Title, Company: "星云科技", PublisherID: "company_001",
			Description: o.Desc, Requirements: js(o.Reqs), Budget: o.Budget, DurationDays: o.Days, Mode: o.Mode,
			Status: model.OpportunityOpen, CreatedAt: at, UpdatedAt: at}); err != nil {
			return err
		}
	}
	return nil
}

type practice struct {
	K              string
	Study, Correct int
	Total          int
}

// seedGrowthHistory 为演示人才写入真实的学习、项目与任务数据（全部经由成长引擎计算）。
func seedGrowthHistory(d EcoDeps) error {
	g := d.Growth
	learn := func(uid string, items []practice) {
		for _, it := range items {
			for i := 0; i < it.Study; i++ {
				_, _ = g.RecordEvent(uid, model.LearningEventRequest{KnowledgeID: it.K, Type: model.EventStudy, Source: "seed"})
			}
			if it.Total > 0 {
				_, _ = g.RecordEvent(uid, model.LearningEventRequest{KnowledgeID: it.K, Type: model.EventPractice, Correct: it.Correct, Total: it.Total, Source: "seed"})
			}
		}
	}
	mastered := func(ids ...string) []practice {
		out := make([]practice, 0, len(ids))
		for _, id := range ids {
			out = append(out, practice{id, 2, 14, 15})
		}
		return out
	}

	// 周航：资深 Java 方向
	_, _ = g.SetGoal("talent_bob", model.GoalRequest{RoleID: "senior-java", WeeklyHours: 12, TargetWeeks: 20, Motivation: "希望半年内成长为能独立负责核心系统的高级工程师"})
	learn("talent_bob", mastered("java-oop", "java-collections", "thread", "thread-safety", "jmm", "volatile", "synchronized", "cas",
		"spring-ioc", "spring-rest", "spring-transaction", "sql-basics", "mysql-index", "mysql-transaction", "slow-query",
		"redis-datatypes", "cache-patterns", "distributed-lock", "aqs", "thread-pool"))
	learn("talent_bob", []practice{{"load-testing", 1, 6, 8}, {"message-queue", 1, 4, 6}, {"service-split", 1, 3, 4}})

	// 林晓：AI 应用方向
	_, _ = g.SetGoal("talent_alice", model.GoalRequest{RoleID: "ai-app-engineer", WeeklyHours: 15, TargetWeeks: 16, Motivation: "从数据分析转型 AI 应用开发"})
	learn("talent_alice", mastered("llm-basics", "prompt-structure", "few-shot", "embedding", "chunking", "hybrid-retrieval", "python-basics", "function-calling"))
	learn("talent_alice", []practice{{"agent-react", 1, 5, 8}, {"llm-eval", 1, 3, 6}, {"etl", 1, 4, 5}})

	// 陈一：前端方向
	_, _ = g.SetGoal("talent_chen", model.GoalRequest{RoleID: "web-frontend", WeeklyHours: 10, TargetWeeks: 12})
	learn("talent_chen", mastered("html-semantic", "flex-layout", "js-async", "vue-reactivity", "vue-component"))

	// 演示学员：刚开始 AI 应用工程师之路
	_, _ = g.SetGoal("learner_demo", model.GoalRequest{RoleID: "ai-app-engineer", WeeklyHours: 10, TargetWeeks: 16, Motivation: "想用 AI 提升工作效率，并接到第一个 AI 项目"})
	learn("learner_demo", []practice{{"llm-basics", 2, 9, 10}, {"prompt-structure", 1, 7, 10}, {"embedding", 1, 0, 0}})

	ctx := context.Background()
	submit := func(uid, project, task, content string) {
		_, _ = d.Practice.Submit(ctx, uid, project, task, model.SubmissionRequest{Content: content, RepoURL: "https://github.com/demo/" + project})
	}
	submit("talent_bob", "ecommerce-order", "order-schema", bobSchema)
	submit("talent_bob", "ecommerce-order", "stock-deduct", bobStock)
	submit("talent_bob", "ecommerce-order", "order-cache", bobCache)
	submit("talent_alice", "rag-kb", "rag-ingest", aliceIngest)
	submit("talent_alice", "rag-kb", "rag-retrieval", aliceRetrieval)
	submit("talent_chen", "learning-dashboard", "dash-components", chenComponents)

	// 历史已完成的真实任务：形成收入记录与任务表现证据
	past := time.Now().AddDate(0, 0, -20).UnixMilli()
	history := []struct {
		oppID, title, uid, review string
		budget                    float64
		rating                    int
		reqs                      []model.SkillRequirement
	}{
		{"opp-hist-recon", "支付对账服务开发", "talent_bob", "交付质量高，主动补充了对账差异告警。", 8000, 5, []model.SkillRequirement{r("spring-boot", 3), r("mysql", 3)}},
		{"opp-hist-faq", "售后 FAQ 智能问答原型", "talent_alice", "两周内完成原型并给出评估报告，效果超出预期。", 6000, 4, []model.SkillRequirement{r("rag", 2), r("prompt-eng", 2)}},
	}
	for _, h := range history {
		if err := d.Market.Save(&model.Opportunity{ID: h.oppID, Title: h.title, Company: "星云科技", PublisherID: "company_001",
			Description: h.title + "（历史任务）", Requirements: js(h.reqs), Budget: h.budget, DurationDays: 14, Mode: "remote",
			Status: model.OpportunityOpen, CreatedAt: past, UpdatedAt: past}); err != nil {
			return err
		}
		app, err := d.MarketSvc.Apply(h.uid, "learner", h.oppID, model.ApplyRequest{Message: "有相关项目经验，可以按期交付。"})
		if err != nil {
			return err
		}
		if _, err := d.MarketSvc.Decide("company_001", app.ID, false, model.ApplicationDecision{Action: "accept"}); err != nil {
			return err
		}
		if _, err := d.MarketSvc.Decide("company_001", app.ID, false, model.ApplicationDecision{Action: "complete", Rating: h.rating, Review: h.review, Payout: h.budget}); err != nil {
			return err
		}
	}
	_, _ = d.MarketSvc.Apply("talent_bob", "learner", "opp-stock-query", model.ApplyRequest{Message: "做过订单与库存相关项目，熟悉缓存与索引优化。"})
	_, _ = d.MarketSvc.Apply("talent_alice", "learner", "opp-rag-assistant", model.ApplyRequest{Message: "完成过 RAG 知识库项目与 FAQ 问答原型。"})
	return nil
}

const bobSchema = "## 订单表设计\n\n- `t_order`：id BIGINT 趋势递增主键、order_no 唯一索引、user_id + created_at 联合索引用于用户订单列表\n- `t_order_item`：冗余商品名称与下单价格作为快照，金额字段使用 DECIMAL(12,2)\n- `t_payment`：支付流水，order_no 索引\n\n```sql\nCREATE TABLE t_order (\n  id BIGINT PRIMARY KEY,\n  order_no VARCHAR(32) NOT NULL UNIQUE,\n  user_id BIGINT NOT NULL,\n  amount DECIMAL(12,2) NOT NULL,\n  status TINYINT NOT NULL,\n  created_at DATETIME NOT NULL,\n  KEY idx_user_time (user_id, created_at)\n);\n```\n\n考虑到大促订单量，按 user_id 分库分表预留扩展，状态变更记录日志便于排查；查询走覆盖索引避免回表，并通过压测验证列表接口性能。"

const bobStock = "## 库存扣减方案对比\n\n1. 数据库乐观锁：`UPDATE stock SET num = num - 1 WHERE sku = ? AND num > 0`，实现简单但热点行竞争激烈\n2. Redis Lua 原子扣减：预热库存到 Redis，Lua 脚本判断并扣减，性能最高，通过 MQ 异步落库保证最终一致\n3. Redisson 分布式锁：通用但吞吐较低\n\n```lua\nlocal n = tonumber(redis.call('GET', KEYS[1]))\nif n and n >= tonumber(ARGV[1]) then\n  return redis.call('DECRBY', KEYS[1], ARGV[1])\nend\nreturn -1\n```\n\n最终选择 Redis Lua + 事务消息落库，消费端按订单号幂等。用 JMeter 并发 2000 线程压测 10 万次无超卖，并发场景下对比了三种方案的 QPS 与异常率，增加了库存对账任务与监控告警。"

const bobCache = "## 缓存设计\n\n- 商品详情：本地缓存 Caffeine + Redis 两级缓存，Cache Aside 模式，先更新数据库再删除缓存，延迟双删兜底\n- 穿透：缓存空值 + 布隆过滤器\n- 击穿：热点 key 互斥锁重建\n- 雪崩：过期时间增加随机值\n\n```java\nString v = redis.get(key);\nif (v == null) { v = lockAndLoad(key); }\n```\n\n通过压测验证缓存命中率 96%，接口 P99 从 180ms 降到 35ms，并补充了缓存异常的降级与日志监控。"

const aliceIngest = "## 切分与向量化\n\n- 采用结构感知切分：按 Markdown 标题层级切分，单块 400-600 字，相邻块重叠 50 字\n- PDF 先转 Markdown 保留标题结构\n- Embedding 使用 bge-m3，批量向量化入库，记录模型版本，模型变更时全量重建索引\n\n```python\nchunks = split_by_heading(doc, max_len=600, overlap=50)\nvectors = embed([c.text for c in chunks])\n```\n\n对比固定长度切分，结构感知切分在 50 条测试问题上的召回率从 71% 提升到 84%，并补充了异常文档的日志与重试。"

const aliceRetrieval = "## 混合检索\n\n- 召回：向量检索 Top-20 + BM25 Top-20，RRF 融合\n- 精排：bge-reranker 交叉编码器取 Top-3\n- 多轮对话加入查询重写\n\n```python\ncands = rrf(vector_search(q, 20), bm25(q, 20))\ntop = rerank(q, cands)[:3]\n```\n\n测试集召回率 84% → 92%，专有名词类问题提升最明显；同时评估了 rerank 的延迟开销，缓存热门问题结果。"

const chenComponents = "## 组件拆分\n\n- `StatCard`：props 接收 label/value/trend，纯展示组件\n- `TrendChart`：接收数据数组，emit 选择区间事件\n- `RankList`：使用插槽自定义行内容\n- 数据获取逻辑抽到 `useDashboardData` composable\n\n```vue\n<StatCard label=\"学习时长\" :value=\"total\" @click=\"open\" />\n```\n\n遵循单向数据流，并为组件编写了基础测试。"

func seedCommunity(d EcoDeps) error {
	now := time.Now()
	at := func(hoursAgo int) int64 { return now.Add(-time.Duration(hoursAgo) * time.Hour).UnixMilli() }
	posts := []model.CommunityPost{
		{ID: "post-volatile", AuthorID: "learner_demo", Type: model.PostQuestion, Title: "volatile 能保证原子性吗？count++ 为什么还是有问题？",
			Content: "看了 JMM 的内容，知道 volatile 保证可见性，但用 volatile 修饰的 count 在多线程 count++ 后结果还是不对，是哪里理解错了？",
			Tags: js([]string{"Java并发", "volatile"}), KnowledgeIDs: js([]string{"volatile", "cas", "jmm"}), Likes: 12, Views: 186, CreatedAt: at(30)},
		{ID: "post-rag-recall", AuthorID: "talent_alice", Type: model.PostQuestion, Title: "RAG 召回率低，应该先优化切分还是检索？",
			Content: "知识库问答经常找不到相关段落，测试集召回率只有 70% 左右。时间有限，优先优化切分策略还是检索策略？",
			Tags: js([]string{"RAG", "检索"}), KnowledgeIDs: js([]string{"chunking", "hybrid-retrieval"}), Likes: 21, Views: 342, CreatedAt: at(52)},
		{ID: "post-evidence", AuthorID: "creator_001", Type: model.PostArticle, Title: "从知识到技能：如何用项目证据证明你的能力",
			Content: "很多同学刷了大量题目，却在面试和真实任务中表现一般。原因是：**知道 ≠ 会做**。\n\n## 三层证据\n\n1. 知识状态：掌握度与置信度\n2. 项目证据：可评审的方案与代码\n3. 任务表现：真实业务中的交付与评价\n\n平台的技能等级正是按这三层证据计算的：只有知识最多 L2，有了项目证据才能到 L3，真实任务表现优秀才能达到 L4、L5。\n\n建议每学完一个技能的核心知识，就立刻选一个项目任务练手。",
			Tags: js([]string{"成长方法", "技能评价"}), KnowledgeIDs: js([]string{}), Likes: 48, Views: 913, CreatedAt: at(80)},
		{ID: "post-fe-to-ai", AuthorID: "talent_chen", Type: model.PostDiscussion, Title: "前端转 AI 应用开发，需要补哪些能力？",
			Content: "目前 Vue 和 JS 比较熟，想往 AI 应用工程师方向发展，大家觉得应该先补 Python、RAG 还是 Prompt？",
			Tags: js([]string{"转型", "AI应用"}), KnowledgeIDs: js([]string{"llm-basics", "prompt-structure"}), Likes: 9, Views: 128, CreatedAt: at(10)},
		{ID: "post-first-income", AuthorID: "talent_bob", Type: model.PostShare, Title: "我的第一笔任务收入：支付对账服务开发复盘",
			Content: "通过任务市场接到星云科技的对账服务开发，两周交付拿到 8000 元。复盘几点：\n\n- 先把需求中的边界条件列清楚（跨日、退款、重复回调）\n- 对账差异必须可追踪，我加了差异告警\n- 项目实践里做过的幂等与索引设计直接派上用场\n\n学习 → 实践 → 产出 → 收益，这个飞轮真的转起来了。",
			Tags: js([]string{"任务复盘", "收入"}), KnowledgeIDs: js([]string{"mysql-index", "spring-rest"}), Likes: 35, Views: 520, CreatedAt: at(120)},
	}
	answers := []model.CommunityAnswer{
		{ID: "ans-v1", PostID: "post-volatile", AuthorID: "talent_bob", Content: "volatile 只保证可见性和有序性，不保证原子性。count++ 实际是读取、加一、写回三步，两个线程可能同时读到相同的旧值。解决方案：用 AtomicInteger（基于 CAS）或者 synchronized。", Likes: 15, Accepted: true, CreatedAt: at(29)},
		{ID: "ans-v2", PostID: "post-volatile", AuthorID: "creator_001", Content: "补充一点：判断是否需要原子性，看操作是否依赖当前值。像 `running = false` 这种单纯赋值用 volatile 就够了；依赖旧值计算新值的复合操作就需要 CAS 或锁。", Likes: 9, CreatedAt: at(27)},
		{ID: "ans-r1", PostID: "post-rag-recall", AuthorID: "creator_001", Content: "建议先看切分。召回问题 60% 以上来自切分不合理：块太大噪声多、块太小上下文缺失。先改成按标题结构切分并加少量重叠，通常能提升 10 个点以上；之后再上混合检索解决专有名词问题。", Likes: 18, Accepted: true, CreatedAt: at(50)},
		{ID: "ans-r2", PostID: "post-rag-recall", AuthorID: "talent_bob", Content: "可以先抽 20 个失败用例看原因，如果是专有名词、编号类查询找不到，就是缺 BM25 关键词检索；如果是答案被截断，就是切分问题。", Likes: 7, CreatedAt: at(46)},
		{ID: "ans-f1", PostID: "post-fe-to-ai", AuthorID: "talent_alice", Content: "我的经验是先 Prompt（见效最快），再 RAG，Python 够用即可。前端同学做 AI 应用有优势：交互和流式输出体验是很多 AI 产品的短板。", Likes: 6, CreatedAt: at(8)},
	}
	for i := range posts {
		posts[i].Status, posts[i].UpdatedAt = "published", posts[i].CreatedAt
		for _, a := range answers {
			if a.PostID == posts[i].ID {
				posts[i].AnswerCount++
			}
		}
		if err := d.Community.SavePost(&posts[i]); err != nil {
			return err
		}
	}
	for i := range answers {
		if err := d.Community.SaveAnswer(&answers[i]); err != nil {
			return err
		}
	}
	resources := []model.KnowledgeResource{
		{ID: "res-review-prompt", Type: "prompt", Title: "代码评审 Prompt 模板", Summary: "让大模型按正确性、可读性、性能、安全四个维度评审代码。",
			Content: "```\n你是资深代码评审专家。请从以下维度评审 <code> 中的代码：\n1. 正确性：边界条件、空值、并发\n2. 可读性：命名、结构、注释\n3. 性能：复杂度、IO、缓存\n4. 安全：注入、越权、敏感信息\n输出格式：问题列表（严重程度/位置/建议）\n<code>{{code}}</code>\n```",
			AuthorID: "creator_001", KnowledgeIDs: js([]string{"prompt-structure"}), SkillIDs: js([]string{"prompt-eng"}), Tags: js([]string{"Prompt", "代码评审"}), Likes: 32, Views: 410},
		{ID: "res-slow-sql-sop", Type: "sop", Title: "线上慢查询排查 SOP", Summary: "从告警到根因定位的 6 步标准流程。",
			Content: "1. 确认告警时间窗口与影响接口\n2. 从慢查询日志定位 Top SQL\n3. `EXPLAIN` 分析执行计划（type/key/rows/Extra）\n4. 检查索引失效：函数、隐式转换、前置模糊匹配\n5. 改写 SQL 或补充索引，在预发环境验证\n6. 回归压测并记录复盘",
			AuthorID: "creator_001", KnowledgeIDs: js([]string{"slow-query", "mysql-index"}), SkillIDs: js([]string{"mysql"}), Tags: js([]string{"MySQL", "SOP"}), Likes: 27, Views: 356},
		{ID: "res-redis-lock", Type: "code", Title: "Redis 分布式锁安全释放 Lua 脚本", Summary: "compare-and-delete，防止误删他人持有的锁。",
			Content: "```lua\nif redis.call('GET', KEYS[1]) == ARGV[1] then\n  return redis.call('DEL', KEYS[1])\nend\nreturn 0\n```",
			AuthorID: "talent_bob", KnowledgeIDs: js([]string{"distributed-lock"}), SkillIDs: js([]string{"redis"}), Tags: js([]string{"Redis", "分布式锁"}), Likes: 19, Views: 233},
		{ID: "res-oversell-case", Type: "case", Title: "电商大促库存超卖事故复盘", Summary: "一次因缓存与数据库扣减不一致导致的超卖事故及改进方案。",
			Content: "## 现象\n大促开始 3 分钟内某爆款超卖 214 件。\n\n## 根因\n先扣 Redis 再异步扣库，MQ 积压期间 Redis 预热脚本被重复执行，库存被重置。\n\n## 改进\n- 预热脚本加幂等标记\n- 库存扣减改为 Lua 原子扣减 + 事务消息\n- 增加实时对账告警",
			AuthorID: "creator_001", KnowledgeIDs: js([]string{"distributed-lock", "cache-patterns", "distributed-transaction"}), SkillIDs: js([]string{"redis", "distributed"}), Tags: js([]string{"事故复盘", "高并发"}), Likes: 41, Views: 688},
		{ID: "res-aqs-video", Type: "video", Title: "图解 AQS 原理（20 分钟）", Summary: "用动画讲清 state、CLH 队列与独占/共享模式。", URL: "https://www.bilibili.com/",
			AuthorID: "creator_001", KnowledgeIDs: js([]string{"aqs", "cas"}), SkillIDs: js([]string{"java-concurrency"}), Tags: js([]string{"视频", "并发"}), Likes: 56, Views: 1203},
		{ID: "res-rag-metrics", Type: "article", Title: "RAG 系统评估指标速查", Summary: "召回率、MRR、忠实度、答案相关性的定义与计算方法。",
			Content: "| 指标 | 衡量什么 |\n| --- | --- |\n| Recall@K | 正确文档是否被召回 |\n| MRR | 正确文档排得是否靠前 |\n| 忠实度 | 回答是否基于检索内容 |\n| 答案相关性 | 是否回答了问题 |",
			AuthorID: "talent_alice", KnowledgeIDs: js([]string{"llm-eval", "hybrid-retrieval"}), SkillIDs: js([]string{"rag", "ai-app-dev"}), Tags: js([]string{"RAG", "评估"}), Likes: 23, Views: 301},
		{ID: "res-course-prompt", Type: "course", Title: "提示词工程入门（课程）", Summary: "平台精品课程：从结构化提示词到 Few-shot 与注入防御。", URL: "#/courses/prompt-engineering",
			AuthorID: "creator_001", KnowledgeIDs: js([]string{"prompt-structure", "few-shot"}), SkillIDs: js([]string{"prompt-eng"}), Tags: js([]string{"课程"}), Likes: 64, Views: 1520},
		{ID: "res-project-order", Type: "project", Title: "电商订单系统（实践项目）", Summary: "5 个任务覆盖表设计、接口、库存扣减、缓存与性能优化。", URL: "#/projects/ecommerce-order",
			AuthorID: "creator_001", KnowledgeIDs: js([]string{"sql-basics", "distributed-lock"}), SkillIDs: js([]string{"mysql", "redis", "spring-boot"}), Tags: js([]string{"项目"}), Likes: 38, Views: 760},
	}
	for i := range resources {
		resources[i].Status, resources[i].SourceType = "published", "creator"
		resources[i].CreatedAt = at(24 * (i + 1))
		if err := d.Community.SaveResource(&resources[i]); err != nil {
			return err
		}
	}
	return nil
}
