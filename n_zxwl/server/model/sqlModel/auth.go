package sqlModel

// ChangePasswordRequest 请求数据结构
type ChangePasswordRequest struct {
	User        string `json:"user"`
	Email       string `json:"email"`
	Passwd      string `json:"passwd"`
	VerifyCode  string `json:"varifycode"`
	NewPassword string `json:"new_password"` // 如果需要确认新密码
}
