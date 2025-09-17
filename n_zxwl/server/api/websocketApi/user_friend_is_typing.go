package websocketApi

import (
	"encoding/json"
	"fmt"
	"time"
)

// handleTyping 处理正在输入状态
func (h *WebSocketHandler) handleTyping(userID int64, data json.RawMessage) error {
	var typingMsg struct {
		To     int64 `json:"to"`        // 接收者ID
		Typing bool  `json:"is_typing"` // 是否正在输入
	}

	if err := json.Unmarshal(data, &typingMsg); err != nil {
		return fmt.Errorf("解析输入状态消息失败: %v", err)
	}

	// 转发输入状态给接收者
	return h.forwardTypingStatus(typingMsg.To, userID, typingMsg.Typing)
}

// forwardTypingStatus 转发输入状态给接收者
func (h *WebSocketHandler) forwardTypingStatus(toUserID int64, fromUserID int64, isTyping bool) error {
	// 构建输入状态消息
	typingMsg := map[string]interface{}{
		"type": "typing_status",
		"data": map[string]interface{}{
			"from_user_id": fromUserID,
			"is_typing":    isTyping,
			"timestamp":    time.Now().Unix(),
		},
	}

	// 发送给目标用户
	return h.sendToUser(toUserID, typingMsg)
}
