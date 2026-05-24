import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const request = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

request.interceptors.request.use(
  config => {
    const whiteList = ['/user/login', '/user/register']
    if (!whiteList.includes(config.url)) {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    return config
  },
  error => Promise.reject(error)
)

request.interceptors.response.use(
  res => res.data,
  error => {
    if (error.response?.status === 401) {
      if (error.config.url.includes('/user/login')) {
        ElMessage.error('账号或密码错误')
      } else {
        ElMessage.error('登录已过期，请重新登录')
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/user/login')
      }
    } else if (error.response?.status === 403) {
      ElMessage.error('权限不足')
    } else {
      ElMessage.error(error.response?.data?.error || '请求失败')
    }
    return Promise.reject(error)
  }
)

export default request