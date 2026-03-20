/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-17 09:48:50
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-20 16:24:22
 * @FilePath: \gin-vue-admin-main\server\model\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysTeam struct {
	global.GVA_MODEL
	TeamNum  string `json:"teamNum" gorm:"comment:团队编号"`
	TeamName string `json:"teamName" gorm:"comment:团队名称"`
	Organize string `json:"organize" gorm:"comment:组织"`
	// 当前项目中 team 的 employee 关联尚未做成可迁移的关系/JSON列，
	// 这里先忽略给 GORM AutoMigrate，避免启动时报 invalid field。
	// 后续如果你要持久化并返回员工列表，可再补齐对应的存储字段和填充逻辑。
	EmployeeList []SysEmployee `json:"employeeList" gorm:"-"`
	AdminName    string        `json:"adminName" gorm:"comment:管理员名称"`
	Status       *int          `json:"status" gorm:"comment:团队状态"`
}

func (SysTeam) TableName() string {
	return "sys_team"
}
