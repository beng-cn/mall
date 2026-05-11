<template>
  <div id="app">
    <!-- 全局顶部导航 -->
    <div v-if="isLoggedIn" style="display: flex; gap: 20px; padding: 15px; background: #f5f5f5;">
      <router-link to="/product/list">首页</router-link>
      <router-link to="/cart">购物车</router-link>
      <router-link to="/order/list">我的订单</router-link>
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
</style>