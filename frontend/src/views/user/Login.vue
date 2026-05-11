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
    const res = await login(form.value)
    if (res.user) {
      // 保存登录状态
      localStorage.setItem('user', JSON.stringify(res.user))
      ElMessage.success('登录成功')
      router.push('/product/list')
    } else {
      ElMessage.error(res.error || '登录失败')
    }
  } catch (err) {
    // 打印错误详情，方便排查
    console.error('登录请求错误:', err)
    ElMessage.error(err.response?.data?.error || '账号或密码错误')
  }
}
</script>
