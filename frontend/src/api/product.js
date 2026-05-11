import request from '../utils/request'

// 商品列表
export function getProductList(params) {
  return request({
    url: '/product/list',
    method: 'get',
    params
  })
}

// 商品详情
export function getProductDetail(id) {
  return request({
    url: `/product/${id}`,
    method: 'get'
  })
}