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

/** 成长生态治理 */
export const ecoAdminApi = {
  dashboard: () => api.get('/manage/dashboard'),
  careers: () => api.get('/manage/careers'),
  saveCareer: (data, id) => (id ? api.put(`/manage/careers/${id}`, data) : api.post('/manage/careers', data)),
  removeCareer: (id) => api.delete(`/manage/careers/${id}`),
  role: (id) => api.get(`/manage/roles/${id}`),
  saveRole: (data, id) => (id ? api.put(`/manage/roles/${id}`, data) : api.post('/manage/roles', data)),
  removeRole: (id) => api.delete(`/manage/roles/${id}`),
  saveCapability: (data, id) => (id ? api.put(`/manage/capabilities/${id}`, data) : api.post('/manage/capabilities', data)),
  removeCapability: (id) => api.delete(`/manage/capabilities/${id}`),
  skills: () => api.get('/manage/skills'),
  saveSkill: (data, id) => (id ? api.put(`/manage/skills/${id}`, data) : api.post('/manage/skills', data)),
  removeSkill: (id) => api.delete(`/manage/skills/${id}`),
  knowledge: () => api.get('/manage/knowledge'),
  saveKnowledge: (data, id) => (id ? api.put(`/manage/knowledge/${id}`, data) : api.post('/manage/knowledge', data)),
  removeKnowledge: (id) => api.delete(`/manage/knowledge/${id}`),
  relations: () => api.get('/manage/relations'),
  addRelation: (data) => api.post('/manage/relations', data),
  removeRelation: (id) => api.delete(`/manage/relations/${id}`),
  chapterMappings: (courseId) => api.get('/manage/chapter-knowledge', { params: { courseId } }),
  setChapterKnowledge: (courseId, chapterId, knowledgeIds) => api.put(`/manage/chapter-knowledge/${courseId}/${chapterId}`, { knowledgeIds }),
  projects: () => api.get('/manage/projects'),
  project: (id) => api.get(`/manage/projects/${id}`),
  saveProject: (data, id) => (id ? api.put(`/manage/projects/${id}`, data) : api.post('/manage/projects', data)),
  removeProject: (id) => api.delete(`/manage/projects/${id}`),
  submissions: (params) => api.get('/manage/submissions', { params }),
  opportunities: (params) => api.get('/manage/opportunities', { params }),
  updateOpportunity: (id, data) => api.put(`/manage/opportunities/${id}`, data),
  removeOpportunity: (id) => api.delete(`/manage/opportunities/${id}`),
  candidates: (id) => api.get(`/manage/opportunities/${id}/candidates`),
  decide: (applicationId, data) => api.post(`/manage/applications/${applicationId}/decision`, data),
  talents: (keyword) => api.get('/manage/talents', { params: { keyword } }),
  talent: (id) => api.get(`/manage/talents/${id}`),
  posts: (params) => api.get('/manage/posts', { params }),
  removePost: (id) => api.delete(`/manage/posts/${id}`),
  resources: (params) => api.get('/manage/resources', { params }),
  removeResource: (id) => api.delete(`/manage/resources/${id}`),
  kernelRuns: (params) => api.get('/manage/kernel/runs', { params }),
  kernelStages: () => api.get('/manage/kernel/stages'),
  agents: () => api.get('/manage/agents')
}

export { settingsApi } from './settings.js'
