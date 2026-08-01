import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'
import { loadDb, saveDb, resetCache } from './utils/db.js'
import { hashPassword } from './utils/auth.js'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const ROOT = path.resolve(__dirname, '../../')
const CONTENT_DIR = path.join(ROOT, 'src/data/content')

function readMd(relPath) {
  const full = path.join(CONTENT_DIR, relPath)
  if (!fs.existsSync(full)) return ''
  return fs.readFileSync(full, 'utf-8')
}

const courseDefs = [
  {
    id: 'prompt-engineering',
    title: '提示词工程入门',
    description: '掌握与大语言模型高效协作的核心技能：从基础结构到 Few-shot、思维链与注入防御，建立系统化的提示词设计方法论。',
    level: '入门',
    tags: ['提示词', 'LLM', '基础'],
    icon: '✍️',
    accent: '#6366f1',
    estimatedMinutes: 45,
    status: 'published',
    chapters: [
      { id: 'basics', title: '提示词工程基础', minutes: 12, file: 'prompt-engineering/01-basics.md' },
      { id: 'structured', title: '结构化提示词设计', minutes: 15, file: 'prompt-engineering/02-structured.md' },
      { id: 'advanced', title: '高级技巧与常见陷阱', minutes: 18, file: 'prompt-engineering/03-advanced.md' }
    ]
  },
  {
    id: 'rag-in-action',
    title: 'RAG 检索增强生成实战',
    description: '从原理到生产落地：文档切分、向量化、混合检索与重排序，构建高召回质量的企业级知识库问答系统。',
    level: '进阶',
    tags: ['RAG', '向量数据库', '知识库'],
    icon: '🔍',
    accent: '#0ea5e9',
    estimatedMinutes: 60,
    status: 'published',
    chapters: [
      { id: 'overview', title: 'RAG 原理与整体链路', minutes: 15, file: 'rag-in-action/01-overview.md' },
      { id: 'chunking', title: '文档切分与向量化策略', minutes: 20, file: 'rag-in-action/02-chunking.md' },
      { id: 'retrieval', title: '检索优化：混合检索与重排序', minutes: 25, file: 'rag-in-action/03-retrieval.md' }
    ]
  },
  {
    id: 'ai-native-dev',
    title: 'AI 驱动的应用开发实战',
    description: '从提示词工程到 AI 原生架构的全链路指南：智能体设计、Function Calling、生产级容错与性能优化，完成从"写代码"到"调度智能"的跃迁。',
    level: '高级',
    tags: ['AI原生', '架构设计', 'Agent', '实战'],
    icon: '🤖',
    accent: '#f59e0b',
    estimatedMinutes: 180,
    status: 'published',
    chapters: [
      { id: 'abstract', title: '摘要', minutes: 5, file: 'ai-native-dev/01-abstract.md' },
      { id: 'introduction', title: '引言与背景', minutes: 10, file: 'ai-native-dev/02-introduction.md' },
      { id: 'challenges', title: '问题与挑战', minutes: 15, file: 'ai-native-dev/03-challenges.md' },
      { id: 'core-concepts', title: '核心概念', minutes: 20, file: 'ai-native-dev/04-core-concepts.md' },
      { id: 'architecture', title: '架构设计', minutes: 20, file: 'ai-native-dev/05-architecture.md' },
      { id: 'implementation', title: '实现细节', minutes: 25, file: 'ai-native-dev/06-implementation.md' },
      { id: 'code-practice', title: '代码实践：智能客服 Agent', minutes: 25, file: 'ai-native-dev/07-code-practice.md' },
      { id: 'best-practices', title: '最佳实践与经验总结', minutes: 15, file: 'ai-native-dev/08-best-practices.md' },
      { id: 'pitfalls', title: '常见坑与排错指南', minutes: 15, file: 'ai-native-dev/09-pitfalls.md' },
      { id: 'performance', title: '性能优化', minutes: 15, file: 'ai-native-dev/10-performance.md' },
      { id: 'summary', title: '总结与展望', minutes: 8, file: 'ai-native-dev/11-summary.md' },
      { id: 'references', title: '参考资料与延伸阅读', minutes: 7, file: 'ai-native-dev/12-references.md' }
    ]
  }
]

const quizDefs = [
  {
    id: 'prompt-engineering',
    courseId: 'prompt-engineering',
    title: '提示词工程入门测验',
    description: '检验你对提示词构成要素、结构化设计与常见陷阱的掌握程度。',
    status: 'published',
    questions: [
      { text: '提示词工程的本质最准确的描述是？', options: ['一种向 AI 提问的礼貌用语规范', '面向概率性模型的接口编程', '一种新的编程语言语法', '仅用于聊天机器人的对话技巧'], answer: 1, explanation: '提示词通过重塑模型的条件概率输出来约束其行为，本质是面向概率性模型的接口编程。' },
      { text: '一个结构完整的提示词通常不包含以下哪个要素？', options: ['角色设定', '任务描述', '模型的内部权重参数', '输出约束'], answer: 2, explanation: '提示词的四大要素是角色设定、任务描述、上下文信息与输出约束；模型权重是训练产物，无法通过提示词修改。' },
      { text: '使用 XML 标签划分提示词区域的主要好处是？', options: ['让提示词看起来更专业', '减少 Token 消耗', '明确划分语义区域，隔离不可信内容并降低解析歧义', '强制模型输出 XML 格式'], answer: 2, explanation: '标签划分语义区域，能将指令与动态注入的数据隔离，既降低歧义又能防御提示词注入。' },
      { text: '关于 Few-shot 示例，下列说法正确的是？', options: ['示例越多越好，没有上限', '给出"输入→期望输出"样例比冗长的文字描述更有效', 'Few-shot 只适用于代码生成任务', '示例的格式与最终输出格式无关'], answer: 1, explanation: '示例优于描述是提示词设计的关键原则，模型会模仿示例的格式与风格；但示例过多会占用上下文并可能引入偏差。' },
      { text: '长提示词中"中间迷失"现象的推荐对策是？', options: ['把所有内容都写在中间', '将关键指令放在提示词的开头和结尾', '增大模型的 temperature', '删除所有上下文信息'], answer: 1, explanation: '模型对提示词首尾位置的信息更敏感，关键指令应放在开头和结尾以对抗中间迷失。' },
      { text: '在生产环境中管理提示词的正确做法是？', options: ['硬编码在业务代码里，方便查看', '外置化存储、绑定模型版本、支持 A/B 测试与回滚', '每次上线前手工修改，不留记录', '只保留最新版本，删除历史版本'], answer: 1, explanation: 'Prompt 就是业务逻辑，必须像代码一样进行版本管理，并绑定验证通过的模型版本。' }
    ]
  },
  {
    id: 'rag-in-action',
    courseId: 'rag-in-action',
    title: 'RAG 实战测验',
    description: '覆盖 RAG 链路、切分策略、混合检索与重排序的核心知识点。',
    status: 'published',
    questions: [
      { text: 'RAG 解决大模型哪两个先天缺陷？', options: ['推理速度慢和显存占用高', '知识截止和私有知识缺失', '不支持多语言和多模态', '上下文窗口小和 Token 太贵'], answer: 1, explanation: 'RAG 通过推理时动态检索外部知识，解决模型知识截止与不了解私有数据的问题。' },
      { text: 'RAG 离线阶段的正确处理顺序是？', options: ['向量化 → 切分 → 入库', '切分 → 向量化 → 存入向量数据库', '入库 → 切分 → 向量化', '切分 → 入库 → 检索'], answer: 1, explanation: '离线阶段先将文档切分为语义连贯的 Chunk，再经 Embedding 模型向量化，最后存入向量数据库。' },
      { text: '生产环境推荐的默认文档切分策略是？', options: ['固定长度切分，不设重叠', '整篇文档不切分直接入库', '结构感知切分（利用标题、段落等天然结构）', '每个句子切成一个块'], answer: 2, explanation: '结构感知切分利用文档天然结构保证语义完整性，是效果与成本平衡的推荐默认策略。' },
      { text: '更换 Embedding 模型后必须做什么？', options: ['什么都不用做，向量可以混用', '只需重建最近一个月的数据', '全量重建向量索引', '把向量维度调成一样即可'], answer: 2, explanation: '不同 Embedding 模型的向量空间互不兼容，混用会导致检索完全失效，必须全量重建索引。' },
      { text: '混合检索指的是？', options: ['同时使用两个大模型生成回答', '稠密向量检索与稀疏关键词检索（BM25）结合', '把多个向量数据库的数据合并', '先检索图片再检索文本'], answer: 1, explanation: '混合检索结合语义相似度（向量）与精确匹配（BM25）两路召回，能同时覆盖语义与专有名词场景。' },
      { text: '重排序（Rerank）阶段通常使用什么模型结构，为什么？', options: ['双编码器，因为速度快', '交叉编码器，将查询与候选文档拼接逐一打分，精度更高', '决策树，因为可解释', '与召回阶段完全相同的模型'], answer: 1, explanation: '召回用双编码器保证速度，精排用交叉编码器保证精度，"召回 Top-20 → 精排 Top-3"是典型配置。' },
      { text: '多轮对话中用户问"那它性能怎么样？"导致检索失败，应引入什么机制？', options: ['查询重写', '增大 Top-K', '更换向量数据库', '提高 temperature'], answer: 0, explanation: '查询重写利用对话历史将代词补全为具体实体，生成独立完整的搜索查询。' }
    ]
  },
  {
    id: 'ai-native-dev',
    courseId: 'ai-native-dev',
    title: 'AI 原生应用开发测验',
    description: '综合考察 AI 原生架构、Agent、Function Calling、容错与性能优化。',
    status: 'published',
    questions: [
      { text: 'AI 辅助开发与 AI 原生开发的核心区别是？', options: ['是否使用 Python 语言', 'AI 介入的是研发过程还是作为应用运行时的核心驱动', '是否需要联网', '代码量的多少'], answer: 1, explanation: 'AI 辅助开发中 AI 是生产力工具，交付的软件不依赖大模型运行；AI 原生开发中大模型是系统的推理核心与逻辑路由。' },
      { text: '传统 CRUD 架构与 AI 原生架构在异常处理上的差异是？', options: ['两者完全相同', 'AI 原生架构不需要异常处理', '传统架构依赖 Try-Catch 与事务回滚，AI 原生架构依赖约束输出、重试机制与护栏', 'AI 原生架构只需要打日志'], answer: 2, explanation: '面对概率性输出，AI 原生架构用 JSON Schema 约束、重试与护栏机制取代传统的确定性异常处理。' },
      { text: 'Function Calling 机制中，大模型实际做的事情是？', options: ['直接在模型内部执行函数代码', '输出包含函数名和参数的 JSON 对象，由应用层执行真实代码', '把函数编译成机器码', '直接连接数据库执行 SQL'], answer: 1, explanation: '模型只输出调用意图（函数名+参数），应用层拦截执行后将结果回传给模型继续推理，实现推理与执行解耦。' },
      { text: 'ReAct 智能体范式的循环是？', options: ['编码 → 测试 → 部署', '感知输入 → 思考规划 → 调用工具 → 观察结果 → 再次思考', '训练 → 验证 → 上线', '提问 → 回答 → 结束'], answer: 1, explanation: 'ReAct（Reason + Act）通过"思考-行动-观察"循环让模型自主完成多步任务，直到达成目标。' },
      { text: '防止 Agent 陷入工具调用死循环的有效手段不包括？', options: ['强制最大迭代次数', '重复调用熔断（相同工具+相同参数连续两次则中断）', '将工具报错信息显式、清晰地反馈给模型', '无限增加 API 额度让它一直试'], answer: 3, explanation: '死循环防护靠最大步数限制、重复调用熔断和清晰的错误反馈；增加额度只会放大损失。' },
      { text: '流式输出（SSE）优化的核心指标是？', options: ['总生成时间', '首字延迟（TTFB）', 'Token 单价', '模型参数量'], answer: 1, explanation: '流式输出不缩短总生成时间，但将首字延迟降到几百毫秒，大幅改善用户体感。' },
      { text: '语义缓存与传统缓存的关键区别是？', options: ['语义缓存基于向量相似度命中，可处理字面不同但语义相同的请求', '语义缓存速度更快', '传统缓存不能存字符串', '两者没有区别'], answer: 0, explanation: '传统缓存要求 Key 完全匹配；语义缓存通过 Embedding 相似度匹配语义相同的变体表达，但需防范"假阳性"。' },
      { text: '评估驱动开发（EDD）的核心实践是？', options: ['不写测试，靠人工体验', '构建黄金数据集，每次迭代自动运行评估并以量化指标判断是否退化', '只在上线前评估一次', '让用户当测试员'], answer: 1, explanation: '由于 LLM 输出非确定性，传统断言失效，必须用黄金数据集+量化指标（准确率、相关性等）持续评估。' },
      { text: '人在回路（HITL）设计适用于什么场景？', options: ['所有 AI 请求都需要人工审核', '涉及资金交易、隐私或高风险决策的关键流程', '仅用于测试环境', '模型响应太快时用来限速'], answer: 1, explanation: 'HITL 在高风险场景插入人工审批节点：AI 负责意图理解与信息提取，最终执行由风险策略路由到自动执行或人工审核。' }
    ]
  }
]

const now = Date.now()

resetCache()
const db = loadDb()

db.users = [
  {
    id: 'admin_001',
    username: 'admin',
    nickname: '系统管理员',
    passwordHash: hashPassword('admin123'),
    role: 'admin',
    status: 'active',
    avatar: '管',
    avatarColor: '#6366f1',
    joinedAt: now
  },
  {
    id: 'reviewer_001',
    username: 'reviewer',
    nickname: '内容审核员',
    passwordHash: hashPassword('review123'),
    role: 'reviewer',
    status: 'active',
    avatar: '审',
    avatarColor: '#0ea5e9',
    joinedAt: now
  },
  {
    id: 'learner_demo',
    username: 'demo',
    nickname: '演示学员',
    passwordHash: hashPassword('demo123'),
    role: 'learner',
    status: 'active',
    avatar: '学',
    avatarColor: '#10b981',
    joinedAt: now
  }
]

db.courses = courseDefs.map((c) => ({
  ...c,
  chapters: c.chapters.map((ch) => ({
    id: ch.id,
    title: ch.title,
    minutes: ch.minutes,
    content: readMd(ch.file),
    status: 'published',
    updatedAt: now
  })),
  createdAt: now,
  updatedAt: now
}))

db.quizzes = quizDefs.map((q) => ({ ...q, createdAt: now, updatedAt: now }))

db.reviews = [
  {
    id: 'rev_001',
    type: 'chapter',
    courseId: 'prompt-engineering',
    targetId: 'basics',
    title: '提示词工程基础（修订版）',
    content: '## 修订内容\n\n新增 Few-shot 示例说明...',
    submitter: 'demo',
    status: 'pending',
    aiScore: 82,
    aiFeedback: '内容结构完整，建议补充更多实际案例。',
    createdAt: now - 86400000
  },
  {
    id: 'rev_002',
    type: 'quiz',
    courseId: 'rag-in-action',
    targetId: 'rag-in-action',
    title: 'RAG 测验新增题目',
    content: JSON.stringify({ text: '向量数据库选型时首要考虑的因素是？', options: ['价格', '向量维度兼容性', '品牌知名度', '界面美观'], answer: 1 }),
    submitter: 'demo',
    status: 'pending',
    aiScore: 75,
    aiFeedback: '题目质量尚可，选项区分度需加强。',
    createdAt: now - 43200000
  }
]

saveDb()
console.log('✅ 数据种子已写入 server/data/store.json')
console.log('   管理员账号: admin / admin123')
console.log('   审核员账号: reviewer / review123')
