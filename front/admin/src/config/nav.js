import { PERM } from '../constants/permissions'

/** 后台导航统一配置 — 侧边栏、顶栏、命令面板、面包屑、路由守卫共用 */
export const NAV_GROUPS = [
  {
    key: 'operations',
    title: '运营管理',
    items: [
      {
        path: '/dashboard',
        name: 'dashboard',
        title: '数据看板',
        icon: 'DataAnalysis',
        permission: PERM.DASHBOARD,
        keywords: ['看板', '统计', 'dashboard']
      }
    ]
  },
  {
    key: 'users',
    title: '用户管理',
    items: [
      {
        path: '/users',
        name: 'users',
        title: '用户管理',
        icon: 'User',
        permission: PERM.USER_READ,
        keywords: ['用户', '账号', 'users']
      }
    ]
  },
  {
    key: 'content',
    title: '内容管理',
    items: [
      {
        path: '/courses',
        name: 'courses',
        title: '课程管理',
        icon: 'Reading',
        permission: PERM.COURSE_READ,
        activePrefix: '/courses',
        keywords: ['课程', 'course']
      },
      {
        path: '/quizzes',
        name: 'quizzes',
        title: '题库管理',
        icon: 'EditPen',
        permission: PERM.QUIZ_READ,
        activePrefix: '/quizzes',
        keywords: ['题库', '测验', 'quiz']
      },
      {
        path: '/reviews',
        name: 'reviews',
        title: '内容审核',
        icon: 'DocumentChecked',
        permission: PERM.REVIEW_READ,
        keywords: ['审核', 'review']
      }
    ]
  },
  {
    key: 'service',
    title: '客户服务',
    items: [
      {
        path: '/customers',
        name: 'customers',
        title: '客户咨询',
        icon: 'Service',
        permission: PERM.CUSTOMER_READ,
        keywords: ['客户', '咨询', '客服']
      }
    ]
  },
  {
    key: 'documents',
    title: '单据管理',
    items: [
      {
        path: '/documents',
        name: 'documents',
        title: '单据管理',
        icon: 'Ticket',
        permission: PERM.DOCUMENT_READ,
        keywords: ['单据', 'document']
      }
    ]
  },
  {
    key: 'knowledge',
    title: '知识库',
    items: [
      {
        path: '/knowledge',
        name: 'knowledge',
        title: '知识库',
        icon: 'Collection',
        permission: PERM.KNOWLEDGE_READ,
        keywords: ['知识', 'knowledge']
      }
    ]
  },
  {
    key: 'ai',
    title: 'AI 服务',
    items: [
      {
        path: '/chatgpt',
        name: 'chatgpt',
        title: 'ChatGPT 对话',
        icon: 'ChatDotRound',
        permission: PERM.AI_CHAT,
        keywords: ['ChatGPT', '对话', '聊天', 'AI']
      }
    ]
  },
  {
    key: 'system',
    title: '系统管理',
    items: [
      {
        path: '/roles',
        name: 'roles',
        title: '权限管理',
        icon: 'Lock',
        permission: PERM.ROLE_MANAGE,
        keywords: ['权限', '角色', 'role']
      },
      {
        path: '/settings',
        name: 'settings',
        title: '系统设置',
        icon: 'Setting',
        permission: PERM.SETTINGS_MANAGE,
        keywords: ['设置', '配置', 'settings']
      }
    ]
  }
]

/** 子页面面包屑父级映射 */
export const NAV_CHILDREN = {
  'course-edit': { parent: 'courses', title: '课程编辑' },
  'quiz-edit': { parent: 'quizzes', title: '测验编辑' }
}

/** 扁平化所有导航项 */
export function flattenNavItems(groups = NAV_GROUPS) {
  return groups.flatMap((group) =>
    group.items.map((item) => ({ ...item, group: group.key, groupTitle: group.title }))
  )
}

/** 按权限过滤后的首个可访问路由 */
export function getFirstAccessibleRoute(auth, groups = NAV_GROUPS) {
  for (const item of flattenNavItems(groups)) {
    if (!item.permission || auth.hasPermission(item.permission)) {
      return { name: item.name }
    }
  }
  return { name: 'login' }
}
