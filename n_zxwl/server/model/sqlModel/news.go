package sqlModel

import "time"

// NewsInfo 资讯数据表
type NewsInfo struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvinceID     string    `json:"province_id"`
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `json:"description"`
	Keywords       string    `json:"keywords"`
	Content        string    `json:"content"`
	VideoDetail    string    `json:"video_detail"`
	VideoType      string    `json:"video_type"`
	VideoImg       string    `json:"video_img"`
	FromSource     string    `json:"from_source"`
	NewsNum        string    `json:"news_num"`
	IsPush         int       `json:"is_push"`
	IsTop          int       `json:"is_top"`
	StyleType      string    `json:"style_type"`
	StyleURL       string    `json:"style_url"`
	PublishTime    time.Time `gorm:"column:publish_time" json:"publish_time"`
	AddTime        time.Time `json:"add_time"`
	CardSchoolID   string    `json:"card_school_id"`
	CardLiveID     string    `json:"card_live_id"`
	ClassName      string    `json:"class_name"`
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
	PublisherEmail string    `json:"publisher_email"`
	LikeCount      int       `json:"like_count"`
	ShareCount     int       `json:"share_count"`
	FavoriteCount  int       `json:"favorite_count"`
	CommentCount   int       `json:"comment_count"`
}

// NewsComment 资讯评论表
type NewsComment struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	NewsID         uint64    `gorm:"not null" json:"news_id"`
	ParentID       uint64    `gorm:"default:0" json:"parent_id"`
	CommenterName  string    `gorm:"size:100;not null" json:"commenter_name"`
	CommenterEmail string    `gorm:"size:255;not null" json:"commenter_email"`
	CommentContent string    `gorm:"type:text;not null" json:"comment_content"`
	CommentTime    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"comment_time"`
	IsApproved     int       `gorm:"default:0" json:"is_approved"`
	LikeCount      int       `gorm:"default:0" json:"like_count"`
	ReplyCount     int       `gorm:"default:0" json:"reply_count"`
	IpAddress      string    `gorm:"size:45" json:"ip_address"`
	UserAgent      string    `gorm:"size:500" json:"user_agent"`

	// 非表字段: 用于存储回复列表
	Replies []NewsComment `gorm:"-" json:"replies,omitempty"`
}

// TableName 指定表名
func (NewsComment) TableName() string {
	return "news_comments"
}

func (NewsInfo) TableName() string {
	return "news_info"
}
