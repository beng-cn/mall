<template>
  <div style="padding: 40px">
    <el-card v-if="info">
      <div style="display: flex; gap: 40px">
        <div>
          <h2>{{ info.name }}</h2>
          <p style="color: red; font-size: 24px; margin: 20px 0">¥{{ info.price }}</p>
          <p>库存：{{ info.stock }}</p>
          <el-button type="primary" style="margin-top: 20px" @click="handleAddCart">
            加入购物车
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getProductDetail } from '../../api/product'
import { addCart } from '../../api/cart'
import { ElMessage } from 'element-plus'

const route = useRoute()
const info = ref(null)

const getDetail = async () => {
  const res = await getProductDetail(route.params.id)
  info.value = res
}

const handleAddCart = async () => {
  await addCart({
    user_id: 1,
    product_id: info.value.id,
    quantity: 1
  })
  ElMessage.success('已加入购物车')
}

onMounted(() => getDetail())
</script>
