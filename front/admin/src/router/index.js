import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    component: () => import('../layouts/AdminLayout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'dashboard', component: () => import('../views/DashboardView.vue'), meta: { title: '数据看板' } },
      { path: 'users', name: 'users', component: () => import('../views/UsersView.vue'), meta: { title: '用户管理', roles: ['admin'] } },
      { path: 'courses', name: 'courses', component: () => import('../views/CoursesView.vue'), meta: { title: '课程管理' } },
      { path: 'courses/:id', name: 'course-edit', component: () => import('../views/CourseEditView.vue'), meta: { title: '课程编辑' } },
      { path: 'quizzes', name: 'quizzes', component: () => import('../views/QuizzesView.vue'), meta: { title: '题库管理' } },
      { path: 'quizzes/:id', name: 'quiz-edit', component: () => import('../views/QuizEditView.vue'), meta: { title: '测验编辑' } },
      { path: 'reviews', name: 'reviews', component: () => import('../views/ReviewView.vue'), meta: { title: '内容审核', roles: ['admin', 'reviewer'] } },
      { path: 'customers', name: 'customers', component: () => import('../views/CustomersView.vue'), meta: { title: '客户咨询', permissions: [PERM.CUSTOMER_READ] } },
      { path: 'documents', name: 'documents', component: () => import('../views/DocumentsView.vue'), meta: { title: '单据管理', permissions: [PERM.DOCUMENT_READ] } },
      { path: 'knowledge', name: 'knowledge', component: () => import('../views/KnowledgeView.vue'), meta: { title: '知识库管理', permissions: [PERM.KNOWLEDGE_READ] } },
      { path: 'roles', name: 'roles', component: () => import('../views/RolesView.vue'), meta: { title: '权限管理', roles: ['admin'] } },
      { path: 'settings', name: 'settings', component: () => import('../views/SettingsView.vue'), meta: { title: '系统设置', roles: ['admin'] } }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

function firstAccessibleRoute(auth) {
  const order = ['dashboard', 'users', 'customers', 'documents', 'knowledge', 'courses', 'quizzes', 'reviews', 'roles']
  for (const name of order) {
    const route = routes[1].children.find((r) => r.name === name)
    if (route?.meta?.permissions && auth.hasAnyPermission(route.meta.permissions)) {
      return { name }
    }
  }
  return { name: 'login' }
}

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.isLoggedIn) {
    return firstAccessibleRoute(auth)
  }
  if (to.meta.permissions && !auth.hasAnyPermission(to.meta.permissions)) {
    return firstAccessibleRoute(auth)
  }
})

export default router
