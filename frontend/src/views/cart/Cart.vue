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
      <!-- 选择框 -->
      <el-table-column type="selection" width="55" />
      <!-- 商品名称 -->
      <el-table-column prop="name" label="商品名称" />
      <!-- 商品价格 -->
      <el-table-column prop="price" label="单价" />
      <el-table-column prop="quantity" label="数量" />
      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="danger" @click="delCart(scope.row.ID)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 右下角去结算 -->
    <div v-if="cartList.length > 0" style="margin-top: 20px; text-align: right">
      <el-button
        type="primary"
        size="large"
        :disabled="selectedList.length === 0"
        @click="goToCheckout"
      >
        去结算
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const cartList = ref([])
const selectedList = ref([])

// 获取购物车 + 商品信息
const getCartList = async () => {
  try {
    const cartRes = await fetch("http://localhost:8080/api/cart/list?user_id=1")
    const carts = await cartRes.json()

    const finalList = []
    for (let cart of carts) {
      const pRes = await fetch(`http://localhost:8080/api/product/${cart.product_id}`)
      const product = await pRes.json()
      finalList.push({
        ...cart,
        name: product.name,
        price: product.price
      })
    }
    cartList.value = finalList
  } catch (e) {
    console.log(e)
  }
}

// 选择商品
const handleSelectionChange = (val) => {
  selectedList.value = val
}

// 删除
const delCart = async (id) => {
  await fetch(`http://localhost:8080/api/cart/${id}`, { method: "DELETE" })
  ElMessage.success("删除成功")
  getCartList()
}

// =============================================
// 🔥 对接你的后端订单接口：CreateOrder
// =============================================
const goToCheckout = async () => {
  const cartIds = selectedList.value.map(item => item.ID)

  try {
    const res = await fetch("http://localhost:8080/api/order/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        user_id: 1,
        cart_ids: cartIds
      })
    })

    const data = await res.json()
    ElMessage.success(`订单创建成功！订单号：${data.order.orderNo}`)
    getCartList() // 刷新购物车（已清空）
  } catch (err) {
    ElMessage.error("订单创建失败")
  }
}

onMounted(() => {
  getCartList()
})
</script>