import api from './index.js'

export const aiModelsApi = {
  overview: () => api.get('/ai-models/overview'),
  setDefault: (code) => api.put('/ai-models/default', { code }),
  resolve: (virtualModel) => api.get('/ai-models/resolve', { params: { virtualModel } }),

  listCanonicalModels: (params) => api.get('/ai-models/canonical-models', { params }),
  createCanonicalModel: (data) => api.post('/ai-models/canonical-models', data),
  updateCanonicalModel: (id, data) => api.put(`/ai-models/canonical-models/${id}`, data),
  deleteCanonicalModel: (id) => api.delete(`/ai-models/canonical-models/${id}`),

  listCapabilities: (params) => api.get('/ai-models/capabilities', { params }),
  createCapability: (data) => api.post('/ai-models/capabilities', data),
  updateCapability: (id, data) => api.put(`/ai-models/capabilities/${id}`, data),
  deleteCapability: (id) => api.delete(`/ai-models/capabilities/${id}`),

  listCapabilityModels: (params) => api.get('/ai-models/capability-models', { params }),
  createCapabilityModel: (data) => api.post('/ai-models/capability-models', data),
  deleteCapabilityModel: (id) => api.delete(`/ai-models/capability-models/${id}`),

  listProviders: (params) => api.get('/ai-models/providers', { params }),
  getProvider: (id) => api.get(`/ai-models/providers/${id}`),
  createProvider: (data) => api.post('/ai-models/providers', data),
  updateProvider: (id, data) => api.put(`/ai-models/providers/${id}`, data),
  deleteProvider: (id) => api.delete(`/ai-models/providers/${id}`),

  listProviderModels: (params) => api.get('/ai-models/provider-models', { params }),
  createProviderModel: (data) => api.post('/ai-models/provider-models', data),
  updateProviderModel: (id, data) => api.put(`/ai-models/provider-models/${id}`, data),
  deleteProviderModel: (id) => api.delete(`/ai-models/provider-models/${id}`),

  listVirtualModels: (params) => api.get('/ai-models/virtual-models', { params }),
  createVirtualModel: (data) => api.post('/ai-models/virtual-models', data),
  updateVirtualModel: (id, data) => api.put(`/ai-models/virtual-models/${id}`, data),
  deleteVirtualModel: (id) => api.delete(`/ai-models/virtual-models/${id}`),

  listVirtualModelMappings: (params) => api.get('/ai-models/virtual-model-mappings', { params }),
  createVirtualModelMapping: (data) => api.post('/ai-models/virtual-model-mappings', data),
  updateVirtualModelMapping: (id, data) => api.put(`/ai-models/virtual-model-mappings/${id}`, data),
  deleteVirtualModelMapping: (id) => api.delete(`/ai-models/virtual-model-mappings/${id}`)
}
