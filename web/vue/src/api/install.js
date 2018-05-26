import httpClient from '../utils/httpClient'

export default {
  store (data, callback) {
    httpClient.post('/install/store', data, callback)
  },
  status (callback) {
    httpClient.get('/install/status', {}, callback)
  }
}
