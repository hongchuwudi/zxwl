package param

// Response 通用响应数据结构
type Response struct {
	Error   int         `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
