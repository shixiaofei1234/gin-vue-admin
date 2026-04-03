package config

type AI struct {
	APIKey         string `mapstructure:"api-key" json:"apiKey" yaml:"api-key"`
	BaseURL        string `mapstructure:"base-url" json:"baseURL" yaml:"base-url"`
	Model          string `mapstructure:"model" json:"model" yaml:"model"`
	TimeoutSeconds int    `mapstructure:"timeout-seconds" json:"timeoutSeconds" yaml:"timeout-seconds"`
	// RAG：从 knowledge-dir 读取 *.md，分块后调用 DashScope 兼容接口 /v1/embeddings 建索引；对话时检索相关片段注入上下文
	RAGEnabled        bool   `mapstructure:"rag-enabled" json:"ragEnabled" yaml:"rag-enabled"`
	KnowledgeDir      string `mapstructure:"knowledge-dir" json:"knowledgeDir" yaml:"knowledge-dir"`
	EmbeddingModel    string `mapstructure:"embedding-model" json:"embeddingModel" yaml:"embedding-model"`
	RAGTopK           int    `mapstructure:"rag-top-k" json:"ragTopK" yaml:"rag-top-k"`
	RAGChunkMaxRunes  int    `mapstructure:"rag-chunk-max-runes" json:"ragChunkMaxRunes" yaml:"rag-chunk-max-runes"`
	RAGInjectMaxRunes int    `mapstructure:"rag-inject-max-runes" json:"ragInjectMaxRunes" yaml:"rag-inject-max-runes"`
}
