<template>
  <div class="login-box">
    <el-card style="width: 400px; margin: 100px auto; padding: 20px">
      <h2 style="text-align: center; margin-bottom: 20px">用户登录</h2>
      <el-form :model="form" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" style="width: 100%" @click="handleLogin">登录</el-button>
        </el-form-item>
        <el-form-item>
          <el-button link @click="$router.push('/user/register')">
            没有账号？去注册
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../../api/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const form = ref({
  username: '',
  password: ''
})

const handleLogin = async () => {
  try {
    localStorage.clear()

    console.log('📤 发送登录请求：', form.value)
    const res = await login(form.value)
    console.log('📥 登录响应完整数据：', res)
    
    if (res && res.token) {
      console.log('✅ 登录成功，准备跳转')
      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(res.user))
      localStorage.setItem('role_id', res.user.role_id) 
      
      ElMessage.success('登录成功')
      await router.push('/product/list')
    } else {
      ElMessage.error(res?.error || '登录失败，服务器返回数据异常')
    }
  } catch (err) {
    console.error('❌ 登录请求异常：', err)
    if (err.response) {
      ElMessage.error(err.response.data?.error || '登录失败')
    } else if (err.request) {
      ElMessage.error('网络异常，请检查连接')
    } else {
      ElMessage.error('登录失败，请重试')
    }
  }
}
</script>
