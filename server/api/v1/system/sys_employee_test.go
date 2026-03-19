package system

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
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

// ========== CreateEmployee 测试 ==========

func TestEmployeeApi_CreateEmployee_Success(t *testing.T) {
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
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.CreateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "创建成功", resp["msg"])
}

func TestEmployeeApi_CreateEmployee_BindJSONError(t *testing.T) {
	// 无效的 JSON
	invalidJSON := `{"employeeName": "test", employeePhone:}`

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.CreateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestEmployeeApi_CreateEmployee_VerifyError_EmptyName(t *testing.T) {
	// 测试数据 - 缺少员工姓名
	emp := map[string]interface{}{
		"employeeName":     "", // 空姓名
		"employeePhone":    "13800138000",
		"employeeEmail":    "test@example.com",
		"employeeGender":   1,
		"employeePosition": "工程师",
	}
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.CreateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestEmployeeApi_CreateEmployee_VerifyError_EmptyPhone(t *testing.T) {
	// 测试数据 - 缺少员工电话
	emp := map[string]interface{}{
		"employeeName":     "张三",
		"employeePhone":    "", // 空电话
		"employeeEmail":    "test@example.com",
		"employeeGender":   1,
		"employeePosition": "工程师",
	}
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.CreateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestEmployeeApi_CreateEmployee_VerifyError_EmptyEmail(t *testing.T) {
	// 测试数据 - 缺少员工邮箱
	emp := map[string]interface{}{
		"employeeName":     "张三",
		"employeePhone":    "13800138000",
		"employeeEmail":    "", // 空邮箱
		"employeeGender":   1,
		"employeePosition": "工程师",
	}
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.CreateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

// ========== UpdateEmployee 测试 ==========

func TestEmployeeApi_UpdateEmployee_Success(t *testing.T) {
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
	updateData := map[string]interface{}{
		"id":               emp.ID,
		"employeeName":     "更新姓名",
		"employeePhone":    "13900139000",
		"employeeEmail":    "updated@example.com",
		"employeeGender":   2,
		"employeePosition": "更新职位",
	}
	jsonData, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/employee/updateEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.UpdateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "更新成功", resp["msg"])

	// 验证数据确实已更新
	var updatedEmp system.SysEmployee
	global.GVA_DB.First(&updatedEmp, emp.ID)
	assert.Equal(t, "更新姓名", updatedEmp.EmployeeName)
}

func TestEmployeeApi_UpdateEmployee_NotFound(t *testing.T) {
	// 准备更新数据 - 不存在的员工
	updateData := map[string]interface{}{
		"id":               99999,
		"employeeName":     "测试姓名",
		"employeePhone":    "13800138000",
		"employeeEmail":    "test@example.com",
		"employeeGender":   1,
		"employeePosition": "工程师",
	}
	jsonData, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/employee/updateEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.UpdateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
	assert.Contains(t, resp["msg"].(string), "更新失败")
}

func TestEmployeeApi_UpdateEmployee_BindJSONError(t *testing.T) {
	// 无效的 JSON
	invalidJSON := `{"id": 1, employeeName:}`

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/employee/updateEmployee", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.UpdateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestEmployeeApi_UpdateEmployee_VerifyError_EmptyName(t *testing.T) {
	// 先创建一个员工
	emp := system.SysEmployee{
		EmployeeName:   "原始姓名",
		EmployeePhone:  "13800138000",
		EmployeeEmail:  "original@example.com",
		EmployeeGender: 1,
	}
	global.GVA_DB.Create(&emp)

	// 准备更新数据 - 空姓名
	updateData := map[string]interface{}{
		"id":               emp.ID,
		"employeeName":     "", // 空姓名
		"employeePhone":    "13800138000",
		"employeeEmail":    "test@example.com",
		"employeeGender":   1,
		"employeePosition": "工程师",
	}
	jsonData, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/employee/updateEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := EmployeeApi{}
	api.UpdateEmployee(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

// ========== GetEmployeeDetail 测试 ==========

func TestEmployeeApi_GetEmployeeDetail_Success(t *testing.T) {
	// 先创建一个员工
	emp := system.SysEmployee{
		EmployeeName:     "测试员工",
		EmployeePhone:    "13800138000",
		EmployeeEmail:    "test@example.com",
		EmployeeGender:   1,
		EmployeePosition: "工程师",
	}
	global.GVA_DB.Create(&emp)

	// 创建 gin.Context - 使用 Query 参数
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeDetail?id="+string(rune('0'+emp.ID)), nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeDetail(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证返回的数据
	data := resp["data"].(map[string]interface{})
	employee := data["employee"].(map[string]interface{})
	assert.Equal(t, "测试员工", employee["employeeName"])
}

func TestEmployeeApi_GetEmployeeDetail_NotFound(t *testing.T) {
	// 创建 gin.Context - 查询不存在的员工
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeDetail?id=99999", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeDetail(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
	assert.Contains(t, resp["msg"].(string), "获取详情失败")
}

func TestEmployeeApi_GetEmployeeDetail_BindQueryError(t *testing.T) {
	// 创建 gin.Context - 无效的 Query 参数
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// 使用无效的 id 参数
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeDetail?id=invalid", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeDetail(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestEmployeeApi_GetEmployeeDetail_MissingID(t *testing.T) {
	// 创建 gin.Context - 缺少必需的 ID 参数
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeDetail", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeDetail(c)

	// 验证响应 - ID 验证失败
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestEmployeeApi_GetEmployeeDetail_ID_Zero(t *testing.T) {
	// 创建 gin.Context - ID 为 0
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeDetail?id=0", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeDetail(c)

	// 验证响应 - ID 验证失败
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

// ========== GetEmployeeList 测试 ==========

func TestEmployeeApi_GetEmployeeList_Success(t *testing.T) {
	// 先创建一些员工
	employees := []system.SysEmployee{
		{EmployeeName: "员工1", EmployeePhone: "13800138001", EmployeeEmail: "emp1@example.com", EmployeeGender: 1},
		{EmployeeName: "员工2", EmployeePhone: "13800138002", EmployeeEmail: "emp2@example.com", EmployeeGender: 2},
		{EmployeeName: "员工3", EmployeePhone: "13800138003", EmployeeEmail: "emp3@example.com", EmployeeGender: 0},
	}
	for _, emp := range employees {
		global.GVA_DB.Create(&emp)
	}

	// 创建 gin.Context - 正常的分页参数
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeList?page=1&pageSize=10", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeList(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证分页数据
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["total"])
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(10), data["pageSize"])
}

func TestEmployeeApi_GetEmployeeList_Pagination(t *testing.T) {
	// 先创建多个员工
	for i := 1; i <= 15; i++ {
		emp := system.SysEmployee{
			EmployeeName:   "员工" + string(rune('0'+i)),
			EmployeePhone:  "1380013800" + string(rune('0'+(i%10))),
			EmployeeEmail:  "emp" + string(rune('0'+i)) + "@example.com",
			EmployeeGender: i % 3,
		}
		global.GVA_DB.Create(&emp)
	}

	// 创建 gin.Context - 第二页，每页5条
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeList?page=2&pageSize=5", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeList(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证分页数据
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(15), data["total"])
	assert.Equal(t, float64(2), data["page"])
	assert.Equal(t, float64(5), data["pageSize"])
}

func TestEmployeeApi_GetEmployeeList_Empty(t *testing.T) {
	// 不创建任何员工
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeList?page=1&pageSize=10", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeList(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证分页数据
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["total"])
}

func TestEmployeeApi_GetEmployeeList_BindQueryError(t *testing.T) {
	// 创建 gin.Context - 无效的分页参数
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// pageSize 使用非数字
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeList?page=1&pageSize=invalid", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeList(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestEmployeeApi_GetEmployeeList_DefaultPage(t *testing.T) {
	// 创建 gin.Context - 缺少分页参数
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeList", nil)

	// 调用被测试方法
	api := EmployeeApi{}
	api.GetEmployeeList(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

// ========== 边界值测试 ==========

func TestEmployeeApi_CreateEmployee_AllFields(t *testing.T) {
	// 测试所有字段
	emp := system.SysEmployee{
		EmployeeName:       "完整员工",
		EmployeePhone:      "13800138000",
		EmployeeEmail:      "full@example.com",
		EmployeeAddress:    "北京市朝阳区某街道",
		EmployeeBirthday:   "1990-01-01",
		EmployeeGender:     2,
		EmployeePosition:   "高级工程师",
		EmployeeDepartment: "研发部",
	}
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	api := EmployeeApi{}
	api.CreateEmployee(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestEmployeeApi_CreateEmployee_GenderUnknown(t *testing.T) {
	// 测试性别为0（未知）
	emp := map[string]interface{}{
		"employeeName":     "测试员工",
		"employeePhone":    "13800138000",
		"employeeEmail":    "test@example.com",
		"employeeGender":   0, // 未知
		"employeePosition": "工程师",
	}
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	api := EmployeeApi{}
	api.CreateEmployee(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestEmployeeApi_CreateEmployee_GenderMale(t *testing.T) {
	// 测试性别为1（男）
	emp := map[string]interface{}{
		"employeeName":     "测试员工",
		"employeePhone":    "13800138000",
		"employeeEmail":    "test@example.com",
		"employeeGender":   1, // 男
		"employeePosition": "工程师",
	}
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	api := EmployeeApi{}
	api.CreateEmployee(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestEmployeeApi_CreateEmployee_GenderFemale(t *testing.T) {
	// 测试性别为2（女）
	emp := map[string]interface{}{
		"employeeName":     "测试员工",
		"employeePhone":    "13800138000",
		"employeeEmail":    "test@example.com",
		"employeeGender":   2, // 女
		"employeePosition": "工程师",
	}
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	api := EmployeeApi{}
	api.CreateEmployee(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestEmployeeApi_CreateEmployee_InvalidContentType(t *testing.T) {
	// 测试无效的 Content-Type
	emp := system.SysEmployee{
		EmployeeName:     "测试员工",
		EmployeePhone:    "13800138000",
		EmployeeEmail:    "test@example.com",
		EmployeeGender:   1,
		EmployeePosition: "工程师",
	}
	jsonData, err := json.Marshal(emp)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/employee/createEmployee", bytes.NewBuffer(jsonData))
	// 不设置 Content-Type

	api := EmployeeApi{}
	api.CreateEmployee(c)

	// 应该失败，因为无法解析 JSON
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestEmployeeApi_UpdateEmployee_AllFields(t *testing.T) {
	// 先创建一个员工
	emp := system.SysEmployee{
		EmployeeName:       "原始",
		EmployeePhone:      "13800138000",
		EmployeeEmail:      "original@example.com",
		EmployeeAddress:    "原始地址",
		EmployeeBirthday:   "1990-01-01",
		EmployeeGender:     1,
		EmployeePosition:   "原始职位",
		EmployeeDepartment: "原始部门",
	}
	global.GVA_DB.Create(&emp)

	// 准备更新数据 - 所有字段
	updateData := map[string]interface{}{
		"id":                 emp.ID,
		"employeeName":       "更新后",
		"employeePhone":      "13900139000",
		"employeeEmail":      "updated@example.com",
		"employeeAddress":    "新地址",
		"employeeBirthday":   "1991-01-01",
		"employeeGender":     2,
		"employeePosition":   "新职位",
		"employeeDepartment": "新部门",
	}
	jsonData, err := json.Marshal(updateData)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/employee/updateEmployee", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	api := EmployeeApi{}
	api.UpdateEmployee(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestEmployeeApi_GetEmployeeList_MaxPageSize(t *testing.T) {
	// 先创建一些员工
	for i := 1; i <= 5; i++ {
		emp := system.SysEmployee{
			EmployeeName:  "员工" + string(rune('0'+i)),
			EmployeePhone: "1380013800" + string(rune('0'+i)),
			EmployeeEmail: "emp" + string(rune('0'+i)) + "@example.com",
		}
		global.GVA_DB.Create(&emp)
	}

	// 创建 gin.Context - 最大的 pageSize
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeList?page=1&pageSize=100", nil)

	api := EmployeeApi{}
	api.GetEmployeeList(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestEmployeeApi_GetEmployeeList_PageOutOfRange(t *testing.T) {
	// 先创建一些员工
	for i := 1; i <= 3; i++ {
		emp := system.SysEmployee{
			EmployeeName:  "员工" + string(rune('0'+i)),
			EmployeePhone: "1380013800" + string(rune('0'+i)),
			EmployeeEmail: "emp" + string(rune('0'+i)) + "@example.com",
		}
		global.GVA_DB.Create(&emp)
	}

	// 创建 gin.Context - 超出范围的页码
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/employee/getEmployeeList?page=100&pageSize=10", nil)

	api := EmployeeApi{}
	api.GetEmployeeList(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	// 验证返回空列表
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["total"])
}

// 为了保持兼容性，需要添加 request 包别名
// 这里的 request.PageInfo 实际上来自 model/common/request
func init() {
	gin.SetMode(gin.TestMode)
}
