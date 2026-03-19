import service from '@/utils/request'

// 员工列表（GET，分页参数 page、pageSize）
export const getEmployeeList = (params) => {
  return service({
    url: '/employee/getEmployeeList',
    method: 'get',
    params
  })
}

// 创建员工（POST）
export const createEmployee = (data) => {
  return service({
    url: '/employee/createEmployee',
    method: 'post',
    data
  })
}

// 更新员工（PUT）
export const updateEmployee = (data) => {
  return service({
    url: '/employee/updateEmployee',
    method: 'put',
    data
  })
}

// 员工详情（GET，params: { ID: number }）
export const getEmployeeDetail = (params) => {
  return service({
    url: '/employee/getEmployeeDetail',
    method: 'get',
    params
  })
}

