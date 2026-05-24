import request from '../utils/request'

// 获取购物车列表
export function getCartList() {
  return request({
    url: '/auth/cart/list',
    method: 'get'
  })
}

// 添加商品到购物车
export function addCart(data) {
  return request({
    url: '/auth/cart/add',
    method: 'post',
    data
  })
}

// 更新购物车商品数量
export function updateCartQuantity(id, data) {
  return request({
    url: `/auth/cart/${id}`,
    method: 'put',
    data
  })
}

// 删除购物车商品
export function deleteCartItem(id) {
  return request({
    url: `/auth/cart/${id}`,
    method: 'delete'
  })
}