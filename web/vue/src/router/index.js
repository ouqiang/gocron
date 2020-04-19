import Vue from 'vue'
import Router from 'vue-router'
import store from '../store/index'
import NotFound from '../components/common/notFound'

import TaskList from '../pages/task/list'
import TaskLog from '../pages/taskLog/list'

import HostList from '../pages/host/list'

import UserList from '../pages/user/list'
import UserLogin from '../pages/user/login'

import NotificationTab from '../pages/system/notification/tab'
import NotificationEmail from '../pages/system/notification/email'
import NotificationSlack from '../pages/system/notification/slack'
import NotificationWebhook from '../pages/system/notification/webhook'

import Install from '../pages/install/index'
import LoginLog from '../pages/system/loginLog'

Vue.use(Router)

const router = new Router({
  routes: [
    {
      path: '*',
      component: NotFound,
      meta: {
        noLogin: true,
        noNeedAdmin: true
      }
    },
    {
      path: '/',
      redirect: '/task'
    },
    {
      path: '/install',
      name: 'install',
      component: Install,
      meta: {
        noLogin: true,
        noNeedAdmin: true
      }
    },
    {
      path: '/task',
      name: 'task-list',
      component: TaskList,
      meta: {
        noNeedAdmin: true
      }
    },
    {
      path: '/task-log',
      name: 'task-log',
      component: TaskLog,
      meta: {
        noNeedAdmin: true
      }
    },
    {
      path: '/host',
      name: 'host-list',
      component: HostList,
      meta: {
        noNeedAdmin: true
      }
    },
    {
      path: '/user',
      name: 'user-list',
      component: UserList
    },
    {
      path: '/user/login',
      name: 'user-login',
      component: UserLogin,
      meta: {
        noLogin: true
      }
    },
    {
      path: '/system/notification',
      component: NotificationTab,
      children: [
        {
          path: '/system/notification/email',
          name: 'system-notification-email',
          component: NotificationEmail
        },
        {
          path: '/system/notification/slack',
          name: 'system-notification-slack',
          component: NotificationSlack
        },
        {
          path: '/system/notification/webhook',
          name: 'system-notification-webhook',
          component: NotificationWebhook
        }
      ]
    },
    {
      path: '/system/login-log',
      name: 'login-log',
      component: LoginLog
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.meta.noLogin) {
    next()
    return
  }
  if (store.getters.user.token) {
    if ((store.getters.user.isAdmin || to.meta.noNeedAdmin)) {
      next()
      return
    }
    if (!store.getters.user.isAdmin) {
      next(
        {
          path: '/404.html'
        }
      )
      return
    }
  }

  next({
    path: '/user/login',
    query: {redirect: to.fullPath}
  })
})

export default router
