import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// VirtualServices
export const virtualServiceApi = {
  list: (namespace = '') => api.get(`/virtualservices?namespace=${namespace}`),
  get: (namespace, name) => api.get(`/virtualservices/${namespace}/${name}`),
  create: (data) => api.post('/virtualservices', data),
  update: (namespace, name, data) => api.put(`/virtualservices/${namespace}/${name}`, data),
  delete: (namespace, name) => api.delete(`/virtualservices/${namespace}/${name}`),
};

// DestinationRules
export const destinationRuleApi = {
  list: (namespace = '') => api.get(`/destinationrules?namespace=${namespace}`),
  get: (namespace, name) => api.get(`/destinationrules/${namespace}/${name}`),
  create: (data) => api.post('/destinationrules', data),
  update: (namespace, name, data) => api.put(`/destinationrules/${namespace}/${name}`, data),
  delete: (namespace, name) => api.delete(`/destinationrules/${namespace}/${name}`),
};

// Gateways
export const gatewayApi = {
  list: (namespace = '') => api.get(`/gateways?namespace=${namespace}`),
  get: (namespace, name) => api.get(`/gateways/${namespace}/${name}`),
  create: (data) => api.post('/gateways', data),
  update: (namespace, name, data) => api.put(`/gateways/${namespace}/${name}`, data),
  delete: (namespace, name) => api.delete(`/gateways/${namespace}/${name}`),
};

// ServiceEntries
export const serviceEntryApi = {
  list: (namespace = '') => api.get(`/serviceentries?namespace=${namespace}`),
  get: (namespace, name) => api.get(`/serviceentries/${namespace}/${name}`),
  create: (data) => api.post('/serviceentries', data),
  update: (namespace, name, data) => api.put(`/serviceentries/${namespace}/${name}`, data),
  delete: (namespace, name) => api.delete(`/serviceentries/${namespace}/${name}`),
};

// AuthorizationPolicies
export const authorizationPolicyApi = {
  list: (namespace = '') => api.get(`/authorizationpolicies?namespace=${namespace}`),
  get: (namespace, name) => api.get(`/authorizationpolicies/${namespace}/${name}`),
  create: (data) => api.post('/authorizationpolicies', data),
  update: (namespace, name, data) => api.put(`/authorizationpolicies/${namespace}/${name}`, data),
  delete: (namespace, name) => api.delete(`/authorizationpolicies/${namespace}/${name}`),
};

// PeerAuthentications
export const peerAuthenticationApi = {
  list: (namespace = '') => api.get(`/peerauthentications?namespace=${namespace}`),
  get: (namespace, name) => api.get(`/peerauthentications/${namespace}/${name}`),
  create: (data) => api.post('/peerauthentications', data),
  update: (namespace, name, data) => api.put(`/peerauthentications/${namespace}/${name}`, data),
  delete: (namespace, name) => api.delete(`/peerauthentications/${namespace}/${name}`),
};

// RequestAuthentications
export const requestAuthenticationApi = {
  list: (namespace = '') => api.get(`/requestauthentications?namespace=${namespace}`),
  get: (namespace, name) => api.get(`/requestauthentications/${namespace}/${name}`),
  create: (data) => api.post('/requestauthentications', data),
  update: (namespace, name, data) => api.put(`/requestauthentications/${namespace}/${name}`, data),
  delete: (namespace, name) => api.delete(`/requestauthentications/${namespace}/${name}`),
};

// Schedules
export const scheduleApi = {
  list: () => api.get('/schedules'),
  create: (data) => api.post('/schedules', data),
  delete: (id) => api.delete(`/schedules/${id}`),
};

// Kubernetes
export const kubernetesApi = {
  listNamespaces: () => api.get('/namespaces'),
  listServices: (namespace = '') => api.get(`/services?namespace=${namespace}`),
  getService: (namespace, name) => api.get(`/services/${namespace}/${name}`),
};

// Topology
export const topologyApi = {
  get: (namespace = '') => api.get(`/topology?namespace=${namespace}`),
};

export default api;
