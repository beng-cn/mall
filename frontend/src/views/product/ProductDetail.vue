<template>
  <div style="padding: 20px; max-width: 1200px; margin: 0 auto;">
    <!-- 加载状态 -->
    <div v-if="product === null" style="text-align: center; padding: 50px;">
      <el-spinner size="large" />
    </div>

    <!-- 商品详情 -->
    <div v-else>
      <el-button @click="$router.back()" style="margin-bottom: 20px;">
        返回商品列表
      </el-button>

      <el-row :gutter="40">
        <el-col :span="12">
          <div style="border: 1px solid #eee; border-radius: 8px; padding: 20px; text-align: center;">
            <div style="width: 100%; height: 400px; background: #f5f5f5; display: flex; align-items: center; justify-content: center;">
              <span style="color: #999;">商品图片</span>
            </div>
          </div>
        </el-col>

        <el-col :span="12">
          <h1 style="font-size: 24px; margin-bottom: 20px;">{{ product.name }}</h1>
          <p style="color: #e53935; font-size: 32px; font-weight: bold; margin-bottom: 20px;">
            ¥{{ product.price }}
          </p>
          <p style="color: #666; margin-bottom: 10px;">库存：{{ product.stock }} 件</p>
          <p style="color: #666; margin-bottom: 10px;">商品编号：{{ product.id }}</p>
          <p style="color: #666; margin-bottom: 30px;">分类：{{ product.category.name }}</p>

          <div style="display: flex; gap: 10px; align-items: center; margin-bottom: 30px;">
            <span>数量：</span>
            <el-input-number v-model="quantity" :min="1" :max="product.stock" />
          </div>

          <el-button type="primary" size="large" @click="handleAddCart">
            加入购物车
          </el-button>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

const route = useRoute()
const product = ref(null)
const quantity = ref(1) 

const getProductDetail = async () => {
  const id = route.params.id
  
  if (!id || id === 'undefined') {
    ElMessage.error("无效的商品ID")
    return
  }
  
  try {
    const res = await fetch(`http://localhost:8080/api/product/${id}`)
    if (!res.ok) {
      throw new Error("商品不存在")
    }
    
    product.value = await res.json()
    console.log("商品详情数据:", product.value) 
  } catch (e) {
    console.error("获取商品详情失败:", e)
    ElMessage.error("获取商品详情失败")
  }
}


const handleAddCart = async () => {
  console.log("加入购物车的商品ID:", product.value.id)
  
  try {
    const res = await fetch("http://localhost:8080/api/cart/add", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        user_id: 1,
        product_id: product.value.id,
        quantity: quantity.value
      })
    })
    
    if (res.ok) {
      ElMessage.success("加入购物车成功")
    } else {
      const data = await res.json()
      ElMessage.error(data.error || "加入购物车失败")
    }
  } catch (e) {
    console.error("加入购物车失败:", e)
    ElMessage.error("网络错误，加入购物车失败")
  }
}

onMounted(() => {
  getProductDetail()
})
</script>