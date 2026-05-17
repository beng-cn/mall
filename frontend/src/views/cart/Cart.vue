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

      <!-- 数量：- 数字 + -->
      <el-table-column label="数量" width="180">
        <template #default="scope">
          <div style="display: flex; align-items: center; justify-content: center; gap: 5px;">
            <el-button size="small" @click="handleMinus(scope.row)">-</el-button>
            <el-input
              v-model.number="scope.row.quantity"
              style="width: 60px; text-align: center"
              @blur="handleBlur(scope.row)"
              @keyup.enter="handleBlur(scope.row)"
            />
            <el-button size="small" @click="handlePlus(scope.row)">+</el-button>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="danger" @click="delCart(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="margin-top: 20px; text-align: right">
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
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const cartList = ref([])
const selectedList = ref([])

// 获取购物车
const getCartList = async () => {
  try {
    const res = await fetch("http://localhost:8080/api/cart/list?user_id=1")
    const carts = await res.json()
    
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
    ElMessage.error("获取购物车失败")
  }
}

// +
const handlePlus = (row) => {
  if (row.quantity >= 99) {
    ElMessage.warning("数量不能超过99")
    return
  }
  row.quantity += 1
  updateQuantity(row)
}

// -
const handleMinus = (row) => {
  if (row.quantity <= 1) {
    ElMessage.warning("数量不能小于1")
    return
  }
  row.quantity -= 1
  updateQuantity(row)
}

// 失去焦点校验
const handleBlur = async (row) => {
  let qty = parseInt(row.quantity)

  if (isNaN(qty) || qty < 1 || qty > 99) {
    ElMessage.warning("请输入 1‑99 之间的数字")
    row.quantity = 1
    await updateQuantity(row)
    return
  }

  row.quantity = qty
  await updateQuantity(row)
}

// 同步后端
const updateQuantity = async (row) => {
  try {
    await fetch(`http://localhost:8080/api/cart/${row.id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ quantity: row.quantity })
    })
  } catch (e) {
    ElMessage.error("更新数量失败")
  }
}

const handleSelectionChange = (val) => {
  selectedList.value = val
}

// 删除
const delCart = async (id) => {
  await fetch(`http://localhost:8080/api/cart/${id}`, { method: "DELETE" })
  ElMessage.success("删除成功")
  getCartList()
}

// 生成订单
const goToCheckout = async () => {
  const cartIds = selectedList.value.map(item => item.id)
  const res = await fetch("http://localhost:8080/api/order/create", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user_id: 1, cart_ids: cartIds })
  })
  const data = await res.json()
  ElMessage.success("订单创建成功：" + data.order.order_no)
  getCartList()
}

onMounted(() => {
  getCartList()
})
</script>