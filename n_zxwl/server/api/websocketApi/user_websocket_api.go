// Package websocketApi api/websocket_api.go
package websocketApi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"mymod/ws"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"mymod/config"
	"mymod/service"
)

type WebSocketHandler struct {
	service           *service.WebSocketService
	upgrader          websocket.Upgrader
	wsManager         *ws.WebSocketManager
	clientConnections map[int64]*websocket.Conn // 存储用户ID到连接的映射
	mu                sync.RWMutex              // 保护并发访问
}

func NewWebSocketHandler() *WebSocketHandler {
	wsService := service.NewWebSocketService(config.GetDB())
	return &WebSocketHandler{
		service:           &wsService,
		wsManager:         ws.GetInstance(),
		clientConnections: make(map[int64]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 获取用户ID
	id := r.URL.Query().Get("userID")
	if id == "" {
		http.Error(w, "用户ID不能为空", http.StatusBadRequest)
		return
	}
	userID, _ := strconv.ParseInt(id, 10, 64)

	// 更新用户为在线状态
	if err := h.service.UpdateUserOnlineStatus(userID, true); err != nil {
		log.Printf("更新用户状态失败: %v", err)
		http.Error(w, "更新用户状态失败", http.StatusInternalServerError)
		return
	}

	// 建立WebSocket连接
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// 连接失败，更新用户为离线状态
		h.service.UpdateUserOnlineStatus(userID, false)
		http.Error(w, "升级到WebSocket失败", http.StatusInternalServerError)
		return
	}

	// 存储连接
	h.mu.Lock()
	h.clientConnections[userID] = conn
	h.mu.Unlock()

	log.Printf("用户 %d WebSocket连接建立成功", userID)

	// 通知好友该用户上线
	go h.notifyFriendsOnlineStatus(userID, true)

	// 设置连接关闭时的回调
	conn.SetCloseHandler(func(code int, text string) error {
		h.handleConnectionClose(userID)
		return nil
	})

	// 启动消息处理循环
	go h.handleMessages(userID, conn)
}

// handleConnectionClose 处理连接关闭
func (h *WebSocketHandler) handleConnectionClose(userID int64) {
	// 更新用户为离线状态
	if err := h.service.UpdateUserOnlineStatus(userID, false); err != nil {
		log.Printf("更新用户离线状态失败: %v", err)
	}

	// 移除连接记录
	h.mu.Lock()
	delete(h.clientConnections, userID)
	h.mu.Unlock()

	// 通知好友该用户下线
	go h.notifyFriendsOnlineStatus(userID, false)

	log.Printf("用户 %d WebSocket连接已关闭", userID)
}

// handleMessages 处理接收到的ws消息
func (h *WebSocketHandler) handleMessages(userID int64, conn *websocket.Conn) {
	defer func() {
		// 确保连接关闭时清理资源
		h.handleConnectionClose(userID)
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		// 设置读取超时
		err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("用户 %d WebSocket读取错误: %v", userID, err)
			}
			break
		}

		// 处理接收到的消息
		if err := h.processMessage(userID, message); err != nil {
			log.Printf("处理用户 %d 的消息失败: %v", userID, err)
		}
	}
}

// handlePing 处理ws心跳消息
func (h *WebSocketHandler) handlePing(userID int64) error {
	h.mu.RLock()
	conn, exists := h.clientConnections[userID]
	h.mu.RUnlock()

	if !exists {
		return fmt.Errorf("用户连接不存在")
	}

	// 发送pong响应
	pongMsg := map[string]interface{}{
		"type": "pong",
		"data": map[string]interface{}{
			"timestamp": time.Now().Unix(),
		},
	}

	return conn.WriteJSON(pongMsg)
}

// isUserOnline 检查用户是否在线
func (h *WebSocketHandler) isUserOnline(userID int64) bool {
	h.mu.RLock()
	_, exists := h.clientConnections[userID]
	h.mu.RUnlock()
	return exists
}

// sendToUser 向指定用户发送消息
func (h *WebSocketHandler) sendToUser(userID int64, message interface{}) error {
	h.mu.RLock()
	conn, exists := h.clientConnections[userID]
	h.mu.RUnlock()

	if !exists {
		return fmt.Errorf("用户 %d 不在线", userID)
	}

	return conn.WriteJSON(message)
}

// GetOnlineUsers 获取在线用户列表（HTTP API）
func (h *WebSocketHandler) GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	onlineUsers := make([]int64, 0, len(h.clientConnections))
	for userID := range h.clientConnections {
		onlineUsers = append(onlineUsers, userID)
	}
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    onlineUsers,
	})
}

// SendMessage 发送消息API（HTTP API）
func (h *WebSocketHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserIDStr := vars["userID"]
	targetUserID, _ := strconv.ParseInt(targetUserIDStr, 10, 64)

	var message struct {
		Content string `json:"content"`
		Type    string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 这里可以获取当前用户ID（从session或token中）
	fromUserID := int64(0) // 示例，实际应该从认证信息中获取

	if err := h.forwardMessageToUser(targetUserID, fromUserID, message.Content, message.Type); err != nil {
		http.Error(w, "发送消息失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "消息发送成功",
	})
}

/**
 * 处理消息已读回执
 */
// handleReadReceipt 处理已读回执
func (h *WebSocketHandler) handleReadReceipt(userID int64, data json.RawMessage) error {
	var readMsg struct {
		MessageID string `json:"message_id"` // 消息ID
		From      int64  `json:"from"`       // 消息发送者
	}

	if err := json.Unmarshal(data, &readMsg); err != nil {
		return fmt.Errorf("解析已读回执失败: %v", err)
	}

	// 通知消息发送者消息已被阅读
	return h.notifyMessageRead(readMsg.From, readMsg.MessageID, userID)
}

// notifyMessageRead 通知消息发送者消息已被阅读
func (h *WebSocketHandler) notifyMessageRead(senderUserID int64, messageID string, readerUserID int64) error {
	// 构建已读回执消息
	readReceipt := map[string]interface{}{
		"type": "message_read",
		"data": map[string]interface{}{
			"message_id":    messageID,
			"reader_id":     readerUserID,
			"read_time":     time.Now().Unix(),
			"read_time_str": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	// 发送给消息的原始发送者
	return h.sendToUser(senderUserID, readReceipt)
}
