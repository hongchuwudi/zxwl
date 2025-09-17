// model/websocket.go
package sqlModel

import (
	"time"
)

// WSConnection WebSocket 连接记录
type WSConnection struct {
	ID           string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID       int        `gorm:"not null;uniqueIndex" json:"userId"` // 唯一索引确保单设备
	ConnectionID string     `gorm:"size:255;not null" json:"connectionId"`
	IsActive     bool       `gorm:"default:true" json:"isActive"`
	IPAddress    string     `gorm:"size:45" json:"ipAddress"`
	CreatedAt    time.Time  `json:"createdAt"`
	LastActivity time.Time  `json:"lastActivity"`
	ClosedAt     *time.Time `json:"closedAt"`
}

// WSMessage WebSocket 消息结构
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	To   int         `json:"to,omitempty"`
	From int         `json:"from,omitempty"`
}

func (WSConnection) TableName() string {
	return "ws_connection"
}
