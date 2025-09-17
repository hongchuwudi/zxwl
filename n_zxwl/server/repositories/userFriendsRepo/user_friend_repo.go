// Package userFriendsRepo repository/friend_repository.go
package userFriendsRepo

import (
	"errors"
	"gorm.io/gorm"
	"mymod/model/sqlModel"
)

type FriendRepository struct {
	Db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) *FriendRepository {
	return &FriendRepository{Db: db}
}

// CheckFriendship 检查是否已经是好友
func (r *FriendRepository) CheckFriendship(userID, friendID int64) (bool, error) {
	var count int64
	err := r.Db.Model(&sqlModel.UserFriend{}).
		Where("(user_a_id = ? AND user_b_id = ?) OR (user_a_id = ? AND user_b_id = ?)",
			userID, friendID, friendID, userID).
		Count(&count).Error
	return count > 0, err
}

// CheckFriendRequest 检查是否存在待处理的好友请求
func (r *FriendRepository) CheckFriendRequest(fromUserID, toUserID int64) (bool, error) {
	var count int64
	err := r.Db.Model(&sqlModel.FriendRequest{}).
		Where("from_user_id = ? AND to_user_id = ? AND status = 'pending'", fromUserID, toUserID).
		Count(&count).Error
	return count > 0, err
}

// CreateFriendRequest 创建好友请求
func (r *FriendRepository) CreateFriendRequest(request *sqlModel.FriendRequest) (error, int64) {
	return r.Db.Create(request).Error, request.ID
}

// GetFriendRequestByID 根据ID获取好友请求
func (r *FriendRepository) GetFriendRequestByID(requestID int64) (*sqlModel.FriendRequest, error) {
	var request sqlModel.FriendRequest
	err := r.Db.Preload("FromUser").First(&request, requestID).Error
	return &request, err
}

// UpdateFriendRequestStatus 更新好友请求状态
func (r *FriendRepository) UpdateFriendRequestStatus(requestID int64, status string) error {
	return r.Db.Model(&sqlModel.FriendRequest{}).
		Where("id = ?", requestID).
		Update("status", status).Error
}

// CreateFriendship 创建好友关系
func (r *FriendRepository) CreateFriendship(friendship *sqlModel.UserFriend) error {
	return r.Db.Create(friendship).Error
}

// DeleteFriendship 删除好友关系
func (r *FriendRepository) DeleteFriendship(userID, friendID int64) error {
	return r.Db.Where("(user_a_id = ? AND user_b_id = ?) OR (user_a_id = ? AND user_b_id = ?)",
		userID, friendID, friendID, userID).
		Delete(&sqlModel.UserFriend{}).Error
}

// SetFriendNickname 设置好友昵称
func (r *FriendRepository) SetFriendNickname(nickname *sqlModel.FriendNickname) error {
	return r.Db.Save(nickname).Error
}

// GetFriendNickname 获取好友昵称
func (r *FriendRepository) GetFriendNickname(userID, friendID int64) (*sqlModel.FriendNickname, error) {
	var nickname sqlModel.FriendNickname
	err := r.Db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&nickname).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &nickname, err
}

// GetUserFriends 获取用户的所有好友
func (r *FriendRepository) GetUserFriends(userID int64) ([]sqlModel.UserFriend, error) {
	var friendships []sqlModel.UserFriend
	err := r.Db.Preload("UserA").Preload("UserB").
		Where("user_a_id = ? OR user_b_id = ?", userID, userID).
		Find(&friendships).Error
	for i := range friendships {
		if friendships[i].UserAID == userID {
			friendships[i].UserA = friendships[i].UserB
		}
		friendships[i].UserB = sqlModel.UserProfile{} // 清空当前用户信息
	}
	return friendships, err
}

// GetFriendRequestsToMe 获取发给我的好友请求
func (r *FriendRepository) GetFriendRequestsToMe(userID int64) ([]sqlModel.FriendRequest, error) {
	var requests []sqlModel.FriendRequest
	err := r.Db.Preload("FromUser").
		Where("to_user_id = ? AND status = 'pending'", userID).
		Find(&requests).Error
	return requests, err
}

// GetFriendRequestsFromMe 获取我发出的好友请求
func (r *FriendRepository) GetFriendRequestsFromMe(userID int64) ([]sqlModel.FriendRequest, error) {
	var requests []sqlModel.FriendRequest
	err := r.Db.Preload("ToUser").
		Where("from_user_id = ? AND status = 'pending'", userID).
		Find(&requests).Error
	return requests, err
}
