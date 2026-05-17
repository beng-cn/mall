<template>
  <div id="app">
    <!-- 全局顶部导航 -->
    <div v-if="isLoggedIn" class="navbar">
      <router-link to="/product/list" class="nav-btn">首页</router-link>

      <!-- 商品分类弹窗 -->
      <div class="nav-dropdown">
        <div class="nav-btn" @click="showCategoryDialog = true">商品分类 ▾</div>

        <div v-if="showCategoryDialog" class="category-modal">
          <div class="modal-box">
            <h3>选择商品分类</h3>

            <!-- 父分类复选框（标准小写id）-->
            <div class="parent-list">
              <label class="checkbox-item" v-for="p in parentCategories" :key="p.id">
                <input
                  type="checkbox"
                  :value="p.id"
                  v-model="checkedParentIds"
                  @change="handleParentCheck(p.id)"
                >
                {{ p.name }}
              </label>
            </div>

            <!-- 子分类复选框（标准小写id）-->
            <div v-for="pid in checkedParentIds" :key="pid" class="child-list">
              <div v-if="childMap[pid] && childMap[pid].length > 0">
                <p class="child-title">{{ getParentName(pid) }} ：</p>
                <label
                  class="checkbox-item child"
                  v-for="c in childMap[pid]"
                  :key="c.id"
                >
                  <input
                    type="checkbox"
                    :value="String(c.id)"
                    v-model="checkedChildIds"
                  >
                  {{ c.name }}
                </label>
              </div>
              <div v-else-if="childMap[pid] && childMap[pid].length === 0" class="child-empty">
                暂无子分类
              </div>
              <div v-else class="child-empty">
                加载中...
              </div>
            </div>

            <div class="btns">
              <el-button size="small" @click="resetCategory">重置</el-button>
              <el-button type="primary" size="small" @click="confirmCategory">确定筛选</el-button>
              <el-button size="small" @click="showCategoryDialog = false">取消</el-button>
            </div>
          </div>
        </div>
      </div>

      <router-link to="/cart" class="nav-btn">购物车</router-link>
      <router-link to="/order/list" class="nav-btn">我的订单</router-link>
      <el-button type="danger" size="small" @click="logout">退出登录</el-button>
    </div>

    <router-view />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const isLoggedIn = ref(false)
const showCategoryDialog = ref(false)

// 分类数据
const parentCategories = ref([])
const childMap = ref({})

// 复选框选中
const checkedParentIds = ref([])
const checkedChildIds = ref([])

// 判断是否登录
const checkLogin = () => {
  isLoggedIn.value = !!localStorage.getItem('user')
}

// 退出登录
const logout = () => {
  localStorage.removeItem('user')
  ElMessage.success('退出成功')
  router.push('/user/login')
}

// 加载父分类
const loadParents = async () => {
  try {
    const res = await fetch('http://localhost:8080/api/product/category/parents')
    const data = await res.json()
    parentCategories.value = data
    console.log('✅ 加载父分类成功：', data)
  } catch (e) {
    console.error('❌ 加载父分类失败', e)
    ElMessage.error('加载分类失败')
  }
}

// 父分类勾选处理
const handleParentCheck = async (pid) => {
  if (!pid && pid !== 0) {
    console.error('❌ 无效的父分类ID：', pid)
    return
  }

  // 取消勾选时清空对应子分类
  if (!checkedParentIds.value.includes(pid)) {
    delete childMap.value[pid]
    checkedChildIds.value = checkedChildIds.value.filter(id => {
      const childList = childMap.value[pid] || []
      return !childList.some(c => String(c.id) === id)
    })
    return
  }

  // 勾选时加载子分类
  await loadChildren(pid)
}

// 加载子分类
const loadChildren = async (pid) => {
  if (childMap.value[pid]) return
  
  childMap.value[pid] = null
  
  try {
    const res = await fetch(`http://localhost:8080/api/product/category/children?parent_id=${pid}`)
    if (!res.ok) throw new Error(`HTTP错误：${res.status}`)
    
    const list = await res.json()
    console.log(`✅ 加载父分类${pid}的子分类成功：`, list)
    
    if (Array.isArray(list)) {
      childMap.value[pid] = list.filter(c => c && c.id)
    } else {
      console.error('❌ 子分类数据格式错误', list)
      childMap.value[pid] = []
      ElMessage.error('子分类数据格式错误')
    }
  } catch (e) {
    console.error(`❌ 加载父分类${pid}的子分类失败`, e)
    childMap.value[pid] = []
    ElMessage.error(`加载${getParentName(pid)}的子分类失败`)
  }
}

// 根据父分类ID获取名称
const getParentName = (pid) => {
  const parent = parentCategories.value.find(p => p.id === pid)
  return parent ? parent.name : `分类${pid}`
}

// 确认筛选（支持多分类同时筛选）
const confirmCategory = () => {
  // 合并所有选中的子分类ID和父分类ID
  const allSelectedIds = [...checkedChildIds.value, ...checkedParentIds.value.map(String)]
  
  // 去重（避免同时选中父分类和子分类时重复）
  const uniqueIds = [...new Set(allSelectedIds)]

  if (uniqueIds.length === 0) {
    ElMessage.warning('请至少选择一个分类')
    return
  }

  // 用逗号拼接成字符串传递给后端
  const categoryIdParam = uniqueIds.join(',')
  
  router.push({ 
    path: '/product/list', 
    query: { category_id: categoryIdParam } 
  })
  
  showCategoryDialog.value = false
}

// 重置
const resetCategory = () => {
  checkedParentIds.value = []
  checkedChildIds.value = []
  childMap.value = {}
}

onMounted(() => {
  checkLogin()
  loadParents()
})

router.afterEach(() => {
  checkLogin()
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.navbar {
  display: flex;
  gap: 15px;
  padding: 15px 25px;
  background: #db590d;
  align-items: center;
  position: relative;
  z-index: 9999;
}

.nav-btn {
  background: #e68200;
  color: #fff !important;
  padding: 6px 14px;
  border-radius: 4px;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  transition: 0.2s;
}
.nav-btn:hover {
  background: #d36f00;
  color: #fff !important;
}

.category-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}
.modal-box {
  background: #fff;
  padding: 25px 30px;
  border-radius: 10px;
  width: 420px;
}
.modal-box h3 {
  margin-bottom: 15px;
  font-size: 16px;
}

.checkbox-item {
  display: block;
  margin: 8px 0;
  cursor: pointer;
}
.child {
  margin-left: 20px;
  color: #666;
}
.child-title {
  margin: 10px 0 5px 10px;
  font-size: 13px;
  color: #999;
}
.child-empty {
  margin-left: 20px;
  font-size: 12px;
  color: #999;
  padding: 5px 0;
}

.btns {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}
</style>