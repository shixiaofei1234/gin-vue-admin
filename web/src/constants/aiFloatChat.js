/**
 * 悬浮窗「小飞」与独立 AI 页共用的文案、路由名、快捷操作配置（避免散落在多个 .vue 里写死）
 */

export const AI_ASSISTANT = Object.freeze({
  name: '小飞',
  floatButtonChar: '飞',
  brandSubtitle: 'shop-test 项目助手',
  userLabel: '我',
  assistantLabel: '小飞'
})

/** 独立聊天页停用说明（与悬浮窗助手名一致） */
export const AI_DISABLED_PAGE = Object.freeze({
  title: '独立聊天页已停用',
  /** 正文里会插入助手名与悬浮按钮上的字 */
  templateLines: (assistantName, floatChar) => [
    `当前项目仅保留右侧悬浮窗「${assistantName}」，请在任意页面点击右侧`,
    `「${floatChar}」按钮进行对话。`
  ]
})

export const AI_FLOAT_STORAGE = Object.freeze({
  openEmployee: 'ai_open_employee'
})

export const AI_FLOAT_UI = Object.freeze({
  greeting: '你好，我是 shop-test 项目助手。',
  inputPlaceholder: '有问题尽管问小飞…',
  send: '发送',
  loadingText: '正在生成回复',
  modeHintLabel: '链路',
  emptyReply: '抱歉，本次没有拿到有效回复。',
  requestFail: '请求失败，请稍后重试',
  requestFailContent: '请求失败，请稍后重试。',
  streamSingleChunkWarn: '当前虽走流式接口，但仅收到1个分片，可能被网关/服务缓冲',
  fallbackModeWarn: '流式接口未生效，已自动回退普通模式',
  bufferedModeWarn: '流式接口已命中，但响应被缓冲，体感会接近普通模式'
})

/** 与路由 name 一致（gin-vue-admin 动态路由里超级管理员下常见 name） */
export const AI_ROUTE_NAMES = Object.freeze({
  apiManage: 'api',
  roleManage: 'authority'
})

/** API 流程卡片底部按钮：跳转 + 失败时菜单提示 */
export const API_FLOW_QUICK_ACTIONS = Object.freeze([
  {
    routeName: AI_ROUTE_NAMES.apiManage,
    label: '去创建 / 登记 API',
    navigateFail: '无法跳转，请从左侧菜单进入：超级管理员 → API 管理'
  },
  {
    routeName: AI_ROUTE_NAMES.roleManage,
    label: '去角色管理授权',
    navigateFail: '无法跳转，请从左侧菜单进入：超级管理员 → 角色管理'
  }
])

export const API_FLOW_DEFAULT_EMBED = () => ({
  kind: 'apiFlow',
  title: '新增 API 流程',
  steps: [
    '后端：写路由、Handler，在 router 里注册',
    '后台：超级管理员 → API 管理，登记路径与方法',
    '权限：角色管理 → 为角色勾选该 API'
  ]
})

export const EMPLOYEE_DETAIL_EMBED = Object.freeze({
  sectionTitle: '快捷操作',
  descBefore:
    '我无法代替你操作系统，但可以在本页已打开「员工列表」时，直接弹出该员工的详情窗口。若当前不在列表页，请先通过左侧菜单进入 ',
  menuPathStrong: '员工管理 → 员工列表',
  descAfter: '，再点下方按钮（或进入列表后会自动尝试打开）。',
  toastAfterEmit:
    '已尝试打开详情：若当前在员工列表页会立即弹出；否则请先进入「员工管理 → 员工列表」。',
  buttonDefault: '打开员工详情',
  buttonById: (id) => `打开员工详情（ID ${id}）`,
  buttonByName: (name) => `打开「${name}」详情`
})
