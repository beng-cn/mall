import request from '../utils/request'

// 创建订单
export function createOrder(data) {
  return request({
    url: '/auth/order/create',
    method: 'post',
    data
  })
}

// 获取订单列表
export function getOrderList() {
  return request({
    url: '/auth/order/list',
    method: 'get'
  })
}

// 支付订单
export function payOrder(id) {
  return request({
    url: `/auth/order/pay/${id}`,
    method: 'post'
  })
}

// 获取订单商品
export function getOrderItems(id) {
  return request({
    url: `/auth/order/items/${id}`,
    method: 'get'
  })
}

// 删除订单
export function deleteOrder(id) {
  return request({
    url: `/auth/order/delete/${id}`,
    method: 'delete'
  })
}