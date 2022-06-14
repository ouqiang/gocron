import {createRouter, createWebHashHistory} from 'vue-router'
import store from '../store/index'


const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            redirect: '/dashboard'
        }, {
            path: '/dashboard',
            component: () => import('@/views/dashboard')
        },
        {
            path: '/task',
            name: '任务管理',
            component: () => import('@/views/task/list')
        },
        {
            path: '/task/create',
            name: '新建任务',
            component: () => import('@/views/task/edit')
        },
        {
            path: '/task/edit/:id',
            name:'任务编辑',
            component: () => import('@/views/task/edit')
        },
        {
            path: '/task/logs',
            name:'任务日志',
            component: () => import('@/views/task/log')
        },
        {
            path: '/process/create',
            name: '进程创建',
            component: () => import('@/views/process/edit')
        },
        {
            path: '/process/edit/:id',
            name:'进程编辑',
            component: () => import('@/views/process/edit')
        },
        {
            path: '/process/index',
            name:'进程管理',
            component: () => import('@/views/process/list')
        },
        {
            path: '/host/index',
            name:'任务节点',
            component: () => import('@/views/host/list')
        },
        {
            path: '/host/create',
            name: '创建节点',
            component: () => import('@/views/host/edit')
        },
        {
            path: '/host/edit/:id',
            name:'节点编辑',
            component: () => import('@/views/host/edit')
        },
        {
            path: '/project',
            name: '项目管理',
            children: [
                {
                    path: '/project/index',
                    name: '项目管理',
                    component: () => import('@/views/project/list'),
                },
            ]
        },
        {
            path: '/user',
            name: '用户管理',
            children: [
                {
                    path: '/user/index',
                    name: '用户列表',
                    component: () => import('@/views/user/list'),
                },
                {
                    path: '/user/create',
                    name: '新增用户',
                    component: () => import('@/views/user/edit'),
                },
                {
                    path: '/user/edit/:id',
                    name: '编辑用户',
                    component: () => import('@/views/user/edit'),
                },
                {
                    path: '/user/edit-password/:id',
                    name: '修改密码',
                    component: () => import('@/views/user/editPassword'),
                },
                {
                    path: '/user/edit-my-password',
                    name: '修改我的密码',
                    component: () => import('@/views/user/editMyPassword'),
                },
                {
                    path: '/user/login',
                    name: '用户登录',
                    component: () => import('@/views/user/login'),
                    meta: {
                        noLogin: true
                    }
                }
            ]
        },{
            path:'/system',
            name:'系统管理',
            children:[
                {
                    path: '/system/notification/email',
                    name: '通知配置-Email',
                    component: () => import('@/views/system/notification/email'),
                },
                {
                    path: '/system/notification/slack',
                    name: '通知配置-Slack',
                    component: () => import('@/views/system/notification/slack'),
                },
                {
                    path: '/system/notification/webhook',
                    name: '通知配置-Webhook',
                    component: () => import('@/views/system/notification/webhook'),
                },
                {
                    path: '/system/login-log',
                    name: '登录日志',
                    component: () => import('@/views/system/loginLog'),
                },
                {
                    path: '/system/login-setting',
                    name: '登录认证',
                    component: () => import('@/views/system/loginSetting'),
                }
            ]
        },
        {
            path: '/:pathMatch(.*)*',
            name: '404',
            component: () => import('@/components/NotFound')
        }
    ]
})

router.beforeEach((to, from, next) => {
    if (to.meta.noLogin) {
        return next()
    }
    // console.log(to, from, next, store, store.getters.user.token, {redirect: to.fullPath})
    if (store.getters.user.token) {
        return next()
    }
    next({
        path: '/user/login',
        query: {redirect: to.fullPath}
    })
})

let routes = router.getRoutes()
let breadcrumb = []
let pathResolve = function (routes, path) {
    let isFind = false
    routes.forEach(route => {
        if (route.children) {
            if(pathResolve(route.children, path)){
                breadcrumb.push(route)
            }
        }
        if (route.path) {
            let re = new RegExp('^' + route.path.replace(':id', '.*') + '$')
            if (re.test(path)) {
                breadcrumb = []
                breadcrumb.push(route)
                isFind = true
            }
        }
    })
    return isFind
}

router.afterEach((to, from) => {
    if (to.fullPath === from.fullPath) {
        return
    }

    store.commit('setCurrentPath', to.fullPath === '/' ? '/dashboard' : to.fullPath)

    breadcrumb = []
    pathResolve(routes, to.fullPath)
    breadcrumb = breadcrumb.reverse()
    store.commit('setBreadcrumb', breadcrumb)
})

export default router