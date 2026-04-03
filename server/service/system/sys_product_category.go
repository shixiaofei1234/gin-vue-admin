package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

type ProductCategoryService struct{}

func (productCategoryService *ProductCategoryService) CreateProductCategory(productCategory system.SysProductCategory) (err error) {
	return global.GVA_DB.Create(&productCategory).Error
}

func (productCategoryService *ProductCategoryService) GetProductCategoryList(pageInfo systemReq.SysProductCategorySearch) (list []system.SysProductCategory, total int64, err error) {
	return []system.SysProductCategory{}, 0, nil
}
