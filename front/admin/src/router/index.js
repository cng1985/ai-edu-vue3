import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import { getFirstAccessibleRoute } from '../config/nav'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { public: true, title: '登录' }
  },
  {
    path: '/',
    component: () => import('../layouts/AdminLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('../views/DashboardView.vue'),
        meta: { title: '数据看板', permissions: [PERM.DASHBOARD] }
      },
      {
        path: 'users',
        name: 'users',
        component: () => import('../views/UsersView.vue'),
        meta: { title: '用户管理', permissions: [PERM.USER_READ] }
      },
      {
        path: 'courses',
        name: 'courses',
        component: () => import('../views/CoursesView.vue'),
        meta: { title: '课程管理', permissions: [PERM.COURSE_READ] }
      },
      {
        path: 'courses/:id',
        name: 'course-edit',
        component: () => import('../views/CourseEditView.vue'),
        meta: { title: '课程编辑', permissions: [PERM.COURSE_READ] }
      },
      {
        path: 'quizzes',
        name: 'quizzes',
        component: () => import('../views/QuizzesView.vue'),
        meta: { title: '题库管理', permissions: [PERM.QUIZ_READ] }
      },
      {
        path: 'quizzes/:id',
        name: 'quiz-edit',
        component: () => import('../views/QuizEditView.vue'),
        meta: { title: '测验编辑', permissions: [PERM.QUIZ_READ] }
      },
      {
        path: 'reviews',
        name: 'reviews',
        component: () => import('../views/ReviewView.vue'),
        meta: { title: '内容审核', permissions: [PERM.REVIEW_READ] }
      },
      {
        path: 'customers',
        name: 'customers',
        component: () => import('../views/CustomersView.vue'),
        meta: { title: '客户咨询', permissions: [PERM.CUSTOMER_READ] }
      },
      {
        path: 'documents',
        name: 'documents',
        component: () => import('../views/DocumentsView.vue'),
        meta: { title: '单据管理', permissions: [PERM.DOCUMENT_READ] }
      },
      {
        path: 'knowledge',
        name: 'knowledge',
        component: () => import('../views/KnowledgeView.vue'),
        meta: { title: '知识库管理', permissions: [PERM.KNOWLEDGE_READ] }
      },
      {
        path: 'roles',
        name: 'roles',
        component: () => import('../views/RolesView.vue'),
        meta: { title: '权限管理', permissions: [PERM.ROLE_MANAGE] }
      },
      {
        path: 'ai-models',
        name: 'ai-models',
        component: () => import('../views/AiModelsView.vue'),
        meta: { title: 'AI 大模型配置', permissions: [PERM.AI_MODEL_READ] }
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('../views/SettingsView.vue'),
        meta: { title: '系统设置', permissions: [PERM.SETTINGS_MANAGE] }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.isLoggedIn) {
    return getFirstAccessibleRoute(auth)
  }
  if (to.meta.permissions && !auth.hasAnyPermission(to.meta.permissions)) {
    return getFirstAccessibleRoute(auth)
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} - AI 学习管理` : 'AI 学习管理'
})

export default router
