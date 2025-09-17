// repositories/websocket_repo.go
package websocketRepo

import (
	"gorm.io/gorm"
	"mymod/model/sqlModel"
)

type WebSocketRepository struct {
	db *gorm.DB
}

func NewWebSocketRepository(db *gorm.DB) WebSocketRepository {
	return WebSocketRepository{db: db}
}

// CreateConnection 创建连接记录
func (r *WebSocketRepository) CreateConnection(conn *sqlModel.WSConnection) error {
	return r.db.Create(conn).Error
}
