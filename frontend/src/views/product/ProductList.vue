<template>
  <div class="product-list">
    <div style="padding: 20px; display: flex; justify-content: space-between;">
      <div>
        <el-input
          v-model="keyword"
          placeholder="搜索商品"
          style="width: 300px; margin-right: 10px"
          @keyup.enter="getList"
        />
        <el-button type="primary" @click="getList">搜索</el-button>
      </div>
      <el-button type="danger" @click="logout">退出登录</el-button>
    </div>

    <el-row :gutter="20" style="padding: 0 20px">
      <!-- 全部统一用 item.ID 大写 -->
      <el-col :span="6" v-for="item in list" :key="item.ID">
        <el-card shadow="hover" @click="toDetail(item.ID)">
          <div style="font-size: 16px; font-weight: bold">{{ item.name }}</div>
          <div style="color: red; font-size: 18px; margin: 10px 0">¥{{ item.price }}</div>
          <div style="font-size: 12px; color: #999">库存：{{ item.stock }}</div>
          <el-button
            type="primary"
            style="margin-top: 10px"
            @click.stop="handleAddCart(item.ID)"
          >
            加入购物车
          </el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getProductList } from '../../api/product'
import { addCart } from '../../api/cart'
import { ElMessage } from 'element-plus'

const router = useRouter()
const keyword = ref('')
const list = ref([])

const getList = async () => {
  const res = await getProductList({ keyword: keyword.value })
  console.log("商品列表:", res)
  list.value = res
}

// 最终修复！！！
const handleAddCart = async (pid) => {
  console.log("加入购物车的商品ID:", pid)
  
  if (!pid || pid <= 0) {
    ElMessage.error("商品ID无效")
    return
  }

  try {
    await addCart({
      user_id: 1,
      product_id: pid,
      quantity: 1
    })
    ElMessage.success("加入购物车成功")
  } catch (e) {
    ElMessage.error("加入失败")
  }
}

const toDetail = (id) => {
  router.push(`/product/detail/${id}`)
}

const logout = () => {
  localStorage.removeItem('user')
  ElMessage.success('退出成功')
  router.push('/user/login')
}

onMounted(() => getList())
</script>