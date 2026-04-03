/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-17 14:24:24
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-23 15:01:19
 * @FilePath: \gin-vue-admin-main\server\service\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"gorm.io/gorm"
)

type TeamService struct{}

func nextTeamNum(tx *gorm.DB) (string, error) {
	var last system.SysTeam
	err := tx.Model(&system.SysTeam{}).Where("team_num LIKE ?", "TD%").Order("team_num DESC").Take(&last).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "TD0001", nil
	}
	if err != nil {
		return "", err
	}
	s := strings.TrimPrefix(strings.ToUpper(strings.TrimSpace(last.TeamNum)), "TD")
	n, aerr := strconv.Atoi(s)
	if aerr != nil {
		return "", aerr
	}
	return fmt.Sprintf("TD%04d", n+1), nil
}

// CreateTeam 创建团队
func (teamService *TeamService) CreateTeam(team system.SysTeam) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		num, err := nextTeamNum(tx)
		if err != nil {
			return err
		}
		team.TeamNum = num
		st := false
		team.Status = st
		return tx.Create(&team).Error
	})
	return err
}

// SwitchTeamStatus 切换团队状态
func (teamService *TeamService) SwitchTeamStatus(id uint, status bool) (err error) {
	err = global.GVA_DB.Model(&system.SysTeam{}).Where("id = ?", id).Update("status", status).Error
	return err
}

// SetAdmin 设置管理员（校验员工属于该团队，并同步管理员姓名）
func (teamService *TeamService) SetAdmin(teamId uint, adminId uint) (err error) {
	var emp system.SysEmployee
	if err = global.GVA_DB.First(&emp, "id = ?", adminId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("员工不存在")
		}
		return err
	}
	if emp.TeamID != teamId {
		return errors.New("该员工不属于当前团队")
	}
	return global.GVA_DB.Model(&system.SysTeam{}).Where("id = ?", teamId).Updates(map[string]interface{}{
		"admin_id":   adminId,
		"admin_name": emp.EmployeeName,
	}).Error
}

// GetTeamEmployeeList 分页获取某团队下的员工
func (teamService *TeamService) GetTeamEmployeeList(info systemReq.SysTeamEmployeeSearch) (list []system.SysEmployee, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysEmployee{}).Where("team_id = ?", info.TeamID)
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		list[i].FillGenderText()
	}
	return list, total, nil
}

// CreateTeam 获取团队列表
func (teamService *TeamService) GetTeamList(pageInfo systemReq.SysTeamSearch) (teamList []system.SysTeam, total int64, err error) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	db := global.GVA_DB.Model(&system.SysTeam{})
	if pageInfo.TeamName != "" {
		db = db.Where("team_name LIKE ?", "%"+pageInfo.TeamName+"%")
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&teamList).Error
	for i := range teamList {
		if err = global.GVA_DB.Model(&system.SysEmployee{}).Where("team_id = ?", teamList[i].ID).Find(&teamList[i].EmployeeList).Error; err != nil {
			return nil, 0, err
		}
	}
	return teamList, total, err
}
