import httpClient from '../utils/httpClient'

export default {
  // 任务列表
  list (query, callback) {
    httpClient.get(`/process`, query, callback)
  },
  store (data, callback) {
    httpClient.post('/process/store', data, callback)
  },
  get (id, callback) {
    httpClient.get(`/process/${id}`, {}, callback)
  },
  start (id, callback) {
    httpClient.post(`/process/start/${id}`, {}, callback)
  },
  stop (id, callback) {
    httpClient.post(`/process/stop/${id}`, {}, callback)
  },
  restart (id, callback) {
    httpClient.post(`/process/stop/${id}`, {}, callback)
  },
  enable (id, callback) {
    httpClient.instance
        .post(`/process/enable/${id}`)
        .then(callback)
  },

  disable (id, callback) {
    httpClient.instance
        .post(`/process/disable/${id}`)
        .then(callback)
  }
}
