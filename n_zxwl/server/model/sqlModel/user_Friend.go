package sqlModel

import "time"

// 好友状态常量
const (
	FriendStatusPending  = "pending"
	FriendStatusAccepted = "accepted"
	FriendStatusBlocked  = "blocked"
)

// UserFriend 好友关系模型
type UserFriend struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         int       `gorm:"not null" json:"userId"`
	FriendID       int       `gorm:"not null" json:"friendId"`
	Salutation     *string   `gorm:"column:salutation" json:"salutation"`           // 对好友的称呼
	RequestMessage *string   `gorm:"column:request_message" json:"request_message"` // 打招呼信息
	Status         string    `gorm:"type:enum('pending','accepted','blocked');default:'pending'" json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// 关联关系
	User   UserProfile `gorm:"foreignKey:UserID" json:"user"`
	Friend UserProfile `gorm:"foreignKey:FriendID" json:"friend"`
}

// AddFriendRequest 添加好友请求结构体
type AddFriendRequest struct {
	FriendID       int64   `json:"friend_id" binding:"required"`
	Salutation     *string `json:"salutation"`
	RequestMessage *string `json:"request_message"`
}

func (UserFriend) TableName() string {
	return "user_friend"
}

// FriendResponse 好友响应结构体
type FriendResponse struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	FriendID       int64       `json:"friend_id"`
	Status         string      `json:"status"`
	Salutation     *string     `json:"salutation"`
	RequestMessage *string     `json:"request_message"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	FriendInfo     UserProfile `json:"friend_info"`
}

// IsValidStatus 验证状态有效性
func (uf *UserFriend) IsValidStatus() bool {
	switch uf.Status {
	case FriendStatusPending, FriendStatusAccepted, FriendStatusBlocked:
		return true
	default:
		return false
	}
}
