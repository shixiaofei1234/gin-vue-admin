/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-19 14:52:21
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-19 17:44:04
 * @FilePath: \serverd:\project\gin-vue-admin\web\src\api\employee.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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

// 删除员工（DELETE，data: { ID: number }）
export const deleteEmployee = (data) => {
  return service({
    url: '/employee/deleteEmployee',
    method: 'delete',
    data
  })
}
