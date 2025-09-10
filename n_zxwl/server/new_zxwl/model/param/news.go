package param

import "time"

// NewsQueryRequest 资讯分页查询请求
type NewsQueryRequest struct {
	Title      string `json:"title,omitempty"`       // 标题关键词
	Keywords   string `json:"keywords,omitempty"`    // 关键词
	FromSource string `json:"from_source,omitempty"` // 来源
	Page       int    `json:"page,omitempty"`        // 页码
	PageSize   int    `json:"page_size,omitempty"`   // 每页大小
}

// NewsVO 资讯响应VO
type NewsVO struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Keywords    string    `json:"keywords"`
	FromSource  string    `json:"from_source"`
	NewsFrom    string    `json:"news_from"`
	ClassName   string    `json:"class_name"`
	PublishTime time.Time `gorm:"column:publish_time" json:"publish_time"`
	CreateTime  time.Time `json:"create_time"`
	StyleURL    string    `json:"style_url,omitempty"`
}
