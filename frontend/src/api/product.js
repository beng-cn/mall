import request from '../utils/request'

// 商品列表（公开接口，游客也能看）
export function getProductList(params) {
  return request({
    url: '/product/list',   
    method: 'get',
    params
  })
}

// 商品详情（公开接口，游客也能看）
export function getProductDetail(id) {
  return request({
    url: `/product/${id}`,  
    method: 'get'
  })
}

// 获取所有父分类（公开接口，游客也能看）
export function getParentCategories() {
  return request({
    url: '/product/category/parents',  
    method: 'get'
  })
}

// 根据父分类ID获取子分类（公开接口，游客也能看）
export function getChildCategories(parentId) {
  return request({
    url: '/product/category/children',  
    method: 'get',
    params: { parent_id: parentId }
  })
}

// 新增商品（管理员接口，需要加/admin前缀）
export function createProduct(data) {
  return request({
    url: '/admin/product',  // 管理员接口加/admin前缀
    method: 'post',
    data
  })
}

// 更新商品（管理员接口）
export function updateProduct(id, data) {
  return request({
    url: `/admin/product/${id}`,  // 管理员接口加/admin前缀
    method: 'put',
    data
  })
}

// 删除商品（管理员接口）
export function deleteProduct(id) {
  return request({
    url: `/admin/product/${id}`,  // 管理员接口加/admin前缀
    method: 'delete'
  })
}