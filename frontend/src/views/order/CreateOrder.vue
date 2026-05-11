<template>
  <div style="padding: 20px">
    <h2>创建订单</h2>
    <el-card style="margin-top: 20px">
      <p>订单商品：示例商品 x1</p>
      <p>总价：¥99</p>
      <el-button type="primary" style="margin-top: 20px" @click="submitOrder">
        提交订单
      </el-button>
    </el-card>
  </div>
</template>

<script setup>
import { createOrder } from '../../api/order'
import { ElMessage } from 'element-plus'

const submitOrder = async () => {
  const res = await createOrder({
    user_id: 1,
    cart_ids: [1]
  })
  if (res.order) {
    ElMessage.success('订单创建成功，订单号：' + res.order.order_no)
  } else {
    ElMessage.error('订单创建失败')
  }
}
</script>
