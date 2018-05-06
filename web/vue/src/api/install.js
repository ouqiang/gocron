import httpClient from '../utils/httpClient'

export default {
  store (data, callback) {
    httpClient.post('/install/store', data, callback)
  }
}
