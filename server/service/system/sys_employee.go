/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-17 14:24:24
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-19 14:23:17
 * @FilePath: \gin-vue-admin-main\server\service\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"gorm.io/gorm"
)

type EmployeeService struct{}

// CreateEmployee 创建员工
func (employeeService *EmployeeService) CreateEmployee(emp system.SysEmployee) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&emp).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateEmployee 更新员工
func (employeeService *EmployeeService) UpdateEmployee(emp system.SysEmployee) (employee system.SysEmployee, err error) {
	var employeeBox system.SysEmployee
	err = global.GVA_DB.Where("id = ?", emp.ID).First(&employeeBox).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return system.SysEmployee{}, errors.New("查询角色数据失败")
	}
	err = global.GVA_DB.Model(&employeeBox).Updates(&emp).Error
	return employeeBox, err
}

// GetEmployeeDetail 获取员工详情
func (employeeService *EmployeeService) GetEmployeeDetail(id uint) (employee system.SysEmployee, err error) {
	err = global.GVA_DB.First(&employee, "id = ?", id).Error
	return employee, err
}

func (employeeService *EmployeeService) GetEmployeeInfoList(info systemReq.SysEmployeeSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysEmployee{})
	var sysEmployees []system.SysEmployee
	if err = db.Count(&total).Error; err != nil {
		return
	}
	if err = db.Order("id desc").Limit(limit).Offset(offset).Find(&sysEmployees).Error; err != nil {
		return
	}
	return sysEmployees, total, nil
}
