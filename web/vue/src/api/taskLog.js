import httpClient from '../utils/httpClient'

export default {
  list (query, callback) {
    httpClient.get('/task/log', query, callback)
  },

  clear (taskId, callback) {
    httpClient.post('/task/log/clear', { task_id: taskId }, callback)
  },

  stop (id, taskId, callback) {
    httpClient.post('/task/log/stop', {id, task_id: taskId}, callback)
  }
}
