import axios from 'axios'

const API_BASE_URL = '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Services
export const servicesAPI = {
  list: () => api.get('/services'),
  listInNamespace: (namespace: string) => api.get(`/services/${namespace}`),
}

// Namespaces
export const namespacesAPI = {
  list: () => api.get('/namespaces'),
}

// VirtualServices
export const virtualServicesAPI = {
  list: (namespace: string) => api.get(`/virtualservices/${namespace}`),
  get: (namespace: string, name: string) => api.get(`/virtualservices/${namespace}/${name}`),
  create: (data: any) => api.post('/virtualservices', data),
  update: (namespace: string, name: string, data: any) => api.put(`/virtualservices/${namespace}/${name}`, data),
  delete: (namespace: string, name: string) => api.delete(`/virtualservices/${namespace}/${name}`),
}

// DestinationRules
export const destinationRulesAPI = {
  list: (namespace: string) => api.get(`/destinationrules/${namespace}`),
  get: (namespace: string, name: string) => api.get(`/destinationrules/${namespace}/${name}`),
  create: (data: any) => api.post('/destinationrules', data),
  update: (namespace: string, name: string, data: any) => api.put(`/destinationrules/${namespace}/${name}`, data),
  delete: (namespace: string, name: string) => api.delete(`/destinationrules/${namespace}/${name}`),
}

// Gateways
export const gatewaysAPI = {
  list: (namespace: string) => api.get(`/gateways/${namespace}`),
  create: (data: any) => api.post('/gateways', data),
  delete: (namespace: string, name: string) => api.delete(`/gateways/${namespace}/${name}`),
}

// PeerAuthentications
export const peerAuthAPI = {
  list: (namespace: string) => api.get(`/peerauthentications/${namespace}`),
  create: (data: any) => api.post('/peerauthentications', data),
  delete: (namespace: string, name: string) => api.delete(`/peerauthentications/${namespace}/${name}`),
}

// AuthorizationPolicies
export const authzPolicyAPI = {
  list: (namespace: string) => api.get(`/authorizationpolicies/${namespace}`),
  create: (data: any) => api.post('/authorizationpolicies', data),
  delete: (namespace: string, name: string) => api.delete(`/authorizationpolicies/${namespace}/${name}`),
}

// Traffic Management
export const trafficAPI = {
  updateWeights: (data: any) => api.post('/traffic/weights', data),
  canary: (data: any) => api.post('/traffic/canary', data),
}

// Snapshots
export const snapshotsAPI = {
  list: () => api.get('/snapshots'),
  get: (id: number) => api.get(`/snapshots/${id}`),
  create: (data: any) => api.post('/snapshots', data),
  delete: (id: number) => api.delete(`/snapshots/${id}`),
  restore: (id: number) => api.post(`/snapshots/${id}/restore`),
}

// Scheduled Actions
export const scheduledActionsAPI = {
  list: () => api.get('/scheduled-actions'),
  get: (id: number) => api.get(`/scheduled-actions/${id}`),
  create: (data: any) => api.post('/scheduled-actions', data),
  update: (id: number, data: any) => api.put(`/scheduled-actions/${id}`, data),
  delete: (id: number) => api.delete(`/scheduled-actions/${id}`),
}

// Audit Logs
export const auditLogsAPI = {
  list: (limit?: number) => api.get('/audit-logs', { params: { limit } }),
}

// YAML Preview
export const yamlAPI = {
  preview: (type: string, spec: any) => api.post('/yaml/preview', { type, spec }),
}

export default api
