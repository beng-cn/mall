import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 1. 引入 Element Plus 中文语言包
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import axios from 'axios'

if (import.meta.env.DEV) {
  localStorage.clear() 
}

const app = createApp(App)

// 2. 注册 Element Plus 并配置中文语言包（关键！）
app.use(ElementPlus, {
  locale: zhCn,
})

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// axios 全局配置
axios.defaults.baseURL = 'http://localhost:8080'
axios.defaults.timeout = 10000

axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

app.config.globalProperties.$axios = axios

app.use(router)
app.mount('#app')