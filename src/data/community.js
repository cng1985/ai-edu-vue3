// 首页社区板块的演示数据：动态、排行榜与热门话题。
// 接入真实后端时，替换为 API 返回的数据即可，界面结构无需改动。

export const communityPosts = [
  {
    id: 'p1',
    author: '林语汐',
    avatar: '汐',
    avatarColor: '#8b5cf6',
    level: 'Lv.6',
    title: 'RAG 检索质量提升 30% 的三个关键改动',
    excerpt:
      '折腾了两周的混合检索终于跑通了。核心结论：BM25 + 向量召回做粗排，再用 Cross-Encoder 精排，比单一向量检索的命中率高出一大截，尤其是专有名词多的场景……',
    tags: ['RAG', '重排序'],
    time: '12 分钟前',
    likes: 128,
    comments: 32,
    hot: true
  },
  {
    id: 'p2',
    author: '陈默',
    avatar: '默',
    avatarColor: '#0ea5e9',
    level: 'Lv.4',
    title: '提示词里加一句话，输出稳定性肉眼可见地变好',
    excerpt:
      '在系统提示词末尾加上「若不确定请回答不知道，禁止编造」，幻觉率直接降了一半。课程第 3 章讲的注入防御也亲测有效，附上我的对比实验记录……',
    tags: ['提示词', '实验'],
    time: '1 小时前',
    likes: 96,
    comments: 21,
    hot: true
  },
  {
    id: 'p3',
    author: 'Kepler',
    avatar: 'K',
    avatarColor: '#10b981',
    level: 'Lv.5',
    title: '智能客服 Agent 实战踩坑：Function Calling 超时重试',
    excerpt:
      '跟着「AI 驱动的应用开发实战」第 7 章搭了个客服 Agent，生产环境里工具调用偶发超时。分享一下我的指数退避 + 降级回复方案，代码已贴在评论区……',
    tags: ['Agent', 'Function Calling'],
    time: '3 小时前',
    likes: 74,
    comments: 18,
    hot: false
  },
  {
    id: 'p4',
    author: '阿栗',
    avatar: '栗',
    avatarColor: '#f59e0b',
    level: 'Lv.2',
    title: '零基础两周学完提示词工程，测验 95 分的笔记全公开',
    excerpt:
      '从完全不懂 LLM 到能独立设计结构化提示词，我把每一章的思维导图和易错点都整理出来了。新手推荐先做完测验再回头看笔记，效果翻倍……',
    tags: ['学习笔记', '入门'],
    time: '昨天 21:40',
    likes: 215,
    comments: 47,
    hot: false
  }
]

export const leaderboard = [
  { name: '苏子瞻', avatar: '苏', avatarColor: '#6366f1', chapters: 18, streak: 21 },
  { name: 'Nova', avatar: 'N', avatarColor: '#0ea5e9', chapters: 17, streak: 14 },
  { name: '江晚吟', avatar: '江', avatarColor: '#10b981', chapters: 15, streak: 9 },
  { name: '半仙', avatar: '仙', avatarColor: '#f97316', chapters: 12, streak: 12 },
  { name: 'Debbie', avatar: 'D', avatarColor: '#ef4444', chapters: 11, streak: 6 }
]

export const hotTopics = [
  { name: 'RAG 实战', posts: 342, trend: '+28%' },
  { name: '提示词工程', posts: 287, trend: '+16%' },
  { name: 'Agent 设计', posts: 201, trend: '+41%' },
  { name: '向量数据库', posts: 156, trend: '+9%' },
  { name: 'AI 原生架构', posts: 118, trend: '+22%' },
  { name: 'Function Calling', posts: 94, trend: '+12%' }
]

export const dailyQuote = {
  text: '未来的程序员不是写更多的代码，而是更好地调度智能。',
  from: '《AI 驱动的应用开发实战》· 引言'
}
