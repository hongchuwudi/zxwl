// Package chatApi api/chat_api.go
package chatApi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/repositories/chatRepo"
	"mymod/service"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	service *service.ChatService
}

// NewChatHandler 使用默认配置创建handler
func NewChatHandler() *ChatHandler {
	repo := chatRepo.NewChatRepository(config.GetDB())
	serv := service.NewChatService(repo)
	return &ChatHandler{
		service: &serv,
	}
}

// SendFriendMessage 发送好友消息
func (h *ChatHandler) SendFriendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	var request sqlModel.SendFriendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"code": 400, "message": "无效的请求数据"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.SendFriendMessage(userID, &request); err != nil {
		http.Error(w, `{"code": 500, "message": "发送消息失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "消息发送成功",
	})
}

// SendRoomMessage 发送聊天室消息
func (h *ChatHandler) SendRoomMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	var request sqlModel.SendRoomMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"code": 400, "message": "无效的请求数据"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.SendRoomMessage(userID, &request); err != nil {
		http.Error(w, `{"code": 500, "message": "发送消息失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "消息发送成功",
	})
}

// GetFriendChatHistory 获取好友聊天历史
func (h *ChatHandler) GetFriendChatHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	friendIDStr := vars["friendID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的好友ID"}`, http.StatusBadRequest)
		return
	}

	// 获取查询参数
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // 默认50条
	if limitStr != "" {
		if limitInt, err := strconv.Atoi(limitStr); err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	messages, err := h.service.GetFriendChatHistory(userID, friendID, limit)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取聊天记录失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    messages,
	})
}

// GetRoomChatHistory 获取聊天室聊天历史
func (h *ChatHandler) GetRoomChatHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	meetingID := vars["meetingID"]

	// 获取查询参数
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // 默认50条
	if limitStr != "" {
		if limitInt, err := strconv.Atoi(limitStr); err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	messages, err := h.service.GetRoomChatHistory(meetingID, limit)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取聊天记录失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    messages,
	})
}

// MarkMessagesAsRead 标记消息为已读
func (h *ChatHandler) MarkMessagesAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	friendIDStr := vars["friendID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的好友ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.MarkFriendMessagesAsRead(userID, friendID); err != nil {
		http.Error(w, `{"code": 500, "message": "标记已读失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "消息已标记为已读",
	})
}

// GetUnreadCount 获取未读消息数量
func (h *ChatHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	count, err := h.service.GetUnreadMessageCount(userID)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取未读消息数量失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    count,
	})
}

// GetRecentChats 获取最近聊天列表
func (h *ChatHandler) GetRecentChats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	// 获取查询参数
	limitStr := r.URL.Query().Get("limit")
	limit := 20 // 默认20条
	if limitStr != "" {
		if limitInt, err := strconv.Atoi(limitStr); err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	chats, err := h.service.GetRecentChats(userID, limit)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取最近聊天列表失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    chats,
	})
}

// DeleteFriendMessage 删除好友消息
func (h *ChatHandler) DeleteFriendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	messageIDStr := vars["messageID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的消息ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteFriendMessage(messageID, userID); err != nil {
		http.Error(w, `{"code": 500, "message": "删除消息失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "消息删除成功",
	})
}

// ClearChatHistory 清空聊天记录
func (h *ChatHandler) ClearChatHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	friendIDStr := vars["friendID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的好友ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.ClearChatHistory(userID, friendID); err != nil {
		http.Error(w, `{"code": 500, "message": "清空聊天记录失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "聊天记录已清空",
	})
}

// GetUnreadCountByFriend 获取指定好友的未读消息数量
func (h *ChatHandler) GetUnreadCountByFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	friendIDStr := vars["friendID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的好友ID"}`, http.StatusBadRequest)
		return
	}

	count, err := h.service.GetUnreadMessageCountByFriend(userID, friendID)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取未读消息数量失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    count,
	})
}
