/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-04-01 10:11:47
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-04-01 10:14:01
 * @FilePath: \server\model\system\request\sys_product_category.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SysProductCategorySearch struct {
	request.PageInfo
	CategoryName     string `json:"categoryName" form:"categoryName"`
	CategoryCode     string `json:"categoryCode" form:"categoryCode"`
	CategoryParentID uint   `json:"categoryParentID" form:"categoryParentID"`
	CategoryStatus   bool   `json:"categoryStatus" form:"categoryStatus"`
}
