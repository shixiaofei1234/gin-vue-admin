/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-17 09:48:50
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-04-01 11:33:37
 * @FilePath: \gin-vue-admin-main\server\model\system\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysProductCategory struct {
	global.GVA_MODEL
	CategoryName     string `json:"categoryName" gorm:"comment:分类名称"`
	CategoryCode     string `json:"categoryCode" gorm:"comment:分类编码"`
	CategoryParentID uint   `json:"categoryParentID" gorm:"comment:父分类ID"`
	CategoryStatus   bool   `json:"categoryStatus" gorm:"comment:分类状态"`
	CategorySort     int    `json:"categorySort" gorm:"comment:分类排序"`
	CategoryRemark   string `json:"categoryRemark" gorm:"comment:分类备注"`
	CategoryIcon     string `json:"categoryIcon" gorm:"comment:分类图标"`
}

func (SysProductCategory) TableName() string {
	return "sys_product_category"
}
