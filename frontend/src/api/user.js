import request from '../utils/request'

// 登录（公开接口）
export function login(data) {
  return request({
    url: '/user/login',  
    method: 'post',
    data
  })
}

// 注册（公开接口）
export function register(data) {
  return request({
    url: '/user/register',  
    method: 'post',
    data
  })
}