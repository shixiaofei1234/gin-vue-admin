/*
 * @Author: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @Date: 2026-03-19 14:19:19
 * @LastEditors: shixiaofei1234 31613391+shixiaofei1234@users.noreply.github.com
 * @LastEditTime: 2026-03-19 14:19:46
 * @FilePath: \gin-vue-admin-main\server\model\system\response\sys_employee.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/system"

type SysEmployeeResponse struct {
	Employee system.SysEmployee `json:"employee"`
}
