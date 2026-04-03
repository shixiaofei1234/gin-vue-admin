package system

import (
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelpkg "github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

const maxMenuTreeLines = 150

// buildDynamicMenuContext 从数据库读取当前登录角色**实际可见**侧边栏菜单树（与 /menu/getMenu 一致），含自定义菜单。
func buildDynamicMenuContext(authorityID uint) string {
	if authorityID == 0 || global.GVA_DB == nil {
		return ""
	}
	var ms MenuService
	menus, err := ms.GetMenuTree(authorityID)
	if err != nil || len(menus) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("【当前登录角色可见菜单（来自数据库权限，与左侧栏一致；回答「在哪个菜单」时以此为最高优先级）】\n")
	n := 0
	var walk func([]modelpkg.SysMenu, int)
	walk = func(nodes []modelpkg.SysMenu, depth int) {
		for i := range nodes {
			if n >= maxMenuTreeLines {
				return
			}
			menu := nodes[i]
			if menu.Hidden {
				continue
			}
			title := strings.TrimSpace(menu.Meta.Title)
			if title == "" {
				title = strings.TrimSpace(menu.Name)
			}
			if title == "" {
				continue
			}
			pad := strings.Repeat("  ", depth)
			b.WriteString(fmt.Sprintf("%s- %s\n", pad, title))
			n++
			if len(menu.Children) > 0 {
				walk(menu.Children, depth+1)
			}
		}
	}
	walk(menus, 0)
	if n == 0 {
		return ""
	}
	if n >= maxMenuTreeLines {
		b.WriteString("…（菜单过多已截断）\n")
	}
	return strings.TrimSpace(b.String())
}
