package websocketApi

import (
	"encoding/json"
	"fmt"
	"log"
	"mymod/model/sqlModel"
	"time"
)

// handleChatMessage 处理聊天消息
func (h *WebSocketHandler) handleChatMessage(userID int64, data json.RawMessage) error {
	var chatMsg sqlModel.SendFriendMessageRequest
	if err := json.Unmarshal(data, &chatMsg); err != nil {
		return fmt.Errorf("解析聊天消息失败: %v", err)
	}
	// 转发消息给接收者
	return h.forwardMessageToUser(chatMsg.ReceiverID, userID, chatMsg.Content, chatMsg.MessageType)
}

// forwardMessageToUser 转发消息给指定用户
func (h *WebSocketHandler) forwardMessageToUser(toUserID int64, fromUserID int64, content, msgType string) error {
	h.mu.RLock()
	conn, exists := h.clientConnections[toUserID]
	h.mu.RUnlock()

	if !exists {
		log.Printf("用户 %d 不在线，消息无法实时送达", toUserID)
		return nil
	}

	// 构建与前端的handleNewMessage期望格式匹配的消息
	forwardMsg := map[string]interface{}{
		"type": "new_message", // 事件类型
		"data": map[string]interface{}{
			"from_user_id": fromUserID,
			"content":      content,
			"message_type": msgType,
			"created_at":   time.Now().Format(time.RFC3339),
		},
	}

	return conn.WriteJSON(forwardMsg)
}
