<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="openDrawer"
        >
          创建团队
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
          label="团队编号"
          prop="teamNum"
          min-width="120"
        />
        <el-table-column
          align="left"
          label="团队名称"
          prop="teamName"
          min-width="150"
        />
        <el-table-column
          align="left"
          label="组织"
          prop="organize"
          min-width="120"
        />
        <el-table-column
          align="left"
          label="员工列表"
          min-width="220"
        >
          <template #default="scope">
            {{ formatEmployeeList(scope.row.employeeList) }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="管理员名称"
          prop="adminName"
          min-width="140"
        />
        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          min-width="140"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              @click="openAdminDialog(scope.row)"
            >
              设置管理员
            </el-button>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="团队状态"
          prop="status"
          min-width="140"
        >
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              inline-prompt
              active-text="开启"
              inactive-text="关闭"
              :active-value="true"
              :inactive-value="false"
              :loading="!!switchLoadingMap[scope.row.ID]"
              @change="handleSwitchStatus(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="创建时间"
          min-width="180"
        >
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) || '-' }}
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
      title="创建团队"
      width="520px"
      :before-close="closeDrawer"
      :show-close="true"
    >
      <el-form
        :model="form"
        label-width="110px"
        ref="formRef"
        :rules="rules"
      >
        <el-form-item
          label="团队名称"
          prop="teamName"
        >
          <el-input
            v-model="form.teamName"
            placeholder="请输入团队名称"
          />
        </el-form-item>

        <el-form-item>
          <div class="flex gap-2">
            <el-button @click="closeDrawer">取 消</el-button>
            <el-button
              type="primary"
              @click="handleCreateTeam"
            >
              确 定
            </el-button>
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-dialog
      v-model="adminDialogVisible"
      title="设置团队管理员"
      width="520px"
      destroy-on-close
      @closed="resetAdminDialog"
    >
      <div
        v-if="currentTeamRow"
        class="admin-dialog-meta"
      >
        团队：{{ currentTeamRow.teamName || '-' }}
        <span class="admin-dialog-meta-gap">当前管理员：{{ currentTeamRow.adminName || '未设置' }}</span>
      </div>
      <div
        v-loading="empLoading"
        class="admin-emp-scroll"
        @scroll.passive="onEmpScroll"
      >
        <el-empty
          v-if="!empLoading && empList.length === 0"
          description="该团队暂无员工"
        />
        <el-radio-group
          v-else
          v-model="selectedAdminId"
          class="admin-radio-group"
        >
          <el-radio
            v-for="e in empList"
            :key="e.ID"
            :value="e.ID"
            class="admin-radio-item"
          >
            <span>{{ e.employeeName }}</span>
            <span
              v-if="e.employeePhone"
              class="admin-radio-phone"
            >{{ e.employeePhone }}</span>
          </el-radio>
        </el-radio-group>
      </div>
      <template #footer>
        <el-button @click="adminDialogVisible = false">取 消</el-button>
        <el-button
          type="primary"
          :loading="adminSubmitLoading"
          @click="handleConfirmAdmin"
        >
          确 定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import service from '@/utils/request'
import { formatDate } from '@/utils/format'
import { getTeamEmployeeList, setTeamAdmin } from '@/api/team'
import { emitter } from '@/utils/bus.js'

defineOptions({
  name: 'TeamList'
})

// 你后端自行添加路由/接口时，可以按需修改这两个 URL。
const TEAM_LIST_URL = '/team/getTeamList'
const TEAM_CREATE_URL = '/team/createTeam'
const TEAM_SWITCH_STATUS_URL = '/team/switchTeamStatus'

const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const drawerVisible = ref(false)
const switchLoadingMap = ref({})

const adminDialogVisible = ref(false)
const currentTeamRow = ref(null)
const adminTeamId = ref(null)
const empList = ref([])
const empPage = ref(1)
const empTotal = ref(0)
const empLoading = ref(false)
const selectedAdminId = ref(null)
const adminSubmitLoading = ref(false)

const emptyForm = () => ({
  teamName: '',
})

const form = ref(emptyForm())
const formRef = ref(null)

const rules = {
  teamName: [
    {
      validator: (rule, value, callback) => {
        if (!value || !String(value).trim()) {
          callback(new Error('团队名称必填'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const formatEmployeeList = (val) => {
  if (!val) return '-'
  if (Array.isArray(val)) {
    const arr = val
      .map((v) => {
        if (v === null || v === undefined) return ''
        if (typeof v === 'string' || typeof v === 'number') return String(v)
        if (v.employeeName) return v.employeeName
        if (v.ID !== undefined) return String(v.ID)
        return ''
      })
      .filter(Boolean)
    return arr.length ? arr.join(',') : '-'
  }
  if (typeof val === 'string') return val
  if (val.employeeName) return val.employeeName
  return '-'
}

const getTeamList = async() => {
  const res = await service({
    url: TEAM_LIST_URL,
    method: 'get',
    params: {
      page: page.value,
      pageSize: pageSize.value
    }
  })

  if (res?.code === 0) {
    tableData.value = res.data?.list ?? []
    total.value = res.data?.total ?? 0
    page.value = res.data?.page ?? page.value
    pageSize.value = res.data?.pageSize ?? pageSize.value
  }
}

const createTeam = async(data) => {
  return service({
    url: TEAM_CREATE_URL,
    method: 'post',
    data
  })
}

const switchTeamStatus = async(data) => {
  return service({
    url: TEAM_SWITCH_STATUS_URL,
    method: 'post',
    data
  })
}

const handleSwitchStatus = async(row) => {
  if (!row || !row.ID) {
    ElMessage.error('团队ID缺失，无法切换状态')
    return
  }
  const id = row.ID
  const nextStatus = Boolean(row.status)
  const prevStatus = !nextStatus
  switchLoadingMap.value[id] = true
  try {
    const res = await switchTeamStatus({
      ID: id,
      status: nextStatus
    })
    if (res?.code !== 0) {
      row.status = prevStatus
      ElMessage.error(res?.msg || '切换团队状态失败')
      return
    }
    ElMessage.success(nextStatus ? '团队已开启' : '团队已关闭')
  } catch (e) {
    row.status = prevStatus
    ElMessage.error('切换团队状态失败')
  } finally {
    switchLoadingMap.value[id] = false
  }
}

const resetAdminDialog = () => {
  currentTeamRow.value = null
  adminTeamId.value = null
  empList.value = []
  empPage.value = 1
  empTotal.value = 0
  selectedAdminId.value = null
}

const ensureCurrentAdminInList = () => {
  const row = currentTeamRow.value
  if (!row?.adminID) return
  if (empList.value.some((e) => e.ID === row.adminID)) return
  empList.value.unshift({
    ID: row.adminID,
    employeeName: row.adminName ? `${row.adminName}（当前）` : '当前管理员',
    employeePhone: ''
  })
}

const loadTeamEmployees = async(reset) => {
  if (!adminTeamId.value) return
  if (empLoading.value) return
  if (!reset && empList.value.length >= empTotal.value && empTotal.value > 0) return
  empLoading.value = true
  try {
    const page = reset ? 1 : empPage.value
    const res = await getTeamEmployeeList({
      teamID: adminTeamId.value,
      page,
      pageSize: 10
    })
    if (res?.code !== 0) {
      ElMessage.error(res?.msg || '加载员工失败')
      return
    }
    const list = res.data?.list ?? []
    empTotal.value = res.data?.total ?? 0
    if (reset) {
      empList.value = list
      empPage.value = 2
      ensureCurrentAdminInList()
    } else {
      empList.value = [...empList.value, ...list]
      empPage.value += 1
    }
  } finally {
    empLoading.value = false
  }
}

const onEmpScroll = (e) => {
  const el = e.target
  if (!el || empLoading.value) return
  const nearBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 40
  if (nearBottom) {
    loadTeamEmployees(false)
  }
}

const openAdminDialog = (row) => {
  currentTeamRow.value = row
  adminTeamId.value = row.ID
  empList.value = []
  empPage.value = 1
  empTotal.value = 0
  selectedAdminId.value = row.adminID > 0 ? row.adminID : null
  adminDialogVisible.value = true
  loadTeamEmployees(true)
}

const handleConfirmAdmin = async() => {
  if (!adminTeamId.value) return
  if (selectedAdminId.value === null || selectedAdminId.value === undefined) {
    ElMessage.warning('请选择一名员工作为管理员')
    return
  }
  adminSubmitLoading.value = true
  try {
    const res = await setTeamAdmin({
      ID: adminTeamId.value,
      adminID: selectedAdminId.value
    })
    if (res?.code !== 0) {
      ElMessage.error(res?.msg || '设置失败')
      return
    }
    ElMessage.success('设置成功')
    adminDialogVisible.value = false
    getTeamList()
  } finally {
    adminSubmitLoading.value = false
  }
}

const openDrawer = () => {
  form.value = emptyForm()
  drawerVisible.value = true
}

const closeDrawer = () => {
  drawerVisible.value = false
}

const handleCreateTeam = async() => {
  if (!formRef.value) {
    ElMessage.error('表单校验失败')
    return
  }

  formRef.value.validate(async(valid) => {
    if (!valid) return

    const payload = { teamName: form.value.teamName }
    const res = await createTeam(payload)
    if (res?.code !== 0) {
      ElMessage.error(res?.msg || '创建团队失败')
      return
    }

    ElMessage.success('创建团队成功')
    closeDrawer()
    getTeamList()
  })
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTeamList()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTeamList()
}

onMounted(() => {
  getTeamList()

  // AI 浮窗里切换团队状态后，需要在本页重新拉取列表刷新状态
  emitter.on('teamListReload', getTeamList)
})

onBeforeUnmount(() => {
  emitter.off('teamListReload', getTeamList)
})
</script>

<style scoped>
.gva-table-box {
  width: 100%;
}

.admin-emp-scroll {
  max-height: 320px;
  overflow-y: auto;
  padding-right: 4px;
}

.admin-dialog-meta {
  margin-bottom: 12px;
  font-size: 13px;
  color: #606266;
}

.admin-dialog-meta-gap {
  margin-left: 12px;
}

.admin-radio-group {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  width: 100%;
}

.admin-radio-item {
  margin-right: 0;
  height: auto;
  align-items: flex-start;
  white-space: normal;
}

.admin-radio-phone {
  margin-left: 8px;
  color: #909399;
  font-size: 13px;
}
</style>

