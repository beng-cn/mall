import request from '../utils/request'

// 添加购物车
export function addCart(data) {
  return request({
    url: '/cart/add',
    method: 'post',
    data
  })
}

// 修改数量
export function updateCart(id, data) {
  return request({
    url: `/cart/${id}`,
    method: 'put',
    data
  })
}

// 删除
export function deleteCart(id) {
  return request({
    url: `/cart/${id}`,
    method: 'delete'
  })
}