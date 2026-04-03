/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-19 14:47:55
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-30 13:52:18
 * @FilePath: \server\router\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TeamRouter struct{}

func (s *TeamRouter) InitTeamRouter(Router *gin.RouterGroup) {
	teamRouter := Router.Group("team").Use(middleware.OperationRecord())
	{
		teamRouter.GET("getTeamList", teamApi.GetTeamList)            // 获取团队列表
		teamRouter.POST("createTeam", teamApi.CreateTeam)             // 创建团队
		teamRouter.POST("switchTeamStatus", teamApi.SwitchTeamStatus) // 切换团队状态
		teamRouter.POST("setAdmin", teamApi.SetAdmin)                  // 设置管理员
		teamRouter.GET("getTeamEmployeeList", teamApi.GetTeamEmployeeList) // 团队下员工分页
	}
}
