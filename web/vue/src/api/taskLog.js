import httpClient from '../utils/httpClient'

export default {
  list (query, callback) {
    httpClient.get('/task/log', query, callback)
  },

  clear (callback) {
    httpClient.post('/task/log/clear', {}, callback)
  },

  stop (id, callback) {
    httpClient.post('/task/log/stop', {id}, callback)
  }
}
