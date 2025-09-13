// Package service/user_friends_service.go
package service

import (
	"errors"
	"fmt"
	"mymod/model/sqlModel"
	"mymod/repositories/userFriendsRepo"
)

type UserFriendsService struct {
	repo userFriendsRepo.UserFriendsRepository
}

func NewUserFriendsService(repo userFriendsRepo.UserFriendsRepository) UserFriendsService {
	return UserFriendsService{repo: repo}
}

// AddFriend 添加好友
func (s *UserFriendsService) AddFriend(userID int64, request *sqlModel.AddFriendRequest) error {
	// 检查是否已经是好友
	exists, err := s.repo.CheckFriendRelation(userID, request.FriendID)
	if err != nil {
		return fmt.Errorf("检查好友关系失败: %v", err)
	}
	if exists {
		return errors.New("已经是好友关系或已发送请求")
	}

	// 检查不能添加自己为好友
	if userID == request.FriendID {
		return errors.New("不能添加自己为好友")
	}

	friendRequest := &sqlModel.UserFriend{
		UserID:         int(userID),
		FriendID:       int(request.FriendID),
		Status:         sqlModel.FriendStatusPending,
		Salutation:     request.Salutation,
		RequestMessage: request.RequestMessage,
	}

	return s.repo.CreateFriendRequest(friendRequest)
}

// AcceptFriend 同意好友请求
func (s *UserFriendsService) AcceptFriend(userID, friendID int64) error {
	// 检查请求是否存在
	request, err := s.repo.GetFriendRequest(friendID, userID)
	if err != nil {
		return errors.New("好友请求不存在")
	}

	if request.Status != sqlModel.FriendStatusPending {
		return errors.New("好友请求状态不正确")
	}

	return s.repo.UpdateFriendStatus(friendID, userID, sqlModel.FriendStatusAccepted)
}

// RejectFriend 拒绝好友请求
func (s *UserFriendsService) RejectFriend(userID, friendID int64) error {
	// 检查请求是否存在
	request, err := s.repo.GetFriendRequest(friendID, userID)
	if err != nil {
		return errors.New("好友请求不存在")
	}

	if request.Status != sqlModel.FriendStatusPending {
		return errors.New("好友请求状态不正确")
	}

	return s.repo.DeleteFriend(friendID, userID)
}

// DeleteFriend 删除好友
func (s *UserFriendsService) DeleteFriend(userID, friendID int64) error {
	// 检查是否是好友
	exists, err := s.repo.CheckFriendRelation(userID, friendID)
	if err != nil {
		return fmt.Errorf("检查好友关系失败: %v", err)
	}
	if !exists {
		return errors.New("不是好友关系")
	}

	return s.repo.DeleteFriend(userID, friendID)
}

// GetFriends 获取好友列表
func (s *UserFriendsService) GetFriends(userID int64) ([]sqlModel.FriendResponse, error) {
	return s.repo.GetUserFriends(userID, sqlModel.FriendStatusAccepted)
}

// GetPendingRequests 获取待处理的好友请求
func (s *UserFriendsService) GetPendingRequests(userID int64) ([]sqlModel.FriendResponse, error) {
	return s.repo.GetPendingRequests(userID)
}

// BlockFriend 拉黑好友
func (s *UserFriendsService) BlockFriend(userID, friendID int64) error {
	// 检查是否是好友
	exists, err := s.repo.CheckFriendRelation(userID, friendID)
	if err != nil {
		return fmt.Errorf("检查好友关系失败: %v", err)
	}
	if !exists {
		return errors.New("不是好友关系")
	}

	return s.repo.UpdateFriendStatus(userID, friendID, sqlModel.FriendStatusBlocked)
}

// GetFriendCount 获取好友数量
func (s *UserFriendsService) GetFriendCount(userID int64) (int64, error) {
	return s.repo.GetFriendCount(userID)
}
