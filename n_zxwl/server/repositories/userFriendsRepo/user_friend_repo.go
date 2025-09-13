// Package userFriendsRepo repository/user_friends_repo.go
package userFriendsRepo

import (
	"gorm.io/gorm"
	"mymod/model/sqlModel"
)

type UserFriendsRepository struct {
	db *gorm.DB
}

// NewUserFriendsRepository 创建好友仓库实例
func NewUserFriendsRepository(db *gorm.DB) UserFriendsRepository {
	return UserFriendsRepository{db: db}
}

// CreateFriendRequest 创建好友请求
func (r *UserFriendsRepository) CreateFriendRequest(friendRequest *sqlModel.UserFriend) error {
	return r.db.Create(friendRequest).Error
}

// GetFriendRequest 获取好友关系
func (r *UserFriendsRepository) GetFriendRequest(userID, friendID int64) (*sqlModel.UserFriend, error) {
	var friendRequest sqlModel.UserFriend
	err := r.db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friendRequest).Error
	if err != nil {
		return nil, err
	}
	return &friendRequest, nil
}

// UpdateFriendStatus 更新好友状态
func (r *UserFriendsRepository) UpdateFriendStatus(userID, friendID int64, status string) error {
	return r.db.Model(&sqlModel.UserFriend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("status", status).Error
}

// DeleteFriend 删除好友关系
func (r *UserFriendsRepository) DeleteFriend(userID, friendID int64) error {
	return r.db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friendID, friendID, userID).Delete(&sqlModel.UserFriend{}).Error
}

// GetUserFriends 获取用户好友列表（包含好友信息）
func (r *UserFriendsRepository) GetUserFriends(userID int64, status string) ([]sqlModel.FriendResponse, error) {
	var friends []sqlModel.FriendResponse

	query := `
		SELECT uf.*, 
		       up.id as friend_info_id, 
		       up.username as friend_info_username,
		       up.email as friend_info_email,
		       up.avatar_url as friend_info_avatar_url,
		       up.display_name as friend_info_display_name,
		       up.gender as friend_info_gender,
		       up.is_online as friend_info_is_online
		FROM user_friend uf
		LEFT JOIN user_profile up ON 
			CASE 
				WHEN uf.user_id = ? THEN uf.friend_id = up.id
				WHEN uf.friend_id = ? THEN uf.user_id = up.id
			END
		WHERE (uf.user_id = ? OR uf.friend_id = ?) AND uf.status = ?
		ORDER BY uf.updated_at DESC
	`

	err := r.db.Raw(query, userID, userID, userID, userID, status).Scan(&friends).Error
	if err != nil {
		return nil, err
	}

	return friends, nil
}

// GetPendingRequests 获取待处理的好友请求
func (r *UserFriendsRepository) GetPendingRequests(userID int64) ([]sqlModel.FriendResponse, error) {
	var requests []sqlModel.FriendResponse

	query := `
		SELECT uf.*, 
		       up.id as friend_info_id, 
		       up.username as friend_info_username,
		       up.email as friend_info_email,
		       up.avatar_url as friend_info_avatar_url,
		       up.display_name as friend_info_display_name,
		       up.gender as friend_info_gender,
		       up.is_online as friend_info_is_online
		FROM user_friend uf
		LEFT JOIN user_profile up ON uf.user_id = up.id
		WHERE uf.friend_id = ? AND uf.status = 'pending'
		ORDER BY uf.created_at DESC
	`

	err := r.db.Raw(query, userID).Scan(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// CheckFriendRelation 检查好友关系是否存在
func (r *UserFriendsRepository) CheckFriendRelation(userID, friendID int64) (bool, error) {
	var count int64
	err := r.db.Model(&sqlModel.UserFriend{}).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			userID, friendID, friendID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetFriendCount 获取好友数量
func (r *UserFriendsRepository) GetFriendCount(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&sqlModel.UserFriend{}).
		Where("(user_id = ? OR friend_id = ?) AND status = 'accepted'", userID, userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
