<template>
  <el-drawer
    v-model="visibleLocal"
    :before-close="handleClose"
    :show-close="false"
    size="40%"
  >
    <template #header>
      <div class="flex justify-between items-center">
        <span class="text-lg">{{ title }}</span>
        <div>
          <el-button @click="handleClose">取 消</el-button>
          <el-button
            type="primary"
            @click="handleSubmit"
          >
            {{ submitText }}
          </el-button>
        </div>
      </div>
    </template>
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item
        label="员工姓名"
        prop="employeeName"
      >
        <el-input
          v-model="form.employeeName"
          maxlength="10"
          show-word-limit
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="员工电话"
        prop="employeePhone"
      >
        <el-input
          v-model="form.employeePhone"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="员工邮箱"
        prop="employeeEmail"
      >
        <el-input
          v-model="form.employeeEmail"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="员工地址"
        prop="employeeAddress"
      >
        <el-input
          v-model="form.employeeAddress"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="员工生日"
        prop="employeeBirthday"
      >
        <el-date-picker
          v-model="form.employeeBirthday"
          type="date"
          value-format="YYYY-MM-DD"
          placeholder="请选择日期"
        />
      </el-form-item>
      <el-form-item
        label="员工性别"
        prop="employeeGender"
      >
        <el-radio-group v-model="form.employeeGender">
          <el-radio :label="1">男</el-radio>
          <el-radio :label="2">女</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item
        label="员工职位"
        prop="employeePosition"
      >
        <el-input
          v-model="form.employeePosition"
          autocomplete="off"
        />
      </el-form-item>
      <el-form-item
        label="员工部门"
        prop="employeeDepartment"
      >
        <el-input
          v-model="form.employeeDepartment"
          autocomplete="off"
        />
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '创建员工'
  },
  submitText: {
    type: String,
    default: '确 定'
  },
  initialValue: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue', 'submit'])

const visibleLocal = ref(props.modelValue)

watch(
  () => props.modelValue,
  (val) => {
    visibleLocal.value = val
    if (val) {
      resetForm()
      applyInitialValue()
    }
  }
)

watch(
  () => visibleLocal.value,
  (val) => {
    emit('update:modelValue', val)
  }
)

const formRef = ref(null)
const form = ref({
  employeeName: '',
  employeePhone: '',
  employeeEmail: '',
  employeeAddress: '',
  employeeBirthday: '',
  employeeGender: 1,
  employeePosition: '',
  employeeDepartment: ''
})

const rules = {
  employeeName: [
    { required: true, message: '请输入员工姓名', trigger: 'blur' },
    { min: 1, max: 10, message: '长度在 1 到 10 个字符', trigger: 'blur' }
  ],
  employeePhone: [
    { required: true, message: '请输入员工电话', trigger: 'blur' },
    {
      pattern: /^1[3-9]\d{9}$/,
      message: '请输入正确的手机号',
      trigger: 'blur'
    }
  ],
  employeeEmail: [
    { required: true, message: '请输入员工邮箱', trigger: 'blur' },
    {
      type: 'email',
      message: '请输入正确的邮箱地址',
      trigger: ['blur', 'change']
    }
  ],
  employeeGender: [
    { required: true, message: '请选择性别', trigger: 'change' }
  ]
}

const resetForm = () => {
  form.value = {
    employeeName: '',
    employeePhone: '',
    employeeEmail: '',
    employeeAddress: '',
    employeeBirthday: '',
    employeeGender: 1,
    employeePosition: '',
    employeeDepartment: ''
  }
}

const applyInitialValue = () => {
  if (!props.initialValue) return
  const v = props.initialValue
  form.value = {
    employeeName: v.employeeName ?? '',
    employeePhone: v.employeePhone ?? '',
    employeeEmail: v.employeeEmail ?? '',
    employeeAddress: v.employeeAddress ?? '',
    employeeBirthday: v.employeeBirthday ?? '',
    employeeGender: v.employeeGender ?? 1,
    employeePosition: v.employeePosition ?? '',
    employeeDepartment: v.employeeDepartment ?? ''
  }
}

const handleClose = () => {
  visibleLocal.value = false
}

const handleSubmit = () => {
  if (!formRef.value) return
  formRef.value.validate((valid) => {
    if (!valid) return
    emit('submit', { ...form.value })
    visibleLocal.value = false
  })
}
</script>

