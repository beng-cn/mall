import { createRouter, createWebHistory } from 'vue-router'

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

// 登录拦截路由守卫
router.beforeEach((to, from, next) => {
  const userInfo = localStorage.getItem('user')
  const isLoggedIn = !!userInfo

  if (to.meta.requiresAuth && !isLoggedIn) {
    next('/user/login')
  } else {
    next()
  }
})

export default router
