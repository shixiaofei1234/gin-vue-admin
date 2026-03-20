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
          label="团队状态"
          prop="status"
          min-width="100"
        >
          <template #default="scope">
            {{ scope.row.status ?? '-' }}
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import service from '@/utils/request'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'TeamList'
})

// 你后端自行添加路由/接口时，可以按需修改这两个 URL。
const TEAM_LIST_URL = '/team/getTeamList'
const TEAM_CREATE_URL = '/team/createTeam'

const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const drawerVisible = ref(false)

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
})
</script>

<style scoped>
.gva-table-box {
  width: 100%;
}
</style>

