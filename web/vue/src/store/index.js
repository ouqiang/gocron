import vuex from 'vuex'
import userStorage from '../storage/user'

export default new vuex.Store({
  state: {
    hiddenNavMenu: false,
    user: userStorage.get(),
    currentPath: '',
    breadcrumb: []
  },
  getters: {
    user (state) {
      return state.user
    },
    login (state) {
      return state.user.token !== ''
    }
  },
  mutations: {
    hiddenNavMenu (state) {
      state.hiddenNavMenu = true
    },
    showNavMenu (state) {
      state.hiddenNavMenu = false
    },
    setUser (state, user) {
      userStorage.setToken(user.token)
      userStorage.setUid(user.uid)
      userStorage.setUsername(user.username)
      userStorage.setIsAdmin(user.isAdmin)
      state.user = user
    },
    setCurrentPath(state, path) {
      state.currentPath = path
    },
    setBreadcrumb(state, breadcrumb) {
      this.state.breadcrumb = breadcrumb
    },
    logout (state) {
      userStorage.clear()
      state.user = userStorage.get()
    }
  }
})
