import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin-token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (res) => {
    const { code, message, data } = res.data
    if (code !== 0) {
      ElMessage.error(message || '请求失败')
      return Promise.reject(new Error(message))
    }
    return data
  },
  (err) => {
    const msg = err.response?.data?.message || err.message || '网络错误'
    if (err.response?.status === 401) {
      localStorage.removeItem('admin-token')
      localStorage.removeItem('admin-user')
      if (!window.location.hash.includes('/login')) {
        window.location.hash = '#/login'
      }
    }
    ElMessage.error(msg)
    return Promise.reject(err)
  }
)

export default api

export const authApi = {
  login: (username, password) => api.post('/auth/login', { username, password }),
  me: () => api.get('/auth/me')
}

export const dashboardApi = {
  stats: () => api.get('/dashboard/stats')
}

export const usersApi = {
  list: (params) => api.get('/users', { params }),
  get: (id) => api.get(`/users/${id}`),
  create: (data) => api.post('/users', data),
  update: (id, data) => api.put(`/users/${id}`, data),
  remove: (id) => api.delete(`/users/${id}`)
}

export const coursesApi = {
  list: (params) => api.get('/courses', { params }),
  get: (id) => api.get(`/courses/${id}`),
  create: (data) => api.post('/courses', data),
  update: (id, data) => api.put(`/courses/${id}`, data),
  remove: (id) => api.delete(`/courses/${id}`),
  addChapter: (courseId, data) => api.post(`/courses/${courseId}/chapters`, data),
  updateChapter: (courseId, chapterId, data) => api.put(`/courses/${courseId}/chapters/${chapterId}`, data),
  removeChapter: (courseId, chapterId) => api.delete(`/courses/${courseId}/chapters/${chapterId}`)
}

export const quizzesApi = {
  list: (params) => api.get('/quizzes', { params }),
  get: (id) => api.get(`/quizzes/${id}`),
  create: (data) => api.post('/quizzes', data),
  update: (id, data) => api.put(`/quizzes/${id}`, data),
  remove: (id) => api.delete(`/quizzes/${id}`)
}

export const reviewsApi = {
  list: (params) => api.get('/reviews', { params }),
  get: (id) => api.get(`/reviews/${id}`),
  approve: (id, comment) => api.post(`/reviews/${id}/approve`, { comment }),
  reject: (id, comment) => api.post(`/reviews/${id}/reject`, { comment })
}

export { settingsApi } from './settings.js'
