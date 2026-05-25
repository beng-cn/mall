<template>
  <div class="user-management-container">
    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="请输入用户名搜索"
        style="width: 240px; margin-right: 10px"
        @keyup.enter="getUserList"
      >
        <template #append>
          <el-icon @click="getUserList"><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="getUserList">搜索</el-button>
      <el-button @click="resetSearch">重置</el-button>
    </div>

    <el-table
      :data="userList"
      border
      stripe
      style="width: 100%; margin-top: 15px"
      v-loading="loading"
    >
      <el-table-column prop="id" label="ID" width="80" align="center" />
      <el-table-column prop="username" label="用户名" min-width="120" />
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column prop="phone" label="手机号" min-width="120" />
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
            {{ scope.row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="role_id" label="角色" width="100" align="center">
        <template #default="scope">
          <el-tag type="primary">
            {{ scope.row.role_id === 1 ? '管理员' : '普通用户' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" min-width="180" />
      <el-table-column label="操作" width="200" align="center">
        <template #default="scope">
          <el-button
            size="small"
            :type="scope.row.status === 1 ? 'warning' : 'success'"
            @click="toggleUserStatus(scope.row)"
            :disabled="scope.row.role_id === 1"
          >
            {{ scope.row.status === 1 ? '禁用' : '启用' }}
          </el-button>
          <el-button
            size="small"
            type="danger"
            @click="deleteUser(scope.row.id)"
            :disabled="scope.row.role_id === 1"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

   <div class="pagination-container">
  <el-pagination
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
    :current-page="pageNum"
    :page-sizes="[10, 20, 50, 100]"
    :page-size="pageSize"
    layout="total, sizes, prev, pager, next, jumper"
    :total="total"
    style="margin-top: 15px; text-align: right"
  />
</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import axios from 'axios'

// 状态管理
const loading = ref(false)
const userList = ref([])
const total = ref(0)
const pageNum = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')

// 获取用户列表
const getUserList = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/admin/user/list', {
      params: {
        keyword: searchKeyword.value,
        page_num: pageNum.value,
        page_size: pageSize.value
      }
    })
    
    if (res.data.code === 200 || res.status === 200) {
      userList.value = res.data.list || []
      total.value = res.data.total || 0
    } else {
      ElMessage.error(res.data.message || '获取用户列表失败')
    }
  } catch (error) {
    console.error('获取用户列表失败：', error)
    ElMessage.error('网络错误，获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 切换用户状态（禁用/启用）
const toggleUserStatus = async (user) => {
  try {
    const confirmText = user.status === 1 ? '禁用' : '启用'
    await ElMessageBox.confirm(
      `确定要${confirmText}用户【${user.username}】吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: user.status === 1 ? 'warning' : 'info'
      }
    )

    const res = await axios.put(
      `/api/admin/user/${user.id}/status`,
      { status: user.status === 1 ? 0 : 1 }
    )

    if (res.data.code === 200 || res.status === 200) {
      ElMessage.success(`${confirmText}成功`)
      getUserList()
    } else {
      ElMessage.error(res.data.message || `${confirmText}失败`)
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('切换用户状态失败：', error)
      ElMessage.error('操作失败，请重试')
    }
  }
}

// 删除用户
const deleteUser = async (userId) => {
  try {
    await ElMessageBox.confirm(
      '此操作将永久删除该用户，是否继续？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )

    const res = await axios.delete(`/api/admin/user/${userId}`)

    if (res.data.code === 200 || res.status === 200) {
      ElMessage.success('删除成功')
      getUserList()
    } else {
      ElMessage.error(res.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用户失败：', error)
      ElMessage.error('删除失败，请重试')
    }
  }
}

// 分页相关
const handleSizeChange = (val) => {
  pageSize.value = val
  getUserList()
}

const handleCurrentChange = (val) => {
  pageNum.value = val
  getUserList()
}

// 重置搜索
const resetSearch = () => {
  searchKeyword.value = ''
  pageNum.value = 1
  getUserList()
}

// 页面挂载时加载数据
onMounted(() => {
  getUserList()
})
</script>

<style scoped>
.user-management-container {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.search-bar {
  display: flex;
  align-items: center;
  padding: 10px 0;
}

.pagination-container {
  margin-top: 20px;
}

:deep(.el-table .el-table-column--selection) {
  width: 60px;
}

:deep(.el-tag) {
  cursor: default;
}
</style>