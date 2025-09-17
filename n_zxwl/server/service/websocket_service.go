// service/websocket_service.go
package service

import (
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/repositories/websocketRepo"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type WebSocketService struct {
	repo      websocketRepo.WebSocketRepository
	clients   map[int]*Client
	broadcast chan sqlModel.WSMessage
	mutex     sync.RWMutex
}

type Client struct {
	ID           string
	UserID       int
	Conn         *websocket.Conn
	Send         chan sqlModel.WSMessage
	LastActivity time.Time
}

func NewWebSocketService(db *gorm.DB) WebSocketService {
	repo := websocketRepo.NewWebSocketRepository(db)
	return WebSocketService{
		repo:      repo,
		clients:   make(map[int]*Client),
		broadcast: make(chan sqlModel.WSMessage, 100),
	}
}

// UpdateUserOnlineStatus 更新用户在线状态
func (s *WebSocketService) UpdateUserOnlineStatus(userID int64, isOnline bool) error {
	db := config.GetDB()

	updateData := map[string]interface{}{
		"is_online":        isOnline,
		"last_online_time": time.Now(),
	}

	result := db.Model(&sqlModel.UserProfile{}).Where("id = ?", userID).Updates(updateData)
	return result.Error
}
