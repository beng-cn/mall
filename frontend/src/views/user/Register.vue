<template>
  <div class="register-box">
    <el-card style="width: 400px; margin: 100px auto; padding: 20px">
      <h2 style="text-align: center; margin-bottom: 20px">用户注册</h2>
      <el-form :model="form" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" style="width: 100%" @click="handleRegister">注册</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../../api/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const form = ref({
  username: '',
  password: '',
  nickname: ''
})

const handleRegister = async () => {
  const res = await register(form.value)
  if (res.message === '注册成功') {
    ElMessage.success('注册成功，请登录')
    router.push('/user/login')
  } else {
    ElMessage.error(res.error || '注册失败')
  }
}
</script>
