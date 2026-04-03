<template>
  <div class="ai-float-chat">
    <el-button
      v-show="!drawerVisible"
      class="ai-float-btn"
      circle
      @click="drawerVisible = true"
    >
      <span class="ai-float-btn-inner">{{ AI_ASSISTANT.floatButtonChar }}</span>
    </el-button>

    <el-drawer
      v-model="drawerVisible"
      class="ai-drawer"
      direction="rtl"
      size="420px"
      :with-header="false"
      append-to-body
    >
      <div class="chat-wrap">
        <div class="chat-header">
          <div class="chat-header-brand">
            <div class="chat-avatar">{{ AI_ASSISTANT.floatButtonChar }}</div>
            <div class="chat-header-text">
              <div class="chat-title">{{ AI_ASSISTANT.name }}</div>
              <div class="chat-sub">{{ AI_ASSISTANT.brandSubtitle }}</div>
            </div>
          </div>
          <div
            class="mode-tag"
            :class="{
              'mode-stream': streamMode === 'stream',
              'mode-fallback': streamMode === 'fallback',
              'mode-buffered': streamMode === 'buffered',
              'mode-connecting': streamMode === 'connecting'
            }"
          >
            {{ modeText }}
          </div>
        </div>
        <div class="mode-hint">
          <span class="mode-hint-label">{{ AI_FLOAT_UI.modeHintLabel }}</span>
          传输 {{ streamStats.transportChunkCount }} · 内容 {{ streamStats.contentChunkCount }} · 首包 {{ streamStats.firstChunkMs }}ms
        </div>

        <div
          ref="chatBoxRef"
          class="chat-box"
        >
          <div
            v-for="(item, idx) in messages"
            :key="idx"
            :class="['msg-row', item.role === 'user' ? 'user' : 'ai']"
          >
            <div
              class="msg-bubble"
              :class="{ 'msg-bubble--await': isAssistantAwaitingContent(idx, item) }"
            >
              <div class="msg-role">{{ item.role === 'user' ? '我' : '小飞' }}</div>
              <div
                v-if="item.content"
                class="msg-content"
              >{{ item.content }}</div>
              <div
                v-else-if="isAssistantAwaitingContent(idx, item)"
                class="msg-loading"
                role="status"
                aria-live="polite"
                aria-busy="true"
              >
                <div class="msg-loading-pulse" />
                <div class="msg-loading-main">
                  <span class="msg-loading-dots" aria-hidden="true">
                    <span class="msg-loading-dot" />
                    <span class="msg-loading-dot" />
                    <span class="msg-loading-dot" />
                  </span>
                  <span class="msg-loading-text">{{ AI_FLOAT_UI.loadingText }}</span>
                </div>
                <div class="msg-loading-skeleton">
                  <span class="msg-loading-bar" />
                  <span class="msg-loading-bar msg-loading-bar--short" />
                </div>
              </div>
              <div
                v-if="item.role === 'assistant' && item.embed === 'apiFlow' && item.embedMeta"
                class="api-flow-inline"
              >
                <!--
                  API 流程示例卡片：
                  embedMeta 来自后端的二次短调用生成 JSON（或前端兜底 ensureApiFlowEmbed）。
                  这里只负责“把 title + steps 渲染成可点击的流程卡”。
                -->
                <div class="api-flow-inline-title">{{ item.embedMeta.title || apiFlowDefaultTitle }}</div>
                <div class="api-flow-chart">
                  <template
                    v-for="(stepText, si) in (item.embedMeta.steps || [])"
                    :key="si"
                  >
                    <div class="flow-step">
                      <span class="flow-num">{{ si + 1 }}</span>
                      <span class="flow-text">{{ stepText }}</span>
                    </div>
                    <div
                      v-if="si < (item.embedMeta.steps || []).length - 1"
                      class="flow-arrow"
                    >
                      ↓
                    </div>
                  </template>
                </div>
                <div class="api-flow-actions">
                  <el-button
                    v-for="(act, ai) in API_FLOW_QUICK_ACTIONS"
                    :key="act.routeName"
                    :type="ai === 0 ? 'primary' : 'default'"
                    size="small"
                    @click="goEmbedRoute(act)"
                  >
                    {{ act.label }}
                  </el-button>
                </div>
              </div>
              <div
                v-if="item.role === 'assistant' && item.embed === 'employeeDetail' && item.embedMeta"
                class="employee-detail-inline"
              >
                <div class="employee-detail-inline-title">{{ EMPLOYEE_DETAIL_EMBED.sectionTitle }}</div>
                <p class="employee-detail-inline-desc">
                  {{ EMPLOYEE_DETAIL_EMBED.descBefore }}<strong>{{ EMPLOYEE_DETAIL_EMBED.menuPathStrong }}</strong>{{ EMPLOYEE_DETAIL_EMBED.descAfter }}
                </p>
                <div class="employee-detail-inline-actions">
                  <el-button
                    type="primary"
                    size="small"
                    @click="goOpenEmployeeDetail(item)"
                  >
                    {{ employeeDetailButtonText(item.embedMeta) }}
                  </el-button>
                </div>
              </div>
              <div
                v-if="item.role === 'assistant' && item.embed === 'teamStatusToggle' && item.embedMeta"
                class="team-status-inline"
              >
                <div class="team-status-inline-title">快捷操作</div>
                <p class="team-status-inline-desc">
                  团队「{{ item.embedMeta.teamName }}」当前状态：
                  <strong>{{ item.embedMeta.currentStatus ? '开启' : '关闭' }}</strong>
                  ，点击按钮将切换为：
                  <strong>{{ item.embedMeta.targetStatus ? '开启' : '关闭' }}</strong>
                </p>
                <div class="team-status-inline-actions">
                  <el-button
                    type="primary"
                    size="small"
                    :loading="teamToggleLoadingMap[item.embedMeta.teamID]"
                    @click="toggleTeamStatus(item, idx)"
                  >
                    {{ item.embedMeta.targetStatus ? '开启团队' : '关闭团队' }}
                  </el-button>
                </div>
              </div>

              <div
                v-if="item.role === 'assistant' && item.embed === 'teamAdminSet' && item.embedMeta"
                class="team-admin-inline"
              >
                <div class="team-admin-inline-title">快捷设置</div>
                <p class="team-admin-inline-desc">
                  团队「{{ item.embedMeta.teamName }}」管理员当前为：
                  <strong>{{ item.embedMeta.currentAdminName || item.embedMeta.targetAdminName }}</strong>
                  ，点击按钮将设置为：
                  <strong>{{ item.embedMeta.targetAdminName }}</strong>
                </p>
                <div class="team-admin-inline-actions">
                  <el-button
                    type="primary"
                    size="small"
                    :loading="teamAdminLoadingMap[item.embedMeta.teamID]"
                    @click="setTeamAdminQuick(item)"
                  >
                    设置管理员
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="chat-input-bar">
          <el-input
            v-model="inputValue"
            class="chat-input-field"
            :placeholder="AI_FLOAT_UI.inputPlaceholder"
            clearable
            @keyup.enter="handleSend"
          />
          <el-button
            class="chat-send-btn"
            type="primary"
            :loading="sending"
            @click="handleSend"
          >
            {{ AI_FLOAT_UI.send }}
          </el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, nextTick, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { chatWithAIStream } from '@/api/aiChat'
import { emitter } from '@/utils/bus.js'
import { getTeamEmployeeList, setTeamAdmin, switchTeamStatus } from '@/api/team'
import { getTeamList } from '@/api/employee'
import {
  AI_ASSISTANT,
  AI_FLOAT_STORAGE,
  AI_FLOAT_UI,
  API_FLOW_DEFAULT_EMBED,
  API_FLOW_QUICK_ACTIONS,
  EMPLOYEE_DETAIL_EMBED
} from '@/constants/aiFloatChat'

const router = useRouter()

const apiFlowDefaultTitle = API_FLOW_DEFAULT_EMBED().title

defineOptions({
  name: 'AIFloatChat'
})

const drawerVisible = ref(false)
const inputValue = ref('')
const sending = ref(false)
const chatBoxRef = ref()
const MAX_HISTORY_MESSAGES = 4
const streamMode = ref('stream')
const streamStats = ref({
  transportChunkCount: 0,
  contentChunkCount: 0,
  firstChunkMs: 0
})
const modeTextMap = {
  stream: '流式',
  fallback: '普通(回退)',
  buffered: '疑似缓冲',
  connecting: '连接中'
}
const modeText = computed(() => modeTextMap[streamMode.value] || '未知')
const messages = ref([
  { role: 'assistant', content: AI_FLOAT_UI.greeting }
])
// team 状态切换按钮 loading：按 teamID 细粒度控制，避免全局禁用
const teamToggleLoadingMap = ref({})
const teamAdminLoadingMap = ref({})
const cotPattern = /(hypothesis|final decision|self-correction|let'?s|wait,|i should|rule\s*\d|internal plan|thought process|reasoning)/i

const sanitizeChunk = (text) => {
  if (!text) return ''
  const lines = text.split('\n')
  const filtered = lines.filter((line) => !cotPattern.test(line.trim()))
  return filtered.join('\n').trim()
}

// 前端兜底：当后端没有下发 `teamStatusToggle` embed，
// 但回复文本里包含「团队（ID:x）当前状态为【启用/关闭】」，则直接解析并渲染可点击按钮。
const parseTeamStatusFromAssistantText = (text) => {
  if (!text) return null
  const raw = String(text)

  // teamID：支持 ID:2 / ID：2
  const idMatch = raw.match(/ID\s*[:：]\s*(\d+)/i)
  const teamID = idMatch ? Number(idMatch[1]) : null

  // 当前状态：兼容【启用】、「启用」、“启用”、'启用'、(启用) 以及无括号写法
  const statusMatch = raw.match(/当前状态为(?:\s*[「『“"'【\(\[]\s*)?(启用|开启|关闭|停用|禁用|已启用|已开启|已关闭|已禁用)(?:\s*[」』”"'】\)\]])?/i)
  if (!statusMatch) return null
  const statusText = statusMatch[1]
  const currentStatus = /(启用|开启|已启用|已开启)/i.test(statusText)
  const targetStatus = !currentStatus

  // 团队名：优先从 “马里奥团队（ID:2）” 提取“马里奥”
  const nameMatch = raw.match(/([^\n。；;:\(（)）]{1,30})团队\s*[（(]\s*ID\s*[:：]\s*\d+\s*[）)]/)
  const teamNameWithID = nameMatch ? String(nameMatch[1]).trim() : ''
  const nameMatchNoID = raw.match(/([^\n。；;:\(（)）]{1,30})团队.*?当前状态为/i)
  const teamNameNoID = nameMatchNoID ? String(nameMatchNoID[1]).trim() : ''
  const teamName = teamNameWithID || teamNameNoID

  if (!teamName) return null

  return {
    kind: 'teamStatusToggle',
    teamID: teamID || 0,
    teamName,
    currentStatus,
    targetStatus
  }
}

const resolveTeamIDByName = async(teamName) => {
  const name = String(teamName || '').trim()
  if (!name) return 0

  try {
    const res = await getTeamList({
      page: 1,
      pageSize: 50,
      teamName: name
    })
    if (res?.code !== 0) return 0

    const list = res.data?.list ?? []
    const exact = list.find((t) => String(t.teamName || '').trim() === name)
    if (exact?.ID) return exact.ID

    // 兜底模糊：包含任意关键词
    const fuzzy = list.find((t) => String(t.teamName || '').includes(name) || name.includes(String(t.teamName || '')))
    return fuzzy?.ID || 0
  } catch (_) {
    return 0
  }
}

// 前端兜底：当后端没有下发 teamAdminSet embed，
// 但回复文本里包含“团队/团 + 管理员为（或是）某人”，则渲染可点击按钮完成管理员变更。
const parseTeamAdminSetFromAssistantText = (text) => {
  if (!text) return null
  const raw = String(text)

  // 团队名：优先从 “X团队...管理员” 里取 X
  const teamMatch = raw.match(/(.+?)(团队|队伍|组|团).*?管理员/i)
  let teamName = teamMatch ? String(teamMatch[1]).trim() : ''

  // 清理常见前缀，避免把“我想设置/请您前往”等带进去
  const prefixes = ['我想设置', '我想', '请您', '你', '该', '该团队', '马里奥团', '马里奥团队']
  for (const p of prefixes) {
    if (teamName.startsWith(p)) teamName = teamName.slice(p.length).trim()
  }
  teamName = teamName.replace(/^(帮我|请把|请|将|我要|设置|更改)\s*/i, '').replace(/的$/, '').trim()

  if (!teamName) return null

  // 目标管理员：匹配 “管理员为X” / “管理员是X”
  const adminMatch = raw.match(/管理员(?:为|是)\s*[「『“"'【\(\[]?\s*([^。；;，,、\n]{1,24})/)
  const targetAdminName = adminMatch ? String(adminMatch[1]).trim() : ''
  if (!targetAdminName) return null

  return {
    kind: 'teamAdminSet',
    teamName,
    targetAdminName
  }
}

const resolveTeamAdminIDByTeamAndName = async(teamName, targetAdminName) => {
  const teamID = await resolveTeamIDByName(teamName)
  if (!teamID) return { teamID: 0, adminID: 0 }

  try {
    const res = await getTeamEmployeeList({
      teamID,
      page: 1,
      pageSize: 50
    })
    if (res?.code !== 0) return { teamID, adminID: 0 }

    const list = res.data?.list ?? []
    const normTarget = String(targetAdminName || '').trim()
    const exact = list.find((e) => String(e?.employeeName || '').trim() === normTarget)
    const best = exact || list.find((e) => String(e?.employeeName || '').includes(normTarget) || normTarget.includes(String(e?.employeeName || '')))

    return { teamID, adminID: best?.ID ? Number(best.ID) : 0 }
  } catch (_) {
    return { teamID, adminID: 0 }
  }
}

/** 用户问题是否与「新建/登记 API」相关，命中则在助手气泡内展示流程图例 */
/** 请求已发出、尚未收到可展示的正文时（首包前或流式开始前） */
const isAssistantAwaitingContent = (idx, item) => {
  if (item.role !== 'assistant' || !sending.value) return false
  if (idx !== messages.value.length - 1) return false
  return !String(item.content || '').trim()
}

/**
 * 是否展示「API 流程」示例卡片：
 * - 前端与后端 DetectApiFlowIntent 会尽量保持一致，但前端做“略宽”的兜底，
 *   目的：避免模型/翻译导致命中意图边界不稳，从而卡片不出现。
 * - 最终渲染仍依赖 embedMeta 的归一化（ensureApiFlowEmbed），避免 JSON 结构不对导致空白卡片。
 */
const shouldEmbedApiFlow = (userText) => {
  const t = (userText || '').trim()
  if (t.length < 4) return false
  // 允许用户输入带空格/点号的情况：如 "a p i"、"a.p.i"
  const hasApi =
    /接口/i.test(t) ||
    /api/i.test(t) ||
    /a[\s\.\-]*p[\s\.\-]*i/i.test(t)
  if (!hasApi) return false
  const intent = /创建|新增|添加|建立|注册|登记|写.*路由|怎么|如何|怎样|步骤|流程|整一个|做一个/i.test(t)
  const colloquial = /想|要|问|如何|怎么|怎样|操作|教我|弄|搞|整一个/i.test(t)
  return intent || colloquial
}

/**
 * 用新对象替换消息，保证：
 * 1) embed/embedMeta 的响应式更新一定生效（Vue 需要新对象触发）
 * 2) title/steps 做归一化 + 必要时强制回退默认值（避免“只显示按钮不显示示例图”）
 */
const ensureApiFlowEmbed = (userText, idx) => {
  if (idx < 0 || !shouldEmbedApiFlow(userText)) return
  const m = messages.value[idx]
  if (!m || m.role !== 'assistant') return
  const fallback = API_FLOW_DEFAULT_EMBED()

  const rawTitle = m.embedMeta?.title
  const title = typeof rawTitle === 'string' && rawTitle.trim()
    ? rawTitle.trim()
    : fallback.title

  const rawSteps = m.embedMeta?.steps
  const normalizedSteps = Array.isArray(rawSteps)
    ? rawSteps
      .map((s) => String(s || '').trim())
      .filter(Boolean)
      .slice(0, 5)
    : []

  const steps = normalizedSteps.length > 0
    ? normalizedSteps
    : fallback.steps

  messages.value[idx] = {
    ...m,
    embed: 'apiFlow',
    embedMeta: {
      kind: 'apiFlow',
      title,
      steps
    }
  }
}

/** 打开源码/文件/路由 等，不当作「业务里的员工详情」 */
const EMPLOYEE_DETAIL_CODE_CTX = /代码|源码|文件|\.vue|路由|组件|IDE|仓库|git|vscode|cursor/i

/**
 * 从用户话里解析「打开员工详情」：支持员工姓名或员工 ID
 * @returns {{ employeeName?: string, employeeId?: number } | null}
 */
const parseEmployeeDetailIntent = (userText) => {
  const t = (userText || '').trim()
  if (t.length < 4) return null
  if (EMPLOYEE_DETAIL_CODE_CTX.test(t)) return null
  if (!/(详情|资料|打开|查看|进入|显示)/.test(t)) return null
  if (!/(员工|人员)/.test(t)) return null

  const idMatch = t.match(/(?:ID|id|编号)\s*[:：]?\s*(\d+)/)
  if (idMatch) return { employeeId: Number(idMatch[1]) }

  const patterns = [
    /员工\s*[「『]?(.+?)[」』]?(?:\s*的)?(?:详情|资料)/,
    /(?:打开|查看|进入|显示).{0,24}?员工\s*[「『]?(.+?)[」』]?(?:\s*的)?(?:详情|资料)?/,
    /(?:人员|员工)\s*[:：]\s*['「]?(.+?)['」]?\s*(?:的)?(?:详情|资料)/
  ]
  for (const re of patterns) {
    const m = t.match(re)
    if (m && m[1]) {
      const name = m[1].trim().replace(/^(的|员工|人员)/, '').replace(/的$/, '').trim()
      if (name.length >= 1 && name.length <= 24) return { employeeName: name }
    }
  }
  return null
}

const shouldEmbedEmployeeDetail = (userText) => !!parseEmployeeDetailIntent(userText)

const employeeDetailButtonText = (meta) => {
  if (!meta) return EMPLOYEE_DETAIL_EMBED.buttonDefault
  if (meta.employeeId != null) return EMPLOYEE_DETAIL_EMBED.buttonById(meta.employeeId)
  if (meta.employeeName) return EMPLOYEE_DETAIL_EMBED.buttonByName(meta.employeeName)
  return EMPLOYEE_DETAIL_EMBED.buttonDefault
}

const goOpenEmployeeDetail = (item) => {
  const meta = item?.embedMeta
  if (!meta) return
  try {
    sessionStorage.setItem(AI_FLOAT_STORAGE.openEmployee, JSON.stringify({ ...meta, ts: Date.now() }))
  } catch (_) {}
  drawerVisible.value = false
  emitter.emit('openEmployeeDetail', meta)
  ElMessage.info(EMPLOYEE_DETAIL_EMBED.toastAfterEmit)
}

const goEmbedRoute = (act) => {
  drawerVisible.value = false
  router.push({ name: act.routeName }).catch(() => {
    ElMessage.error(act.navigateFail)
  })
}

const toggleTeamStatus = async(item, idx) => {
  const meta = item?.embedMeta
  if (!meta?.teamID) return

  const teamID = meta.teamID
  const targetStatus = Boolean(meta.targetStatus)

  teamToggleLoadingMap.value[teamID] = true
  try {
    const res = await switchTeamStatus({
      ID: teamID,
      status: targetStatus
    })
    if (res?.code !== 0) {
      ElMessage.error(res?.msg || '切换团队状态失败')
      return
    }
    ElMessage.success(targetStatus ? '团队已开启' : '团队已关闭')

    // 如果用户当前在“团队管理”页面，需要刷新列表以保持按钮状态与表格一致
    emitter.emit('teamListReload')

    // 更新按钮：currentStatus=targetStatus；下一次点击则反向切换
    const updatedMeta = {
      ...meta,
      currentStatus: targetStatus,
      targetStatus: !targetStatus
    }
    messages.value[idx] = {
      ...item,
      embedMeta: updatedMeta
    }
  } catch (e) {
    ElMessage.error('切换团队状态失败')
  } finally {
    teamToggleLoadingMap.value[teamID] = false
  }
}

const setTeamAdminQuick = async(item) => {
  const meta = item?.embedMeta
  if (!meta?.teamID || !meta?.targetAdminName) return

  const teamID = meta.teamID
  const targetAdminName = meta.targetAdminName

  teamAdminLoadingMap.value[teamID] = true
  try {
    const teamEmployeeRes = await getTeamEmployeeList({
      teamID,
      page: 1,
      pageSize: 50
    })
    if (teamEmployeeRes?.code !== 0) {
      ElMessage.error(teamEmployeeRes?.msg || '加载团队员工失败')
      return
    }

    const list = teamEmployeeRes.data?.list ?? []
    const normTarget = String(targetAdminName || '').trim()
    const exact = list.find((e) => String(e?.employeeName || '').trim() === normTarget)
    const best = exact || list.find((e) => String(e?.employeeName || '').includes(normTarget) || normTarget.includes(String(e?.employeeName || '')))

    if (!best?.ID) {
      ElMessage.error(`未在团队内找到员工「${targetAdminName}」`)
      return
    }

    const res = await setTeamAdmin({
      ID: teamID,
      adminID: Number(best.ID)
    })
    if (res?.code !== 0) {
      ElMessage.error(res?.msg || '设置管理员失败')
      return
    }

    ElMessage.success(`管理员已设置为「${targetAdminName}」`)
    emitter.emit('teamListReload')
  } catch (e) {
    ElMessage.error('设置管理员失败')
  } finally {
    teamAdminLoadingMap.value[teamID] = false
  }
}

const scrollToBottom = async() => {
  await nextTick()
  if (!chatBoxRef.value) return
  chatBoxRef.value.scrollTop = chatBoxRef.value.scrollHeight
}

const buildHistory = () => {
  const recentMessages = messages.value.slice(-MAX_HISTORY_MESSAGES)
  const chatHistory = recentMessages.map((item) => ({
    role: item.role === 'assistant' ? 'assistant' : 'user',
    content: item.content
  }))
  return chatHistory
}

const handleSend = async() => {
  const message = inputValue.value.trim()
  if (!message || sending.value) return

  const history = buildHistory()
  messages.value.push({ role: 'user', content: message })
  inputValue.value = ''
  await scrollToBottom()

  sending.value = true
  let aiMsgIndex = -1
  try {
    messages.value.push({ role: 'assistant', content: '' })
    aiMsgIndex = messages.value.length - 1
    await chatWithAIStream({
      message,
      history,
      temperature: 0.5
    }, {
      onMode: (mode) => {
        streamMode.value = mode
        if (mode === 'fallback') {
          ElMessage.warning(AI_FLOAT_UI.fallbackModeWarn)
        }
        if (mode === 'buffered') {
          ElMessage.warning(AI_FLOAT_UI.bufferedModeWarn)
        }
      },
      onChunk: (chunk) => {
        const safeChunk = sanitizeChunk(chunk)
        if (!safeChunk) return
        messages.value[aiMsgIndex].content += `${safeChunk}\n`
        scrollToBottom()
      },
      onEmbed: (embed) => {
        // 后端会在流结束后额外推送 embed 信息（例如 apiFlow 卡片 JSON）。
        // 前端只负责“把 embed 写入当前 assistant 消息”，模板再渲染卡片。
        if (!embed) return
        const kind =
          embed.kind ||
          (Array.isArray(embed.steps) && embed.steps.length > 0 ? 'apiFlow' : null)
        if (!kind) return

        const cur = messages.value[aiMsgIndex]
        if (kind === 'apiFlow') {
          messages.value[aiMsgIndex] = {
            ...cur,
            embed: 'apiFlow',
            embedMeta:
              embed.kind === 'apiFlow'
                ? embed
                : { kind: 'apiFlow', title: embed.title, steps: embed.steps }
          }
          return
        }

        messages.value[aiMsgIndex] = {
          ...cur,
          embed: kind,
          embedMeta: embed
        }
      },
      onMetrics: (metrics) => {
        streamStats.value = metrics
      }
    })
    if (streamMode.value === 'stream' && streamStats.value.contentChunkCount <= 1) {
      ElMessage.warning(AI_FLOAT_UI.streamSingleChunkWarn)
    }
    if (!messages.value[aiMsgIndex].content) {
      messages.value[aiMsgIndex].content = AI_FLOAT_UI.emptyReply
    }
  } catch (e) {
    ElMessage.error(AI_FLOAT_UI.requestFail)
    if (aiMsgIndex >= 0 && messages.value[aiMsgIndex] && !String(messages.value[aiMsgIndex].content || '').trim()) {
      messages.value[aiMsgIndex].content = AI_FLOAT_UI.requestFailContent
    }
  } finally {
    if (aiMsgIndex >= 0) {
      ensureApiFlowEmbed(message, aiMsgIndex)
      if (!messages.value[aiMsgIndex].embed) {
        // 1) 团队状态切换按钮兜底：从“助手回复文本”解析
        const assistantText = messages.value[aiMsgIndex]?.content
        const teamMeta = parseTeamStatusFromAssistantText(assistantText)
        if (teamMeta) {
          let resolved = teamMeta
          // 如果回复文本没有携带 ID，则用团队名再查一次 ID
          if (!resolved.teamID) {
            const resolvedTeamID = await resolveTeamIDByName(resolved.teamName)
            if (resolvedTeamID) {
              resolved = { ...resolved, teamID: resolvedTeamID }
            }
          }
          // teamID 还是拿不到时不展示按钮，避免点击失败
          if (resolved.teamID) {
            const cur = messages.value[aiMsgIndex]
            messages.value[aiMsgIndex] = {
              ...cur,
              embed: 'teamStatusToggle',
              embedMeta: resolved
            }
          }
        } else {
          // 2) 团队管理员变更按钮兜底：从“助手回复文本”解析
          const teamAdminMeta = parseTeamAdminSetFromAssistantText(assistantText)
          if (teamAdminMeta) {
            const cur = messages.value[aiMsgIndex]
            let resolved = teamAdminMeta

            // 先解析 teamID
            if (!resolved.teamID) {
              const resolvedTeamID = await resolveTeamIDByName(resolved.teamName)
              if (resolvedTeamID) resolved = { ...resolved, teamID: resolvedTeamID }
            }

            // 再解析 adminID（需要拉取该团队员工列表）
            if (resolved.teamID) {
              const empRes = await getTeamEmployeeList({
                teamID: resolved.teamID,
                page: 1,
                pageSize: 50
              })
              if (empRes?.code === 0) {
                const list = empRes.data?.list ?? []
                const normTarget = String(resolved.targetAdminName || '').trim()
                const exact = list.find((e) => String(e?.employeeName || '').trim() === normTarget)
                const best = exact || list.find((e) => String(e?.employeeName || '').includes(normTarget) || normTarget.includes(String(e?.employeeName || '')))
                if (best?.ID) {
                  resolved = { ...resolved, adminID: Number(best.ID), currentAdminName: best.employeeName }
                }
              }
            }

            // 找不到 adminID 就不展示按钮（避免点击失败）
            if (resolved.teamID && resolved.adminID) {
              messages.value[aiMsgIndex] = {
                ...cur,
                embed: 'teamAdminSet',
                embedMeta: resolved
              }
            }
          } else {
            // 3) 员工详情兜底：从用户问题解析
            const empMeta = parseEmployeeDetailIntent(message)
            if (empMeta) {
              const cur = messages.value[aiMsgIndex]
              messages.value[aiMsgIndex] = {
                ...cur,
                embed: 'employeeDetail',
                embedMeta: empMeta
              }
            }
          }
        }
      }
    }
  sending.value = false
    scrollToBottom()
  }
}
</script>

<style scoped>
/* 抽屉整体：浅色渐变底 */
.ai-drawer :deep(.el-drawer__body) {
  padding: 0;
  background: linear-gradient(165deg, #f8fafc 0%, #eef2ff 38%, #f1f5f9 100%);
}

.ai-float-btn {
  --el-button-bg-color: transparent;
  --el-button-border-color: transparent;
  --el-button-hover-bg-color: transparent;
  --el-button-hover-border-color: transparent;
  --el-button-active-bg-color: transparent;
  --el-button-active-border-color: transparent;
  position: fixed;
  right: 18px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 10020;
  width: 54px;
  height: 54px;
  padding: 0;
  border: none;
  background: linear-gradient(145deg, #4f46e5 0%, #2563eb 50%, #0ea5e9 100%);
  box-shadow:
    0 4px 14px rgba(37, 99, 235, 0.45),
    0 0 0 1px rgba(255, 255, 255, 0.2) inset;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.ai-float-btn:hover {
  transform: translateY(-50%) scale(1.06);
  box-shadow:
    0 8px 24px rgba(37, 99, 235, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.25) inset;
}

.ai-float-btn-inner {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.02em;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
}

.chat-wrap {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px 16px 14px;
  box-sizing: border-box;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.35);
}

.chat-header-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.chat-avatar {
  flex-shrink: 0;
  width: 44px;
  height: 44px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 800;
  color: #fff;
  background: linear-gradient(135deg, #6366f1 0%, #3b82f6 55%, #06b6d4 100%);
  box-shadow:
    0 6px 16px rgba(99, 102, 241, 0.35),
    0 0 0 1px rgba(255, 255, 255, 0.15) inset;
}

.chat-header-text {
  min-width: 0;
}

.chat-title {
  font-weight: 700;
  font-size: 17px;
  letter-spacing: 0.02em;
  color: #0f172a;
  line-height: 1.25;
}

.chat-sub {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
  line-height: 1.3;
}

.mode-tag {
  flex-shrink: 0;
  font-size: 11px;
  font-weight: 600;
  line-height: 1;
  border-radius: 999px;
  padding: 6px 11px;
  border: 1px solid transparent;
  letter-spacing: 0.02em;
}

.mode-stream {
  color: #166534;
  background: linear-gradient(180deg, #ecfdf5 0%, #d1fae5 100%);
  border-color: rgba(34, 197, 94, 0.45);
  box-shadow: 0 1px 2px rgba(22, 101, 52, 0.08);
}

.mode-fallback {
  color: #9a3412;
  background: linear-gradient(180deg, #fff7ed 0%, #ffedd5 100%);
  border-color: rgba(249, 115, 22, 0.4);
  box-shadow: 0 1px 2px rgba(154, 52, 18, 0.08);
}

.mode-buffered {
  color: #7c2d12;
  background: linear-gradient(180deg, #fef2f2 0%, #fee2e2 100%);
  border-color: rgba(239, 68, 68, 0.35);
  box-shadow: 0 1px 2px rgba(124, 45, 18, 0.08);
}

.mode-connecting {
  color: #1d4ed8;
  background: linear-gradient(180deg, #eff6ff 0%, #dbeafe 100%);
  border-color: rgba(59, 130, 246, 0.45);
  box-shadow: 0 1px 2px rgba(29, 78, 216, 0.1);
}

.mode-hint {
  font-size: 11px;
  color: #64748b;
  margin-bottom: 10px;
  padding: 8px 10px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.65);
  border: 1px solid rgba(226, 232, 240, 0.9);
  backdrop-filter: blur(6px);
  font-variant-numeric: tabular-nums;
}

.mode-hint-label {
  display: inline-block;
  margin-right: 6px;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 700;
  color: #475569;
  background: rgba(148, 163, 184, 0.22);
}

.chat-box {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  border-radius: 16px;
  padding: 12px 12px 14px;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(226, 232, 240, 0.95);
  box-shadow:
    0 1px 3px rgba(15, 23, 42, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.85);
}

.chat-box::-webkit-scrollbar {
  width: 6px;
}

.chat-box::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.45);
  border-radius: 999px;
}

.chat-box::-webkit-scrollbar-track {
  background: transparent;
}

/* 嵌在助手气泡内的流程图例 */
.api-flow-inline {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed rgba(148, 163, 184, 0.65);
}

.api-flow-inline-title {
  font-size: 12px;
  font-weight: 700;
  color: #1e40af;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.api-flow-inline-title::before {
  content: '';
  width: 4px;
  height: 14px;
  border-radius: 2px;
  background: linear-gradient(180deg, #3b82f6, #6366f1);
}

.api-flow-chart {
  font-size: 12px;
  color: #334155;
  line-height: 1.5;
}

.flow-step {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.flow-num {
  flex-shrink: 0;
  width: 22px;
  height: 22px;
  border-radius: 999px;
  background: linear-gradient(145deg, #3b82f6, #6366f1);
  color: #fff;
  font-size: 11px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 6px rgba(59, 130, 246, 0.35);
}

.flow-text {
  flex: 1;
}

.flow-arrow {
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
  line-height: 1.2;
  padding: 3px 0;
}

.api-flow-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.employee-detail-inline {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed rgba(148, 163, 184, 0.65);
}

.employee-detail-inline-title {
  font-size: 12px;
  font-weight: 700;
  color: #0f766e;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.employee-detail-inline-title::before {
  content: '';
  width: 4px;
  height: 14px;
  border-radius: 2px;
  background: linear-gradient(180deg, #14b8a6, #0d9488);
}

.employee-detail-inline-desc {
  margin: 0 0 10px;
  font-size: 12px;
  line-height: 1.55;
  color: #475569;
}

.employee-detail-inline-desc strong {
  color: #334155;
  font-weight: 600;
}

.employee-detail-inline-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.team-status-inline {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed rgba(148, 163, 184, 0.65);
}

.team-status-inline-title {
  font-size: 12px;
  font-weight: 700;
  color: #0f766e;
  margin-bottom: 8px;
}

.team-status-inline-desc {
  margin: 0 0 10px;
  font-size: 12px;
  line-height: 1.55;
  color: #475569;
}

.team-status-inline-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.team-admin-inline {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed rgba(99, 102, 241, 0.55);
}

.team-admin-inline-title {
  font-size: 12px;
  font-weight: 700;
  color: #4f46e5;
  margin-bottom: 8px;
}

.team-admin-inline-desc {
  margin: 0 0 10px;
  font-size: 12px;
  line-height: 1.55;
  color: #475569;
}

.team-admin-inline-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.msg-row {
  display: flex;
  margin-bottom: 12px;
}

.msg-row:last-child {
  margin-bottom: 2px;
}

.msg-row.user {
  justify-content: flex-end;
}

.msg-bubble {
  max-width: 85%;
  border-radius: 16px;
  padding: 10px 14px 12px;
  line-height: 1.55;
  white-space: pre-wrap;
  word-break: break-word;
}

.msg-row.user .msg-bubble {
  border-bottom-right-radius: 6px;
  background: linear-gradient(145deg, #3b82f6 0%, #2563eb 55%, #1d4ed8 100%);
  color: #fff;
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.35);
}

.msg-row.ai .msg-bubble {
  border-bottom-left-radius: 6px;
  background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
  color: #0f172a;
  border: 1px solid rgba(226, 232, 240, 0.95);
  box-shadow: 0 2px 10px rgba(15, 23, 42, 0.06);
}

.msg-row.ai .msg-bubble--await {
  border-color: rgba(99, 102, 241, 0.35);
  box-shadow:
    0 2px 12px rgba(99, 102, 241, 0.12),
    0 0 0 1px rgba(99, 102, 241, 0.08) inset;
}

.msg-loading {
  position: relative;
  padding-top: 2px;
  min-height: 52px;
}

.msg-loading-pulse {
  position: absolute;
  inset: -2px -4px -4px -4px;
  border-radius: 14px;
  background: radial-gradient(ellipse 80% 60% at 50% 0%, rgba(99, 102, 241, 0.14), transparent 65%);
  pointer-events: none;
  animation: msg-await-pulse 2s ease-in-out infinite;
}

@keyframes msg-await-pulse {
  0%,
  100% {
    opacity: 0.65;
  }

  50% {
    opacity: 1;
  }
}

.msg-loading-main {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.msg-loading-dots {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 20px;
}

.msg-loading-dot {
  width: 7px;
  height: 7px;
  border-radius: 999px;
  background: linear-gradient(180deg, #6366f1, #3b82f6);
  box-shadow: 0 1px 3px rgba(59, 130, 246, 0.45);
  animation: msg-dot-bounce 1.05s ease-in-out infinite;
}

.msg-loading-dot:nth-child(2) {
  animation-delay: 0.15s;
}

.msg-loading-dot:nth-child(3) {
  animation-delay: 0.3s;
}

@keyframes msg-dot-bounce {
  0%,
  80%,
  100% {
    transform: translateY(0);
    opacity: 0.55;
  }

  40% {
    transform: translateY(-5px);
    opacity: 1;
  }
}

.msg-loading-text {
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  letter-spacing: 0.02em;
}

.msg-loading-skeleton {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.msg-loading-bar {
  display: block;
  height: 8px;
  border-radius: 999px;
  background: linear-gradient(
    90deg,
    rgba(226, 232, 240, 0.85) 0%,
    rgba(241, 245, 249, 0.95) 40%,
    rgba(226, 232, 240, 0.85) 80%
  );
  background-size: 200% 100%;
  animation: msg-shimmer 1.35s ease-in-out infinite;
}

.msg-loading-bar--short {
  width: 62%;
}

@keyframes msg-shimmer {
  0% {
    background-position: 100% 0;
  }

  100% {
    background-position: -100% 0;
  }
}

.msg-role {
  font-size: 11px;
  font-weight: 600;
  text-transform: none;
  letter-spacing: 0.04em;
  margin-bottom: 4px;
}

.msg-row.user .msg-role {
  color: rgba(255, 255, 255, 0.82);
}

.msg-row.ai .msg-role {
  color: #64748b;
}

.msg-content {
  font-size: 14px;
}

.chat-input-bar {
  display: flex;
  align-items: stretch;
  gap: 10px;
  margin-top: 12px;
  padding: 10px 12px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid rgba(226, 232, 240, 0.95);
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.05);
}

.chat-input-field {
  flex: 1;
  min-width: 0;
}

.chat-input-field :deep(.el-input__wrapper) {
  border-radius: 10px;
  box-shadow: none;
  background: rgba(241, 245, 249, 0.9);
  border: 1px solid transparent;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.chat-input-field :deep(.el-input__wrapper:hover) {
  background: #f1f5f9;
}

.chat-input-field :deep(.el-input__wrapper.is-focus) {
  background: #fff;
  border-color: rgba(59, 130, 246, 0.45);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.12);
}

.chat-send-btn {
  flex-shrink: 0;
  border-radius: 10px;
  padding-left: 18px;
  padding-right: 18px;
  font-weight: 600;
  background: linear-gradient(145deg, #3b82f6 0%, #2563eb 100%);
  border: none;
  box-shadow: 0 2px 10px rgba(37, 99, 235, 0.35);
}

.chat-send-btn:hover {
  filter: brightness(1.05);
}

@media (max-width: 768px) {
  .ai-float-btn {
    right: 12px;
    top: auto;
    bottom: 76px;
    transform: none;
    width: 50px;
    height: 50px;
  }

  .ai-float-btn:hover {
    transform: scale(1.06);
  }
}
</style>
