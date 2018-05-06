import httpClient from '../utils/httpClient'

export default {
  loginLogList (query, callback) {
    httpClient.get('/system/login-log', query, callback)
  }
}
