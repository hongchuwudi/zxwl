package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Notification 表示通知消息结构
type Notification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	To   string      `json:"to,omitempty"`
}

// Client 表示 WebSocket 客户端
type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan Notification
	Hub    *Hub
}

// Hub 管理所有客户端和广播消息
type Hub struct {
	Clients    map[string]*Client // userID -> Client
	Broadcast  chan Notification
	Register   chan *Client
	Unregister chan *Client
	mutex      sync.RWMutex
}

// WebSocketManager 管理 WebSocket 连接
type WebSocketManager struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

var (
	instance *WebSocketManager
	once     sync.Once
)

// GetInstance 获取 WebSocketManager 单例实例
func GetInstance() *WebSocketManager {
	once.Do(func() {
		instance = &WebSocketManager{
			hub: &Hub{
				Clients:    make(map[string]*Client),
				Broadcast:  make(chan Notification),
				Register:   make(chan *Client),
				Unregister: make(chan *Client),
			},
			upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		}
		go instance.hub.run()
	})
	return instance
}

// HandleWebSocket 处理 WebSocket 连接
func (wm *WebSocketManager) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := wm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}

	client := &Client{
		ID:     fmt.Sprintf("%s-%d", userID, time.Now().UnixNano()),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan Notification, 256),
		Hub:    wm.hub,
	}

	client.Hub.Register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()

	log.Printf("客户端已连接: %s (用户ID: %s)", client.ID, userID)
}

// SendNotification 发送通知给特定用户
func (wm *WebSocketManager) SendNotification(userID string, notification Notification) {
	wm.hub.mutex.RLock()
	client, exists := wm.hub.Clients[userID]
	wm.hub.mutex.RUnlock()

	if exists {
		select {
		case client.Send <- notification:
		default:
			close(client.Send)
			wm.hub.Unregister <- client
		}
	}
}

// BroadcastNotification 广播通知给所有用户
func (wm *WebSocketManager) BroadcastNotification(notification Notification) {
	wm.hub.Broadcast <- notification
}

// GetConnectedUsers 获取所有已连接的用户ID
func (wm *WebSocketManager) GetConnectedUsers() []string {
	wm.hub.mutex.RLock()
	defer wm.hub.mutex.RUnlock()

	users := make([]string, 0, len(wm.hub.Clients))
	for userID := range wm.hub.Clients {
		users = append(users, userID)
	}
	return users
}

// IsUserConnected 检查用户是否已连接
func (wm *WebSocketManager) IsUserConnected(userID string) bool {
	wm.hub.mutex.RLock()
	defer wm.hub.mutex.RUnlock()

	_, exists := wm.hub.Clients[userID]
	return exists
}

// Hub 的主循环
func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client.UserID] = client
			h.mutex.Unlock()

			// 发送欢迎消息
			welcomeMsg := Notification{
				Type: "connected",
				Data: map[string]interface{}{
					"message": "WebSocket连接已建立",
					"time":    time.Now().Format(time.RFC3339),
				},
			}
			client.Send <- welcomeMsg

		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, exists := h.Clients[client.UserID]; exists {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.mutex.Unlock()

		case notification := <-h.Broadcast:
			h.mutex.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- notification:
				default:
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// 写协程
func (c *Client) writePump() {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			json.NewEncoder(w).Encode(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				json.NewEncoder(w).Encode(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 读协程
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512000)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 错误: %v", err)
			}
			break
		}

		var msg Notification
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("消息解析错误: %v", err)
			continue
		}

		c.handleMessage(msg)
	}
}

// 处理客户端消息
func (c *Client) handleMessage(msg Notification) {
	switch msg.Type {
	case "ping":
		response := Notification{
			Type: "pong",
			Data: map[string]interface{}{
				"timestamp": time.Now().Unix(),
			},
		}
		c.Send <- response

	case "auth":
		log.Printf("用户 %s 身份验证", c.UserID)

	default:
		log.Printf("收到未知消息类型: %s from user: %s", msg.Type, c.UserID)
	}
}
