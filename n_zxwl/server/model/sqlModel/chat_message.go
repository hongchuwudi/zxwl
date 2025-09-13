// Package sqlModel models/chat_messages.go
package sqlModel

import (
	"time"
)

// UserFriendChatMessage 好友私聊消息模型
type UserFriendChatMessage struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID    int64     `gorm:"not null" json:"sender_id"`
	ReceiverID  int64     `gorm:"not null" json:"receiver_id"`
	MessageType string    `gorm:"type:enum('text','image','file','system');default:'text'" json:"message_type"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	IsRead      bool      `gorm:"default:false" json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联关系
	Sender   UserProfile `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver UserProfile `gorm:"foreignKey:ReceiverID" json:"receiver"`
}

// TableName 指定表名
func (UserFriendChatMessage) TableName() string {
	return "user_friend_chat_messages"
}

// ChatRoomChatMessage 聊天室群聊消息模型
type ChatRoomChatMessage struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MeetingID   string    `gorm:"type:varchar(36);not null" json:"meeting_id"`
	SenderID    int64     `gorm:"not null" json:"sender_id"`
	MessageType string    `gorm:"type:enum('text','image','file','system');default:'text'" json:"message_type"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联关系
	Meeting ChatMeetingRoom `gorm:"foreignKey:MeetingID" json:"meeting"`
	Sender  UserProfile     `gorm:"foreignKey:SenderID" json:"sender"`
}

// TableName 指定表名
func (ChatRoomChatMessage) TableName() string {
	return "chat_room_chat_messages"
}

// SendFriendMessageRequest 发送好友消息请求
type SendFriendMessageRequest struct {
	ReceiverID  int64  `json:"receiver_id" binding:"required"`
	MessageType string `json:"message_type" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

// SendRoomMessageRequest 发送聊天室消息请求
type SendRoomMessageRequest struct {
	MeetingID   string `json:"meeting_id" binding:"required"`
	MessageType string `json:"message_type" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

// ChatMessageResponse 聊天消息响应
type ChatMessageResponse struct {
	ID          int64       `json:"id"`
	SenderID    int64       `json:"sender_id"`
	ReceiverID  *int64      `json:"receiver_id,omitempty"`
	MeetingID   *string     `json:"meeting_id,omitempty"`
	MessageType string      `json:"message_type"`
	Content     string      `json:"content"`
	IsRead      bool        `json:"is_read"`
	CreatedAt   time.Time   `json:"created_at"`
	SenderInfo  UserProfile `json:"sender_info"`
}
