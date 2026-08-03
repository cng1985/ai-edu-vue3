import api from './index.js'

export const knowledgeApi = {
  status: () => api.get('/knowledge/status'),
  listChunks: (params) => api.get('/knowledge/chunks', { params }),
  reindex: () => api.post('/knowledge/reindex'),
  search: (q) => api.get('/knowledge/search', { params: { q } })
}
