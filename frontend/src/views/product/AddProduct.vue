<template>
  <div class="add-product-container">
    <div class="card">
      <h2>{{ isEditMode ? '编辑商品' : '新增商品' }}</h2>

      <!-- 分类选择区域 -->
      <div class="category-section">
        <h3>选择商品分类</h3>
        
        <div class="form-item">
          <label>父分类：</label>
          <el-select 
            v-model="selectedParentId" 
            placeholder="请选择父分类" 
            @change="handleParentChange"
            style="width: 220px"
          >
            <el-option
              v-for="p in parentCategories"
              :key="p.id"
              :label="p.name"
              :value="p.id"
            ></el-option>
          </el-select>
          <el-button size="small" type="primary" @click="openAddParentCategory">添加</el-button>
          <el-button size="small" type="warning" @click="openEditParentCategory" :disabled="!selectedParentId">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteParentCategory" :disabled="!selectedParentId">删除</el-button>
        </div>

        <div class="form-item" v-if="selectedParentId">
          <label>子分类：</label>
          <el-select 
            v-model="selectedChildId" 
            placeholder="请选择子分类（可选）" 
            style="width: 220px"
            @change="handleChildChange"
          >
            <el-option
              v-for="c in childCategories"
              :key="c.id"
              :label="c.name"
              :value="c.id"
            ></el-option>
          </el-select>
          <el-button size="small" type="primary" @click="openAddChildCategory">添加</el-button>
          <el-button size="small" type="warning" @click="openEditChildCategory" :disabled="!selectedChildId">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteChildCategory" :disabled="!selectedChildId">删除</el-button>
        </div>

        <!-- 选择已有商品进行编辑/删除 -->
        <div class="form-item" v-if="selectedChildId">
          <label>选择商品：</label>
          <el-select 
            v-model="selectedProductId" 
            placeholder="选择已有商品进行编辑" 
            style="width: 220px"
            @change="handleProductSelect"
            clearable
          >
            <el-option
              v-for="p in productOptions"
              :key="p.id"
              :label="p.name"
              :value="p.id"
            ></el-option>
          </el-select>
          <el-button 
            size="small" 
            type="danger" 
            @click="deleteProduct" 
            :disabled="!selectedProductId"
          >
            删除商品
          </el-button>
          <span style="color: #999; font-size: 12px; margin-left: 10px;">选择后自动填充商品信息</span>
        </div>
      </div>

      <!-- 商品信息表单 -->
      <div class="form-section">
        <h3>商品基本信息</h3>
        
        <div class="form-item">
          <label>商品名称：</label>
          <el-input 
            v-model="productForm.name" 
            placeholder="请输入商品名称" 
            style="width: 300px"
            maxlength="50"
            show-word-limit
          ></el-input>
        </div>

        <div class="form-item">
          <label>商品价格：</label>
          <el-input-number 
            v-model="productForm.price" 
            :min="0.01" 
            :step="0.01" 
            :precision="2"
            placeholder="请输入商品价格"
          ></el-input-number>
          <span class="unit">元</span>
        </div>

        <div class="form-item">
          <label>库存数量：</label>
          <el-input-number 
            v-model="productForm.stock" 
            :min="0" 
            :step="1"
            placeholder="请输入库存数量"
          ></el-input-number>
          <span class="unit">件</span>
        </div>

        <div class="form-item">
          <label>商品图片：</label>
          <input 
            type="file" 
            ref="fileInput" 
            style="display: none" 
            accept="image/jpeg,image/jpg,image/png,image/gif"
            @change="handleFileUpload"
          />
          <el-button type="primary" size="small" @click="triggerFileSelect">点击上传图片</el-button>
          
          <div v-if="productForm.image" class="image-preview">
            <img :src="productForm.image" alt="商品预览" />
            <el-button type="danger" size="small" @click="removeImage">删除图片</el-button>
          </div>
        </div>

        <div class="form-item">
          <label>是否上架：</label>
          <el-switch 
            v-model="productForm.status" 
            active-value="1" 
            inactive-value="0"
            active-text="上架"
            inactive-text="下架"
          ></el-switch>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="submit-btn">
        <el-button 
          type="primary" 
          size="large" 
          @click="submitProduct" 
          :loading="submitting"
        >
          {{ isEditMode ? '更新商品' : '提交新增' }}
        </el-button>
        <el-button size="large" @click="resetForm">重置表单</el-button>
      </div>
    </div>

    <!-- 添加/编辑父分类弹窗 -->
    <el-dialog 
      :title="isCategoryEditMode ? '编辑父分类' : '添加父分类'" 
      v-model="showParentCategoryDialog" 
      width="400px"
      @close="closeCategoryDialog"
    >
      <el-form :model="categoryForm" label-width="80px">
        <el-form-item label="分类名称：" required>
          <el-input v-model="categoryForm.name" placeholder="请输入分类名称" maxlength="20"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showParentCategoryDialog = false">取消</el-button>
        <el-button type="primary" @click="saveParentCategory" :loading="saving">确定</el-button>
      </template>
    </el-dialog>

    <!-- 添加/编辑子分类弹窗 -->
    <el-dialog 
      :title="isCategoryEditMode ? '编辑子分类' : '添加子分类'" 
      v-model="showChildCategoryDialog" 
      width="400px"
      @close="closeCategoryDialog"
    >
      <el-form :model="categoryForm" label-width="80px">
        <el-form-item label="分类名称：" required>
          <el-input v-model="categoryForm.name" placeholder="请输入分类名称" maxlength="20"></el-input>
        </el-form-item>
        <el-form-item label="所属父分类：">
          <el-input :value="getParentName(selectedParentId)" disabled></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showChildCategoryDialog = false">取消</el-button>
        <el-button type="primary" @click="saveChildCategory" :loading="saving">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import request from '../../utils/request'

const router = useRouter()

// 分类数据
const parentCategories = ref([])
const childCategories = ref([])
const selectedParentId = ref('')
const selectedChildId = ref('')

// 商品选择相关
const productOptions = ref([])
const selectedProductId = ref('')
const isEditMode = ref(false)

// 弹窗控制
const showParentCategoryDialog = ref(false)
const showChildCategoryDialog = ref(false)
const isCategoryEditMode = ref(false)
const editCategoryId = ref('')

// 表单数据
const categoryForm = ref({ name: '' })
const productForm = ref({
  name: '',
  price: 0.01,
  stock: 0,
  image: '',
  status: '1'
})

// 原生文件上传引用
const fileInput = ref(null)

// 加载状态
const submitting = ref(false)
const saving = ref(false)

// 加载父分类列表
const loadParentCategories = async () => {
  try {
    parentCategories.value = await request({
      url: '/product/category/parents',
      method: 'get'
    })
  } catch (e) {
    console.error('加载父分类失败：', e)
    ElMessage.error('加载分类列表失败')
  }
}

// 父分类变更时加载对应子分类
const handleParentChange = async (pid) => {
  selectedChildId.value = ''
  selectedProductId.value = '' // 清空商品选择
  productOptions.value = []
  if (!pid) {
    childCategories.value = []
    return
  }
  try {
    childCategories.value = await request({
      url: '/product/category/children',
      method: 'get',
      params: { parent_id: pid }
    })
  } catch (e) {
    console.error('加载子分类失败：', e)
    ElMessage.error('加载子分类失败')
  }
}

// 子分类变更时加载该分类下的所有商品
const handleChildChange = async (cid) => {
  selectedProductId.value = ''
  productOptions.value = []
  if (!cid) return
  
  try {
    productOptions.value = await request({
      url: '/product/list',
      method: 'get',
      params: { category_id: cid }
    })
  } catch (e) {
    console.error('加载商品列表失败：', e)
    ElMessage.error('加载该分类下的商品失败')
  }
}

// 选择商品后自动填充表单
const handleProductSelect = async (productId) => {
  if (!productId) {
    // 清空选择，回到新增模式
    isEditMode.value = false
    resetForm()
    return
  }

  try {
    const product = await request({
      url: `/product/${productId}`,
      method: 'get'
    })

    // 自动填充所有表单字段
    productForm.value = {
      name: product.name,
      price: product.price,
      stock: product.stock,
      image: product.image,
      status: String(product.status) // 转换为字符串匹配switch
    }

    isEditMode.value = true
    ElMessage.success('商品信息已加载，可进行编辑')
  } catch (e) {
    console.error('获取商品详情失败：', e)
    ElMessage.error('加载商品信息失败')
  }
}

// 删除商品
const deleteProduct = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该商品吗？删除后无法恢复',
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await request({
      url: `/admin/product/${selectedProductId.value}`,
      method: 'delete'
    })

    ElMessage.success('商品删除成功')
    
    // 重置表单并重新加载商品列表
    resetForm()
    handleChildChange(selectedChildId.value)
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message)
    }
  }
}

// 打开添加父分类弹窗
const openAddParentCategory = () => {
  isCategoryEditMode.value = false
  categoryForm.value = { name: '' }
  showParentCategoryDialog.value = true
}

// 打开编辑父分类弹窗
const openEditParentCategory = () => {
  const parent = parentCategories.value.find(p => p.id === selectedParentId.value)
  if (!parent) return
  isCategoryEditMode.value = true
  editCategoryId.value = parent.id
  categoryForm.value = { name: parent.name }
  showParentCategoryDialog.value = true
}

// 保存父分类（新增/编辑）
const saveParentCategory = async () => {
  if (!categoryForm.value.name.trim()) {
    ElMessage.warning('请输入分类名称')
    return
  }
  saving.value = true
  try {
    if (isCategoryEditMode.value) {
      // 编辑模式
      await request({
        url: `/admin/category/${editCategoryId.value}`,
        method: 'put',
        data: { name: categoryForm.value.name }
      })
    } else {
      // 新增模式
      await request({
        url: '/admin/category/add',
        method: 'post',
        data: {
          name: categoryForm.value.name,
          parent_id: 0,
          status: 1
        }
      })
    }
    ElMessage.success(isCategoryEditMode.value ? '父分类编辑成功' : '父分类添加成功')
    showParentCategoryDialog.value = false
    loadParentCategories()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

// 删除父分类
const deleteParentCategory = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该父分类吗？删除后无法恢复',
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await request({
      url: `/admin/category/${selectedParentId.value}`,
      method: 'delete'
    })
    ElMessage.success('父分类删除成功')
    selectedParentId.value = ''
    childCategories.value = []
    loadParentCategories()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message)
    }
  }
}

// 打开添加子分类弹窗
const openAddChildCategory = () => {
  isCategoryEditMode.value = false
  categoryForm.value = { name: '' }
  showChildCategoryDialog.value = true
}

// 打开编辑子分类弹窗
const openEditChildCategory = () => {
  const child = childCategories.value.find(c => c.id === selectedChildId.value)
  if (!child) return
  isCategoryEditMode.value = true
  editCategoryId.value = child.id
  categoryForm.value = { name: child.name }
  showChildCategoryDialog.value = true
}

// 保存子分类（新增/编辑）
const saveChildCategory = async () => {
  if (!categoryForm.value.name.trim()) {
    ElMessage.warning('请输入分类名称')
    return
  }
  saving.value = true
  try {
    if (isCategoryEditMode.value) {
      // 编辑模式
      await request({
        url: `/admin/category/${editCategoryId.value}`,
        method: 'put',
        data: { name: categoryForm.value.name }
      })
    } else {
      // 新增模式
      await request({
        url: '/admin/category/add',
        method: 'post',
        data: {
          name: categoryForm.value.name,
          parent_id: selectedParentId.value,
          status: 1
        }
      })
    }
    ElMessage.success(isCategoryEditMode.value ? '子分类编辑成功' : '子分类添加成功')
    showChildCategoryDialog.value = false
    handleParentChange(selectedParentId.value)
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

// 删除子分类
const deleteChildCategory = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该子分类吗？删除后无法恢复',
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await request({
      url: `/admin/category/${selectedChildId.value}`,
      method: 'delete'
    })
    ElMessage.success('子分类删除成功')
    selectedChildId.value = ''
    handleParentChange(selectedParentId.value)
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message)
    }
  }
}

// 关闭分类弹窗时重置状态
const closeCategoryDialog = () => {
  isCategoryEditMode.value = false
  editCategoryId.value = ''
  categoryForm.value = { name: '' }
}

// 根据父分类ID获取名称
const getParentName = (pid) => {
  const parent = parentCategories.value.find(p => p.id === pid)
  return parent ? parent.name : ''
}

// 触发文件选择对话框
const triggerFileSelect = () => {
  fileInput.value.click()
}

// 处理文件上传
const handleFileUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return

  // 校验文件类型
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    ElMessage.error('仅支持上传 jpg、png、gif 格式的图片')
    e.target.value = ''
    return
  }

  // 校验文件大小（最大10MB）
  if (file.size > 10 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过10MB')
    e.target.value = ''
    return
  }

  // 创建FormData对象
  const formData = new FormData()
  formData.append('image', file)

  try {
    const data = await request({
      url: '/admin/upload',
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    productForm.value.image = data.url
    ElMessage.success('图片上传成功')
  } catch (e) {
    console.error('图片上传失败:', e)
    ElMessage.error(e.message || '图片上传失败')
  } finally {
    // 清空文件输入框，允许重复上传同一张图片
    e.target.value = ''
  }
}

// 删除已上传的图片
const removeImage = () => {
  productForm.value.image = ''
}

// 提交商品（支持新增/编辑两种模式）
const submitProduct = async () => {
  // 表单验证
  if (!productForm.value.name.trim()) {
    return ElMessage.warning('请输入商品名称')
  }
  if (productForm.value.price < 0.01) {
    return ElMessage.warning('商品价格不能低于0.01元')
  }
  if (productForm.value.stock < 0) {
    return ElMessage.warning('库存数量不能为负数')
  }
  if (!selectedParentId.value) {
    return ElMessage.warning('请至少选择父分类')
  }

  submitting.value = true
  try {
    // 优先使用子分类ID，没有则使用父分类ID
    const categoryId = parseInt(selectedChildId.value || selectedParentId.value)
    const requestData = {
      name: productForm.value.name,
      price: productForm.value.price,
      stock: productForm.value.stock,
      image: productForm.value.image,
      status: parseInt(productForm.value.status),
      category_id: categoryId
    }

    if (isEditMode.value) {
      // 编辑模式：调用更新接口
      await request({
        url: `/admin/product/${selectedProductId.value}`,
        method: 'put',
        data: requestData
      })
    } else {
      // 新增模式：调用新增接口
      await request({
        url: '/admin/product',
        method: 'post',
        data: requestData
      })
    }
    
    ElMessage.success(isEditMode.value ? '商品更新成功' : '商品新增成功')
    router.push('/product/list')
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    submitting.value = false
  }
}

// 重置表单
const resetForm = () => {
  productForm.value = {
    name: '',
    price: 0.01,
    stock: 0,
    image: '',
    status: '1'
  }
  selectedParentId.value = ''
  selectedChildId.value = ''
  selectedProductId.value = ''
  childCategories.value = []
  productOptions.value = []
  isEditMode.value = false
}

onMounted(() => {
  loadParentCategories()
})
</script>

<style scoped>
.add-product-container {
  max-width: 850px;
  margin: 20px auto;
  padding: 0 20px;
}
.card {
  background: #fff;
  padding: 35px;
  border-radius: 10px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}
h2 { 
  margin-bottom: 25px; 
  font-size: 20px; 
  color: #333; 
  border-bottom: 2px solid #e68200; 
  padding-bottom: 10px;
}
h3 { 
  margin: 20px 0 15px; 
  font-size: 16px; 
  color: #666; 
  border-bottom: 1px solid #eee; 
  padding-bottom: 8px;
}
.form-item {
  margin-bottom: 18px;
  display: flex;
  align-items: center;
  gap: 10px;
}
.form-item label {
  width: 90px;
  text-align: right;
  color: #666;
  font-weight: 500;
}
.form-item .el-input, .form-item .el-select { 
  flex: 1; 
  max-width: 350px;
}
.unit {
  color: #999;
  margin-left: 5px;
}
.image-preview {
  margin-left: 100px;
  margin-top: 10px;
}
.image-preview img {
  width: 120px;
  height: 120px;
  object-fit: cover;
  border: 1px solid #eee;
  border-radius: 4px;
  display: block;
  margin-bottom: 8px;
}
.submit-btn { 
  margin-top: 30px; 
  display: flex; 
  gap: 15px; 
  justify-content: center; 
}
.category-section {
  margin-bottom: 30px;
}
</style>