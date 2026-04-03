package system

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"go.uber.org/zap"
)

// 复用连接、关闭压缩，减轻上游/中间层对流式响应的缓冲
var dashScopeTransport = &http.Transport{
	MaxIdleConns:        32,
	MaxIdleConnsPerHost: 8,
	IdleConnTimeout:     90 * time.Second,
	DisableCompression:  true,
}

func newDashScopeClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: dashScopeTransport,
	}
}

func capAITemperature(t float64) float64 {
	// 为了避免温度过大导致回答结构不稳定，这里对 temperature 做了裁剪。
	if t <= 0 {
		return 0.5
	}
	if t > 0.6 {
		return 0.6
	}
	return t
}

// shopTestSystemPrompt 服务端强制注入，保证回答范围仅限 shop-test（忽略客户端伪造的 system）
// 菜单层级与 server/source/system/menu.go 初始化数据一致，便于回答「在哪个菜单」类问题。
const shopTestSystemPrompt = `你是「shop-test」项目的专属助手（基于 gin-vue-admin），知识边界仅限本项目。

必须遵守：
1. 只回答与 shop-test / 本后台相关的内容（业务、功能位置、操作说明等）。
2. 无关问题用 1～2 句礼貌拒绝，引导对方问本项目内问题。
3. 不要编造不存在的菜单名；不要提「构建管理」「插件中心」等与本项目默认菜单不符的泛称。
4. 回答简短（约 200 字内），禁止思考过程、英文推理。

【菜单类问题】用户问「在哪个菜单」「去哪找」时：
- 若上下文中提供了「当前登录角色可见菜单」树，必须**只依据该树**回答路径，不要猜测。
- 回答格式：「一级菜单名 → 子菜单名」。
- 只有在未提供动态菜单树时，才可参考下面「初始化默认」列表。

【本项目默认侧边栏（仅作兜底；动态菜单树优先）】
- 打包插件、插件市场、插件安装、邮件插件 等：均在「插件系统」下；其中「打包插件」即：插件系统 → 打包插件。
- 代码生成器、自动化代码管理、表单生成器、系统配置、模板配置、表格模板 等：在「系统工具」下。
- 角色管理、菜单管理、API 管理、用户管理、字典、操作历史：在「超级管理员」下。
- 「服务器状态」「仪表盘」「关于我们」等为顶级菜单。
- 「示例文件」下含客户列表示例、上传下载等子菜单。`

type AIChatService struct{}

type dashScopeChatRequest struct {
	Model       string                    `json:"model"`
	Messages    []systemReq.AIChatMessage `json:"messages"`
	Temperature float64                   `json:"temperature"`
	MaxTokens   int                       `json:"max_tokens,omitempty"`
	ExtraBody   map[string]any            `json:"extra_body,omitempty"`
	Stream      bool                      `json:"stream"`
}

type dashScopeChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type dashScopeStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var spacesRegexp = regexp.MustCompile(`\s+`)
var cotLineRegexp = regexp.MustCompile(`(?i)(hypothesis|final decision|self-correction|let'?s|wait,|i should|rule\s*[0-9]|internal plan|thought process|reasoning)`)

func sanitizeAIText(content string) string {
	// 处理大模型返回的“思考链/推理模板类”文本行，
	// 保证前端展示更短、更符合产品要求。
	if strings.TrimSpace(content) == "" {
		return ""
	}
	lines := strings.Split(content, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// 过滤典型思考链文本行
		if cotLineRegexp.MatchString(trimmed) {
			continue
		}
		filtered = append(filtered, line)
	}
	result := strings.TrimSpace(strings.Join(filtered, "\n"))
	// 兜底：如果仍出现“Final Decision”等片段，按标记截断到后半段业务内容
	lowered := strings.ToLower(result)
	if idx := strings.Index(lowered, "final decision"); idx >= 0 {
		if hIdx := strings.Index(result[idx:], "###"); hIdx >= 0 {
			result = strings.TrimSpace(result[idx+hIdx:])
		}
	}
	return result
}

func extractString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []any:
		var parts []string
		for _, item := range val {
			if m, ok := item.(map[string]any); ok {
				if text, ok := m["text"].(string); ok && strings.TrimSpace(text) != "" {
					parts = append(parts, text)
				}
			}
		}
		return strings.Join(parts, "")
	default:
		return ""
	}
}

func extractContentFromStreamPayload(payload string) string {
	// 不同模型/网关对流式事件的 JSON 结构可能不同：
	// 这里做兼容解析，尽量抽出 delta/message.content。
	var raw map[string]any
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return ""
	}
	tryChoices := func(root map[string]any) string {
		choices, ok := root["choices"].([]any)
		if !ok || len(choices) == 0 {
			return ""
		}
		first, ok := choices[0].(map[string]any)
		if !ok {
			return ""
		}
		if delta, ok := first["delta"].(map[string]any); ok {
			if c := extractString(delta["content"]); strings.TrimSpace(c) != "" {
				return c
			}
		}
		if msg, ok := first["message"].(map[string]any); ok {
			if c := extractString(msg["content"]); strings.TrimSpace(c) != "" {
				return c
			}
		}
		return ""
	}
	if c := tryChoices(raw); c != "" {
		return c
	}
	if output, ok := raw["output"].(map[string]any); ok {
		if c := tryChoices(output); c != "" {
			return c
		}
	}
	return ""
}

func (s *AIChatService) buildMessages(req systemReq.AIChatRequest, authorityID uint) ([]systemReq.AIChatMessage, error) {
	// 组装发送给大模型的 messages：
	// 1) system：强制注入本项目助手规则（shop-testSystemPrompt）
	// 2) system：动态菜单树（用于回答“去哪找/在哪个菜单”）
	// 3) system：RAG 检索片段（如开启）
	// 4) history：保留 user/assistant 的最近上下文
	// 5) user：当前用户问题
	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		return nil, errors.New("消息不能为空")
	}
	var history []systemReq.AIChatMessage
	for _, m := range req.History {
		role := strings.ToLower(strings.TrimSpace(m.Role))
		if role == "system" {
			continue
		}
		if role != "user" && role != "assistant" {
			continue
		}
		history = append(history, m)
	}
	out := make([]systemReq.AIChatMessage, 0, len(history)+4)
	out = append(out, systemReq.AIChatMessage{
		Role:    "system",
		Content: shopTestSystemPrompt,
	})
	if dyn := buildDynamicMenuContext(authorityID); strings.TrimSpace(dyn) != "" {
		out = append(out, systemReq.AIChatMessage{
			Role:    "system",
			Content: dyn,
		})
	}
	// 业务数据直查上下文：让模型回答前可拿到“当前库里的真实团队/员工信息”。
	// 这样对于“某团队配置/管理员/成员是谁”这类问题，不再只给操作路径。
	if bizCtx := buildBusinessDataContext(msg); strings.TrimSpace(bizCtx) != "" {
		out = append(out, systemReq.AIChatMessage{
			Role:    "system",
			Content: bizCtx,
		})
	}
	if ragCtx := retrieveRAGContext(msg); strings.TrimSpace(ragCtx) != "" {
		out = append(out, systemReq.AIChatMessage{
			Role: "system",
			Content: `【项目知识库 RAG】以下为从 resource/ai_knowledge 下 Markdown 文档检索到的相关片段。请优先据此作答；若与固定规则冲突，以固定规则为准；检索内容不足时可结合你对 shop-test 的常识简要补充。` +
				"\n\n" + ragCtx,
		})
	}
	out = append(out, history...)
	out = append(out, systemReq.AIChatMessage{
		Role:    "user",
		Content: msg,
	})
	return out, nil
}

func shouldQueryBusinessContext(userQuery string) bool {
	q := strings.ToLower(strings.TrimSpace(userQuery))
	if q == "" {
		return false
	}
	return strings.Contains(q, "团队") ||
		strings.Contains(q, "管理员") ||
		strings.Contains(q, "队员") ||
		strings.Contains(q, "成员") ||
		strings.Contains(q, "员工") ||
		strings.Contains(q, "team")
}

func buildBusinessDataContext(userQuery string) string {
	if !shouldQueryBusinessContext(userQuery) {
		return ""
	}

	// 先取最近团队，再按“问题文本是否命中团队名/编号/管理员/成员名”筛选出更相关项。
	var teams []system.SysTeam
	if err := global.GVA_DB.Model(&system.SysTeam{}).Order("id desc").Limit(20).Find(&teams).Error; err != nil {
		global.GVA_LOG.Warn("AI 业务上下文：查询团队失败", zap.Error(err))
		return ""
	}
	if len(teams) == 0 {
		return ""
	}

	type teamHit struct {
		team      system.SysTeam
		employees []system.SysEmployee
		hit       bool
	}
	q := strings.ToLower(strings.TrimSpace(userQuery))
	hits := make([]teamHit, 0, len(teams))
	for i := range teams {
		var employees []system.SysEmployee
		if err := global.GVA_DB.Model(&system.SysEmployee{}).
			Where("team_id = ?", teams[i].ID).
			Order("id desc").
			Limit(30).
			Find(&employees).Error; err != nil {
			global.GVA_LOG.Warn("AI 业务上下文：查询团队成员失败", zap.Error(err))
		}
		matched := false
		tname := strings.ToLower(strings.TrimSpace(teams[i].TeamName))
		tnum := strings.ToLower(strings.TrimSpace(teams[i].TeamNum))
		admin := strings.ToLower(strings.TrimSpace(teams[i].AdminName))
		if (tname != "" && strings.Contains(q, tname)) ||
			(tnum != "" && strings.Contains(q, tnum)) ||
			(admin != "" && strings.Contains(q, admin)) {
			matched = true
		}
		if !matched {
			for _, e := range employees {
				en := strings.ToLower(strings.TrimSpace(e.EmployeeName))
				if en != "" && strings.Contains(q, en) {
					matched = true
					break
				}
			}
		}
		hits = append(hits, teamHit{
			team:      teams[i],
			employees: employees,
			hit:       matched,
		})
	}

	// 优先命中项；如果一个都没命中，取最近 5 个团队做概览兜底。
	selected := make([]teamHit, 0, 6)
	for _, h := range hits {
		if h.hit {
			selected = append(selected, h)
		}
		if len(selected) >= 6 {
			break
		}
	}
	if len(selected) == 0 {
		for i := 0; i < len(hits) && i < 5; i++ {
			selected = append(selected, hits[i])
		}
	}
	if len(selected) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("【实时业务数据（团队/员工）】以下信息来自当前数据库，请优先依据这些真实数据回答；若用户问某团队/管理员/成员，请直接给结论，不必只讲操作路径。\n")
	for _, s := range selected {
		t := s.team
		status := "停用"
		if t.Status {
			status = "启用"
		}
		admin := strings.TrimSpace(t.AdminName)
		if admin == "" {
			admin = "未设置"
		}
		b.WriteString(fmt.Sprintf("- 团队[%s | %s]：名称=%s，状态=%s，管理员=%s，成员数=%d\n",
			strings.TrimSpace(t.TeamNum),
			fmt.Sprintf("ID:%d", t.ID),
			strings.TrimSpace(t.TeamName),
			status,
			admin,
			len(s.employees),
		))
		if len(s.employees) > 0 {
			b.WriteString("  成员：")
			for i, e := range s.employees {
				if i >= 12 {
					b.WriteString(" ...")
					break
				}
				if i > 0 {
					b.WriteString("、")
				}
				b.WriteString(strings.TrimSpace(e.EmployeeName))
			}
			b.WriteString("\n")
		}
	}
	// 控制注入长度，避免把太多数据库内容塞进 prompt。
	ctx := strings.TrimSpace(b.String())
	const maxRunes = 1800
	r := []rune(ctx)
	if len(r) > maxRunes {
		ctx = string(r[:maxRunes])
	}
	return ctx
}

func (s *AIChatService) Chat(req systemReq.AIChatRequest, authorityID uint) (string, error) {
	// 对“团队管理员是谁/叫什么”这类强结构化问题，
	// 直接从数据库取 admin_name，避免模型漏字段导致回答不完整。
	if reply, ok := s.tryDirectTeamAdminReply(req.Message); ok {
		return reply, nil
	}

	conf := global.GVA_CONFIG.AI
	if strings.TrimSpace(conf.APIKey) == "" {
		return "", errors.New("请先在 config.yaml 配置 ai.api-key")
	}
	if strings.TrimSpace(conf.BaseURL) == "" {
		return "", errors.New("请先在 config.yaml 配置 ai.base-url")
	}
	if strings.TrimSpace(conf.Model) == "" {
		return "", errors.New("请先在 config.yaml 配置 ai.model")
	}

	messages, err := s.buildMessages(req, authorityID)
	if err != nil {
		return "", err
	}

	temperature := capAITemperature(req.Temperature)

	requestBody := dashScopeChatRequest{
		Model:       conf.Model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   128,
		ExtraBody: map[string]any{
			"enable_thinking": false,
		},
		Stream: false,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	timeout := time.Duration(conf.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	client := newDashScopeClient(timeout)

	httpReq, err := http.NewRequest(http.MethodPost, conf.BaseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+conf.APIKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return "", fmt.Errorf("大模型请求超时（%ds），请稍后重试或在 config.yaml 中调大 ai.timeout-seconds", conf.TimeoutSeconds)
		}
		if strings.Contains(strings.ToLower(err.Error()), "context deadline exceeded") {
			return "", fmt.Errorf("大模型请求超时（%ds），请稍后重试或在 config.yaml 中调大 ai.timeout-seconds", conf.TimeoutSeconds)
		}
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("大模型调用失败: " + string(respBytes))
	}

	var result dashScopeChatResponse
	if err = json.Unmarshal(respBytes, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", errors.New("大模型返回为空")
	}

	return sanitizeAIText(result.Choices[0].Message.Content), nil
}

func (s *AIChatService) ChatStream(req systemReq.AIChatRequest, authorityID uint, onChunk func(string) error) error {
	// 同 Chat：流式接口也对强结构化问题做直查直回。
	if reply, ok := s.tryDirectTeamAdminReply(req.Message); ok {
		// 前端把每个 onChunk 当作正文分片展示，这里只推送一段即可。
		return onChunk(reply + "\n")
	}

	conf := global.GVA_CONFIG.AI
	if strings.TrimSpace(conf.APIKey) == "" {
		return errors.New("请先在 config.yaml 配置 ai.api-key")
	}
	if strings.TrimSpace(conf.BaseURL) == "" {
		return errors.New("请先在 config.yaml 配置 ai.base-url")
	}
	if strings.TrimSpace(conf.Model) == "" {
		return errors.New("请先在 config.yaml 配置 ai.model")
	}

	messages, err := s.buildMessages(req, authorityID)
	if err != nil {
		return err
	}

	temperature := capAITemperature(req.Temperature)

	requestBody := dashScopeChatRequest{
		Model:       conf.Model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   128,
		ExtraBody: map[string]any{
			"enable_thinking": false,
		},
		Stream: true,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	timeout := time.Duration(conf.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 90 * time.Second
	}
	client := newDashScopeClient(timeout)

	httpReq, err := http.NewRequest(http.MethodPost, conf.BaseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Authorization", "Bearer "+conf.APIKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return fmt.Errorf("大模型请求超时（%ds），请稍后重试或在 config.yaml 中调大 ai.timeout-seconds", conf.TimeoutSeconds)
		}
		if strings.Contains(strings.ToLower(err.Error()), "context deadline exceeded") {
			return fmt.Errorf("大模型请求超时（%ds），请稍后重试或在 config.yaml 中调大 ai.timeout-seconds", conf.TimeoutSeconds)
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return errors.New("大模型调用失败")
		}
		return errors.New("大模型调用失败: " + string(respBytes))
	}

	reader := bufio.NewReader(resp.Body)
	// pending 是一个“尚未确认属于最终渲染”的缓冲区：
	// 用于在收到流式碎片时，把疑似被拆分的内容尽量拼完整后再 onChunk 给前端。
	var pending string
	flushPending := func(force bool) error {
		if pending == "" {
			return nil
		}
		parts := strings.Split(pending, "\n")
		if !force && len(parts) <= 1 {
			return nil
		}
		limit := len(parts)
		if !force {
			limit = len(parts) - 1
		}
		var out []string
		for i := 0; i < limit; i++ {
			line := parts[i]
			if cotLineRegexp.MatchString(strings.TrimSpace(line)) {
				continue
			}
			out = append(out, line)
		}
		if force {
			pending = ""
		} else {
			pending = parts[len(parts)-1]
		}
		// 拼出最终要渲染的 chunk：如果空则不输出。
		chunk := strings.TrimSpace(strings.Join(out, "\n"))
		if chunk != "" {
			return onChunk(chunk + "\n")
		}
		return nil
	}
	for {
		line, readErr := reader.ReadString('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return readErr
		}
		// dashScope 的流式事件通常是 line-by-line 的 SSE 风格：
		// 每行以 data: 开头（json），并在最后用 [DONE] 结束。
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data:") {
			payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if payload == "[DONE]" {
				if err = flushPending(true); err != nil {
					return err
				}
				return nil
			}
			if payload != "" {
				content := extractContentFromStreamPayload(payload)
				content = spacesRegexp.ReplaceAllString(content, " ")
				if strings.TrimSpace(content) != "" {
					pending += content
					if cbErr := flushPending(false); cbErr != nil {
						return cbErr
					}
				}
			}
		}
		if errors.Is(readErr, io.EOF) {
			if err = flushPending(true); err != nil {
				return err
			}
			return nil
		}
	}
}

// tryDirectTeamAdminReply：
// 专门处理“团队管理员是谁/叫什么”的强结构化问题：
// - 命中则直接查询 sys_team.admin_name 返回（不走大模型）。
// - 避免模型可能漏字段导致“查不出管理员名称”的体验问题。
func (s *AIChatService) tryDirectTeamAdminReply(userQuery string) (string, bool) {
	q := strings.ToLower(strings.TrimSpace(userQuery))
	if q == "" {
		return "", false
	}
	// DEBUG: 用于确认直查分支是否触发（需要结合终端日志观察）
	global.GVA_LOG.Info("AI tryDirectTeamAdminReply", zap.String("rawQuery", userQuery), zap.String("q", q))
	// 必须包含“管理员/负责人/管理者”等强结构化关键词，并且尽可能要包含“团队/队伍/组/team”
	// 以定位到具体团队，否则很容易误命中。
	isAdminIntent := strings.Contains(q, "管理员") ||
		strings.Contains(q, "负责人") ||
		strings.Contains(q, "管理者") ||
		strings.Contains(q, "admin")
	if !isAdminIntent {
		return "", false
	}
	if !(strings.Contains(q, "团队") || strings.Contains(q, "队伍") || strings.Contains(q, "组") || strings.Contains(q, "team")) {
		return "", false
	}

	// 1) 优先从用户问题中“抽取团队名/团队编号”，做精确查询。
	// 典型问题：`马思平团队的管理员是谁`
	// 典型问题：`TD0002团队的管理员是谁`
	var extractedTeamNum string
	{
		// SysTeam.TeamNum 形如 TD0001
		reTeamNum := regexp.MustCompile(`(?i)\btd\d+\b`)
		if m := reTeamNum.FindStringSubmatch(q); len(m) >= 1 && strings.TrimSpace(m[0]) != "" {
			extractedTeamNum = strings.ToUpper(strings.TrimSpace(m[0]))
		}
	}
	var extractedTeamName string
	{
		// 先用正则“尽量”提取（保留兜底逻辑：如果正则失败，就走 deterministic 截取）。
		// 注意：Go regexp(RE2) 不支持懒惰量词的行为可能导致提取不稳，所以这里只做尝试。
		reTeamName := regexp.MustCompile(`(.+)(团队|队伍|组).*?(管理员|负责人|管理者|admin)`)
		if m := reTeamName.FindStringSubmatch(q); len(m) >= 2 && strings.TrimSpace(m[1]) != "" {
			extractedTeamName = strings.TrimSpace(m[1])
		}

		// 再用更确定的方式兜底：
		// 对“马里奥团队的管理员是谁”，直接取“团队”关键字之前的前缀。
		if strings.TrimSpace(extractedTeamName) == "" {
			kwPos := -1
			for _, cand := range []string{"团队", "队伍", "组"} {
				if p := strings.Index(q, cand); p >= 0 && (kwPos == -1 || p < kwPos) {
					kwPos = p
				}
			}
			if kwPos > 0 {
				extractedTeamName = strings.TrimSpace(q[:kwPos])
			}
		}

		// 清理尾缀/连接词：比如 “马里奥团队的” -> “马里奥”
		extractedTeamName = strings.TrimSpace(extractedTeamName)
		extractedTeamName = strings.TrimSuffix(extractedTeamName, "团队")
		extractedTeamName = strings.TrimSuffix(extractedTeamName, "队伍")
		extractedTeamName = strings.TrimSuffix(extractedTeamName, "组")
		extractedTeamName = strings.TrimSuffix(extractedTeamName, "的")
		extractedTeamName = strings.TrimSpace(extractedTeamName)
	}

	// 命中团队编号：直接按 team_num 精确查
	if extractedTeamNum != "" {
		var team system.SysTeam
		if err := global.GVA_DB.Where("team_num = ?", extractedTeamNum).First(&team).Error; err == nil {
			admin := strings.TrimSpace(team.AdminName)
			if admin == "" {
				admin = "未设置"
			}
			teamName := strings.TrimSpace(team.TeamName)
			if teamName == "" {
				teamName = "未命名团队"
			}
			return fmt.Sprintf("当前数据库中，团队「%s」（%s）的管理员是：%s。", teamName, strings.TrimSpace(team.TeamNum), admin), true
		}
		// team_num 明确给了但查不到：直接返回明确提示，避免让模型接管输出“无法访问业务数据”。
		return fmt.Sprintf("未找到团队编号「%s」，请确认编号是否正确后再问。", extractedTeamNum), true
	}

	// 命中团队名：按包含关系查询（更宽容匹配）
	if extractedTeamName != "" {
		var team system.SysTeam
		if err := global.GVA_DB.Where("team_name LIKE ?", "%"+extractedTeamName+"%").Order("id desc").Limit(1).First(&team).Error; err == nil {
			admin := strings.TrimSpace(team.AdminName)
			if admin == "" {
				admin = "未设置"
			}
			teamNum := strings.TrimSpace(team.TeamNum)
			if teamNum == "" {
				return fmt.Sprintf("当前数据库中，团队「%s」的管理员是：%s。", strings.TrimSpace(team.TeamName), admin), true
			}
			return fmt.Sprintf("当前数据库中，团队「%s」（%s）的管理员是：%s。", strings.TrimSpace(team.TeamName), teamNum, admin), true
		}
		// team_name 抽取到了但查不到：直接返回明确提示，避免模型拒答。
		return fmt.Sprintf("未找到团队「%s」，请确认团队名称是否完全一致，或提供团队编号（如 TD0001）。", extractedTeamName), true
	}

	// 识别到了“团队 + 管理员”，但提取不出团队名/编号：
	// 给一个可执行的明确引导，而不是让模型拒答“无法访问业务数据”。
	return "请提供要查询的团队名称或团队编号（例如：‘马里奥团队的管理员是谁’或‘TD0001团队的管理员是谁’）。", true
}

var (
	apiFlowIntentHasAPI     = regexp.MustCompile(`(?i)api|接口`)
	apiFlowIntentKeywords   = regexp.MustCompile(`创建|新增|添加|建立|注册|登记|写.*路由|怎么|如何|怎样|步骤|流程|整一个|做一个`)
	apiFlowIntentColloquial = regexp.MustCompile(`想|要|问|如何|怎么|怎样|操作|教我|弄|搞|整一个`)
)

// DetectApiFlowIntent 与前端 shouldEmbedApiFlow 尽量保持一致：
// 命中时才会在流式结束后再生成 apiFlow 卡片（标题+步骤）。
func DetectApiFlowIntent(msg string) bool {
	t := strings.TrimSpace(msg)
	if len([]rune(t)) < 4 {
		return false
	}
	if !apiFlowIntentHasAPI.MatchString(t) {
		return false
	}
	return apiFlowIntentKeywords.MatchString(t) || apiFlowIntentColloquial.MatchString(t)
}

func truncateRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}

func stripMarkdownJSONFence(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return s
	}
	end := len(lines) - 1
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) == "```" {
			end = i
			break
		}
	}
	start := 1
	if end <= start {
		return s
	}
	return strings.TrimSpace(strings.Join(lines[start:end], "\n"))
}

func defaultApiFlowEmbedMap() map[string]interface{} {
	// 默认兜底：当模型二次生成失败/解析失败/steps 不足时使用。
	return map[string]interface{}{
		"kind":  "apiFlow",
		"title": "新增 API 流程",
		"steps": []string{
			"后端：写路由、Handler，在 router 里注册",
			"后台：超级管理员 → API 管理，登记路径与方法",
			"权限：角色管理 → 为角色勾选该 API",
		},
	}
}

// chatCompletionSingle 单次非流式对话（不注入 RAG/菜单）：
// 只让模型输出一个“严格 JSON”，用于生成结构化卡片内容。
func (s *AIChatService) chatCompletionSingle(systemPrompt, userPrompt string, maxTokens int) (string, error) {
	conf := global.GVA_CONFIG.AI
	if strings.TrimSpace(conf.APIKey) == "" {
		return "", errors.New("请先在 config.yaml 配置 ai.api-key")
	}
	if strings.TrimSpace(conf.BaseURL) == "" {
		return "", errors.New("请先在 config.yaml 配置 ai.base-url")
	}
	if strings.TrimSpace(conf.Model) == "" {
		return "", errors.New("请先在 config.yaml 配置 ai.model")
	}

	messages := []systemReq.AIChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
	requestBody := dashScopeChatRequest{
		Model:       conf.Model,
		Messages:    messages,
		Temperature: 0.35,
		MaxTokens:   maxTokens,
		ExtraBody: map[string]any{
			"enable_thinking": false,
		},
		Stream: false,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	timeout := time.Duration(conf.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	client := newDashScopeClient(timeout)

	httpReq, err := http.NewRequest(http.MethodPost, conf.BaseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+conf.APIKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("大模型调用失败: " + string(respBytes))
	}

	var result dashScopeChatResponse
	if err = json.Unmarshal(respBytes, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", errors.New("大模型返回为空")
	}
	return sanitizeAIText(result.Choices[0].Message.Content), nil
}

// GenerateApiFlowEmbedCardIfNeeded：
// 命中意图后，在“主回答流式结束”后再做一次短调用，
// 让模型生成标题与 steps 的结构化 JSON，返回给前端做卡片渲染。
func (s *AIChatService) GenerateApiFlowEmbedCardIfNeeded(userQuestion, assistantReply string) map[string]interface{} {
	if !DetectApiFlowIntent(userQuestion) {
		return nil
	}
	fallback := defaultApiFlowEmbedMap()
	reply := truncateRunes(strings.TrimSpace(assistantReply), 900)
	systemPrompt := `你是 shop-test（gin-vue-admin）后台「帮助卡片」生成器。
只输出一个 JSON 对象，不要 markdown、不要代码块、不要解释。
格式严格为：
{"title":"string","steps":["string","string",...]}
要求：
- title 为 8～20 字内的中文标题，概括「新增业务 API 并登记、授权」。
- steps 为 3～5 条中文短句，每条一行说明一步，顺序为：后端路由 → 后台 API 登记 → 角色授权（或等价合理顺序）。
- 结合用户问题与助手已回答的摘要，不要重复长段落，不要编造不存在的菜单名；超级管理员、API 管理、角色管理 等须与 gin-vue-admin 一致。`

	userPrompt := fmt.Sprintf("用户问题：\n%s\n\n助手已回答（摘要）：\n%s\n\n请输出 JSON。", strings.TrimSpace(userQuestion), reply)

	raw, err := s.chatCompletionSingle(systemPrompt, userPrompt, 384)
	if err != nil {
		global.GVA_LOG.Warn("API 流程卡片生成失败，使用默认", zap.Error(err))
		return fallback
	}
	raw = stripMarkdownJSONFence(sanitizeAIText(raw))
	var parsed struct {
		Title string   `json:"title"`
		Steps []string `json:"steps"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		global.GVA_LOG.Warn("API 流程卡片 JSON 解析失败，使用默认", zap.Error(err))
		return fallback
	}
	if len(parsed.Steps) < 2 {
		return fallback
	}
	steps := make([]string, 0, len(parsed.Steps))
	for _, st := range parsed.Steps {
		t := strings.TrimSpace(st)
		if t != "" {
			steps = append(steps, t)
		}
	}
	if len(steps) < 2 {
		return fallback
	}
	title := strings.TrimSpace(parsed.Title)
	if title == "" {
		title = "新增 API 流程"
	}
	return map[string]interface{}{
		"kind":  "apiFlow",
		"title": title,
		"steps": steps,
	}
}

// GenerateTeamStatusToggleEmbedCardIfNeeded：
// 解析“把某团队关闭/开启”意图，直接从数据库取 team.ID 与当前/目标状态，
// 返回给前端一个可点击的队伍状态切换按钮配置。
//
// 说明：
// - 这里不走二次大模型调用，而是做确定性逻辑，保证操作按钮可靠可用。
// - 返回的 embed.kind 固定为 'teamStatusToggle'，前端据此渲染按钮。
func (s *AIChatService) GenerateTeamStatusToggleEmbedCardIfNeeded(userQuestion string) map[string]interface{} {
	q := strings.ToLower(strings.TrimSpace(userQuestion))
	if q == "" {
		return nil
	}

	// 1) 判断目标：关闭/启用
	wantClose := regexp.MustCompile(`关闭|停用|禁用|关掉|下线`).MatchString(q)
	wantOpen := regexp.MustCompile(`开启|启用|恢复|打开`).MatchString(q)
	// 支持“切换/反转/来回”，即自动取反当前状态
	wantToggle := regexp.MustCompile(`切换|反转|倒过来|来回|互换|开关`).MatchString(q)
	if !wantClose && !wantOpen && !wantToggle {
		return nil
	}
	if wantClose && wantOpen {
		// 歧义时不做确定性操作，避免误切换
		return nil
	}
	// targetStatus 在查询出当前状态后再计算（toggle 情况需要取反）

	// 2) 抽取团队名：取“团队/队伍/组”前面的片段
	var teamName string
	{
		re := regexp.MustCompile(`(.+?)(团队|队伍|组)`)
		if m := re.FindStringSubmatch(userQuestion); len(m) >= 2 {
			teamName = strings.TrimSpace(m[1])
		}
		// 清理常见前缀，如“帮我把/请把/把”
		for _, p := range []string{
			"帮我把", "帮忙把", "请把", "把", "将", "我要把",
			"切换一下", "切换", "反转一下", "反转", "倒过来", "来回", "互换", "开关", "开关一下",
		} {
			if strings.HasPrefix(teamName, p) {
				teamName = strings.TrimSpace(strings.TrimPrefix(teamName, p))
			}
		}
		teamName = strings.TrimSuffix(teamName, "的")
		teamName = strings.TrimSpace(teamName)
	}
	if teamName == "" {
		return nil
	}

	// 3) 查询团队：优先按名称模糊匹配（teamName 可能有“马里奥”这种短文本）
	var team system.SysTeam
	if err := global.GVA_DB.Where("team_name LIKE ?", "%"+teamName+"%").Order("id desc").First(&team).Error; err != nil {
		return nil
	}
	if team.ID == 0 {
		return nil
	}

	cur := team.Status
	var targetStatus bool
	if wantOpen {
		targetStatus = true
	} else if wantClose {
		targetStatus = false
	} else if wantToggle {
		targetStatus = !cur
	} else {
		return nil
	}
	return map[string]interface{}{
		"kind":          "teamStatusToggle",
		"teamID":        team.ID,
		"teamName":      team.TeamName,
		"currentStatus": cur,
		"targetStatus":  targetStatus,
	}
}
