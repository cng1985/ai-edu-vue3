// 内置 AI 学习助手引擎（本地模拟版）
//
// 采用与真实 RAG 系统一致的"检索 → 组装 → 流式生成"三段式结构：
// 1. 将全部课程内容按标题切分为知识块，构建轻量倒排检索
// 2. 对用户问题做关键词召回与打分，取 Top-K 知识块
// 3. 组装回答并以流式（逐字符）回调输出，模拟 SSE 体验
//
// 如需接入真实大模型 API，只需替换 generateAnswer 的实现，
// 保持 ask(question, handlers) 的接口不变即可。

import { courses } from '../data/courses'

// ---------- 1. 知识库构建：结构感知切分 ----------

function buildKnowledgeBase() {
  const chunks = []
  for (const course of courses) {
    for (const chapter of course.chapters) {
      const sections = splitByHeadings(chapter.content)
      for (const section of sections) {
        const plain = stripMarkdown(section.body)
        if (plain.length < 30) continue
        chunks.push({
          courseId: course.id,
          courseTitle: course.title,
          chapterId: chapter.id,
          chapterTitle: chapter.title,
          heading: section.heading,
          text: plain
        })
      }
    }
  }
  return chunks
}

function splitByHeadings(markdown) {
  const lines = markdown.split('\n')
  const sections = []
  let current = { heading: '', body: '' }
  let inCode = false
  for (const line of lines) {
    if (line.trim().startsWith('```')) inCode = !inCode
    if (!inCode && /^#{1,3}\s/.test(line)) {
      if (current.body.trim()) sections.push(current)
      current = { heading: line.replace(/^#+\s*/, '').trim(), body: '' }
    } else {
      current.body += line + '\n'
    }
  }
  if (current.body.trim()) sections.push(current)
  return sections
}

function stripMarkdown(text) {
  return text
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/\|.*\|/g, ' ')
    .replace(/[#>*`_[\]()!-]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}

const knowledgeBase = buildKnowledgeBase()

// ---------- 2. 检索：关键词召回与打分 ----------

const STOP_WORDS = new Set([
  '什么', '怎么', '如何', '为什么', '哪些', '一个', '可以', '需要', '应该',
  '是不是', '吗', '呢', '的', '了', '和', '与', '或', '在', '是', '有',
  '请问', '请', '介绍', '一下', '讲讲', '说说', '解释', 'the', 'a', 'is', 'what', 'how'
])

function tokenize(query) {
  const tokens = new Set()
  // 英文与数字词
  for (const m of query.toLowerCase().matchAll(/[a-z0-9-]{2,}/g)) {
    tokens.add(m[0])
  }
  // 中文按 2~4 字滑窗组合，近似分词效果
  const zh = query.replace(/[^\u4e00-\u9fa5]/g, '#')
  for (const seg of zh.split('#')) {
    if (seg.length < 2) continue
    for (let n = 2; n <= Math.min(4, seg.length); n++) {
      for (let i = 0; i + n <= seg.length; i++) {
        const word = seg.slice(i, i + n)
        if (!STOP_WORDS.has(word)) tokens.add(word)
      }
    }
  }
  return [...tokens].filter((t) => !STOP_WORDS.has(t))
}

export function retrieve(query, k = 3) {
  const tokens = tokenize(query)
  if (tokens.length === 0) return []
  const scored = knowledgeBase.map((chunk) => {
    let score = 0
    const haystack = (chunk.heading + ' ' + chunk.chapterTitle + ' ' + chunk.text).toLowerCase()
    for (const token of tokens) {
      if (!haystack.includes(token)) continue
      // 长词与标题命中权重更高
      let weight = token.length
      if ((chunk.heading + chunk.chapterTitle).toLowerCase().includes(token)) weight *= 3
      score += weight
    }
    return { chunk, score }
  })
  return scored
    .filter((s) => s.score >= 4)
    .sort((a, b) => b.score - a.score)
    .slice(0, k)
}

// ---------- 3. 生成：组装回答 ----------

const GREETING_PATTERN = /^(你好|您好|hi|hello|嗨|哈喽|在吗)[!！。~\s]*$/i

function generateAnswer(question) {
  if (GREETING_PATTERN.test(question.trim())) {
    return {
      text:
        '你好！我是本系统内置的 AI 学习助手 🤖\n\n' +
        '我基于平台全部课程内容构建了本地知识库，可以帮你：\n\n' +
        '- **解答课程相关问题**，例如"什么是 RAG？"、"如何防止 Agent 死循环？"\n' +
        '- **定位知识点出处**，告诉你答案来自哪门课程的哪个章节\n' +
        '- **推荐学习路径**，试试问我"我该从哪门课开始学？"\n\n' +
        '请直接输入你的问题吧！',
      sources: []
    }
  }

  if (/学习路径|从哪|怎么学|先学|顺序|入门建议/.test(question)) {
    return {
      text:
        '根据平台的课程体系，推荐按以下路径循序渐进：\n\n' +
        '1. **提示词工程入门**（入门 · 约 45 分钟）——先掌握与大模型协作的基础语言，这是后续一切的地基。\n' +
        '2. **RAG 检索增强生成实战**（进阶 · 约 60 分钟）——学会为模型外挂知识库，解决幻觉与私有知识问题。\n' +
        '3. **AI 驱动的应用开发实战**（高级 · 约 180 分钟）——全链路整合：架构设计、Agent、Function Calling、生产容错与性能优化。\n\n' +
        '每完成一门课程，建议到「知识测验」板块做对应的测验巩固效果。',
      sources: []
    }
  }

  const hits = retrieve(question, 3)
  if (hits.length === 0) {
    return {
      text:
        '抱歉，我在当前课程知识库中没有检索到与这个问题足够相关的内容。\n\n' +
        '我目前的知识范围覆盖以下课程：\n\n' +
        courses.map((c) => `- **${c.title}**（${c.tags.join(' / ')}）`).join('\n') +
        '\n\n你可以换个提问方式，或使用更具体的关键词（如"混合检索"、"Function Calling"、"语义缓存"）再试一次。',
      sources: []
    }
  }

  const top = hits[0].chunk
  const intro = `关于「${question.trim().replace(/[？?]+$/, '')}」，我在课程知识库中找到了以下内容：\n\n`
  const parts = hits.map(({ chunk }, i) => {
    const snippet = chunk.text.length > 220 ? chunk.text.slice(0, 220) + '……' : chunk.text
    const title = chunk.heading || chunk.chapterTitle
    return `**${i + 1}. ${title}**\n\n${snippet}`
  })
  const outro =
    `\n\n---\n\n以上内容主要来自课程 **《${top.courseTitle}》** 的 **「${top.chapterTitle}」** 章节，` +
    '建议前往该章节阅读完整讲解与代码示例。'

  return {
    text: intro + parts.join('\n\n') + outro,
    sources: hits.map(({ chunk }) => ({
      courseId: chunk.courseId,
      courseTitle: chunk.courseTitle,
      chapterId: chunk.chapterId,
      chapterTitle: chunk.chapterTitle
    }))
  }
}

// ---------- 4. 流式输出 ----------

/**
 * 向助手提问，回答以流式回调返回。
 * @param {string} question 用户问题
 * @param {{ onToken?: (chunk: string) => void, onDone?: (result: { text: string, sources: Array }) => void }} handlers
 * @returns {() => void} 取消函数
 */
export function ask(question, handlers = {}) {
  const { onToken, onDone } = handlers
  const result = generateAnswer(question)
  let cancelled = false
  let index = 0

  // 模拟真实 LLM 的"思考延迟 + 逐块流式输出"
  const thinkDelay = 350 + Math.random() * 400

  function pump() {
    if (cancelled) return
    if (index >= result.text.length) {
      onDone && onDone(result)
      return
    }
    const step = 2 + Math.floor(Math.random() * 4)
    const chunk = result.text.slice(index, index + step)
    index += step
    onToken && onToken(chunk)
    setTimeout(pump, 12 + Math.random() * 24)
  }

  const timer = setTimeout(pump, thinkDelay)
  return () => {
    cancelled = true
    clearTimeout(timer)
  }
}
