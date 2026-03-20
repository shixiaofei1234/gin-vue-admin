/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-17 10:06:23
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-20 17:39:27
 * @FilePath: \gin-vue-admin-main\server\api\v1\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TeamApi struct{}

// CreateEmployee
// @Tags      Employee
// @Summary   创建角色
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysEmployee                                                true  "权限id, 权限名, 父角色id"
// @Success   200   {object}   "创建员工"
// @Router    /employee/createEmployee [post]
func (a *TeamApi) CreateTeam(c *gin.Context) {
	var team system.SysTeam
	var err error

	if err = c.ShouldBindJSON(&team); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = utils.Verify(team, utils.EmployeeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = teamService.CreateTeam(team); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetEmployeeList
// @Tags      Employee
// @Summary   分页获取员工列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SysEmployeeSearch                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取员工列表,返回包括列表,总数,页码,每页数量"
// @Router    /employee/GetEmployeeList [get]

func (a *TeamApi) GetTeamList(c *gin.Context) {
	var pageInfo systemReq.SysEmployeeSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := employeeService.GetEmployeeInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
