/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-19 14:47:55
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-04-01 10:22:56
 * @FilePath: \server\router\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProductCategoryRouter struct{}

func (s *ProductCategoryRouter) InitProductCategoryRouter(Router *gin.RouterGroup) {
	productCategoryRouter := Router.Group("productCategory").Use(middleware.OperationRecord())
	{
		// productCategoryRouter.GET("getTeamList", productCategoryApi.GetProductCategoryList)           // 获取分类列表
		productCategoryRouter.POST("createProductCategory", productCategoryApi.CreateProductCategory) // 创建产品分类
	}
}
