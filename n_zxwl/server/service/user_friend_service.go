// Package service service/friend_service.go
package service

import (
	"errors"
	"gorm.io/gorm"
	"mymod/model/sqlModel"
	"mymod/repositories/userFriendsRepo"
	"time"
)

type FriendService struct {
	repo *userFriendsRepo.FriendRepository
}

func NewFriendService(repo *userFriendsRepo.FriendRepository) *FriendService {
	return &FriendService{repo: repo}
}

// SendFriendRequest 发送好友请求
func (s *FriendService) SendFriendRequest(fromUserID, toUserID int64, salutation, message *string) (error, int64) {
	// 检查是否已经是好友
	isFriend, err := s.repo.CheckFriendship(fromUserID, toUserID)
	if err != nil {
		return err, 0
	}
	if isFriend {
		return errors.New("已经是好友关系"), 0
	}

	// 检查是否已经发送过请求
	hasRequest, err := s.repo.CheckFriendRequest(fromUserID, toUserID)
	if err != nil {
		return err, 0
	}
	if hasRequest {
		return errors.New("已经发送过好友请求"), 0
	}

	// 创建好友请求
	request := &sqlModel.FriendRequest{
		FromUserID:     fromUserID,
		ToUserID:       toUserID,
		Salutation:     salutation,
		RequestMessage: message,
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return s.repo.CreateFriendRequest(request)
}

// AcceptFriendRequest 接受好友请求
func (s *FriendService) AcceptFriendRequest(requestID, userID int64) error {
	// 获取请求
	request, err := s.repo.GetFriendRequestByID(requestID)
	if err != nil {
		return err
	}

	// 验证权限
	if request.ToUserID != userID {
		return errors.New("无权操作此请求")
	}

	if request.Status != "pending" {
		return errors.New("请求状态不正确")
	}

	// 开始事务
	tx := s.repo.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取请求
	request, err = s.repo.GetFriendRequestByID(requestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 验证权限
	if request.ToUserID != userID {
		tx.Rollback()
		return errors.New("无权操作此请求")
	}

	if request.Status != "pending" {
		tx.Rollback()
		return errors.New("请求状态不正确")
	}

	// 更新请求状态
	if err := tx.Model(&sqlModel.FriendRequest{}).Where("id = ?", requestID).Update("status", "accepted").Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建好友关系
	friendship := &sqlModel.UserFriend{
		UserAID:      request.FromUserID,
		UserBID:      request.ToUserID,
		RelationType: "normal",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := tx.Create(friendship).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// RejectFriendRequest 拒绝好友请求
func (s *FriendService) RejectFriendRequest(requestID, userID int64) error {
	request, err := s.repo.GetFriendRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.ToUserID != userID {
		return errors.New("无权操作此请求")
	}

	return s.repo.UpdateFriendRequestStatus(requestID, "rejected")
}

// DeleteFriend 删除好友
func (s *FriendService) DeleteFriend(userID, friendID int64) error {
	// 检查是否是好友
	isFriend, err := s.repo.CheckFriendship(userID, friendID)
	if err != nil {
		return err
	}
	if !isFriend {
		return errors.New("不是好友关系")
	}

	return s.repo.DeleteFriendship(userID, friendID)
}

// SetFriendNickname 设置好友昵称
func (s *FriendService) SetFriendNickname(userID, friendID int64, nickname string) error {
	// 检查是否是好友
	isFriend, err := s.repo.CheckFriendship(userID, friendID)
	if err != nil {
		return err
	}
	if !isFriend {
		return errors.New("不是好友关系，无法设置昵称")
	}

	nicknameRecord := &sqlModel.FriendNickname{
		UserID:    userID,
		FriendID:  friendID,
		Nickname:  nickname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.SetFriendNickname(nicknameRecord)
}

// GetUserFriends 获取用户的所有好友
func (s *FriendService) GetUserFriends(userID int64) ([]sqlModel.UserFriend, error) {
	// 获取好友列表
	friendships, err := s.repo.GetUserFriends(userID)
	if err != nil {
		return nil, err
	}

	// 为每个好友查询昵称
	for i := range friendships {
		var friendID int64
		// 确定好友ID（好友总是在UserA字段中）
		friendID = int64(friendships[i].UserA.ID)

		// 查询昵称
		nickname, err := s.repo.GetFriendNickname(userID, friendID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不是"记录不存在"的错误，返回错误
			return nil, err
		}

		// 设置昵称
		if nickname != nil {
			friendships[i].NickName = nickname.Nickname
		} else {
			friendships[i].NickName = "" // 或者可以设置为默认值
		}
	}

	return friendships, nil
}

// GetFriendRequestsToMe 获取发给我的好友请求
func (s *FriendService) GetFriendRequestsToMe(userID int64) ([]sqlModel.FriendRequest, error) {
	return s.repo.GetFriendRequestsToMe(userID)
}

// GetFriendRequestsFromMe 获取我发出的好友请求
func (s *FriendService) GetFriendRequestsFromMe(userID int64) ([]sqlModel.FriendRequest, error) {
	return s.repo.GetFriendRequestsFromMe(userID)
}
