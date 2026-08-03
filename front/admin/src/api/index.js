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
    if (res.config.responseType === 'blob') {
      return res
    }
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
      localStorage.removeItem('admin-permissions')
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
  login: (username, password, portal) => api.post('/auth/login', { username, password, portal }),
  me: () => api.get('/auth/me'),
  permissions: () => api.get('/auth/permissions'),
  refreshPermissions: () => api.post('/auth/permissions/refresh')
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

export const customersApi = {
  stats: () => api.get('/customers/stats'),
  listTickets: (params) => api.get('/customers/tickets', { params }),
  getTicket: (id) => api.get(`/customers/tickets/${id}`),
  listMessages: (id, params) => api.get(`/customers/tickets/${id}/messages`, { params }),
  reply: (id, content) => api.post(`/customers/tickets/${id}/reply`, { content }),
  updateStatus: (id, status) => api.put(`/customers/tickets/${id}/status`, { status })
}

function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

export const documentsApi = {
  list: (params) => api.get('/documents', { params }),
  get: (id) => api.get(`/documents/${id}`),
  create: (data) => api.post('/documents', data),
  update: (id, data) => api.put(`/documents/${id}`, data),
  remove: (id) => api.delete(`/documents/${id}`),
  exportExcel: async (params) => {
    const res = await api.get('/documents/export', { params, responseType: 'blob' })
    const disposition = res.headers?.['content-disposition'] || ''
    const match = disposition.match(/filename\*=UTF-8''(.+)/)
    const filename = match ? decodeURIComponent(match[1]) : '单据导出.xlsx'
    downloadBlob(res.data, filename)
  },
  downloadTemplate: async () => {
    const res = await api.get('/documents/import/template', { responseType: 'blob' })
    downloadBlob(res.data, '单据导入模板.xlsx')
  },
  importExcel: (file) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/documents/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },
  importProgress: (taskId) => api.get(`/documents/import/${taskId}/progress`)
}

export { settingsApi } from './settings.js'
