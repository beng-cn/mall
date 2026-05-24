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
      <el-table-column prop="price" label="单价" :formatter="formatPrice" />
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
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCartList, updateCartQuantity, deleteCartItem } from '@/api/cart'
import { createOrder } from '@/api/order'
import { getProductDetail } from '@/api/product'

const router = useRouter()
const cartList = ref([])
const selectedList = ref([])

// 价格格式化
const formatPrice = (row) => {
  return `¥${row.price.toFixed(2)}`
}

// 获取购物车
const loadCartList = async () => {
  try {
    const carts = await getCartList()
    
    const finalList = []
    for (let cart of carts) {
      const product = await getProductDetail(cart.product_id)
      finalList.push({
        ...cart,
        name: product.name,
        price: product.price
      })
    }
    cartList.value = finalList
    console.log('✅ 购物车加载成功：', finalList)
  } catch (e) {
    console.error('❌ 获取购物车失败：', e)
    ElMessage.error(e.response?.data?.error || "获取购物车失败")
    if (e.response?.status === 401) {
      router.push('/user/login')
    }
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
    ElMessage.warning("请输入 1-99 之间的数字")
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
    await updateCartQuantity(row.id, { quantity: row.quantity })
    console.log('✅ 购物车数量更新成功')
  } catch (e) {
    console.error('❌ 更新数量失败：', e)
    ElMessage.error(e.response?.data?.error || "更新数量失败")
  }
}

const handleSelectionChange = (val) => {
  selectedList.value = val
}

// 删除
const delCart = async (id) => {
  try {
    await deleteCartItem(id)
    ElMessage.success("删除成功")
    loadCartList()
  } catch (e) {
    console.error('❌ 删除购物车失败：', e)
    ElMessage.error(e.response?.data?.error || "删除失败")
  }
}

// 生成订单
const goToCheckout = async () => {
  if (selectedList.value.length === 0) {
    ElMessage.warning("请至少选择一件商品")
    return
  }

  try {
    const cartIds = selectedList.value.map(item => item.id)
    const data = await createOrder({ cart_ids: cartIds })
    ElMessage.success("订单创建成功：" + data.order.order_no)
    loadCartList()
  } catch (e) {
    console.error('❌ 创建订单失败：', e)
    ElMessage.error(e.response?.data?.error || "创建订单失败")
  }
}

onMounted(() => {
  loadCartList()
})
</script>