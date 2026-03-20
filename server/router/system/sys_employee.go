/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-19 14:47:55
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-19 17:43:50
 * @FilePath: \server\router\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type EmployeeRouter struct{}

func (s *EmployeeRouter) InitEmployeeRouter(Router *gin.RouterGroup) {
	employeeRouter := Router.Group("employee").Use(middleware.OperationRecord())
	{
		employeeRouter.GET("getEmployeeList", employeeApi.GetEmployeeList)     // 获取员工列表
		employeeRouter.POST("createEmployee", employeeApi.CreateEmployee)      // 创建员工
		employeeRouter.PUT("updateEmployee", employeeApi.UpdateEmployee)       // 更新员工
		employeeRouter.GET("getEmployeeDetail", employeeApi.GetEmployeeDetail) // 获取员工详情
		employeeRouter.DELETE("deleteEmployee", employeeApi.DeleteEmployee)    // 删除员工
	}
}
