<template>
  <div>
    <div class="gva-search-box">
      <el-form
        :inline="true"
        class="demo-form-inline"
      >
        <el-form-item label="父分类">
          <el-select
            v-model="filterParentId"
            clearable
            placeholder="全部"
            style="width: 200px"
            @clear="filterParentId = undefined"
          >
            <el-option
              label="仅一级分类"
              :value="0"
            />
            <el-option
              v-for="c in categoryOptions"
              :key="c.ID"
              :label="c.categoryName"
              :value="c.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="filterStatus"
            clearable
            placeholder="全部"
            style="width: 140px"
          >
            <el-option
              label="启用"
              :value="true"
            />
            <el-option
              label="禁用"
              :value="false"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSearch"
          >
            查询
          </el-button>
          <el-button
            icon="refresh"
            @click="resetSearch"
          >
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="openDrawer('create')"
        >
          新增分类
        </el-button>
      </div>

      <el-table
        :data="tableData"
        style="width: 100%"
        border
        row-key="ID"
        @sort-change="handleSortChange"
      >
        <el-table-column
          align="left"
          label="分类ID"
          prop="ID"
          min-width="90"
        />
        <el-table-column
          align="left"
          label="分类名称"
          prop="categoryName"
          min-width="140"
        />
        <el-table-column
          align="left"
          label="父分类"
          min-width="120"
        >
          <template #default="scope">
            {{ scope.row.parentName || '-' }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="排序"
          prop="categorySort"
          sortable="custom"
          min-width="100"
        />
        <el-table-column
          align="left"
          label="状态"
          min-width="100"
        >
          <template #default="scope">
            <el-tag :type="scope.row.categoryStatus ? 'success' : 'info'">
              {{ scope.row.categoryStatus ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="商品数量"
          min-width="100"
        >
          <template #default="scope">
            {{ scope.row.productCount ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="创建时间"
          prop="CreatedAt"
          sortable="custom"
          min-width="170"
        >
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) || '-' }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="plus"
              @click="openDrawer('createChild', scope.row)"
            >
              添加子分类
            </el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              @click="openDrawer('update', scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="onDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="drawerVisible"
      :title="drawerTitle"
      size="520px"
      :before-close="closeDrawer"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="110px"
      >
        <el-form-item
          label="分类名称"
          prop="categoryName"
        >
          <el-input
            v-model="form.categoryName"
            placeholder="请输入分类名称"
            maxlength="64"
            show-word-limit
          />
        </el-form-item>
        <el-form-item
          v-if="drawerMode === 'update'"
          label="父分类"
          prop="categoryParentID"
        >
          <el-select
            v-model="form.categoryParentID"
            placeholder="无（一级分类）"
            style="width: 100%"
            clearable
          >
            <el-option
              label="无"
              :value="0"
            />
            <el-option
              v-for="c in parentSelectOptions"
              :key="c.ID"
              :label="c.categoryName"
              :value="c.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="drawerMode === 'createChild'"
          label="父分类"
        >
          <el-input
            :model-value="childParentName"
            readonly
          />
        </el-form-item>
        <el-form-item
          label="排序"
          prop="categorySort"
        >
          <el-input-number
            v-model="form.categorySort"
            :min="0"
            :max="999999"
            controls-position="right"
            style="width: 100%"
          />
          <div class="text-gray-400 text-xs mt-1">
            数字越小越靠前
          </div>
        </el-form-item>
        <el-form-item
          label="分类描述"
          prop="categoryRemark"
        >
          <el-input
            v-model="form.categoryRemark"
            type="textarea"
            :rows="3"
            placeholder="最多100字"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="分类图标">
          <el-upload
            class="category-icon-uploader"
            :action="`${baseUrl}/fileUploadAndDownload/upload`"
            :headers="uploadHeaders"
            :show-file-list="false"
            :on-success="handleIconSuccess"
            :before-upload="beforeIconUpload"
          >
            <img
              v-if="form.categoryIcon"
              :src="displayIconUrl(form.categoryIcon)"
              class="icon-preview"
            >
            <el-icon
              v-else
              class="icon-uploader-placeholder"
            >
              <Plus />
            </el-icon>
          </el-upload>
          <div class="text-gray-400 text-xs mt-1">
            建议 40×40 像素，JPG / PNG
          </div>
        </el-form-item>
        <el-form-item
          label="状态"
          prop="categoryStatus"
        >
          <el-switch
            v-model="form.categoryStatus"
            inline-prompt
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
        <el-form-item>
          <el-button @click="closeDrawer">
            取 消
          </el-button>
          <el-button
            type="primary"
            :loading="submitLoading"
            @click="submitForm"
          >
            保 存
          </el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useUserStore } from '@/pinia/modules/user'
import { formatDate, getBaseUrl } from '@/utils/format'
import { getUrl } from '@/utils/image'
import {
  getProductCategoryList,
  getProductCategorySelectList,
  getProductCategory,
  createProductCategory,
  updateProductCategory,
  deleteProductCategory
} from '@/api/productCategory'

defineOptions({
  name: 'CategoryList'
})

const baseUrl = getBaseUrl()
const userStore = useUserStore()
const uploadHeaders = computed(() => ({ 'x-token': userStore.token }))

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

const sortField = ref('categorySort')
const sortOrder = ref('asc')

const filterParentId = ref(undefined)
const filterStatus = ref(undefined)

const categoryOptions = ref([])

const loadCategoryOptions = async () => {
  const res = await getProductCategorySelectList()
  if (res.code === 0) {
    categoryOptions.value = res.data || []
  }
}

const buildQuery = () => {
  const q = {
    page: page.value,
    pageSize: pageSize.value,
    sortField: sortField.value,
    sortOrder: sortOrder.value
  }
  if (filterParentId.value !== undefined && filterParentId.value !== null && filterParentId.value !== '') {
    q.categoryParentID = filterParentId.value
  }
  if (filterStatus.value !== undefined && filterStatus.value !== null && filterStatus.value !== '') {
    q.categoryStatus = filterStatus.value
  }
  return q
}

const getTableData = async () => {
  const res = await getProductCategoryList(buildQuery())
  if (res.code === 0 && res.data) {
    tableData.value = res.data.list || []
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

const handleSortChange = ({ prop, order }) => {
  if (!order) {
    sortField.value = 'categorySort'
    sortOrder.value = 'asc'
  } else if (prop === 'categorySort') {
    sortField.value = 'categorySort'
    sortOrder.value = order === 'ascending' ? 'asc' : 'desc'
  } else if (prop === 'CreatedAt') {
    sortField.value = 'CreatedAt'
    sortOrder.value = order === 'ascending' ? 'asc' : 'desc'
  }
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const onSearch = () => {
  page.value = 1
  getTableData()
}

const resetSearch = () => {
  filterParentId.value = undefined
  filterStatus.value = undefined
  page.value = 1
  getTableData()
}

const drawerVisible = ref(false)
const drawerMode = ref('create')
const drawerTitle = computed(() => {
  if (drawerMode.value === 'update') return '编辑分类'
  if (drawerMode.value === 'createChild') return '新增子分类'
  return '新增分类'
})
const submitLoading = ref(false)
const formRef = ref(null)
const childParentName = ref('')
const categoryDetailCache = ref({})
const MAX_CATEGORY_LEVEL = 3

const defaultForm = () => ({
  ID: 0,
  categoryName: '',
  categoryCode: '',
  categoryParentID: 0,
  categorySort: 0,
  categoryRemark: '',
  categoryIcon: '',
  categoryStatus: true
})

const form = ref(defaultForm())

const rules = {
  categoryName: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const parentSelectOptions = computed(() => {
  const selfId = form.value.ID
  return (categoryOptions.value || []).filter((c) => c.ID !== selfId)
})

const getCategoryDetailByID = async (id) => {
  if (!id) return null
  if (categoryDetailCache.value[id]) return categoryDetailCache.value[id]
  const res = await getProductCategory({ ID: id })
  if (res.code !== 0 || !res.data) return null
  categoryDetailCache.value[id] = res.data
  return res.data
}

const getCategoryLevel = async (row) => {
  let level = 1
  let parentID = row.categoryParentID
  const visited = new Set([row.ID])
  while (parentID && parentID > 0) {
    if (visited.has(parentID)) break
    visited.add(parentID)
    level++
    const parent = await getCategoryDetailByID(parentID)
    if (!parent) break
    parentID = parent.categoryParentID
  }
  return level
}

const openDrawer = async (mode, row) => {
  drawerMode.value = mode
  childParentName.value = ''
  await loadCategoryOptions()
  if (mode === 'create') {
    form.value = defaultForm()
  } else if (mode === 'createChild' && row) {
    const parentLevel = await getCategoryLevel(row)
    if (parentLevel >= MAX_CATEGORY_LEVEL) {
      ElMessage.warning(`最多支持${MAX_CATEGORY_LEVEL}级分类，当前分类已是第${parentLevel}级`)
      return
    }
    form.value = defaultForm()
    form.value.categoryParentID = row.ID
    childParentName.value = row.categoryName
  } else if (row) {
    const detail = await getCategoryDetailByID(row.ID)
    if (detail) {
      form.value = { ...defaultForm(), ...detail }
    }
  }
  drawerVisible.value = true
}

const closeDrawer = () => {
  drawerVisible.value = false
  childParentName.value = ''
  form.value = defaultForm()
}

const beforeIconUpload = (file) => {
  const ok = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!ok) {
    ElMessage.error('仅支持 JPG、PNG 格式')
    return false
  }
  return new Promise((resolve) => {
    const img = new Image()
    const url = URL.createObjectURL(file)
    img.onload = () => {
      URL.revokeObjectURL(url)
      if (img.width !== 40 || img.height !== 40) {
        ElMessage.warning(`图标建议尺寸 40×40，当前 ${img.width}×${img.height}，仍可上传`)
      }
      resolve(true)
    }
    img.onerror = () => {
      URL.revokeObjectURL(url)
      resolve(true)
    }
    img.src = url
  })
}

/** 上传接口返回的地址多为相对路径，预览需拼接 VITE_FILE_API（见 @/utils/image getUrl） */
const displayIconUrl = (u) => (u ? getUrl(u) : '')

/** el-upload 的 response 可能是 { code, data, msg }，也可能被封装成其它形态 */
const pickUploadFileUrl = (res) => {
  if (!res) return ''
  const payload = res.data !== undefined ? res.data : res
  const file = payload?.file ?? payload?.File
  if (!file) return ''
  return file.url || file.Url || ''
}

const handleIconSuccess = (res, uploadFile) => {
  const url = pickUploadFileUrl(res)
  if (url) {
    form.value.categoryIcon = url
    return
  }
  // 个别环境下后端把地址放在其它字段，做一次兜底
  const fallback = uploadFile?.response && pickUploadFileUrl(uploadFile.response)
  if (fallback) {
    form.value.categoryIcon = fallback
    return
  }
  ElMessage.warning('上传成功，但未解析到图片地址，请检查接口返回结构')
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      let res
      if (drawerMode.value === 'create' || drawerMode.value === 'createChild') {
        res = await createProductCategory(form.value)
      } else {
        res = await updateProductCategory(form.value)
      }
      if (res.code === 0) {
        ElMessage.success(res.msg || '操作成功')
        closeDrawer()
        await loadCategoryOptions()
        getTableData()
      }
    } finally {
      submitLoading.value = false
    }
  })
}

const onDelete = (row) => {
  ElMessageBox.confirm(
    '删除后不可恢复。若存在子分类将无法删除；商品归属校验将在商品模块接入后生效。',
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    const res = await deleteProductCategory({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
      loadCategoryOptions()
    }
  }).catch(() => {})
}

onMounted(() => {
  loadCategoryOptions()
  getTableData()
})
</script>

<style scoped lang="scss">
.category-icon-uploader {
  :deep(.el-upload) {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
.icon-preview {
  width: 40px;
  height: 40px;
  display: block;
  object-fit: cover;
}
.icon-uploader-placeholder {
  font-size: 28px;
  color: #8c939d;
}
</style>
