package userParam

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatarUrl"`
	Gender      int8   `json:"gender"`
	BirthYear   int    `json:"birthYear"`
	Location    string `json:"location"`
	Bio         string `json:"bio"`
}

// UserQueryRequest 用户查询请求
type UserQueryRequest struct {
	Page      int    `form:"page" json:"page"`
	Size      int    `form:"size" json:"size"`
	Username  string `form:"username" json:"username"`
	Email     string `form:"email" json:"email"`
	StartTime string `form:"startTime" json:"start_time"`
	EndTime   string `form:"endTime" json:"end_time"`
	Gender    int    `form:"gender" json:"gender"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Data  interface{} `json:"data"`
}
