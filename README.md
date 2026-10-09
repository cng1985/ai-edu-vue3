# 知航 · AI 教育知识工作生态

> 让每个人通过 AI 学习、实践、工作，持续提升个人价值。

知航不是传统在线教育平台，而是连接 **学习、AI 陪伴、知识社区、项目实践、IT 任务、人才成长与商业价值** 的综合生态——一个「AI 时代的个人成长操作系统」。

```
学习 → 实践 → 产出 → 收益 → 投资学习 → 更强能力
```

## 核心领域模型

完整成长链：

```
Career → Role → Capability → Skill → Knowledge → Course → Project → Task → Evidence → Talent Profile
 职业     岗位     能力          技能     知识        课程     项目      任务    证据        人才画像
```

| 概念 | 说明 |
| :--- | :--- |
| 职业 / 岗位 / 能力 | 能力 = 能完成什么类型的问题，每项能力拆解为技能及要求等级 |
| 技能（L0–L5） | 未掌握 / 了解 / 可以辅助完成 / 可以独立完成 / 可以解决复杂问题 / 可以设计和指导他人 |
| 知识图谱 | 知识点之间有前置、关联、相似、高级四类关系，决定学习顺序 |
| 课程 | 课程 → 章节，通过映射引用知识点：课程结构与知识结构分离，知识点可被多门课程复用 |
| 学习状态 | 学会不是 Boolean：`mastery`（掌握度）+ `confidence`（置信度），含遗忘衰减与 1/3/7/15/30 天间隔复习 |
| 技能评价 | **技能 = 知识状态 + 项目证据 + 任务表现**。只有知识最高 L2；L3 需项目证据；L4/L5 需高质量项目或真实任务；证据也不能替代知识 |
| 项目实践 | Project → Task → Execution → Evidence，提交后由 Review Agent 评审生成能力证据 |
| IT 任务市场 | 企业发布任务 → 技能模型匹配 → 人才画像 → 推荐人员 → 录用 → 验收结算（收入 + 任务表现证据） |
| 知识社区 | 问题 → 社区回答 → AI 总结 → 知识资产 |
| 知识资产 | 统一 Knowledge Resource：文章、课程、视频、代码、案例、项目、任务、Prompt、SOP |
| 人才画像 | 职业目标、技能水平、知识掌握、项目经历、任务记录、社区贡献、收入记录 |

## AI 架构

```
用户 → AI Learning Kernel → AI Execution Kernel → Pipeline Runtime → Resource / Tool Layer
                                                                   (Knowledge / Course / Project / Task / Skill / Talent)
```

**学习 Pipeline**（`server/internal/service/kernel.go`，运行时 `server/pkg/pipeline`）：

```
UserContextStage → CareerAnalysisStage → SkillGapStage → KnowledgeRecommendStage
→ LearningPlanStage → AIResponseStage → EvaluationStage → ProfileUpdateStage
```

结构化分析（差距、推荐、计划、评估）由成长计算引擎 `server/pkg/growth` 完成；大模型只负责把结果讲给学习者听。未配置大模型时 `AIResponseStage` 标记为跳过并使用规则模板，Pipeline 仍可完整运行。每次运行都会持久化 Stage 轨迹，可在学习端与管理端查看。

**AI Agent 体系**：

| Agent | 职责 |
| :--- | :--- |
| Career Agent | 职业规划、能力分析、成长路线 |
| Learning Agent | 学习计划、知识解释、学习监督 |
| Knowledge Agent | 知识理解、知识关联、知识生成（结合知识图谱与课程 RAG 知识库） |
| Project Agent | 项目指导、任务拆解、技术方案 |
| Review Agent | 代码分析、项目评价、能力评估 |
| Opportunity Agent | 任务匹配、人才推荐、机会发现 |

每个 Agent 会注入与职责相关的个人成长上下文（目标、差距、掌握度、证据、任务匹配），学习者可在对话中查看 Agent 依据的数据。

## 产品结构

### 学习端（`front/app`，学习者 / 创作者 / 企业）

| 模块 | 说明 |
| :--- | :--- |
| 成长中心 | 岗位达成度、个人成长飞轮、关键能力差距、AI 推荐、到期复习、推荐实践与匹配任务、学习计划、活跃度 |
| 职业与能力 | 职业 → 岗位 → 能力 → 技能要求，与个人等级对比，设定职业目标 |
| 知识图谱 | 分层可视化（按前置深度分列），节点颜色为个人掌握度；知识点详情含正文、练习、前后置关系、引用课程、知识资产与 Knowledge Agent |
| 课程学习 / 知识测验 | 章节完成与测验成绩通过映射同步到知识状态 |
| AI 学习内核 | 实时推送 8 个 Stage 的执行轨迹与学习建议、差距、计划、评估 |
| AI 伙伴 | 6 个 Agent 流式对话 |
| 项目实践 | 任务链、技能要求与个人匹配度、提交成果、评审结果与能力证据 |
| 任务市场 | 按匹配度排序的真实任务、申请与履约状态 |
| 企业工作台 | 发布任务、处理申请、AI 推荐人才、验收结算（企业账号） |
| 人才库 | 人才画像列表与详情（企业账号） |
| 知识社区 / 知识资产 | 问答采纳与点赞、AI 总结沉淀为知识资产；创作者发布资产 |
| 人才画像 | 技能雷达与三项得分明细、知识掌握、项目经历、任务记录、社区贡献、收入记录、能力证据 |

### 管理端（`front/admin`，管理员 / 运营 / 审核员）

生态看板、职业能力体系、技能与知识图谱（技能 / 知识点与练习题 / 图谱关系（防环路）/ 课程章节映射）、项目实践、IT 任务市场、人才库、社区与知识资产、AI 学习内核，以及用户、课程、题库、内容审核、客户咨询、RAG 知识库、大模型配置、权限与系统设置。

## 启动方式

```bash
# 安装依赖（应用端 + 管理端）
npm run install:all

# 终端 1：启动 Go API 服务（http://localhost:3001）
npm run dev:server

# 终端 2：启动管理端（http://localhost:5174）
npm run dev:admin

# 终端 3：启动学习端（http://localhost:5173）
npm run dev
```

首次启动会自动写入种子数据：5 个职业、9 个岗位、19 项技能、45 个带练习题的知识点与 49 条图谱关系、3 个实践项目、企业任务、社区内容与知识资产，以及 4 位拥有真实成长数据的演示人才（全部经由成长引擎计算生成）。已有数据库升级时，生态数据会在检测到没有职业数据时自动补齐。

## 演示账号

| 角色 | 用户名 | 密码 | 入口 |
| :--- | :--- | :--- | :--- |
| 学习者（AI 应用工程师方向） | demo | demo123 | 学习端 |
| 企业（星云科技） | company | company123 | 学习端 → 企业工作台 |
| 创作者 | expert | expert123 | 学习端 |
| 演示人才 | bob / alice / chen | bob123 / alice123 / chen123 | 学习端 |
| 管理员 | admin | admin123 | 管理端 |
| 审核员 | reviewer | review123 | 管理端 |
| 运营 | operator | oper123 | 管理端 |

## 后端

| 类别 | 选型 |
| :--- | :--- |
| 语言 | Go（Uber Fx 依赖注入、Gin、GORM） |
| 数据库 | SQLite（开发环境，可切换 MySQL） |
| 认证 | JWT + bcrypt，角色与状态以数据库实时数据为准 |
| 计算引擎 | `pkg/growth`：掌握度、技能等级、差距分析、匹配度、知识推荐、学习计划（纯函数，含单测） |
| AI 内核 | `pkg/pipeline`：泛型 Pipeline Runtime（Stage 轨迹、跳过、失败中止、panic 恢复） |
| 大模型 | OpenAI 兼容接口，支持多厂商路由与虚拟模型（管理端「大模型配置」） |

```
server/
├── cmd/server/              入口
├── internal/
│   ├── model/               领域模型（ecosystem.go）与接口视图（ecosystem_view.go）
│   ├── repository/          数据访问
│   ├── service/             catalog / growth / practice / market / community / talent / kernel / agent 等
│   ├── handler/ router/     HTTP 接口与路由（router/eco.go 为生态接口）
│   └── seed/                种子数据（eco_catalog.go 为职业与知识图谱）
└── pkg/
    ├── growth/              成长计算引擎
    ├── pipeline/            Pipeline Runtime
    ├── rbac/                角色与权限（学习者 / 创作者 / 企业 / 管理员 / 审核员 / 运营 / 游客）
    ├── llm/ rag/            大模型客户端与 RAG
    └── ...
```

### 主要接口

| 模块 | 接口 | 说明 |
| :--- | :--- | :--- |
| 职业与知识 | `GET /eco/careers` · `GET /eco/roles/:id` · `GET /eco/skills` | 职业体系、岗位详情（含个人差距） |
| 职业与知识 | `GET /eco/knowledge/graph?skillId=` · `GET /eco/knowledge/:id` | 知识图谱、知识点详情 |
| 个人成长 | `GET /me/overview` · `GET/PUT /me/goal` · `GET /me/gap` · `GET /me/skills` | 成长总览、职业目标、差距、技能状态 |
| 个人成长 | `POST /me/events` · `POST /me/knowledge/:id/practice` · `POST /me/chapters/complete` | 学习事件、练习判分、章节/测验同步 |
| 个人成长 | `GET /me/recommendations` · `GET /me/reviews` · `GET /me/profile` | 知识推荐、到期复习、人才画像 |
| 项目实践 | `GET /eco/projects[/:id]` · `POST /eco/projects/:id/tasks/:taskId/submissions` | 项目与任务、提交评审 |
| 任务市场 | `GET /eco/opportunities[/:id]` · `POST /eco/opportunities/:id/apply` | 任务列表（含匹配度）、申请 |
| 企业 | `GET/POST/PUT/DELETE /enterprise/opportunities` · `GET .../:id/candidates` · `POST /enterprise/applications/:id/decision` | 发布、推荐人才、录用/婉拒/验收 |
| 人才 | `GET /talents[/:id]` | 人才库与人才画像 |
| 社区 | `GET/POST /eco/community/posts` · `POST .../:id/answers` · `POST .../:id/summarize` | 社区内容、回答、AI 总结为知识资产 |
| 知识资产 | `GET/POST /eco/resources` | 知识资产 |
| AI | `GET /ai/agents` · `POST /ai/agents/:code/chat/stream` · `GET /ai/agents/:code/context` | Agent 列表、SSE 对话、上下文预览 |
| AI | `GET /ai/kernel/stages` · `POST /ai/kernel/run/stream` · `GET /ai/kernel/runs` | Pipeline 结构、运行（SSE 推送 Stage 轨迹）、运行记录 |
| 管理 | `/manage/*` | 生态看板与职业/技能/知识/项目/任务/人才/社区/内核治理 |

（以上接口前缀均为 `/api/v1`。）

### 服务端环境变量

| 变量 | 默认值 | 说明 |
| :--- | :--- | :--- |
| `PORT` | `3001` | API 服务端口 |
| `JWT_SECRET` | 内置开发密钥 | JWT 签名密钥，**生产环境必须配置** |
| `JWT_TTL_HOURS` | `168`（7 天） | 登录 token 有效期（小时） |
| `DB_PATH` | `data/ai-learning.db` | SQLite 数据库路径 |
| `CORS_ORIGINS` | `*` | 允许跨域的来源（逗号分隔） |
| `LLM_API_KEY` / `LLM_BASE_URL` / `LLM_MODEL` | — | 启用真实大模型（OpenAI 兼容接口），也可在管理端「大模型配置」中配置 |

未配置大模型时：Agent 对话会提示配置；AI 学习内核使用规则模板生成建议；Review Agent 使用可解释的规则评审（标注为「规则评审」）；社区总结使用抽取式总结。

## 测试

```bash
cd server && go test ./...
```

覆盖成长计算引擎（掌握度、遗忘、技能等级、差距、匹配、推荐、计划）、Pipeline Runtime，以及「目标 → 差距 → 推荐 → 练习 → 真实任务验收 → 收入与证据 → 技能升级」的端到端服务测试。
