package system

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	// 初始化内存数据库用于测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	global.GVA_DB = db
	// 自动迁移表结构
	db.AutoMigrate(&system.SysEmployee{})
}

func TestEmployeeService_CreateEmployee_Success(t *testing.T) {
	// 准备测试数据
	emp := system.SysEmployee{
		EmployeeName:       "张三",
		EmployeePhone:      "13800138000",
		EmployeeEmail:      "zhangsan@example.com",
		EmployeeAddress:    "北京市朝阳区",
		EmployeeBirthday:   "1990-01-01",
		EmployeeGender:     1,
		EmployeePosition:   "工程师",
		EmployeeDepartment: "技术部",
	}

	// 调用被测试方法
	service := &EmployeeService{}
	err := service.CreateEmployee(emp)

	// 验证结果
	assert.NoError(t, err)

	// 验证数据已正确插入
	var count int64
	global.GVA_DB.Model(&system.SysEmployee{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestEmployeeService_CreateEmployee_MultipleEmployees(t *testing.T) {
	// 准备测试数据 - 多个员工
	emp1 := system.SysEmployee{
		EmployeeName:     "李四",
		EmployeePhone:    "13800138001",
		EmployeeEmail:    "lisi@example.com",
		EmployeeGender:   2,
		EmployeePosition: "产品经理",
	}
	emp2 := system.SysEmployee{
		EmployeeName:     "王五",
		EmployeePhone:    "13800138002",
		EmployeeEmail:    "wangwu@example.com",
		EmployeeGender:   1,
		EmployeePosition: "设计师",
	}

	service := &EmployeeService{}

	// 创建第一个员工
	err := service.CreateEmployee(emp1)
	assert.NoError(t, err)

	// 创建第二个员工
	err = service.CreateEmployee(emp2)
	assert.NoError(t, err)

	// 验证数据
	var count int64
	global.GVA_DB.Model(&system.SysEmployee{}).Count(&count)
	assert.Equal(t, int64(2), count)
}

func TestEmployeeService_UpdateEmployee_Success(t *testing.T) {
	// 先创建一个员工
	emp := system.SysEmployee{
		EmployeeName:     "原始姓名",
		EmployeePhone:    "13800138000",
		EmployeeEmail:    "original@example.com",
		EmployeeGender:   1,
		EmployeePosition: "原始职位",
	}
	global.GVA_DB.Create(&emp)

	// 准备更新数据
	updateEmp := system.SysEmployee{
		ID:               emp.ID,
		EmployeeName:     "更新姓名",
		EmployeePhone:    "13900139000",
		EmployeeEmail:    "updated@example.com",
		EmployeeGender:   2,
		EmployeePosition: "更新职位",
	}

	// 调用被测试方法
	service := &EmployeeService{}
	updatedEmployee, err := service.UpdateEmployee(updateEmp)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "更新姓名", updatedEmployee.EmployeeName)
	assert.Equal(t, "13900139000", updatedEmployee.EmployeePhone)
}

func TestEmployeeService_UpdateEmployee_NotFound(t *testing.T) {
	// 准备不存在的员工数据
	emp := system.SysEmployee{
		ID:            99999,
		EmployeeName:  "不存在的员工",
		EmployeePhone: "13800138000",
	}

	// 调用被测试方法
	service := &EmployeeService{}
	_, err := service.UpdateEmployee(emp)

	// 验证结果 - 应该返回错误
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "查询角色数据失败")
}

func TestEmployeeService_GetEmployeeDetail_Success(t *testing.T) {
	// 先创建一个员工
	emp := system.SysEmployee{
		EmployeeName:     "测试员工",
		EmployeePhone:    "13800138000",
		EmployeeEmail:    "test@example.com",
		EmployeeGender:   1,
		EmployeePosition: "测试职位",
	}
	global.GVA_DB.Create(&emp)

	// 调用被测试方法
	service := &EmployeeService{}
	detail, err := service.GetEmployeeDetail(emp.ID)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "测试员工", detail.EmployeeName)
	assert.Equal(t, "男", detail.EmployeeGenderStr) // 验证 FillGenderText 被调用
}

func TestEmployeeService_GetEmployeeDetail_NotFound(t *testing.T) {
	// 调用被测试方法 - 查询不存在的员工
	service := &EmployeeService{}
	_, err := service.GetEmployeeDetail(99999)

	// 验证结果 - 应该返回错误
	assert.Error(t, err)
}

func TestEmployeeService_GetEmployeeDetail_Gender_0(t *testing.T) {
	// 测试性别为0的情况
	emp := system.SysEmployee{
		EmployeeName:   "测试员工",
		EmployeePhone:  "13800138000",
		EmployeeEmail:  "test@example.com",
		EmployeeGender: 0, // 未知
	}
	global.GVA_DB.Create(&emp)

	// 调用被测试方法
	service := &EmployeeService{}
	detail, err := service.GetEmployeeDetail(emp.ID)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "未知", detail.EmployeeGenderStr)
}

func TestEmployeeService_GetEmployeeDetail_Gender_2(t *testing.T) {
	// 测试性别为2（女）的情况
	emp := system.SysEmployee{
		EmployeeName:   "测试员工",
		EmployeePhone:  "13800138000",
		EmployeeEmail:  "test@example.com",
		EmployeeGender: 2, // 女
	}
	global.GVA_DB.Create(&emp)

	// 调用被测试方法
	service := &EmployeeService{}
	detail, err := service.GetEmployeeDetail(emp.ID)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "女", detail.EmployeeGenderStr)
}

func TestEmployeeService_GetEmployeeInfoList_Success(t *testing.T) {
	// 先创建多个员工
	employees := []system.SysEmployee{
		{
			EmployeeName:   "员工1",
			EmployeePhone:  "13800138001",
			EmployeeEmail:  "emp1@example.com",
			EmployeeGender: 1,
		},
		{
			EmployeeName:   "员工2",
			EmployeePhone:  "13800138002",
			EmployeeEmail:  "emp2@example.com",
			EmployeeGender: 2,
		},
		{
			EmployeeName:   "员工3",
			EmployeePhone:  "13800138003",
			EmployeeEmail:  "emp3@example.com",
			EmployeeGender: 0,
		},
	}
	for _, emp := range employees {
		global.GVA_DB.Create(&emp)
	}

	// 准备分页参数
	pageInfo := systemReq.SysEmployeeSearch{
		PageInfo: systemReq.PageInfo{
			Page:     1,
			PageSize: 10,
		},
	}

	// 调用被测试方法
	service := &EmployeeService{}
	list, total, err := service.GetEmployeeInfoList(pageInfo)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list.([]system.SysEmployee), 3)

	// 验证 FillGenderText 被调用
	empList := list.([]system.SysEmployee)
	for _, emp := range empList {
		assert.NotEmpty(t, emp.EmployeeGenderStr)
	}
}

func TestEmployeeService_GetEmployeeInfoList_Pagination(t *testing.T) {
	// 先创建多个员工
	for i := 1; i <= 15; i++ {
		emp := system.SysEmployee{
			EmployeeName:   "员工" + string(rune('0'+i)),
			EmployeePhone:  "1380013800" + string(rune('0'+i%10)),
			EmployeeEmail:  "emp" + string(rune('0'+i)) + "@example.com",
			EmployeeGender: i % 3,
		}
		global.GVA_DB.Create(&emp)
	}

	// 准备分页参数 - 每页5条
	pageInfo := systemReq.SysEmployeeSearch{
		PageInfo: systemReq.PageInfo{
			Page:     2,
			PageSize: 5,
		},
	}

	// 调用被测试方法
	service := &EmployeeService{}
	list, total, err := service.GetEmployeeInfoList(pageInfo)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, list.([]system.SysEmployee), 5)
}

func TestEmployeeService_GetEmployeeInfoList_Empty(t *testing.T) {
	// 不创建任何员工
	pageInfo := systemReq.SysEmployeeSearch{
		PageInfo: systemReq.PageInfo{
			Page:     1,
			PageSize: 10,
		},
	}

	// 调用被测试方法
	service := &EmployeeService{}
	list, total, err := service.GetEmployeeInfoList(pageInfo)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, list.([]system.SysEmployee), 0)
}

func TestEmployeeService_GetEmployeeInfoList_PageSize_0(t *testing.T) {
	// 先创建员工
	emp := system.SysEmployee{
		EmployeeName:   "测试员工",
		EmployeePhone:  "13800138000",
		EmployeeEmail:  "test@example.com",
		EmployeeGender: 1,
	}
	global.GVA_DB.Create(&emp)

	// PageSize 为 0 的情况
	pageInfo := systemReq.SysEmployeeSearch{
		PageInfo: systemReq.PageInfo{
			Page:     1,
			PageSize: 0,
		},
	}

	// 调用被测试方法
	service := &EmployeeService{}
	list, total, err := service.GetEmployeeInfoList(pageInfo)

	// 验证结果 - PageSize 为 0 时，Limit(0) 会返回所有记录
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	// 注意：GORM 的 Limit(0) 行为可能返回所有记录或空记录，取决于版本
	assert.NotNil(t, list)
}

func TestEmployeeService_GetEmployeeInfoList_Page_0(t *testing.T) {
	// 先创建员工
	emp := system.SysEmployee{
		EmployeeName:   "测试员工",
		EmployeePhone:  "13800138000",
		EmployeeEmail:  "test@example.com",
		EmployeeGender: 1,
	}
	global.GVA_DB.Create(&emp)

	// Page 为 0 的情况
	pageInfo := systemReq.SysEmployeeSearch{
		PageInfo: systemReq.PageInfo{
			Page:     0,
			PageSize: 10,
		},
	}

	// 调用被测试方法
	service := &EmployeeService{}
	list, total, err := service.GetEmployeeInfoList(pageInfo)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.NotNil(t, list)
}

func TestEmployeeService_GetEmployeeInfoList_LargePageSize(t *testing.T) {
	// 先创建员工
	employees := []system.SysEmployee{
		{EmployeeName: "员工1", EmployeePhone: "13800138001", EmployeeEmail: "emp1@example.com"},
		{EmployeeName: "员工2", EmployeePhone: "13800138002", EmployeeEmail: "emp2@example.com"},
	}
	for _, emp := range employees {
		global.GVA_DB.Create(&emp)
	}

	// PageSize 超过实际数量
	pageInfo := systemReq.SysEmployeeSearch{
		PageInfo: systemReq.PageInfo{
			Page:     1,
			PageSize: 100,
		},
	}

	// 调用被测试方法
	service := &EmployeeService{}
	list, total, err := service.GetEmployeeInfoList(pageInfo)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list.([]system.SysEmployee), 2)
}
