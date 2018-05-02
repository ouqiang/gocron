class User {
  getToken () {
    return sessionStorage.getItem('token') || ''
  }

  setToken (token) {
    sessionStorage.setItem('token', token)
    return this
  }

  isLogin () {
    return this.getToken() !== ''
  }

  clear () {
    sessionStorage.clear()
  }

  getUid (uid) {
    return sessionStorage.getItem('uid') || ''
  }

  setUid (uid) {
    sessionStorage.setItem('uid', uid)
    return this
  }

  getUsername () {
    return sessionStorage.getItem('username') || ''
  }

  setUsername (username) {
    sessionStorage.setItem('username', username)
    return this
  }

  getIsAdmin () {
    return sessionStorage.getItem('is_admin') || 0
  }

  setIsAdmin (isAdmin) {
    sessionStorage.setItem('is_admin', isAdmin)
    return this
  }
}

export default new User()
