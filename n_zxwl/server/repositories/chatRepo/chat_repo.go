// Package chatRepo repository/chat_repo.go
package chatRepo

import (
	"gorm.io/gorm"
	"mymod/model/sqlModel"
)

type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository 创建聊天仓库实例
func NewChatRepository(db *gorm.DB) ChatRepository {
	return ChatRepository{db: db}
}

// CreateFriendMessage 创建好友消息
func (r *ChatRepository) CreateFriendMessage(message *sqlModel.UserFriendChatMessage) error {
	return r.db.Create(message).Error
}

// CreateRoomMessage 创建聊天室消息
func (r *ChatRepository) CreateRoomMessage(message *sqlModel.ChatRoomChatMessage) error {
	return r.db.Create(message).Error
}

// GetFriendMessages 获取好友聊天记录
func (r *ChatRepository) GetFriendMessages(userID, friendID int64, limit int) ([]sqlModel.UserFriendChatMessage, error) {
	var messages []sqlModel.UserFriendChatMessage

	err := r.db.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, friendID, friendID, userID).
		Order("created_at").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

// GetRoomMessages 获取聊天室消息记录
func (r *ChatRepository) GetRoomMessages(meetingID string, limit int) ([]sqlModel.ChatRoomChatMessage, error) {
	var messages []sqlModel.ChatRoomChatMessage

	err := r.db.Preload("Sender").Preload("Meeting").
		Where("meeting_id = ?", meetingID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

// MarkMessagesAsRead 标记消息为已读
func (r *ChatRepository) MarkMessagesAsRead(userID, friendID int64) error {
	return r.db.Model(&sqlModel.UserFriendChatMessage{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = false", userID, friendID).
		Update("is_read", true).Error
}

// GetUnreadCount 获取未读消息数量
func (r *ChatRepository) GetUnreadCount(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&sqlModel.UserFriendChatMessage{}).
		Where("receiver_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

// GetUnreadCountByFriend 获取指定好友的未读消息数量
func (r *ChatRepository) GetUnreadCountByFriend(userID, friendID int64) (int64, error) {
	var count int64
	err := r.db.Model(&sqlModel.UserFriendChatMessage{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = false", userID, friendID).
		Count(&count).Error
	return count, err
}

// GetRecentChats 获取最近聊天列表
func (r *ChatRepository) GetRecentChats(userID int64, limit int) ([]sqlModel.UserFriendChatMessage, error) {
	var messages []sqlModel.UserFriendChatMessage

	// 修正Group方法调用
	subQuery := r.db.Model(&sqlModel.UserFriendChatMessage{}).
		Select("MAX(created_at) as max_time, sender_id, receiver_id").
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Group("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END") // 移除多余的参数

	err := r.db.Preload("Sender").Preload("Receiver").
		Joins("INNER JOIN (?) as latest ON user_friend_chat_messages.created_at = latest.max_time AND "+
			"(user_friend_chat_messages.sender_id = latest.sender_id OR user_friend_chat_messages.receiver_id = latest.receiver_id)", subQuery).
		Where("user_friend_chat_messages.sender_id = ? OR user_friend_chat_messages.receiver_id = ?", userID, userID).
		Order("user_friend_chat_messages.created_at DESC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

// DeleteFriendMessage 删除好友消息
func (r *ChatRepository) DeleteFriendMessage(messageID, userID int64) error {
	return r.db.Where("id = ? AND (sender_id = ? OR receiver_id = ?)", messageID, userID, userID).
		Delete(&sqlModel.UserFriendChatMessage{}).Error
}

// DeleteRoomMessage 删除聊天室消息
func (r *ChatRepository) DeleteRoomMessage(messageID, userID int64) error {
	return r.db.Where("id = ? AND sender_id = ?", messageID, userID).
		Delete(&sqlModel.ChatRoomChatMessage{}).Error
}

// ClearChatHistory 清空聊天记录
func (r *ChatRepository) ClearChatHistory(userID, friendID int64) error {
	return r.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, friendID, friendID, userID).
		Delete(&sqlModel.UserFriendChatMessage{}).Error
}
