import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'

export const BRAND = '知航'

const routes = [
  { path: '/login', name: 'login', component: () => import('../views/LoginView.vue'), meta: { title: '登录', public: true, blank: true } },
  { path: '/register', name: 'register', component: () => import('../views/RegisterView.vue'), meta: { title: '注册', public: true, blank: true } },

  { path: '/', name: 'growth', component: () => import('../views/GrowthView.vue'), meta: { title: '成长中心', permissions: [PERM.ECO_READ] } },
  { path: '/career', name: 'career', component: () => import('../views/CareerView.vue'), meta: { title: '职业与能力', permissions: [PERM.ECO_READ] } },
  { path: '/knowledge', name: 'knowledge', component: () => import('../views/KnowledgeGraphView.vue'), meta: { title: '知识图谱', permissions: [PERM.ECO_READ] } },
  { path: '/knowledge/:id', name: 'knowledge-detail', component: () => import('../views/KnowledgeDetailView.vue'), meta: { title: '知识点', permissions: [PERM.ECO_READ] } },
  { path: '/kernel', name: 'kernel', component: () => import('../views/KernelView.vue'), meta: { title: 'AI 学习内核', permissions: [PERM.AI_CHAT] } },
  { path: '/agents/:code?', name: 'agents', component: () => import('../views/AgentsView.vue'), meta: { title: 'AI 伙伴', permissions: [PERM.AI_CHAT] } },

  { path: '/courses', name: 'courses', component: () => import('../views/CoursesView.vue'), meta: { title: '课程学习', permissions: [PERM.COURSE_READ] } },
  { path: '/courses/:courseId', name: 'course-detail', component: () => import('../views/CourseDetailView.vue'), meta: { title: '课程详情', permissions: [PERM.COURSE_READ] } },
  { path: '/courses/:courseId/:chapterId', name: 'lesson', component: () => import('../views/LessonView.vue'), meta: { title: '章节学习', permissions: [PERM.COURSE_READ] } },
  { path: '/quiz', name: 'quiz-list', component: () => import('../views/QuizListView.vue'), meta: { title: '知识测验', permissions: [PERM.QUIZ_READ] } },
  { path: '/quiz/:quizId', name: 'quiz', component: () => import('../views/QuizView.vue'), meta: { title: '测验', permissions: [PERM.QUIZ_READ] } },

  { path: '/projects', name: 'projects', component: () => import('../views/ProjectsView.vue'), meta: { title: '项目实践', permissions: [PERM.ECO_READ] } },
  { path: '/projects/:id', name: 'project-detail', component: () => import('../views/ProjectDetailView.vue'), meta: { title: '项目详情', permissions: [PERM.ECO_READ] } },
  { path: '/market', name: 'market', component: () => import('../views/MarketView.vue'), meta: { title: '任务市场', permissions: [PERM.ECO_READ] } },
  { path: '/market/:id', name: 'opportunity', component: () => import('../views/OpportunityDetailView.vue'), meta: { title: '任务详情', permissions: [PERM.ECO_READ] } },
  { path: '/enterprise', name: 'enterprise', component: () => import('../views/EnterpriseView.vue'), meta: { title: '企业工作台', permissions: [PERM.OPPORTUNITY_PUBLISH] } },
  { path: '/talents', name: 'talents', component: () => import('../views/TalentsView.vue'), meta: { title: '人才库', permissions: [PERM.TALENT_READ] } },
  { path: '/talents/:id', name: 'talent', component: () => import('../views/ProfileView.vue'), meta: { title: '人才画像', permissions: [PERM.TALENT_READ] } },

  { path: '/community', name: 'community', component: () => import('../views/CommunityView.vue'), meta: { title: '知识社区', permissions: [PERM.ECO_READ] } },
  { path: '/community/:id', name: 'post', component: () => import('../views/PostDetailView.vue'), meta: { title: '社区内容', permissions: [PERM.ECO_READ] } },
  { path: '/resources', name: 'resources', component: () => import('../views/ResourcesView.vue'), meta: { title: '知识资产', permissions: [PERM.ECO_READ] } },

  { path: '/profile', name: 'profile', component: () => import('../views/ProfileView.vue'), meta: { title: '我的人才画像', permissions: [PERM.ECO_READ] } },
  { path: '/support', name: 'support', component: () => import('../views/CustomerSupportView.vue'), meta: { title: '客户咨询', permissions: [PERM.CUSTOMER_CHAT] } },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

/** 不同参与者的默认首页：企业进入工作台，其余进入成长中心 */
export function homeFor(auth) {
  return auth.user?.role === 'enterprise' ? '/enterprise' : '/'
}

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.meta.public && auth.isLoggedIn) {
    return { path: homeFor(auth) }
  }
  if (to.path === '/' && homeFor(auth) !== '/') {
    return { path: homeFor(auth) }
  }
  if (to.meta.permissions && !auth.hasAnyPermission(to.meta.permissions)) {
    return to.path === '/' ? { path: '/courses' } : { path: homeFor(auth), query: { denied: String(to.name) } }
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} · ${BRAND}` : `${BRAND} · AI 成长生态`
})

export default router
