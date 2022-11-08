import httpClient from '../utils/httpClient'

export default {
  list(query, callback) {
    httpClient.get('/task/log', query, callback)
  },

  clear(callback) {
    httpClient.post('/task/log/clear', {}, callback)
  },

  result(id, callback) {
    httpClient.post('/task/log/result', { id }, callback)
  },

  stop(id, taskId, callback) {
    httpClient.post('/task/log/stop', { id, task_id: taskId }, callback)
  }
}
