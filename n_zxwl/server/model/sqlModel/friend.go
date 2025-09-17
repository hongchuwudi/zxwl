package sqlModel

import (
	"time"
)

type FriendRequest struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserID     int64     `gorm:"not null;index:idx_from_user" json:"from_user_id"`
	ToUserID       int64     `gorm:"not null;index:idx_to_user" json:"to_user_id"`
	Salutation     *string   `gorm:"size:50" json:"salutation"`
	RequestMessage *string   `gorm:"size:500" json:"request_message"`
	Status         string    `gorm:"type:enum('pending','accepted','rejected','canceled');default:'pending';index:idx_status" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联关系
	FromUser UserProfile `gorm:"foreignKey:FromUserID" json:"from_user"`
	ToUser   UserProfile `gorm:"foreignKey:ToUserID" json:"to_user"`
}

func (FriendRequest) TableName() string {
	return "friend_request"
}

type UserFriend struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserAID      int64     `gorm:"not null;index:idx_user_a" json:"user_a_id"`
	UserBID      int64     `gorm:"not null;index:idx_user_b" json:"user_b_id"`
	RelationType string    `gorm:"type:enum('mutual','close','normal');default:'normal'" json:"relation_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联关系
	UserA    UserProfile `gorm:"foreignKey:UserAID" json:"user_a"`
	UserB    UserProfile `gorm:"foreignKey:UserBID" json:"user_b"`
	NickName string      `gorm:"-" json:"nickname"`
}

func (UserFriend) TableName() string {
	return "user_friend"
}

type FriendNickname struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index:idx_user_id" json:"user_id"`
	FriendID  int64     `gorm:"not null;index:idx_friend_id" json:"friend_id"`
	Nickname  string    `gorm:"size:50;not null" json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联关系
	User   UserProfile `gorm:"foreignKey:UserID" json:"user"`
	Friend UserProfile `gorm:"foreignKey:FriendID" json:"friend"`
}

func (FriendNickname) TableName() string {
	return "friend_nickname"
}
