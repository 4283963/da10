import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error.response?.data?.error || error.message)
  }
)

export const vehicleAPI = {
  getByVIN: (vin) => api.get(`/vehicles/vin/${vin}`),
  create: (data) => api.post('/vehicles', data),
  list: () => api.get('/vehicles'),
}

export const nodeAPI = {
  update: (id, data) => api.put(`/nodes/${id}`, data),
  getProgress: (vehicleId) => api.get(`/nodes/vehicle/${vehicleId}/progress`),
}

export const expenseAPI = {
  create: (data) => api.post('/expenses', data),
  getStats: (vehicleId) => api.get(`/expenses/vehicle/${vehicleId}/stats`),
}

export const uploadAPI = {
  uploadImage: (file) => {
    const formData = new FormData()
    formData.append('file', file)
    return axios.post('/api/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }).then(res => res.data)
  },
}

export default api
