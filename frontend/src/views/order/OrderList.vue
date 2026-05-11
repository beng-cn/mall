<template>
  <div style="padding: 20px">
    <h2>我的订单</h2>
    <el-table :data="orderList" border style="width: 100%; margin-top: 20px">
      <!-- 这里用 order_no，和你结构体的 json 标签完全一致 -->
      <el-table-column prop="order_no" label="订单号" />
      <el-table-column prop="total" label="订单金额" />
      <el-table-column prop="status" label="订单状态">
        <template #default="scope">
          {{ scope.row.status === 0 ? '待支付' : scope.row.status === 1 ? '已支付' : '已取消' }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const orderList = ref([])

const getOrderList = async () => {
  try {
    const res = await fetch("http://localhost:8080/api/order/list?user_id=1")
    const data = await res.json()
    orderList.value = data
  } catch (e) {
    console.log("获取订单失败", e)
  }
}

onMounted(() => {
  getOrderList()
})
</script>