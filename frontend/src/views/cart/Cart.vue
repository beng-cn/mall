<template>
  <div style="padding: 20px">
    <div style="display: flex; justify-content: space-between;">
      <h2>我的购物车</h2>
      <el-button type="danger" @click="logout">退出登录</el-button>
    </div>

    <el-table :data="cartList" border style="width: 100%; margin-top: 20px">
      <el-table-column prop="product_id" label="商品ID" />
      <el-table-column prop="name" label="商品名称" />
      <el-table-column prop="price" label="单价" />
      <el-table-column label="数量">
        <template #default="scope">
          <el-input-number
            v-model="scope.row.quantity"
            :min="1"
            @change="(val) => updateQuantity(scope.row.id, val)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="danger" @click="delCart(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="margin-top: 20px; text-align: right">
      <el-button type="primary" @click="toOrder">去结算</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { updateCart, deleteCart } from '../../api/cart'
import { ElMessage } from 'element-plus'

const router = useRouter()
const cartList = ref([
  { id: 1, product_id: 1, name: '示例商品', price: 99, quantity: 1 }
])

const updateQuantity = async (id, val) => {
  await updateCart(id, { quantity: val })
  ElMessage.success('更新成功')
}

const delCart = async (id) => {
  await deleteCart(id)
  cartList.value = cartList.value.filter(item => item.id !== id)
  ElMessage.success('删除成功')
}

const toOrder = () => {
  router.push('/order/create')
}

const logout = () => {
  localStorage.removeItem('user')
  ElMessage.success('退出成功')
  router.push('/user/login')
}
</script>