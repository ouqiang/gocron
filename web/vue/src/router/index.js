import Vue from 'vue'
import Router from 'vue-router'
import TaskList from '../pages/task/list'
import HostList from '../pages/host/list'
import UserList from '../pages/user/list'
import UserLogin from '../pages/user/login'
import System from '../pages/system/index'
import Install from '../pages/install/index'

Vue.use(Router)

const router = new Router({
  routes: [
    {
      path: '/install',
      name: 'install',
      component: Install
    },
    {
      path: '/task',
      name: 'task-list',
      component: TaskList
    },
    {
      path: '/host',
      name: 'host-list',
      component: HostList
    },
    {
      path: '/user',
      name: 'user-list',
      component: UserList
    },
    {
      path: '/user/login',
      name: 'user-login',
      component: UserLogin
    },
    {
      path: '/system',
      name: 'system-index',
      component: System
    }
  ]
})

export default router
