package userParam

// ChangePasswordRequest 更改密码请求
type ChangePasswordRequest struct {
	Email       string `json:"email"`
	VerifyCode  string `json:"verify_code"`
	NewPassword string `json:"password"`
}
