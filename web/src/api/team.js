import service from '@/utils/request'

/** 团队下员工分页（GET：teamID, page, pageSize） */
export const getTeamEmployeeList = (params) => {
  return service({
    url: '/team/getTeamEmployeeList',
    method: 'get',
    params
  })
}

/** 设置团队管理员（POST：ID 团队主键, adminID 员工主键） */
export const setTeamAdmin = (data) => {
  return service({
    url: '/team/setAdmin',
    method: 'post',
    data
  })
}

// 切换团队状态（POST：ID 团队主键, status bool）
export const switchTeamStatus = (data) => {
  return service({
    url: '/team/switchTeamStatus',
    method: 'post',
    data
  })
}
