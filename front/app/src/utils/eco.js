/** 成长生态的展示常量与格式化工具 */

export const LEVELS = [
  { level: 0, name: '未掌握', color: '#b6bad0' },
  { level: 1, name: '了解', color: '#8fb8f5' },
  { level: 2, name: '可以辅助完成', color: '#2a8cf4' },
  { level: 3, name: '可以独立完成', color: '#6b5cff' },
  { level: 4, name: '可以解决复杂问题', color: '#0fb981' },
  { level: 5, name: '可以设计和指导他人', color: '#f59e0b' }
]

export const levelInfo = (l) => LEVELS[Math.max(0, Math.min(5, l || 0))]

export const KNOWLEDGE_STATUS = {
  new: { label: '未学习', tone: 'neutral' },
  learning: { label: '学习中', tone: 'info' },
  mastered: { label: '已掌握', tone: 'success' }
}

export const RELATION_TYPES = {
  prerequisite: '前置知识',
  related: '关联知识',
  similar: '相似知识',
  advanced: '高级知识'
}

export const POST_TYPES = {
  question: { label: '问答', icon: 'message' },
  discussion: { label: '讨论', icon: 'users' },
  article: { label: '文章', icon: 'note' },
  share: { label: '分享', icon: 'gift' }
}

export const RESOURCE_TYPES = {
  article: { label: '文章', icon: 'note' },
  course: { label: '课程', icon: 'book' },
  video: { label: '视频', icon: 'play' },
  code: { label: '代码', icon: 'code' },
  case: { label: '案例', icon: 'flag' },
  project: { label: '项目', icon: 'layers' },
  task: { label: '任务', icon: 'target' },
  prompt: { label: 'Prompt', icon: 'sparkles' },
  sop: { label: 'SOP', icon: 'clipboard' }
}

export const APPLICATION_STATUS = {
  applied: { label: '待处理', tone: 'info' },
  accepted: { label: '进行中', tone: 'warning' },
  rejected: { label: '未录用', tone: 'neutral' },
  completed: { label: '已验收', tone: 'success' }
}

export const OPPORTUNITY_STATUS = {
  open: { label: '招募中', tone: 'success' },
  in_progress: { label: '进行中', tone: 'warning' },
  closed: { label: '已结束', tone: 'neutral' }
}

export const ROLE_NAMES = {
  learner: '学习者',
  creator: '创作者',
  enterprise: '企业',
  admin: '管理员',
  operator: '运营',
  reviewer: '审核员',
  guest: '游客'
}

export const FLYWHEEL = [
  { key: 'learn', label: '学习', icon: 'book' },
  { key: 'practice', label: '实践', icon: 'code' },
  { key: 'output', label: '产出', icon: 'trophy' },
  { key: 'income', label: '收益', icon: 'gift' },
  { key: 'invest', label: '投资学习', icon: 'refresh' },
  { key: 'stronger', label: '更强能力', icon: 'trend' }
]

export const pct = (v) => Math.round((v || 0) * 100)

export function money(v) {
  return '¥' + Number(v || 0).toLocaleString('zh-CN', { maximumFractionDigits: 0 })
}

export function timeAgo(ms) {
  if (!ms) return ''
  const diff = Date.now() - ms
  const m = Math.floor(diff / 60000)
  if (m < 1) return '刚刚'
  if (m < 60) return `${m} 分钟前`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h} 小时前`
  const d = Math.floor(h / 24)
  if (d < 30) return `${d} 天前`
  return new Date(ms).toLocaleDateString('zh-CN')
}

export function dateText(ms) {
  return ms ? new Date(ms).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' }) : ''
}

export function masteryColor(m) {
  if (!m) return '#d5d9e6'
  if (m >= 0.8) return '#0fb981'
  if (m >= 0.5) return '#6b5cff'
  if (m >= 0.25) return '#2a8cf4'
  return '#f59e0b'
}
