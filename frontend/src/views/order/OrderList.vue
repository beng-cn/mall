<template>
  <div style="padding: 20px">
    <h2>我的订单</h2>
    <el-table :data="orderList" border style="width: 100%; margin-top: 20px">
      
      <el-table-column prop="order_no" label="订单号" />

      <el-table-column label="商品" width="260">
        <template #default="scope">
          <span 
            style="cursor: pointer; color: #5677fc;" 
            @click="showOrderItems(scope.row)"
          >
            {{ scope.row.goodsText || '加载中...' }}
          </span>
        </template>
      </el-table-column>

      <el-table-column prop="total" label="订单金额" :formatter="formatPrice" />
      
      <el-table-column prop="status" label="订单状态">
        <template #default="scope">
          {{ scope.row.status === 0 ? '待支付' : scope.row.status === 1 ? '已支付' : '已取消' }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="220">
        <template #default="scope">
         <el-button
  type="success"
  size="small"
  v-if="scope.row.status === 0"
  @click="handlePayOrder(scope.row)"
>
  去支付
</el-button>
<el-tag v-else type="success" size="small">
  {{ scope.row.status === 1 ? '已支付' : '已取消' }}
</el-tag>

<el-button
  type="danger"
  size="small"
  style="margin-left: 5px"
  @click="handleDeleteOrder(scope.row)"
>
  删除
</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 订单商品弹窗 -->
    <el-dialog v-model="dialogVisible" title="订单商品详情" width="500px">
      <div>
        <div 
          v-for="item in currentOrderItems" 
          :key="item.id"
          style="padding:10px 0;border-bottom:1px solid #eee;display:flex;justify-content:space-between"
        >
          <span>{{ item.name }}</span>
          <span style="color:#e53935">¥{{ item.price.toFixed(2) }} × {{ item.quantity }}</span>
        </div>
      </div>
      <template #footer>
        <el-button type="primary" @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getOrderList, payOrder, deleteOrder, getOrderItems } from '@/api/order'

const router = useRouter()
const orderList = ref([])
const dialogVisible = ref(false)
const currentOrderItems = ref([])

// 价格格式化
const formatPrice = (row) => {
  return `¥${row.total.toFixed(2)}`
}

// 获取订单列表
const loadOrderList = async () => {
  try {
    const data = await getOrderList()
    orderList.value = data
    console.log('✅ 订单列表加载成功：', data)

    // 批量获取订单商品
    const itemPromises = orderList.value.map(async (order) => {
      try {
        const items = await getOrderItems(order.id)
        const names = items.map(i => i.name)
        order.goodsText = names.length > 2 ? names.slice(0,2).join('，') + '…' : names.join('，')
        order.allItems = items
      } catch (e) {
        console.error("获取订单商品失败", e)
        order.goodsText = '加载失败'
        order.allItems = []
      }
    })
    await Promise.all(itemPromises)
  } catch (e) { 
    console.error("获取订单失败", e) 
    ElMessage.error(e.response?.data?.error || "获取订单列表失败")
    if (e.response?.status === 401) {
      router.push('/user/login')
    }
  }
}

// 支付订单
const handlePayOrder = async (order) => {
  try {
    await payOrder(order.id)
    ElMessage.success("支付成功！")
    loadOrderList()
  } catch (e) { 
    console.error("支付失败", e)
    ElMessage.error(e.response?.data?.error || "支付失败") 
  }
}

// 删除订单
const handleDeleteOrder = async (order) => {
  try {
    await deleteOrder(order.id)
    ElMessage.success("删除成功")
    loadOrderList()
  } catch (e) { 
    console.error("删除订单失败", e)
    ElMessage.error(e.response?.data?.error || "删除失败") 
  }
}

// 显示订单商品详情
const showOrderItems = (order) => {
  currentOrderItems.value = order.allItems
  dialogVisible.value = true
}

onMounted(() => { 
  loadOrderList() 
})
</script>