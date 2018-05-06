import httpClient from '../utils/httpClient'

export default {
  list (query, callback) {
    httpClient.get('/user', {}, callback)
  },

  detail (id, callback) {
    httpClient.get(`/user/${id}`, {}, callback)
  },

  update (data, callback) {
    httpClient.post('/user/store', data, callback)
  },

  login (username, password, callback) {
    httpClient.post('/user/login', {username, password}, callback)
  },

  enable (id, callback) {
    httpClient.post(`/user/enable/${id}`, {}, callback)
  },

  disable (id, callback) {
    httpClient.post(`/user/disable/${id}`, {}, callback)
  },

  remove (id, callback) {
    httpClient.post(`/user/remove/${id}`, {}, callback)
  },

  editPassword (data, callback) {
    httpClient.post(`/user/editPassword/${data.id}`, {
      'new_password': data.new_password,
      'confirm_new_password': data.confirm_new_password
    }, callback)
  },

  editMyPassword (data, callback) {
    httpClient.post(`/user/editMyPassword`, data, callback)
  }
}
