import httpClient from '../utils/httpClient'

export default {
  list (query, callback) {
    httpClient.get('/user', {}, callback)
  },

  detail (id, callback) {
    httpClient.get(`/user/${id}`, {}, callback)
  },

  update (data, callback, failureCallback) {
    httpClient.post('/user/store', data, callback, failureCallback)
  },

  login (username, password, callback, failureCallback) {
    httpClient.post('/user/login', {username, password}, callback, failureCallback)
  },

  enable (id, callback, failureCallback) {
    httpClient.post(`/user/enable/${id}`, {}, callback, failureCallback)
  },

  disable (id, callback, failureCallback) {
    httpClient.post(`/user/disable/${id}`, {}, callback, failureCallback)
  },

  remove (id, callback, failureCallback) {
    httpClient.post(`/user/remove/${id}`, {}, callback, failureCallback)
  },

  editPassword (data, callback, failureCallback) {
    httpClient.post(`/user/editPassword/${data.id}`, {
      'new_password': data.new_password,
      'confirm_new_password': data.confirm_new_password
    }, callback, failureCallback)
  },

  editMyPassword (data, callback, failureCallback) {
    httpClient.post(`/user/editMyPassword`, data, callback, failureCallback)
  }
}
