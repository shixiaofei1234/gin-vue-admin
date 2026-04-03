package system

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ragChunk struct {
	Source string
	Text   string
	Emb    []float64
}

var (
	ragMu       sync.Mutex
	ragIndex    []ragChunk
	ragLoadErr  error
	ragLoaded   bool
	ragEmbedURL string
)

func chatURLToEmbeddingURL(base string) string {
	// 兼容 chat/completions -> embeddings 的不同部署方式：
	// - 如果用户把 baseURL 指到了 /chat/completions，就替换为 /embeddings
	// - 如果 baseURL 以 /v1 结尾，就补上 /embeddings
	base = strings.TrimSpace(base)
	if base == "" {
		return ""
	}
	if strings.Contains(base, "/chat/completions") {
		return strings.Replace(base, "/chat/completions", "/embeddings", 1)
	}
	if strings.HasSuffix(base, "/v1") {
		return base + "/embeddings"
	}
	return strings.TrimSuffix(base, "/") + "/embeddings"
}

type openAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
}

func cosineSim(a, b []float64) float64 {
	// 余弦相似度，用于把 query 向量与文档块向量做相似度排序。
	// 如果维度不一致或某向量为 0，直接返回 0，保证不产生 NaN。
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func embedTexts(client *http.Client, apiKey, model, url string, texts []string) ([][]float64, error) {
	// 调用 embedding 接口，把多个文本 batch 一次向量化。
	// 返回值顺序按 index 排序，保证与 texts[i] 对齐。
	if len(texts) == 0 {
		return nil, nil
	}
	body := map[string]any{
		"model": model,
		"input": texts,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("embedding http %d: %s", resp.StatusCode, string(raw))
	}
	var out openAIEmbeddingResponse
	if err = json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if len(out.Data) == 0 {
		return nil, errors.New("embedding 返回为空")
	}
	sort.Slice(out.Data, func(i, j int) bool { return out.Data[i].Index < out.Data[j].Index })
	vecs := make([][]float64, len(out.Data))
	for i := range out.Data {
		vecs[i] = out.Data[i].Embedding
	}
	return vecs, nil
}

func splitChunks(text, source string, maxRunes int) []ragChunk {
	// 把 Markdown/文本拆成多段 chunk：
	// - 优先按空行段落分割（\n\n）
	// - 控制每个 chunk 的 rune 数，避免过长影响 embedding 与注入 token
	// - 最终仍可能大于 maxRunes 时，做硬截断
	if maxRunes <= 0 {
		maxRunes = 800
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	parts := strings.Split(text, "\n\n")
	var chunks []ragChunk
	var buf strings.Builder
	flush := func() {
		s := strings.TrimSpace(buf.String())
		if s == "" {
			return
		}
		chunks = append(chunks, ragChunk{Source: source, Text: s})
		buf.Reset()
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(p)
		if utf8.RuneCountInString(buf.String()) >= maxRunes {
			flush()
		}
	}
	flush()
	if len(chunks) == 0 {
		return nil
	}
	var out []ragChunk
	for _, c := range chunks {
		s := c.Text
		for utf8.RuneCountInString(s) > maxRunes {
			runes := []rune(s)
			out = append(out, ragChunk{Source: source, Text: string(runes[:maxRunes])})
			s = string(runes[maxRunes:])
		}
		if strings.TrimSpace(s) != "" {
			out = append(out, ragChunk{Source: source, Text: strings.TrimSpace(s)})
		}
	}
	return out
}

func loadMarkdownKnowledge(dir string, maxRunes int) ([]ragChunk, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("knowledge-dir 不是目录: %s", dir)
	}
	var list []ragChunk
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		chunks := splitChunks(string(raw), filepath.Base(path), maxRunes)
		list = append(list, chunks...)
		return nil
	})
	return list, err
}

func keywordScore(query, chunk string) float64 {
	q := strings.TrimSpace(query)
	if q == "" {
		return 0
	}
	score := 0.0
	if strings.Contains(chunk, q) {
		score += 10
	}
	for _, w := range strings.Fields(q) {
		if len([]rune(w)) >= 2 && strings.Contains(chunk, w) {
			score += 2
		}
	}
	for _, r := range q {
		if strings.ContainsRune(chunk, r) {
			score += 0.01
		}
	}
	return score
}

func ensureRAGIndex() {
	// RAG 索引初始化（懒加载）：
	// - 只加载一次（ragLoaded 控制）
	// - 启用向量检索：对所有 chunk 做 embedding
	// - 失败/未配置：降级为 keywordScore（简单但稳）
	conf := global.GVA_CONFIG.AI
	if !conf.RAGEnabled {
		return
	}
	ragMu.Lock()
	defer ragMu.Unlock()
	if ragLoaded {
		return
	}
	ragLoaded = true
	ragEmbedURL = chatURLToEmbeddingURL(conf.BaseURL)
	dir := strings.TrimSpace(conf.KnowledgeDir)
	if dir == "" {
		dir = "resource/ai_knowledge"
	}
	maxRunes := conf.RAGChunkMaxRunes
	if maxRunes <= 0 {
		maxRunes = 800
	}
	chunks, err := loadMarkdownKnowledge(dir, maxRunes)
	if err != nil {
		ragLoadErr = err
		global.GVA_LOG.Warn("AI RAG: 加载知识库目录失败，将不使用向量检索: " + err.Error())
		return
	}
	if len(chunks) == 0 {
		global.GVA_LOG.Info("AI RAG: 知识库目录无 .md 文件，跳过索引")
		return
	}
	if strings.TrimSpace(conf.APIKey) == "" || ragEmbedURL == "" {
		// 未配置 embedding 所需参数：直接把 chunks 塞进 ragIndex，
		// 相似度计算时会自动走 keywordScore（而非 cosineSim）。
		global.GVA_LOG.Warn("AI RAG: 未配置 api-key 或 base-url，使用关键词检索")
		for i := range chunks {
			ragIndex = append(ragIndex, chunks[i])
		}
		return
	}
	model := strings.TrimSpace(conf.EmbeddingModel)
	if model == "" {
		model = "text-embedding-v3"
	}
	timeout := time.Duration(conf.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	client := newDashScopeClient(timeout)
	texts := make([]string, len(chunks))
	for i := range chunks {
		texts[i] = chunks[i].Text
	}
	const batch = 8
	for i := 0; i < len(texts); i += batch {
		end := i + batch
		if end > len(texts) {
			end = len(texts)
		}
		vecs, err := embedTexts(client, conf.APIKey, model, ragEmbedURL, texts[i:end])
		if err != nil {
			// embedding 失败：为了不让检索功能完全不可用，降级为 keywordScore。
			ragLoadErr = err
			global.GVA_LOG.Warn("AI RAG: 向量索引失败，降级为关键词检索: " + err.Error())
			for j := range chunks {
				ragIndex = append(ragIndex, chunks[j])
			}
			return
		}
		for j := range vecs {
			chunks[i+j].Emb = vecs[j]
			ragIndex = append(ragIndex, chunks[i+j])
		}
	}
	global.GVA_LOG.Info(fmt.Sprintf("AI RAG: 已索引 %d 条文档块", len(ragIndex)))
}

func retrieveRAGContext(userQuery string) string {
	// 把 userQuery 与 ragIndex 中 chunk 做检索：
	// 1) 优先使用 embedding 向量相似度排序
	// 2) 若无向量或向量检索失败：降级 keywordScore
	// 3) 取 topK 条，并按字符数 maxInject 控制注入长度
	// 返回拼接后的上下文片段（会注入到 AI messages 的 system prompt 中）。
	conf := global.GVA_CONFIG.AI
	if !conf.RAGEnabled {
		return ""
	}
	ensureRAGIndex()
	if len(ragIndex) == 0 {
		return ""
	}
	topK := conf.RAGTopK
	if topK <= 0 {
		topK = 4
	}
	maxInject := conf.RAGInjectMaxRunes
	if maxInject <= 0 {
		maxInject = 2500
	}
	q := strings.TrimSpace(userQuery)
	if q == "" {
		return ""
	}

	type scored struct {
		idx   int
		score float64
	}
	var scores []scored

	if len(ragIndex) > 0 && len(ragIndex[0].Emb) > 0 && strings.TrimSpace(conf.APIKey) != "" && ragEmbedURL != "" {
		model := strings.TrimSpace(conf.EmbeddingModel)
		if model == "" {
			model = "text-embedding-v3"
		}
		timeout := time.Duration(conf.TimeoutSeconds) * time.Second
		if timeout <= 0 {
			timeout = 60 * time.Second
		}
		client := newDashScopeClient(timeout)
		vecs, err := embedTexts(client, conf.APIKey, model, ragEmbedURL, []string{q})
		if err == nil && len(vecs) == 1 && len(vecs[0]) > 0 {
			qv := vecs[0]
			for i := range ragIndex {
				if len(ragIndex[i].Emb) == 0 {
					continue
				}
				// cosine 相似度：从高到低排序后取 topK 注入上下文。
				s := cosineSim(qv, ragIndex[i].Emb)
				scores = append(scores, scored{idx: i, score: s})
			}
		}
	}

	if len(scores) == 0 {
		// embeddings 不可用/失败：落回关键词分数。
		for i := range ragIndex {
			scores = append(scores, scored{idx: i, score: keywordScore(q, ragIndex[i].Text)})
		}
	}

	sort.Slice(scores, func(i, j int) bool { return scores[i].score > scores[j].score })
	if len(scores) > topK {
		scores = scores[:topK]
	}
	var b strings.Builder
	runes := 0
	for _, s := range scores {
		// 对于低分块可以在 ragIndex 足够大时跳过，减少无用噪声注入。
		if s.score <= 0 && len(ragIndex) > 3 {
			continue
		}
		c := ragIndex[s.idx]
		block := fmt.Sprintf("【%s】\n%s\n\n", c.Source, c.Text)
		runes += utf8.RuneCountInString(block)
		if runes > maxInject {
			break
		}
		b.WriteString(block)
	}
	return strings.TrimSpace(b.String())
}
