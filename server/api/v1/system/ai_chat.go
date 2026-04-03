package system

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AIChatApi struct{}

// Chat 普通非流式：返回完整 content 字符串给前端。
func (a *AIChatApi) Chat(c *gin.Context) {
	var req systemReq.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	authorityID := utils.GetUserAuthorityId(c)
	content, err := aiChatService.Chat(req, authorityID)
	if err != nil {
		global.GVA_LOG.Error("AI 对话失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"content": content,
	}, "请求成功", c)
}

// ChatStream 流式：通过 SSE(事件流) 持续推送：
// - {"content": "..."}：正文分片
// - {"embed": {...}}：卡片嵌入信息（当前用于 apiFlow）
// - {"done": true}：流结束
//
// 关键点：我们把 embed 合并进同一条 done 帧（避免某些情况下最后一帧只到 done 但 embed 丢失）。
func (a *AIChatApi) ChatStream(c *gin.Context) {
	var req systemReq.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	if _, err := c.Writer.Write([]byte(": connected\n\n")); err == nil {
		if flusher, ok := c.Writer.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	writeEvent := func(data gin.H) error {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if _, err = c.Writer.Write([]byte("data: " + string(b) + "\n\n")); err != nil {
			return err
		}
		if flusher, ok := c.Writer.(http.Flusher); ok {
			flusher.Flush()
		}
		return nil
	}

	authorityID := utils.GetUserAuthorityId(c)
	var full strings.Builder
	// full 用来拼接最终完整回复（不是为了展示，而是为了在流结束后再生成卡片 JSON）。
	// 注意：由于 ChatStream 本身是流式回调，这里需要把每个 chunk 内容累积起来。
	if err := aiChatService.ChatStream(req, authorityID, func(chunk string) error {
		full.WriteString(chunk)
		return writeEvent(gin.H{"content": chunk})
	}); err != nil {
		global.GVA_LOG.Error("AI 流式对话失败", zap.Error(err))
		_ = writeEvent(gin.H{"error": err.Error()})
		return
	}

	// 与 embed 合并进同一条 SSE，避免与 done 分开发送时前端只收到 done、丢失 embed
	// 先尝试团队状态切换 embed（确定性操作按钮）
	embed := aiChatService.GenerateTeamStatusToggleEmbedCardIfNeeded(req.Message)
	if embed == nil {
		embed = aiChatService.GenerateApiFlowEmbedCardIfNeeded(req.Message, full.String())
	}
	donePayload := gin.H{"done": true}
	if embed != nil {
		donePayload["embed"] = embed
	}
	// 只写一条 done（可能携带 embed），前端解析逻辑更简单且更不容易遗漏最后一帧。
	_ = writeEvent(donePayload)
}
