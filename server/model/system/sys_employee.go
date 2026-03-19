/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-17 09:48:50
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-18 10:14:09
 * @FilePath: \gin-vue-admin-main\server\model\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

const (
	GenderUnknown = 0 // 未知
	GenderMale    = 1 // 男
	GenderFemale  = 2 // 女
)

type SysEmployee struct {
	global.GVA_MODEL
	EmployeeName       string `json:"employeeName" gorm:"comment:员工姓名"`
	EmployeePhone      string `json:"employeePhone" gorm:"comment:员工电话"`
	EmployeeEmail      string `json:"employeeEmail" gorm:"comment:员工邮箱"`
	EmployeeAddress    string `json:"employeeAddress" gorm:"comment:员工地址"`
	EmployeeBirthday   string `json:"employeeBirthday" gorm:"comment:员工生日"`
	EmployeeGender     int    `json:"employeeGender" gorm:"type:tinyint;comment:员工性别（0=未知，1=男，2=女）"`
	EmployeeGenderStr  string `json:"employeeGenderStr" gorm:"-"`
	EmployeePosition   string `json:"employeePosition" gorm:"comment:员工职位"`
	EmployeeDepartment string `json:"employeeDepartment" gorm:"comment:员工部门"`
}

func (SysEmployee) TableName() string {
	return "sys_employee"
}

func (e *SysEmployee) FillGenderText() {
	e.EmployeeGenderStr = GetEmployeeGenderText(e.EmployeeGender)
}

func GetEmployeeGenderText(gender int) string {
	switch gender {
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未知"
	}
}
