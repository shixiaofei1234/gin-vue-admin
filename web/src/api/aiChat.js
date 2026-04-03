import service from '@/utils/request'
import { getBaseUrl } from '@/utils/format'
import { useUserStore } from '@/pinia/modules/user'

export const chatWithAI = (data) => {
  return service({
    url: '/ai/chat',
    method: 'post',
    data
  })
}

export const chatWithAIStream = async(data, { onChunk, onDone, onError, onMode, onMetrics, onEmbed } = {}) => {
  const userStore = useUserStore()
  const baseUrl = getBaseUrl()
  const response = await fetch(`${baseUrl}/ai/chat/stream`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'x-token': userStore.token || '',
      'x-user-id': userStore.userInfo?.ID || ''
    },
    body: JSON.stringify(data)
  })

  if (response.status === 404) {
    onMode && onMode('fallback')
    const res = await chatWithAI(data)
    const content = res?.data?.content || ''
    if (content) {
      onChunk && onChunk(content)
    }
    onDone && onDone()
    return
  }

  if (!response.ok || !response.body) {
    const text = await response.text()
    throw new Error(text || '流式请求失败')
  }
  onMode && onMode('connecting')

  // SSE 事件流由浏览器把响应体暴露为 ReadableStream，
  // 需要我们自己按 chunk 解码、再按 SSE 的分隔符拆事件。
  const reader = response.body.getReader()
  const decoder = new TextDecoder('utf-8')
  let buffer = ''
  let transportChunkCount = 0
  let contentChunkCount = 0
  let firstChunkAt = 0
  const startedAt = Date.now()

  const processSseEventBlock = (event) => {
    // 这里假设每个 SSE 事件块长这样：
    // data: {...}
    // data: {...}
    // （以空行分隔）
    // 我们把所有 data: 行拼成一个字符串再 JSON.parse。
    const lines = event.split('\n')
    const dataLines = lines
      .map((line) => line.trim())
      .filter((line) => line.startsWith('data:'))
      .map((line) => line.replace(/^data:\s*/, ''))
    if (!dataLines.length) return
    const raw = dataLines.join('')
    try {
      const json = JSON.parse(raw)
      if (json.error) {
        throw new Error(json.error)
      }
      // 必须先处理 embed，再处理 done：
      // 某些情况下后端把 embed 合并到最后的 done 帧里，
      // 如果在 done 分支 return 之前没有先处理 embed，就会“只看到 done、看不到卡片”。
      if (json.embed) {
        onEmbed && onEmbed(json.embed)
      }
      if (json.done) {
        return
      }
      if (json.content) {
        // content 是正文分片：统计 metrics 用、也驱动 onChunk 逐步把文本写入 UI
        contentChunkCount += 1
        if (!firstChunkAt) {
          firstChunkAt = Date.now()
          onMode && onMode('stream')
        }
        onChunk && onChunk(json.content)
      }
    } catch (e) {
      onError && onError(e)
      throw e
    }
  }

  while (true) {
    const { done, value } = await reader.read()
    // 必须在 done 前解码 value：
    // 有些代理/浏览器会把最后一次 read() 里塞进 done+embed/data，
    // 如果直接 if(done) break，会导致最后 data 块没被加入 buffer。
    if (value) {
      transportChunkCount += 1
      buffer += decoder.decode(value, { stream: true })
    }
    const events = buffer.split(/\r?\n\r?\n/)
    buffer = events.pop() || ''
    for (const event of events) {
      processSseEventBlock(event)
    }
    if (done) {
      break
    }
  }
  // 流结束时 buffer 里可能还有未以 \n\n 结尾的整段 data
  if (buffer.trim()) {
    try {
      processSseEventBlock(buffer)
    } catch (_) {
      /* 末尾半包 JSON 忽略 */
    }
  }
  if (contentChunkCount === 0) {
    onMode && onMode('buffered')
  }
  onMetrics && onMetrics({
    // transportChunkCount：网络层/代理层拆分次数
    // contentChunkCount：业务层（实际收到正文）次数
    // firstChunkMs：从请求开始到首个 content 到达的耗时
    transportChunkCount,
    contentChunkCount,
    firstChunkMs: firstChunkAt ? (firstChunkAt - startedAt) : 0
  })
  onDone && onDone()
}
