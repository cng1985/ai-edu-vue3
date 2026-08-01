import api from './index.js'

export const rolesApi = {
  list: () => api.get('/roles'),
  listPermissions: () => api.get('/permissions'),
  update: (role, permissions) => api.put(`/roles/${role}`, { permissions })
}
