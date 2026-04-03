/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-04-01 10:18:04
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-04-01 11:27:13
 * @FilePath: \server\api\v1\system\sys_product_category.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
)

type ProductCategoryApi struct{}

func (productCategoryService *ProductCategoryApi) CreateProductCategory(c *gin.Context) {
	response.OkWithMessage("创建成功", c)
}

func (productCategoryService *ProductCategoryApi) GetProductCategoryList(pageInfo systemReq.SysProductCategorySearch) (list []system.SysProductCategory, total int64, err error) {
	return productCategoryService.GetProductCategoryList(pageInfo)
}
