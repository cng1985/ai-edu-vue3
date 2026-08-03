import api from './index.js'

export const settingsApi = {
  get: () => api.get('/settings'),
  update: (data) => api.put('/settings', data)
}
