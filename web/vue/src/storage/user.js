class User {
  get () {
    return {
      'token': this.getToken(),
      'uid': this.getUid(),
      'username': this.getUsername(),
      'isAdmin': this.getIsAdmin()
    }
  }

  getToken () {
    return localStorage.getItem('token') || ''
  }

  setToken (token) {
    localStorage.setItem('token', token)
    return this
  }

  clear () {
    localStorage.clear()
  }

  getUid () {
    return localStorage.getItem('uid') || ''
  }

  setUid (uid) {
    localStorage.setItem('uid', uid)
    return this
  }

  getUsername () {
    return localStorage.getItem('username') || ''
  }

  setUsername (username) {
    localStorage.setItem('username', username)
    return this
  }

  getIsAdmin () {
    let isAdmin = localStorage.getItem('is_admin')
    return isAdmin === '1'
  }

  setIsAdmin (isAdmin) {
    localStorage.setItem('is_admin', isAdmin)
    return this
  }
}

export default new User()
