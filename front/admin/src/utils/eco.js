/** 成长生态展示常量（与学习端 front/app/src/utils/eco.js 保持一致） */

export const LEVELS = ['未掌握', '了解', '可以辅助完成', '可以独立完成', '可以解决复杂问题', '可以设计和指导他人']
export const LEVEL_COLORS = ['#b6bad0', '#8fb8f5', '#2a8cf4', '#6b5cff', '#0fb981', '#f59e0b']

export const RELATION_TYPES = {
  prerequisite: '前置知识',
  related: '关联知识',
  similar: '相似知识',
  advanced: '高级知识'
}

export const RESOURCE_TYPES = {
  article: '文章',
  course: '课程',
  video: '视频',
  code: '代码',
  case: '案例',
  project: '项目',
  task: '任务',
  prompt: 'Prompt',
  sop: 'SOP'
}

export const POST_TYPES = { question: '问答', discussion: '讨论', article: '文章', share: '分享' }

export const OPPORTUNITY_STATUS = {
  open: { label: '招募中', type: 'success' },
  in_progress: { label: '进行中', type: 'warning' },
  closed: { label: '已结束', type: 'info' }
}

export const APPLICATION_STATUS = {
  applied: { label: '待处理', type: 'primary' },
  accepted: { label: '进行中', type: 'warning' },
  rejected: { label: '未录用', type: 'info' },
  completed: { label: '已验收', type: 'success' }
}

export const ROLE_LEVELS = ['初级', '中级', '高级', '专家']

export const money = (v) => '¥' + Number(v || 0).toLocaleString('zh-CN', { maximumFractionDigits: 0 })
export const pct = (v) => Math.round((v || 0) * 100)
export const dateTime = (ms) => (ms ? new Date(ms).toLocaleString('zh-CN', { hour12: false }) : '—')
