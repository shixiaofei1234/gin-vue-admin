<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="openDrawer"
        >
          创建员工
        </el-button>
      </div>
      <el-table
        :data="tableData"
        style="width: 100%"
        border
        row-key="ID"
      >
        <el-table-column
          align="left"
          label="ID"
          prop="ID"
          width="80"
        />
        <el-table-column
          align="left"
          label="员工姓名"
          prop="employeeName"
          min-width="100"
        />
        <el-table-column
          align="left"
          label="员工电话"
          prop="employeePhone"
          min-width="120"
        />
        <el-table-column
          align="left"
          label="员工邮箱"
          prop="employeeEmail"
          min-width="160"
        />
        <el-table-column
          align="left"
          label="员工地址"
          prop="employeeAddress"
          min-width="180"
        />
        <el-table-column
          align="left"
          label="员工生日"
          prop="employeeBirthday"
          width="120"
        />
        <el-table-column
          align="left"
          label="员工性别"
          prop="employeeGenderStr"
          width="90"
        >
          <template #default="scope">
            {{ getGenderText(scope.row) }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="员工职位"
          prop="employeePosition"
          min-width="100"
        />
        <el-table-column
          align="left"
          label="员工部门"
          prop="employeeDepartment"
          min-width="100"
        />
        <el-table-column
          align="left"
          label="操作"
          min-width="180"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="view"
              @click="openDetail(scope.row)"
            >
              详情
            </el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              @click="openEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="handleDelete(scope.row)"
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
    <EmployeeFormDrawer
      v-model="drawerVisible"
      @submit="handleCreateEmployee"
      :title="drawerTitle"
      :submit-text="drawerSubmitText"
      :initial-value="editingRow"
    />
    <el-dialog
      v-model="detailVisible"
      title="员工详情"
      width="600px"
    >
      <el-descriptions
        v-if="detailRow"
        :column="2"
        border
      >
        <el-descriptions-item label="ID">{{ detailRow.ID }}</el-descriptions-item>
        <el-descriptions-item label="员工姓名">{{ detailRow.employeeName }}</el-descriptions-item>
        <el-descriptions-item label="员工电话">{{ detailRow.employeePhone }}</el-descriptions-item>
        <el-descriptions-item label="员工邮箱">{{ detailRow.employeeEmail }}</el-descriptions-item>
        <el-descriptions-item label="员工地址">{{ detailRow.employeeAddress || '-' }}</el-descriptions-item>
        <el-descriptions-item label="员工生日">{{ detailRow.employeeBirthday || '-' }}</el-descriptions-item>
        <el-descriptions-item label="员工性别">
          {{ getGenderText(detailRow) }}
        </el-descriptions-item>
        <el-descriptions-item label="员工职位">{{ detailRow.employeePosition || '-' }}</el-descriptions-item>
        <el-descriptions-item label="员工部门">{{ detailRow.employeeDepartment || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关 闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getEmployeeList, createEmployee, updateEmployee, getEmployeeDetail, deleteEmployee } from '@/api/employee'
import EmployeeFormDrawer from './EmployeeFormDrawer.vue'
import { emitter } from '@/utils/bus.js'
import { AI_FLOAT_STORAGE } from '@/constants/aiFloatChat'

// AI 悬浮窗/聊天发出的“打开员工详情”快捷指令会通过两种方式到这里：
// 1) 当前就在员工列表页：直接 emitter.emit('openEmployeeDetail', meta) 触发 handleOpenEmployeeFromAi
// 2) 不在列表页：aiFloatChat 会写入 sessionStorage，列表页 onMounted 时再消费（tryConsumePendingEmployeeFromAi）

defineOptions({
  name: 'EmployeeList'
})

const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const drawerVisible = ref(false)
const drawerTitle = ref('创建员工')
const drawerSubmitText = ref('确 定')
const editingRow = ref(null)
const editMode = ref('create') // create | edit

const detailVisible = ref(false)
const detailRow = ref(null)

const getGenderText = (row) => {
  if (!row) return '-'
  if (row.employeeGenderStr) return row.employeeGenderStr
  if (row.employeeGender === 1) return '男'
  if (row.employeeGender === 2) return '女'
  if (row.employeeGender === 0) return '未知'
  return row.employeeGender ?? '-'
}

const openDrawer = () => {
  editMode.value = 'create'
  editingRow.value = null
  drawerTitle.value = '创建员工'
  drawerSubmitText.value = '确 定'
  drawerVisible.value = true
}

const handleCreateEmployee = async(payload) => {
  if (editMode.value === 'edit' && editingRow.value) {
    const res = await updateEmployee({
      ID: editingRow.value.ID,
      ...payload
    })
    if (res.code === 0) {
      ElMessage.success('更新成功')
      fetchEmployeeList()
    }
    drawerVisible.value = false
    return
  }

  const res = await createEmployee(payload)
  if (res.code === 0) {
    ElMessage.success('创建成功')
    page.value = 1
    fetchEmployeeList()
  }
}

const openDetail = async(row) => {
  const res = await getEmployeeDetail({ ID: row.ID })
  if (res.code === 0) {
    detailRow.value = res.data
    detailVisible.value = true
    return
  }
  ElMessage.error(res.msg || '获取员工详情失败')
}

/** 按姓名在当前列表数据中查找后打开详情弹窗（必要时会额外拉取第一页数据兜底） */
const tryOpenDetailByName = async(employeeName) => {
  const name = (employeeName || '').trim()
  if (!name) return
  const matchRow = (row) => {
    if (!row?.employeeName) return false
    return row.employeeName === name || String(row.employeeName).includes(name)
  }
  let found = tableData.value.find(matchRow)
  if (found) {
    await openDetail(found)
    return
  }
  const res = await getEmployeeList({ page: 1, pageSize: 100 })
  if (res.code !== 0) {
    ElMessage.error(res.msg || '加载员工列表失败')
    return
  }
  const list = res.data.list ?? []
  found = list.find(matchRow)
  if (found) {
    await openDetail(found)
    return
  }
  ElMessage.warning(`未在列表中找到「${name}」，请核对姓名或翻页后再试。`)
}

// 处理 AI 端发来的 meta（可能包含 employeeId 或 employeeName）。
// 这里会复用现有的 getEmployeeDetail + detailVisible 弹窗逻辑。
const handleOpenEmployeeFromAi = async(meta) => {
  try {
    sessionStorage.removeItem(AI_FLOAT_STORAGE.openEmployee)
  } catch (_) {}
  if (!meta) return
  const { employeeId, employeeName } = meta
  if (employeeId != null) {
    const res = await getEmployeeDetail({ ID: employeeId })
    if (res.code === 0) {
      detailRow.value = res.data
      detailVisible.value = true
      return
    }
    ElMessage.error(res.msg || '获取员工详情失败')
    return
  }
  if (employeeName) {
    await tryOpenDetailByName(employeeName)
  }
}

// 消费 AI 写入的待执行指令（通常用于：用户不在员工列表页时先发起“打开详情”）
const tryConsumePendingEmployeeFromAi = () => {
  let raw
  try {
    raw = sessionStorage.getItem(AI_FLOAT_STORAGE.openEmployee)
  } catch (_) {
    return
  }
  if (!raw) return
  try {
    const data = JSON.parse(raw)
    const ts = data?.ts
    if (ts && Date.now() - ts > 120000) {
      sessionStorage.removeItem(AI_FLOAT_STORAGE.openEmployee)
      return
    }
    const { ts: _t, ...meta } = data
    if (meta && (meta.employeeId != null || meta.employeeName)) {
      handleOpenEmployeeFromAi(meta)
    }
  } catch (_) {
    sessionStorage.removeItem(AI_FLOAT_STORAGE.openEmployee)
  }
}

const openEdit = (row) => {
  editMode.value = 'edit'
  editingRow.value = { ...row }
  drawerTitle.value = '编辑员工'
  drawerSubmitText.value = '保 存'
  drawerVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该员工吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await deleteEmployee({ ID: row.ID })
    if (res.code !== 0) {
      ElMessage.error(res.msg || '删除失败')
      return
    }
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && page.value > 1) {
      page.value -= 1
    }
    fetchEmployeeList()
  })
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchEmployeeList()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchEmployeeList()
}

// 查询员工列表（GET，分页：page、pageSize；后端需返回 list、total、page、pageSize）
const fetchEmployeeList = async() => {
  const res = await getEmployeeList({ page: page.value, pageSize: pageSize.value })
  if (res.code === 0) {
    tableData.value = res.data.list ?? []
    total.value = res.data.total ?? 0
    page.value = res.data.page ?? page.value
    pageSize.value = res.data.pageSize ?? pageSize.value
  }
}

onMounted(async() => {
  emitter.on('openEmployeeDetail', handleOpenEmployeeFromAi)
  await fetchEmployeeList()
  tryConsumePendingEmployeeFromAi()
})

onBeforeUnmount(() => {
  emitter.off('openEmployeeDetail', handleOpenEmployeeFromAi)
})
</script>

<style scoped></style>
