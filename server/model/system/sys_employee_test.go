package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSysEmployee_TableName(t *testing.T) {
	emp := SysEmployee{}
	assert.Equal(t, "sys_employee", emp.TableName())
}

func TestGetEmployeeGenderText_Unknown(t *testing.T) {
	// 测试边界值：性别为0（未知）
	result := GetEmployeeGenderText(0)
	assert.Equal(t, "未知", result)
}

func TestGetEmployeeGenderText_Male(t *testing.T) {
	// 测试正常值：性别为1（男）
	result := GetEmployeeGenderText(1)
	assert.Equal(t, "男", result)
}

func TestGetEmployeeGenderText_Female(t *testing.T) {
	// 测试正常值：性别为2（女）
	result := GetEmployeeGenderText(2)
	assert.Equal(t, "女", result)
}

func TestGetEmployeeGenderText_InvalidValue_Negative(t *testing.T) {
	// 测试异常值：性别为负数
	result := GetEmployeeGenderText(-1)
	assert.Equal(t, "未知", result)
}

func TestGetEmployeeGenderText_InvalidValue_Large(t *testing.T) {
	// 测试异常值：性别大于2
	result := GetEmployeeGenderText(100)
	assert.Equal(t, "未知", result)
}

func TestSysEmployee_FillGenderText_Unknown(t *testing.T) {
	// 测试 FillGenderText：性别为0（未知）
	emp := SysEmployee{
		EmployeeGender: 0,
	}
	emp.FillGenderText()
	assert.Equal(t, "未知", emp.EmployeeGenderStr)
}

func TestSysEmployee_FillGenderText_Male(t *testing.T) {
	// 测试 FillGenderText：性别为1（男）
	emp := SysEmployee{
		EmployeeGender: 1,
	}
	emp.FillGenderText()
	assert.Equal(t, "男", emp.EmployeeGenderStr)
}

func TestSysEmployee_FillGenderText_Female(t *testing.T) {
	// 测试 FillGenderText：性别为2（女）
	emp := SysEmployee{
		EmployeeGender: 2,
	}
	emp.FillGenderText()
	assert.Equal(t, "女", emp.EmployeeGenderStr)
}

func TestSysEmployee_FillGenderText_InvalidValue(t *testing.T) {
	// 测试 FillGenderText：无效的性别值
	emp := SysEmployee{
		EmployeeGender: 99,
	}
	emp.FillGenderText()
	assert.Equal(t, "未知", emp.EmployeeGenderStr)
}

func TestSysEmployee_AllFields(t *testing.T) {
	// 测试 SysEmployee 结构体的所有字段
	emp := SysEmployee{
		EmployeeName:       "张三",
		EmployeePhone:      "13800138000",
		EmployeeEmail:      "zhangsan@example.com",
		EmployeeAddress:    "北京市朝阳区",
		EmployeeBirthday:   "1990-01-01",
		EmployeeGender:     1,
		EmployeePosition:   "工程师",
		EmployeeDepartment: "技术部",
	}
	assert.Equal(t, "张三", emp.EmployeeName)
	assert.Equal(t, "13800138000", emp.EmployeePhone)
	assert.Equal(t, "zhangsan@example.com", emp.EmployeeEmail)
	assert.Equal(t, "北京市朝阳区", emp.EmployeeAddress)
	assert.Equal(t, "1990-01-01", emp.EmployeeBirthday)
	assert.Equal(t, 1, emp.EmployeeGender)
	assert.Equal(t, "工程师", emp.EmployeePosition)
	assert.Equal(t, "技术部", emp.EmployeeDepartment)
}

func TestSysEmployee_GenderConstants(t *testing.T) {
	// 测试性别常量
	assert.Equal(t, 0, GenderUnknown)
	assert.Equal(t, 1, GenderMale)
	assert.Equal(t, 2, GenderFemale)
}

func TestSysEmployee_FillGenderText_MultipleCalls(t *testing.T) {
	// 测试多次调用 FillGenderText
	emp := SysEmployee{
		EmployeeGender: 1,
	}
	emp.FillGenderText()
	assert.Equal(t, "男", emp.EmployeeGenderStr)

	emp.EmployeeGender = 2
	emp.FillGenderText()
	assert.Equal(t, "女", emp.EmployeeGenderStr)

	emp.EmployeeGender = 0
	emp.FillGenderText()
	assert.Equal(t, "未知", emp.EmployeeGenderStr)
}
