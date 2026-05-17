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

      <el-table-column prop="total" label="订单金额" />
      
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
            @click="payOrder(scope.row)"
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
            @click="deleteOrder(scope.row.id)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 订单商品弹窗（改用普通div，不使用ElList） -->
    <el-dialog v-model="dialogVisible" title="订单商品详情" width="500px">
      <div>
        <div 
          v-for="item in currentOrderItems" 
          :key="item.id"
          style="padding:10px 0;border-bottom:1px solid #eee;display:flex;justify-content:space-between"
        >
          <span>{{ item.name }}</span>
          <span style="color:#e53935">¥{{ item.price }} × {{ item.quantity }}</span>
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
import { ElMessage, ElDialog } from 'element-plus'

const orderList = ref([])
const dialogVisible = ref(false)
const currentOrderItems = ref([])

const getOrderList = async () => {
  try {
    const res = await fetch("http://localhost:8080/api/order/list?user_id=1")
    const data = await res.json()
    orderList.value = data

    const itemPromises = orderList.value.map(async (order) => {
      try {
        const resp = await fetch(`http://localhost:8080/api/order/items/${order.id}`)
        const items = await resp.json()
        const names = items.map(i => i.name)
        order.goodsText = names.length > 2 ? names.slice(0,2).join('，') + '…' : names.join('，')
        order.allItems = items
      } catch (e) {
        console.log("获取订单商品失败", e)
        order.goodsText = '加载失败'
        order.allItems = []
      }
    })
    await Promise.all(itemPromises)
  } catch (e) { 
    console.log("获取订单失败", e) 
    ElMessage.error("获取订单列表失败")
  }
}

const payOrder = async (order) => {
  try {
    const res = await fetch(`http://localhost:8080/api/order/pay/${order.id}`, { method: 'POST' })
    const data = await res.json()
    if (res.ok) {
      ElMessage.success("支付成功！")
      getOrderList()
    } else {
      ElMessage.error(data.error || "支付失败")
    }
  } catch (e) { ElMessage.error("支付请求失败") }
}

const deleteOrder = async (id) => {
  try {
    await fetch(`http://localhost:8080/api/order/delete/${id}`, { method: 'DELETE' })
    ElMessage.success("删除成功")
    getOrderList()
  } catch (e) { ElMessage.error("删除失败") }
}

const showOrderItems = (order) => {
  currentOrderItems.value = order.allItems
  dialogVisible.value = true
}

onMounted(() => { getOrderList() })
</script>