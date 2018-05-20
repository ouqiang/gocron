import Vue from 'vue'
import Router from 'vue-router'
import store from '../store/index'
import NotFound from '../components/common/notFound'

import TaskList from '../pages/task/list'
import TaskEdit from '../pages/task/edit'
import TaskLog from '../pages/taskLog/list'

import HostList from '../pages/host/list'
import HostEdit from '../pages/host/edit'

import UserList from '../pages/user/list'
import UserEdit from '../pages/user/edit'
import UserLogin from '../pages/user/login'
import UserEditPassword from '../pages/user/editPassword'
import UserEditMyPassword from '../pages/user/editMyPassword'

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
      path: '/task/create',
      name: 'task-create',
      component: TaskEdit
    },
    {
      path: '/task/edit/:id',
      name: 'task-edit',
      component: TaskEdit
    },
    {
      path: '/task/log',
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
      path: '/host/create',
      name: 'host-create',
      component: HostEdit
    },
    {
      path: '/host/edit/:id',
      name: 'host-edit',
      component: HostEdit
    },
    {
      path: '/user',
      name: 'user-list',
      component: UserList
    },
    {
      path: '/user/create',
      name: 'user-create',
      component: UserEdit
    },
    {
      path: '/user/edit/:id',
      name: 'user-edit',
      component: UserEdit
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
      path: '/user/edit-password/:id',
      name: 'user-edit-password',
      component: UserEditPassword
    },
    {
      path: '/user/edit-my-password',
      name: 'user-edit-my-password',
      component: UserEditMyPassword,
      meta: {
        noNeedAdmin: true
      }
    },
    {
      path: '/system',
      redirect: '/system/notification/email'
    },
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
