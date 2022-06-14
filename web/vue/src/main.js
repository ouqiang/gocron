import {createApp} from 'vue'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import App from './App.vue'
import Install from '@/views/install/index'
import router from './router'
import store from './store'


let $appConfirm = function (callback) {
    this.$confirm('确定执行此操作?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        callback()
    })
}

let app
if (window.location.hash.indexOf('/install') > -1) {
    app = createApp(Install) //直接打开安装页面
} else {
    app = createApp(App)
    //使用icon
    for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
        app.component(key, component)
    }

    app.use(router)
        .use(store)
}
app.config.globalProperties.$appConfirm = $appConfirm
app.use(ElementPlus, {
    locale: zhCn
})
    .mount('#app')
