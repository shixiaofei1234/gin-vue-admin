/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-17 14:24:24
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-20 17:28:28
 * @FilePath: \gin-vue-admin-main\server\service\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

type TeamService struct{}

// CreateEmployee 创建员工
func (teamService *TeamService) CreateTeam(emp system.SysTeam) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&emp).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
