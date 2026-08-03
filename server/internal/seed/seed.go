package seed

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cng1985/ai-learning-server/internal/model"
	"github.com/cng1985/ai-learning-server/internal/repository"
	"github.com/cng1985/ai-learning-server/pkg/authutil"
	"gorm.io/datatypes"
)

type chapterDef struct {
	ID      string
	Title   string
	Minutes int
	File    string
}

type courseDef struct {
	ID               string
	Title            string
	Description      string
	Level            string
	Tags             []string
	Icon             string
	Accent           string
	EstimatedMinutes int
	Status           string
	Chapters         []chapterDef
}

func Run(users *repository.UserRepo, courses *repository.CourseRepo, quizzes *repository.QuizRepo, reviews *repository.ReviewRepo, roles *repository.RoleRepo) error {
	seedRolePermissions(roles)
	total, err := users.Total()
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	fmt.Println("📦 首次启动，正在初始化数据...")
	now := time.Now().UnixMilli()
	contentDir := filepath.Join("..", "front", "app", "src", "data", "content")

	adminHash, _ := authutil.HashPassword("admin123")
	reviewerHash, _ := authutil.HashPassword("review123")
	demoHash, _ := authutil.HashPassword("demo123")

	operatorHash, _ := authutil.HashPassword("oper123")

	for _, u := range []model.User{
		{ID: "admin_001", Username: "admin", Nickname: "系统管理员", PasswordHash: adminHash, Role: "admin", Status: "active", Avatar: "管", AvatarColor: "#6366f1", JoinedAt: now},
		{ID: "reviewer_001", Username: "reviewer", Nickname: "内容审核员", PasswordHash: reviewerHash, Role: "reviewer", Status: "active", Avatar: "审", AvatarColor: "#0ea5e9", JoinedAt: now},
		{ID: "operator_001", Username: "operator", Nickname: "内容运营", PasswordHash: operatorHash, Role: "operator", Status: "active", Avatar: "运", AvatarColor: "#f59e0b", JoinedAt: now},
		{ID: "learner_demo", Username: "demo", Nickname: "演示学员", PasswordHash: demoHash, Role: "learner", Status: "active", Avatar: "学", AvatarColor: "#10b981", JoinedAt: now},
	} {
		if err := users.Create(&u); err != nil {
			return err
		}
	}

	courseDefs := []courseDef{
		{
			ID: "prompt-engineering", Title: "提示词工程入门",
			Description: "掌握与大语言模型高效协作的核心技能：从基础结构到 Few-shot、思维链与注入防御，建立系统化的提示词设计方法论。",
			Level: "入门", Tags: []string{"提示词", "LLM", "基础"}, Icon: "✍️", Accent: "#6366f1",
			EstimatedMinutes: 45, Status: "published",
			Chapters: []chapterDef{
				{"basics", "提示词工程基础", 12, "prompt-engineering/01-basics.md"},
				{"structured", "结构化提示词设计", 15, "prompt-engineering/02-structured.md"},
				{"advanced", "高级技巧与常见陷阱", 18, "prompt-engineering/03-advanced.md"},
			},
		},
		{
			ID: "rag-in-action", Title: "RAG 检索增强生成实战",
			Description: "从原理到生产落地：文档切分、向量化、混合检索与重排序，构建高召回质量的企业级知识库问答系统。",
			Level: "进阶", Tags: []string{"RAG", "向量数据库", "知识库"}, Icon: "🔍", Accent: "#0ea5e9",
			EstimatedMinutes: 60, Status: "published",
			Chapters: []chapterDef{
				{"overview", "RAG 原理与整体链路", 15, "rag-in-action/01-overview.md"},
				{"chunking", "文档切分与向量化策略", 20, "rag-in-action/02-chunking.md"},
				{"retrieval", "检索优化：混合检索与重排序", 25, "rag-in-action/03-retrieval.md"},
			},
		},
		{
			ID: "ai-native-dev", Title: "AI 驱动的应用开发实战",
			Description: "从提示词工程到 AI 原生架构的全链路指南：智能体设计、Function Calling、生产级容错与性能优化，完成从\"写代码\"到\"调度智能\"的跃迁。",
			Level: "高级", Tags: []string{"AI原生", "架构设计", "Agent", "实战"}, Icon: "🤖", Accent: "#f59e0b",
			EstimatedMinutes: 180, Status: "published",
			Chapters: []chapterDef{
				{"abstract", "摘要", 5, "ai-native-dev/01-abstract.md"},
				{"introduction", "引言与背景", 10, "ai-native-dev/02-introduction.md"},
				{"challenges", "问题与挑战", 15, "ai-native-dev/03-challenges.md"},
				{"core-concepts", "核心概念", 20, "ai-native-dev/04-core-concepts.md"},
				{"architecture", "架构设计", 20, "ai-native-dev/05-architecture.md"},
				{"implementation", "实现细节", 25, "ai-native-dev/06-implementation.md"},
				{"code-practice", "代码实践：智能客服 Agent", 25, "ai-native-dev/07-code-practice.md"},
				{"best-practices", "最佳实践与经验总结", 15, "ai-native-dev/08-best-practices.md"},
				{"pitfalls", "常见坑与排错指南", 15, "ai-native-dev/09-pitfalls.md"},
				{"performance", "性能优化", 15, "ai-native-dev/10-performance.md"},
				{"summary", "总结与展望", 8, "ai-native-dev/11-summary.md"},
				{"references", "参考资料与延伸阅读", 7, "ai-native-dev/12-references.md"},
			},
		},
	}

	for _, def := range courseDefs {
		tags, _ := json.Marshal(def.Tags)
		course := model.Course{
			ID: def.ID, Title: def.Title, Description: def.Description, Level: def.Level,
			Tags: datatypes.JSON(tags), Icon: def.Icon, Accent: def.Accent,
			EstimatedMinutes: def.EstimatedMinutes, Status: def.Status,
			CreatedAt: now, UpdatedAt: now,
		}
		if err := courses.Create(&course); err != nil {
			return err
		}
		for _, ch := range def.Chapters {
			content := readMd(filepath.Join(contentDir, ch.File))
			chapter := model.Chapter{
				ID: ch.ID, CourseID: def.ID, Title: ch.Title, Minutes: ch.Minutes,
				Content: content, Status: "published", UpdatedAt: now,
			}
			if err := courses.AddChapter(&chapter); err != nil {
				return err
			}
		}
	}

	quizData := []struct {
		id, courseID, title, desc string
		questions                 []model.Question
	}{
		{"prompt-engineering", "prompt-engineering", "提示词工程入门测验", "检验你对提示词构成要素、结构化设计与常见陷阱的掌握程度。", promptQuestions()},
		{"rag-in-action", "rag-in-action", "RAG 实战测验", "覆盖 RAG 链路、切分策略、混合检索与重排序的核心知识点。", ragQuestions()},
		{"ai-native-dev", "ai-native-dev", "AI 原生应用开发测验", "综合考察 AI 原生架构、Agent、Function Calling、容错与性能优化。", aiNativeQuestions()},
	}

	for _, q := range quizData {
		questions, _ := json.Marshal(q.questions)
		quiz := model.Quiz{
			ID: q.id, CourseID: q.courseID, Title: q.title, Description: q.desc,
			Questions: datatypes.JSON(questions), Status: "published",
			CreatedAt: now, UpdatedAt: now,
		}
		if err := quizzes.Create(&quiz); err != nil {
			return err
		}
	}

	for _, r := range []model.Review{
		{ID: "rev_001", Type: "chapter", CourseID: "prompt-engineering", TargetID: "basics", Title: "提示词工程基础（修订版）", Content: "## 修订内容\n\n新增 Few-shot 示例说明...", Submitter: "demo", Status: "pending", AIScore: 82, AIFeedback: "内容结构完整，建议补充更多实际案例。", CreatedAt: now - 86400000},
		{ID: "rev_002", Type: "quiz", CourseID: "rag-in-action", TargetID: "rag-in-action", Title: "RAG 测验新增题目", Content: `{"text":"向量数据库选型时首要考虑的因素是？","options":["价格","向量维度兼容性","品牌知名度","界面美观"],"answer":1,"explanation":"兼容性决定能否正确存储和检索向量。"}`, Submitter: "demo", Status: "pending", AIScore: 75, AIFeedback: "题目质量尚可，选项区分度需加强。", CreatedAt: now - 43200000},
	} {
		if err := reviews.Create(&r); err != nil {
			return err
		}
	}

	fmt.Println("✅ 数据种子已写入 SQLite 数据库")
	fmt.Println("   管理员账号: admin / admin123")
	fmt.Println("   审核员账号: reviewer / review123")
	fmt.Println("   运营账号: operator / oper123")
	return nil
}

func seedRolePermissions(roles *repository.RoleRepo) {
	importPerms := func(role string, perms []string) {
		b, _ := json.Marshal(perms)
		_ = roles.Upsert(&model.RolePermission{
			Role: role, Permissions: datatypes.JSON(b), UpdatedAt: time.Now().UnixMilli(),
		})
	}
	defaults := map[string][]string{
		"admin":    {"user:read", "user:create", "user:update", "user:delete", "course:read", "course:write", "course:delete", "quiz:read", "quiz:write", "quiz:delete", "review:read", "review:approve", "dashboard:read", "role:manage", "settings:manage", "ai:chat", "customer:read", "customer:reply"},
		"reviewer": {"course:read", "quiz:read", "review:read", "review:approve", "dashboard:read", "ai:chat"},
		"operator": {"course:read", "course:write", "quiz:read", "quiz:write", "dashboard:read", "ai:chat", "customer:read", "customer:reply"},
		"learner":  {"course:read", "quiz:read", "ai:chat", "customer:chat"},
		"guest":    {"ai:chat"},
	}
	for role, perms := range defaults {
		existing, err := roles.FindByRole(role)
		if err != nil {
			importPerms(role, perms)
			continue
		}
		if role != "admin" && role != "operator" && role != "learner" {
			continue
		}
		var current []string
		_ = json.Unmarshal(existing.Permissions, &current)
		merged := mergePermissions(current, perms)
		if len(merged) != len(current) {
			importPerms(role, merged)
		}
	}
}

func mergePermissions(current, defaults []string) []string {
	set := map[string]bool{}
	for _, p := range current {
		set[p] = true
	}
	for _, p := range defaults {
		set[p] = true
	}
	out := make([]string, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	return out
}

func readMd(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func promptQuestions() []model.Question {
	return []model.Question{
		{Text: "提示词工程的本质最准确的描述是？", Options: []string{"一种向 AI 提问的礼貌用语规范", "面向概率性模型的接口编程", "一种新的编程语言语法", "仅用于聊天机器人的对话技巧"}, Answer: 1, Explanation: "提示词通过重塑模型的条件概率输出来约束其行为，本质是面向概率性模型的接口编程。"},
		{Text: "一个结构完整的提示词通常不包含以下哪个要素？", Options: []string{"角色设定", "任务描述", "模型的内部权重参数", "输出约束"}, Answer: 2, Explanation: "提示词的四大要素是角色设定、任务描述、上下文信息与输出约束；模型权重是训练产物，无法通过提示词修改。"},
		{Text: "使用 XML 标签划分提示词区域的主要好处是？", Options: []string{"让提示词看起来更专业", "减少 Token 消耗", "明确划分语义区域，隔离不可信内容并降低解析歧义", "强制模型输出 XML 格式"}, Answer: 2, Explanation: "标签划分语义区域，能将指令与动态注入的数据隔离，既降低歧义又能防御提示词注入。"},
		{Text: "关于 Few-shot 示例，下列说法正确的是？", Options: []string{"示例越多越好，没有上限", "给出\"输入→期望输出\"样例比冗长的文字描述更有效", "Few-shot 只适用于代码生成任务", "示例的格式与最终输出格式无关"}, Answer: 1, Explanation: "示例优于描述是提示词设计的关键原则，模型会模仿示例的格式与风格；但示例过多会占用上下文并可能引入偏差。"},
		{Text: "长提示词中\"中间迷失\"现象的推荐对策是？", Options: []string{"把所有内容都写在中间", "将关键指令放在提示词的开头和结尾", "增大模型的 temperature", "删除所有上下文信息"}, Answer: 1, Explanation: "模型对提示词首尾位置的信息更敏感，关键指令应放在开头和结尾以对抗中间迷失。"},
		{Text: "在生产环境中管理提示词的正确做法是？", Options: []string{"硬编码在业务代码里，方便查看", "外置化存储、绑定模型版本、支持 A/B 测试与回滚", "每次上线前手工修改，不留记录", "只保留最新版本，删除历史版本"}, Answer: 1, Explanation: "Prompt 就是业务逻辑，必须像代码一样进行版本管理，并绑定验证通过的模型版本。"},
	}
}

func ragQuestions() []model.Question {
	return []model.Question{
		{Text: "RAG 解决大模型哪两个先天缺陷？", Options: []string{"推理速度慢和显存占用高", "知识截止和私有知识缺失", "不支持多语言和多模态", "上下文窗口小和 Token 太贵"}, Answer: 1, Explanation: "RAG 通过推理时动态检索外部知识，解决模型知识截止与不了解私有数据的问题。"},
		{Text: "RAG 离线阶段的正确处理顺序是？", Options: []string{"向量化 → 切分 → 入库", "切分 → 向量化 → 存入向量数据库", "入库 → 切分 → 向量化", "切分 → 入库 → 检索"}, Answer: 1, Explanation: "离线阶段先将文档切分为语义连贯的 Chunk，再经 Embedding 模型向量化，最后存入向量数据库。"},
		{Text: "生产环境推荐的默认文档切分策略是？", Options: []string{"固定长度切分，不设重叠", "整篇文档不切分直接入库", "结构感知切分（利用标题、段落等天然结构）", "每个句子切成一个块"}, Answer: 2, Explanation: "结构感知切分利用文档天然结构保证语义完整性，是效果与成本平衡的推荐默认策略。"},
		{Text: "更换 Embedding 模型后必须做什么？", Options: []string{"什么都不用做，向量可以混用", "只需重建最近一个月的数据", "全量重建向量索引", "把向量维度调成一样即可"}, Answer: 2, Explanation: "不同 Embedding 模型的向量空间互不兼容，混用会导致检索完全失效，必须全量重建索引。"},
		{Text: "混合检索指的是？", Options: []string{"同时使用两个大模型生成回答", "稠密向量检索与稀疏关键词检索（BM25）结合", "把多个向量数据库的数据合并", "先检索图片再检索文本"}, Answer: 1, Explanation: "混合检索结合语义相似度（向量）与精确匹配（BM25）两路召回，能同时覆盖语义与专有名词场景。"},
		{Text: "重排序（Rerank）阶段通常使用什么模型结构，为什么？", Options: []string{"双编码器，因为速度快", "交叉编码器，将查询与候选文档拼接逐一打分，精度更高", "决策树，因为可解释", "与召回阶段完全相同的模型"}, Answer: 1, Explanation: "召回用双编码器保证速度，精排用交叉编码器保证精度，\"召回 Top-20 → 精排 Top-3\"是典型配置。"},
		{Text: "多轮对话中用户问\"那它性能怎么样？\"导致检索失败，应引入什么机制？", Options: []string{"查询重写", "增大 Top-K", "更换向量数据库", "提高 temperature"}, Answer: 0, Explanation: "查询重写利用对话历史将代词补全为具体实体，生成独立完整的搜索查询。"},
	}
}

func aiNativeQuestions() []model.Question {
	return []model.Question{
		{Text: "AI 辅助开发与 AI 原生开发的核心区别是？", Options: []string{"是否使用 Python 语言", "AI 介入的是研发过程还是作为应用运行时的核心驱动", "是否需要联网", "代码量的多少"}, Answer: 1, Explanation: "AI 辅助开发中 AI 是生产力工具，交付的软件不依赖大模型运行；AI 原生开发中大模型是系统的推理核心与逻辑路由。"},
		{Text: "传统 CRUD 架构与 AI 原生架构在异常处理上的差异是？", Options: []string{"两者完全相同", "AI 原生架构不需要异常处理", "传统架构依赖 Try-Catch 与事务回滚，AI 原生架构依赖约束输出、重试机制与护栏", "AI 原生架构只需要打日志"}, Answer: 2, Explanation: "面对概率性输出，AI 原生架构用 JSON Schema 约束、重试与护栏机制取代传统的确定性异常处理。"},
		{Text: "Function Calling 机制中，大模型实际做的事情是？", Options: []string{"直接在模型内部执行函数代码", "输出包含函数名和参数的 JSON 对象，由应用层执行真实代码", "把函数编译成机器码", "直接连接数据库执行 SQL"}, Answer: 1, Explanation: "模型只输出调用意图（函数名+参数），应用层拦截执行后将结果回传给模型继续推理，实现推理与执行解耦。"},
		{Text: "ReAct 智能体范式的循环是？", Options: []string{"编码 → 测试 → 部署", "感知输入 → 思考规划 → 调用工具 → 观察结果 → 再次思考", "训练 → 验证 → 上线", "提问 → 回答 → 结束"}, Answer: 1, Explanation: "ReAct（Reason + Act）通过\"思考-行动-观察\"循环让模型自主完成多步任务，直到达成目标。"},
		{Text: "防止 Agent 陷入工具调用死循环的有效手段不包括？", Options: []string{"强制最大迭代次数", "重复调用熔断（相同工具+相同参数连续两次则中断）", "将工具报错信息显式、清晰地反馈给模型", "无限增加 API 额度让它一直试"}, Answer: 3, Explanation: "死循环防护靠最大步数限制、重复调用熔断和清晰的错误反馈；增加额度只会放大损失。"},
		{Text: "流式输出（SSE）优化的核心指标是？", Options: []string{"总生成时间", "首字延迟（TTFB）", "Token 单价", "模型参数量"}, Answer: 1, Explanation: "流式输出不缩短总生成时间，但将首字延迟降到几百毫秒，大幅改善用户体感。"},
		{Text: "语义缓存与传统缓存的关键区别是？", Options: []string{"语义缓存基于向量相似度命中，可处理字面不同但语义相同的请求", "语义缓存速度更快", "传统缓存不能存字符串", "两者没有区别"}, Answer: 0, Explanation: "传统缓存要求 Key 完全匹配；语义缓存通过 Embedding 相似度匹配语义相同的变体表达，但需防范\"假阳性\"。"},
		{Text: "评估驱动开发（EDD）的核心实践是？", Options: []string{"不写测试，靠人工体验", "构建黄金数据集，每次迭代自动运行评估并以量化指标判断是否退化", "只在上线前评估一次", "让用户当测试员"}, Answer: 1, Explanation: "由于 LLM 输出非确定性，传统断言失效，必须用黄金数据集+量化指标（准确率、相关性等）持续评估。"},
		{Text: "人在回路（HITL）设计适用于什么场景？", Options: []string{"所有 AI 请求都需要人工审核", "涉及资金交易、隐私或高风险决策的关键流程", "仅用于测试环境", "模型响应太快时用来限速"}, Answer: 1, Explanation: "HITL 在高风险场景插入人工审批节点：AI 负责意图理解与信息提取，最终执行由风险策略路由到自动执行或人工审核。"},
	}
}