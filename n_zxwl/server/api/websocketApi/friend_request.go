package websocketApi

import (
	"encoding/json"
	"fmt"
	"log"
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/repositories/userFriendsRepo"
	"time"
)

// handleRequest 处理用户好友请求
func (h *WebSocketHandler) handleRequest(userID int64, data json.RawMessage) error {
	var friendRequest sqlModel.FriendRequest
	if err := json.Unmarshal(data, &friendRequest); err != nil {
		return fmt.Errorf("解析输入状态消息失败: %v", err)
	}

	friRepo := userFriendsRepo.FriendRepository{config.GetDB()}
	firRequest, _ := friRepo.GetFriendRequestByID(friendRequest.ID)
	friendRequest = *firRequest

	// 转发输入状态给接收者
	return h.forwardRequestStatus(friendRequest, userID)
}

func (h *WebSocketHandler) forwardRequestStatus(friendRequest sqlModel.FriendRequest, userID int64) error {
	h.mu.RLock()
	conn, exists := h.clientConnections[friendRequest.ToUserID]
	h.mu.RUnlock()

	if !exists {
		log.Printf("用户 %d 不在线，好友请求-消息无法实时送达", friendRequest.ToUserID)
		return nil
	}

	// 确保用户信息已加载（防止空指针或零值）
	if friendRequest.FromUser.ID == 0 {
		log.Printf("警告：好友请求 %d 的用户信息未预加载", friendRequest.ID)
		// 这里可以添加从数据库加载用户信息的逻辑
	}

	// 构建回执消息格式
	forwardMsg := map[string]interface{}{
		"type": "new_request",
		"data": map[string]interface{}{
			"request_id":      friendRequest.ID,
			"status":          friendRequest.Status,
			"salutation":      friendRequest.Salutation,
			"request_message": friendRequest.RequestMessage,
			"created_at":      friendRequest.CreatedAt.Format(time.RFC3339),
			"from_user": map[string]interface{}{
				"id":          friendRequest.FromUser.ID,
				"username":    friendRequest.FromUser.Username,
				"email":       friendRequest.FromUser.Email,
				"displayName": friendRequest.FromUser.DisplayName,
				"avatarUrl":   friendRequest.FromUser.AvatarURL,
				"gender":      friendRequest.FromUser.Gender,
			},
		},
	}

	if err := conn.WriteJSON(forwardMsg); err != nil {
		log.Printf("向用户 %d 发送好友请求通知失败: %v", friendRequest.ToUserID, err)
		return err
	}

	log.Printf("好友请求通知已发送给用户 %d (请求ID: %d)", friendRequest.ToUserID, friendRequest.ID)
	return nil
}
