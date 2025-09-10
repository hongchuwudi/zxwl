package sqlModel

import "time"

// User 用户表
type User struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UID       *int      `json:"uid,omitempty" gorm:"column:uid"`
	Name      string    `json:"name" gorm:"column:name;not null"`
	Pwd       string    `json:"-" gorm:"column:pwd;not null"` // 密码字段通常不序列化到JSON
	Email     string    `json:"email" gorm:"column:email;not null"`
	Token     *string   `json:"token,omitempty" gorm:"column:token"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
