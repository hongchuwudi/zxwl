package param

// AIRequest 定义AI请求结构体
type AIRequest struct {
	Model       string      `json:"model"`
	Messages    []AIMessage `json:"messages"`
	Stream      bool        `json:"stream"`
	Temperature float64     `json:"temperature"`
	MaxTokens   int         `json:"max_tokens"`
}

// AIMessage 定义消息结构体
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
