import service from '@/utils/request'

/** 分页列表（GET） */
export const getProductCategoryList = (params) => {
  return service({
    url: '/productCategory/getProductCategoryList',
    method: 'get',
    params
  })
}

/** 下拉（启用分类，GET） */
export const getProductCategorySelectList = () => {
  return service({
    url: '/productCategory/getProductCategorySelectList',
    method: 'get'
  })
}

/** 详情（GET params: { ID }） */
export const getProductCategory = (params) => {
  return service({
    url: '/productCategory/getProductCategory',
    method: 'get',
    params
  })
}

export const createProductCategory = (data) => {
  return service({
    url: '/productCategory/createProductCategory',
    method: 'post',
    data
  })
}

export const updateProductCategory = (data) => {
  return service({
    url: '/productCategory/updateProductCategory',
    method: 'put',
    data
  })
}

export const deleteProductCategory = (data) => {
  return service({
    url: '/productCategory/deleteProductCategory',
    method: 'delete',
    data
  })
}
