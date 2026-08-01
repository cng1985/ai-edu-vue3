// 课程数据目录：章节正文以 Markdown 文件形式存放于 content/ 目录，
// 通过 Vite 的 ?raw 导入为字符串，实现内容与代码解耦。

import aiNative01 from './content/ai-native-dev/01-abstract.md?raw'
import aiNative02 from './content/ai-native-dev/02-introduction.md?raw'
import aiNative03 from './content/ai-native-dev/03-challenges.md?raw'
import aiNative04 from './content/ai-native-dev/04-core-concepts.md?raw'
import aiNative05 from './content/ai-native-dev/05-architecture.md?raw'
import aiNative06 from './content/ai-native-dev/06-implementation.md?raw'
import aiNative07 from './content/ai-native-dev/07-code-practice.md?raw'
import aiNative08 from './content/ai-native-dev/08-best-practices.md?raw'
import aiNative09 from './content/ai-native-dev/09-pitfalls.md?raw'
import aiNative10 from './content/ai-native-dev/10-performance.md?raw'
import aiNative11 from './content/ai-native-dev/11-summary.md?raw'
import aiNative12 from './content/ai-native-dev/12-references.md?raw'

import prompt01 from './content/prompt-engineering/01-basics.md?raw'
import prompt02 from './content/prompt-engineering/02-structured.md?raw'
import prompt03 from './content/prompt-engineering/03-advanced.md?raw'

import rag01 from './content/rag-in-action/01-overview.md?raw'
import rag02 from './content/rag-in-action/02-chunking.md?raw'
import rag03 from './content/rag-in-action/03-retrieval.md?raw'

export const courses = [
  {
    id: 'prompt-engineering',
    title: '提示词工程入门',
    description: '掌握与大语言模型高效协作的核心技能：从基础结构到 Few-shot、思维链与注入防御，建立系统化的提示词设计方法论。',
    level: '入门',
    tags: ['提示词', 'LLM', '基础'],
    icon: '✍️',
    accent: '#6366f1',
    estimatedMinutes: 45,
    chapters: [
      { id: 'basics', title: '提示词工程基础', minutes: 12, content: prompt01 },
      { id: 'structured', title: '结构化提示词设计', minutes: 15, content: prompt02 },
      { id: 'advanced', title: '高级技巧与常见陷阱', minutes: 18, content: prompt03 }
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
    chapters: [
      { id: 'overview', title: 'RAG 原理与整体链路', minutes: 15, content: rag01 },
      { id: 'chunking', title: '文档切分与向量化策略', minutes: 20, content: rag02 },
      { id: 'retrieval', title: '检索优化：混合检索与重排序', minutes: 25, content: rag03 }
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
    chapters: [
      { id: 'abstract', title: '摘要', minutes: 5, content: aiNative01 },
      { id: 'introduction', title: '引言与背景', minutes: 10, content: aiNative02 },
      { id: 'challenges', title: '问题与挑战', minutes: 15, content: aiNative03 },
      { id: 'core-concepts', title: '核心概念', minutes: 20, content: aiNative04 },
      { id: 'architecture', title: '架构设计', minutes: 20, content: aiNative05 },
      { id: 'implementation', title: '实现细节', minutes: 25, content: aiNative06 },
      { id: 'code-practice', title: '代码实践：智能客服 Agent', minutes: 25, content: aiNative07 },
      { id: 'best-practices', title: '最佳实践与经验总结', minutes: 15, content: aiNative08 },
      { id: 'pitfalls', title: '常见坑与排错指南', minutes: 15, content: aiNative09 },
      { id: 'performance', title: '性能优化', minutes: 15, content: aiNative10 },
      { id: 'summary', title: '总结与展望', minutes: 8, content: aiNative11 },
      { id: 'references', title: '参考资料与延伸阅读', minutes: 7, content: aiNative12 }
    ]
  }
]

export function getCourse(courseId) {
  return courses.find((c) => c.id === courseId)
}

export function getChapter(courseId, chapterId) {
  const course = getCourse(courseId)
  if (!course) return null
  const index = course.chapters.findIndex((ch) => ch.id === chapterId)
  if (index === -1) return null
  return {
    course,
    chapter: course.chapters[index],
    index,
    prev: index > 0 ? course.chapters[index - 1] : null,
    next: index < course.chapters.length - 1 ? course.chapters[index + 1] : null
  }
}

export const totalChapterCount = courses.reduce((sum, c) => sum + c.chapters.length, 0)
