// Package service service/chat_service.go
package service

import (
	"errors"
	"mymod/model/sqlModel"
	"mymod/repositories/chatRepo"
)

type ChatService struct {
	repo chatRepo.ChatRepository
}

func NewChatService(repo chatRepo.ChatRepository) ChatService {
	return ChatService{repo: repo}
}

// SendFriendMessage 发送好友消息
func (s *ChatService) SendFriendMessage(senderID int64, request *sqlModel.SendFriendMessageRequest) error {
	message := &sqlModel.UserFriendChatMessage{
		SenderID:    senderID,
		ReceiverID:  request.ReceiverID,
		MessageType: request.MessageType,
		Content:     request.Content,
	}

	return s.repo.CreateFriendMessage(message)
}

// SendRoomMessage 发送聊天室消息
func (s *ChatService) SendRoomMessage(senderID int64, request *sqlModel.SendRoomMessageRequest) error {
	message := &sqlModel.ChatRoomChatMessage{
		MeetingID:   request.MeetingID,
		SenderID:    senderID,
		MessageType: request.MessageType,
		Content:     request.Content,
	}

	return s.repo.CreateRoomMessage(message)
}

// GetFriendChatHistory 获取好友聊天历史
func (s *ChatService) GetFriendChatHistory(userID, friendID int64, limit int) ([]sqlModel.UserFriendChatMessage, error) {
	if userID == friendID {
		return nil, errors.New("不能获取与自己的聊天记录")
	}

	return s.repo.GetFriendMessages(userID, friendID, limit)
}

// GetRoomChatHistory 获取聊天室聊天历史
func (s *ChatService) GetRoomChatHistory(meetingID string, limit int) ([]sqlModel.ChatRoomChatMessage, error) {
	return s.repo.GetRoomMessages(meetingID, limit)
}

// MarkFriendMessagesAsRead 标记好友消息为已读
func (s *ChatService) MarkFriendMessagesAsRead(userID, friendID int64) error {
	return s.repo.MarkMessagesAsRead(userID, friendID)
}

// GetUnreadMessageCount 获取未读消息总数
func (s *ChatService) GetUnreadMessageCount(userID int64) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

// GetUnreadMessageCountByFriend 获取指定好友的未读消息数量
func (s *ChatService) GetUnreadMessageCountByFriend(userID, friendID int64) (int64, error) {
	return s.repo.GetUnreadCountByFriend(userID, friendID)
}

// GetRecentChats 获取最近聊天列表
func (s *ChatService) GetRecentChats(userID int64, limit int) ([]sqlModel.UserFriendChatMessage, error) {
	return s.repo.GetRecentChats(userID, limit)
}

// DeleteFriendMessage 删除好友消息
func (s *ChatService) DeleteFriendMessage(messageID, userID int64) error {
	return s.repo.DeleteFriendMessage(messageID, userID)
}

// DeleteRoomMessage 删除聊天室消息
func (s *ChatService) DeleteRoomMessage(messageID, userID int64) error {
	return s.repo.DeleteRoomMessage(messageID, userID)
}

// ClearChatHistory 清空聊天记录
func (s *ChatService) ClearChatHistory(userID, friendID int64) error {
	return s.repo.ClearChatHistory(userID, friendID)
}
