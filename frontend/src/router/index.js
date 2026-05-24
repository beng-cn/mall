import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'

const routes = [
  // 默认首页 → 登录页
  {
    path: '/',
    redirect: '/user/login'
  },

  // 用户模块（无需登录）
  {
    path: '/user/login',
    name: 'Login',
    component: () => import('../views/user/Login.vue')
  },
  {
    path: '/user/register',
    name: 'Register',
    component: () => import('../views/user/Register.vue')
  },

  // 商品模块（必须登录）
  {
    path: '/product/list',
    name: 'ProductList',
    component: () => import('../views/product/ProductList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/product/detail/:id',
    name: 'ProductDetail',
    component: () => import('../views/product/ProductDetail.vue'),
    meta: { requiresAuth: true }
  },
  // 新增商品路由（管理员专属）
  {
    path: '/product/add',
    name: 'AddProduct',
    component: () => import('../views/product/AddProduct.vue'),
    meta: { requiresAuth: true, isAdmin: true }
  },

  // 购物车（必须登录）
  {
    path: '/cart',
    name: 'Cart',
    component: () => import('../views/cart/Cart.vue'),
    meta: { requiresAuth: true }
  },

  // 订单模块（必须登录）
  {
    path: '/order/create',
    name: 'CreateOrder',
    component: () => import('../views/order/CreateOrder.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/order/list',
    name: 'OrderList',
    component: () => import('../views/order/OrderList.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  const userStr = localStorage.getItem('user')
  const isLoggedIn = !!userStr

  if (to.meta.requiresAuth && !isLoggedIn) {
    ElMessage.warning('请先登录')
    return next('/user/login')
  }

  if (to.meta.isAdmin) {
    try {
      const user = JSON.parse(userStr)
      if (user.role_id !== 1) {
        ElMessage.error('权限不足，仅管理员可访问')
        return next('/product/list')
      }
    } catch (e) {
      ElMessage.error('登录状态异常，请重新登录')
      localStorage.clear()
      return next('/user/login')
    }
  }

  next()
})

export default router