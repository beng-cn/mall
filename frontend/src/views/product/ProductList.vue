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
    </div>

    <el-row :gutter="20" style="padding: 0 20px">
      <el-col :span="6" v-for="item in list" :key="item.id">
        <el-card shadow="hover" @click="toDetail(item.id)">
          <div style="font-size: 16px; font-weight: bold">{{ item.name }}</div>
          <div style="color: red; font-size: 18px; margin: 10px 0">¥{{ item.price }}</div>
          <div style="font-size: 12px; color: #999">库存：{{ item.stock }}</div>
          <el-button
            type="primary"
            style="margin-top: 10px"
            @click.stop="handleAddCart(item.id)"
          >
            加入购物车
          </el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getProductList } from '../../api/product'
import { addCart } from '../../api/cart'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()

const keyword = ref('')
const list = ref([])

const getList = async () => {
  const params = {
    keyword: keyword.value
  }
  
  if (route.query.category_id) {
    params.category_id = route.query.category_id
  }
  
  const res = await getProductList(params)
  console.log("商品列表:", res)
  list.value = res
}

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
    ElMessage.error("商品库存不足，无法加入购物车")
  }
}

const toDetail = (id) => {
  router.push(`/product/detail/${id}`)
}

onMounted(() => getList())

watch(() => route.query.category_id, () => {
  getList()
})
</script>