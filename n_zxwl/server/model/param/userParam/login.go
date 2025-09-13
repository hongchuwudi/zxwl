package userParam

import (
	"mymod/model/sqlModel"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Login    string `json:"login"` // 可以是邮箱或用户名
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User        *sqlModel.UserProfile `json:"user"`
	AccessToken string                `json:"access_token"`
}
