import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue'),
    meta: { title: '学习中心' }
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

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} · AI 学习系统` : 'AI 学习系统'
})

export default router
