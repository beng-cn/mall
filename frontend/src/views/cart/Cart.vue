<template>
  <div style="padding: 20px">
    <h2>我的购物车</h2>

    <div v-if="cartList.length === 0" style="margin-top: 30px">
      <h3>购物车暂无商品</h3>
    </div>

    <el-table
      v-else
      :data="cartList"
      border
      style="width: 100%; margin-top: 20px"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="商品名称" />
      <el-table-column prop="price" label="单价" />
      <el-table-column prop="quantity" label="数量" />
      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="danger" @click="delCart(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="cartList.length > 0" style="margin-top: 20px; text-align: right">
      <el-button
        type="primary"
        size="large"
        :disabled="selectedList.length === 0"
        @click="goToCheckout"
      >
        生成订单
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'

const cartList = ref([])
const selectedList = ref([])
const abortController = new AbortController()

const getCartList = async () => {
  try {
    const cartRes = await fetch("http://localhost:8080/api/cart/list?user_id=1", {
      signal: abortController.signal
    })
    if (!cartRes.ok) {
      throw new Error("获取购物车失败")
    }
    
    const carts = await cartRes.json()
    console.log("购物车原始数据:", carts)

    // 并行获取所有商品信息，大幅提升加载速度
    const productPromises = carts.map(async (cart) => {
      try {
        const pRes = await fetch(`http://localhost:8080/api/product/${cart.product_id}`, {
          signal: abortController.signal
        })
        if (!pRes.ok) {
          console.warn("商品不存在，product_id:", cart.product_id)
          return null
        }
        
        const product = await pRes.json()
        return {
          ...cart,
          name: product.name,
          price: product.price
        }
      } catch (e) {
        if (e.name !== 'AbortError') {
          console.error("获取商品信息失败:", e)
        }
        return null
      }
    })

    // 等待所有商品信息获取完成，并过滤掉无效项
    const results = await Promise.all(productPromises)
    cartList.value = results.filter(item => item !== null)
    console.log("处理后的购物车列表:", cartList.value)
  } catch (e) {
    if (e.name !== 'AbortError') {
      console.error("获取购物车失败:", e)
      ElMessage.error("获取购物车失败")
    }
  }
}

const handleSelectionChange = (val) => {
  selectedList.value = val
}

const delCart = async (id) => {
  console.log("删除购物车项ID:", id)
  
  if (!id) {
    ElMessage.error("无效的购物车项ID")
    return
  }

  try {
    const res = await fetch(`http://localhost:8080/api/cart/${id}`, { 
      method: "DELETE" 
    })
    
    if (res.ok) {
      ElMessage.success("删除成功")
      getCartList()
    } else {
      const data = await res.json()
      ElMessage.error(data.error || "删除失败")
    }
  } catch (e) {
    console.error("删除购物车项失败:", e)
    ElMessage.error("网络错误，删除失败")
  }
}

const goToCheckout = async () => {
  if (selectedList.value.length === 0) {
    ElMessage.warning("请选择要结算的商品")
    return
  }

  const cartIds = selectedList.value.map(item => item.id)

  try {
    const res = await fetch("http://localhost:8080/api/order/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        user_id: 1,
        cart_ids: cartIds
      })
    })

    if (res.ok) {
      const data = await res.json()
      ElMessage.success(`订单创建成功！订单号：${data.order.order_no}`)
      getCartList()
    } else {
      const data = await res.json()
      ElMessage.error(data.error || "订单创建失败")
    }
  } catch (err) {
    console.error("订单创建失败:", err)
    ElMessage.error("网络错误，订单创建失败")
  }
}

onMounted(() => {
  getCartList()
})

onUnmounted(() => {
  // 组件卸载时取消所有未完成的请求，防止内存泄漏
  abortController.abort()
})
</script>