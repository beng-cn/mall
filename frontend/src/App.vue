<template>
  <div id="app">
    <!-- 全局顶部导航 → 橙黄色背景 -->
    <div v-if="isLoggedIn" class="navbar">
      <router-link to="/product/list" class="nav-btn">首页</router-link>
      <router-link to="/cart" class="nav-btn">购物车</router-link>
      <router-link to="/order/list" class="nav-btn">我的订单</router-link>
      <el-button type="danger" size="small" @click="logout">退出登录</el-button>
    </div>

    <router-view />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const isLoggedIn = ref(false)

// 判断是否登录
const checkLogin = () => {
  isLoggedIn.value = !!localStorage.getItem('user')
}

// 退出登录
const logout = () => {
  localStorage.removeItem('user')
  ElMessage.success('退出成功')
  router.push('/user/login')
}

onMounted(() => {
  checkLogin()
})

// 路由变化时重新检查登录状态
router.afterEach(() => {
  checkLogin()
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

/* 导航栏整体：橙黄色背景 */
.navbar {
  display: flex;
  gap: 15px;
  padding: 15px 25px;
  background: #db590d; 
  align-items: center;
}

/* 导航按钮：深黄色背景 + 白色文字 */
.nav-btn {
  background: #e68200; 
  color: #fff !important;
  padding: 6px 14px;
  border-radius: 4px;
  text-decoration: none;
  font-size: 14px;
  transition: 0.2s;
}

/* 鼠标悬浮效果 */
.nav-btn:hover {
  background: #d36f00;
  color: #fff !important;
}
</style>
