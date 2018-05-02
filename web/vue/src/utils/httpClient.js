import axios from 'axios'
import {Loading, Message} from 'element-ui'
import userStorage from '../storage/user'
import router from '../router/index'

let loadingInstance = null
const errorMessage = '加载失败, 请稍后再试'
const SUCCESS_CODE = 0
const AUTH_ERROR_CODE = 401
const APP_NOT_INSTALL_CODE = 801

axios.defaults.baseURL = '/api'
axios.defaults.timeout = 5000
axios.defaults.responseType = 'json'
axios.interceptors.request.use(config => {
  loadingInstance = Loading.service({fullscreen: true})
  const token = userStorage.getToken()
  config.headers = {
    'Auth-Token': token
  }

  return config
}, error => {
  loadingInstance.close()
  Message.error({
    message: errorMessage
  })

  return Promise.reject(error)
})

axios.interceptors.response.use(data => {
  loadingInstance.close()
  return data
}, error => {
  loadingInstance.close()
  Message.error({
    message: errorMessage
  })

  return Promise.reject(error)
})

function handle (promise, next) {
  promise.then((res) => successCallback(res, next))
    .catch((error) => failureCallback(error))
}

function checkResponseCode (code, msg) {
  switch (code) {
    // 应用未安装
    case APP_NOT_INSTALL_CODE:
      router.push('/install')
      return false
    // 认证失败
    case AUTH_ERROR_CODE:
      router.push('/user/login')
      return false
  }
  if (code !== SUCCESS_CODE) {
    Message.error({
      message: msg
    })
    return false
  }

  return true
}

function successCallback (res, next) {
  if (!checkResponseCode(res.data.code)) {
    return
  }
  if (!next) {
    return
  }
  next(res.data.data, res.data.code, res.data.message)
}

function failureCallback (error) {
  Message.error({
    message: '请求失败 - ' + error
  })
}

export default {
  get (uri, next) {
    const promise = axios.get(uri)
    handle(promise, next)
  },

  batchGet (uris, next) {
    const requests = []
    for (let uri of uris) {
      requests.push(axios.get(uri))
    }
    axios.all(requests).then(axios.spread(function (...res) {
      const result = []
      for (let item of res) {
        if (!checkResponseCode(item.data.code, item.data.message)) {
          return
        }
        result.push(item.data.data)
      }
      next(...result)
    })).catch((error) => failureCallback(error))
  },

  post (uri, data, next) {
    const promise = axios.post(uri, data)
    handle(promise, next)
  }
}
