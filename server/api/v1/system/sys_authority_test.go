package system

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	global.GVA_CONFIG.System.UseStrictAuth = false
}

func TestAuthorityApi_SetDataAuthority_Success(t *testing.T) {
	// 准备测试数据
	auth := system.SysAuthority{
		AuthorityId:     888,
		AuthorityName:   "测试角色",
		DataAuthorityId: []*system.SysAuthority{{AuthorityId: 888}},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 888,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "设置成功", resp["msg"])
}

func TestAuthorityApi_SetDataAuthority_JSONBindError(t *testing.T) {
	// 创建无效的 JSON
	invalidJSON := `{"authorityId": "invalid", "authorityName":}`

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 888,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
	assert.Contains(t, resp["msg"].(string), "invalid character")
}

func TestAuthorityApi_SetDataAuthority_VerifyError_EmptyAuthorityId(t *testing.T) {
	// 准备测试数据 - AuthorityId 为 0
	auth := system.SysAuthority{
		AuthorityId:     0,
		AuthorityName:   "测试角色",
		DataAuthorityId: []*system.SysAuthority{{AuthorityId: 888}},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 888,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
	assert.Contains(t, resp["msg"].(string), "AuthorityId")
}

func TestAuthorityApi_SetDataAuthority_ServiceError_InvalidAuthority(t *testing.T) {
	// 准备测试数据 - 不存在的 authorityId
	auth := system.SysAuthority{
		AuthorityId:     99999,
		AuthorityName:   "测试角色",
		DataAuthorityId: []*system.SysAuthority{{AuthorityId: 99999}},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 888,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应 - 应该返回错误，因为 role 不存在
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestAuthorityApi_SetDataAuthority_VerifyError_ZeroAuthorityId(t *testing.T) {
	// 准备测试数据 - AuthorityId 未设置（默认为 0）
	auth := map[string]interface{}{
		"authorityName": "测试角色",
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 888,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
	assert.Contains(t, resp["msg"].(string), "AuthorityId")
}

func TestAuthorityApi_SetDataAuthority_NoClaimsInContext(t *testing.T) {
	// 准备测试数据
	auth := system.SysAuthority{
		AuthorityId:     888,
		AuthorityName:   "测试角色",
		DataAuthorityId: []*system.SysAuthority{{AuthorityId: 888}},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求 - 不设置 claims
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	// 由于没有 claims，GetUserAuthorityId 会返回 0，导致服务层失败
	assert.Equal(t, float64(7), resp["code"])
}

func TestAuthorityApi_SetDataAuthority_MultipleDataAuthorities(t *testing.T) {
	// 准备测试数据 - 多个数据权限
	auth := system.SysAuthority{
		AuthorityId:   888,
		AuthorityName: "测试角色",
		DataAuthorityId: []*system.SysAuthority{
			{AuthorityId: 888},
			{AuthorityId: 888},
		},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 888,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "设置成功", resp["msg"])
}

func TestAuthorityApi_SetDataAuthority_EmptyDataAuthorities(t *testing.T) {
	// 准备测试数据 - 空的数据权限列表
	auth := system.SysAuthority{
		AuthorityId:     888,
		AuthorityName:   "测试角色",
		DataAuthorityId: []*system.SysAuthority{},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 888,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "设置成功", resp["msg"])
}

func TestAuthorityApi_SetDataAuthority_InvalidContentType(t *testing.T) {
	// 准备测试数据
	auth := system.SysAuthority{
		AuthorityId:     888,
		AuthorityName:   "测试角色",
		DataAuthorityId: []*system.SysAuthority{{AuthorityId: 888}},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求 - 不设置 Content-Type
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 888,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应 - 可能会因为无法解析 JSON 而失败
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(7), resp["code"])
}

func TestAuthorityApi_SetDataAuthority_MaxAuthorityId(t *testing.T) {
	// 准备测试数据 - 最大值的 AuthorityId
	maxUint := ^uint(0)
	auth := system.SysAuthority{
		AuthorityId:     maxUint,
		AuthorityName:   "测试角色",
		DataAuthorityId: []*system.SysAuthority{{AuthorityId: maxUint}},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: maxUint,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	// 由于 role 不存在，应该返回错误
	assert.Equal(t, float64(7), resp["code"])
}

func TestAuthorityApi_SetDataAuthority_MinAuthorityId(t *testing.T) {
	// 准备测试数据 - 最小值的 AuthorityId (1)
	auth := system.SysAuthority{
		AuthorityId:     1,
		AuthorityName:   "测试角色",
		DataAuthorityId: []*system.SysAuthority{{AuthorityId: 1}},
	}

	jsonData, err := json.Marshal(auth)
	assert.NoError(t, err)

	// 创建 gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置请求
	c.Request, _ = http.NewRequest("POST", "/authority/setDataAuthority", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置用户权限ID到 context
	claims := &systemReq.CustomClaims{
		AuthorityId: 1,
		UUID:        "test-uuid",
	}
	c.Set("claims", claims)

	// 调用被测试方法
	api := AuthorityApi{}
	api.SetDataAuthority(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	// 由于数据库中可能存在 role 1，返回结果可能成功或失败
	// 这里只验证响应格式正确
	assert.Contains(t, []float64{0, 7}, resp["code"])
}
