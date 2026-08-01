import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { title: '登录', public: true, blank: true }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/RegisterView.vue'),
    meta: { title: '注册', public: true, blank: true }
  },
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue'),
    meta: { title: '学习驾驶舱' }
  },
  {
    path: '/career',
    name: 'career',
    component: () => import('../views/CareerView.vue'),
    meta: { title: 'AI 职业规划' }
  },
  {
    path: '/path',
    name: 'learning-path',
    component: () => import('../views/LearningPathView.vue'),
    meta: { title: '学习路径' }
  },
  {
    path: '/micro/:unitId',
    name: 'micro-unit',
    component: () => import('../views/MicroUnitView.vue'),
    meta: { title: '5 分钟微学习' }
  },
  {
    path: '/review',
    name: 'review',
    component: () => import('../views/ReviewView.vue'),
    meta: { title: '复习与补强' }
  },
  {
    path: '/incentives',
    name: 'incentives',
    component: () => import('../views/IncentiveView.vue'),
    meta: { title: '成长激励' }
  },
  {
    path: '/courses',
    name: 'courses',
    component: () => import('../views/CoursesView.vue'),
    meta: { title: '全部课程' }
  },
  {
    path: '/courses/:courseId',
    name: 'course-detail',
    component: () => import('../views/CourseDetailView.vue'),
    meta: { title: '课程详情' }
  },
  {
    path: '/courses/:courseId/:chapterId',
    name: 'lesson',
    component: () => import('../views/LessonView.vue'),
    meta: { title: '章节学习' }
  },
  {
    path: '/chat',
    name: 'chat',
    component: () => import('../views/ChatView.vue'),
    meta: { title: 'AI 学习助手' }
  },
  {
    path: '/quiz',
    name: 'quiz-list',
    component: () => import('../views/QuizListView.vue'),
    meta: { title: '知识测验' }
  },
  {
    path: '/quiz/:quizId',
    name: 'quiz',
    component: () => import('../views/QuizView.vue'),
    meta: { title: '测验' }
  },
  {
    path: '/stats',
    name: 'stats',
    component: () => import('../views/StatsView.vue'),
    meta: { title: '学习统计' }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  // 已登录用户访问登录/注册页时直接回首页
  if (to.meta.public && auth.isLoggedIn) {
    return { path: '/' }
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} · AI 学习系统` : 'AI 学习系统'
})

export default router
