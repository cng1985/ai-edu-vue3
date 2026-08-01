import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

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
      { path: 'reviews', name: 'reviews', component: () => import('../views/ReviewView.vue'), meta: { title: '内容审核', roles: ['admin', 'reviewer'] } }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.isLoggedIn) {
    return { name: 'dashboard' }
  }
  if (to.meta.roles && !to.meta.roles.includes(auth.user?.role)) {
    return { name: 'dashboard' }
  }
})

export default router
