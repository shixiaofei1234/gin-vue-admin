package request

type AIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIChatRequest struct {
	Message     string          `json:"message"`
	History     []AIChatMessage `json:"history"`
	Temperature float64         `json:"temperature"`
}
