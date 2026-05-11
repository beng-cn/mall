<template>
  <div style="padding: 20px">
    <h2>我的订单</h2>
    <el-table :data="orderList" border style="width: 100%; margin-top: 20px">
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
// 这里先写死模拟数据，你可以后续对接后端订单列表接口
const orderList = ref([
  {
    id: 1,
    order_no: 'ORD123456',
    total: 99,
    status: 0
  }
])

onMounted(() => {
  // 后续可在这里调用后端订单列表接口
})
</script>