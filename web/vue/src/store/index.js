import vue from 'vue'
import vuex from 'vuex'
import userStorage from '../storage/user'

vue.use(vuex)
export default new vuex.Store({
  state: {
    hiddenNavMenu: false,
    user: userStorage.get()
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
    logout (state) {
      userStorage.clear()
      state.user = userStorage.get()
    }
  }
})
