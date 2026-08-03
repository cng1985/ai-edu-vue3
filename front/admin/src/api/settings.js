import api from './index.js'

export const settingsApi = {
  get: () => api.get('/settings'),
  resolve: (virtualModel) => api.get('/settings/resolve', { params: { virtualModel } }),
  setDefaultVirtualModel: (code) => api.put('/settings/default-virtual-model', { code }),
  createProvider: (data) => api.post('/settings/providers', data),
  updateProvider: (id, data) => api.put(`/settings/providers/${id}`, data),
  quickSetup: (data) => api.post('/settings/quick-setup', data),
  reindexKnowledge: () => api.post('/settings/knowledge/reindex')
}
