// Package sqlModel models/user_choose.go
package sqlModel

import (
	"time"
)

// UserChoose 用户志愿选择表
type UserChoose struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"column:user_id;not null" json:"user_id"`
	SchoolName *string   `gorm:"column:school_name" json:"school_name"` // 修改字段名
	MajorName  *string   `gorm:"column:major_name" json:"major_name"`   // 修改字段名
	Priority   int8      `gorm:"column:priority;not null;default:0" json:"priority"`
	CreateTime time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
}

// TableName 指定表名
func (UserChoose) TableName() string {
	return "user_choose"
}

// UserChooseResponse 响应结构体
type UserChooseResponse struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	SchoolName *string   `json:"school_name"` // 修改字段名
	MajorName  *string   `json:"major_name"`  // 修改字段名
	Priority   int8      `json:"priority"`
	CreateTime time.Time `json:"create_time"`
}

// 优先级常量
const (
	PriorityHighest = 0
	PriorityHigh    = 1
	PriorityMedium  = 2
	PriorityLow     = 3
)

// IsValid 验证数据有效性
func (uc *UserChoose) IsValid() bool {
	if uc.SchoolName == nil && uc.MajorName == nil {
		return false
	}
	if uc.Priority < 0 || uc.Priority > 3 {
		return false
	}
	return true
}

// PriorityText 获取优先级文本
func (uc *UserChoose) PriorityText() string {
	switch uc.Priority {
	case PriorityHighest:
		return "最高优先级"
	case PriorityHigh:
		return "高优先级"
	case PriorityMedium:
		return "中优先级"
	case PriorityLow:
		return "低优先级"
	default:
		return "未知优先级"
	}
}
