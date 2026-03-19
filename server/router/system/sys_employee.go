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
	}
}
