package websocketApi

import (
	"encoding/json"
	"fmt"
	"log"
)

var msg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// processMessage 处理不同类型的ws消息
func (h *WebSocketHandler) processMessage(userID int64, message []byte) error {
	log.Printf("用户 %d 发送的原始消息: %s", userID, string(message))
	if err := json.Unmarshal(message, &msg); err != nil {
		return fmt.Errorf("解析消息失败: %v", err)
	}
	switch msg.Type {
	case "ping":
		return h.handlePing(userID)
	case "chat_message":
		return h.handleChatMessage(userID, msg.Data)
	case "typing":
		return h.handleTyping(userID, msg.Data)
	case "read_receipt":
		return h.handleReadReceipt(userID, msg.Data)
	case "friend_request":
		return h.handleRequest(userID, msg.Data)
	default:
		log.Printf("用户 %d 发送了未知消息类型: %s", userID, msg.Type)
		return nil
	}
}
