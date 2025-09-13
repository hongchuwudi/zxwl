package userParam

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	VerifyCode  string `json:"verify_code"`
	DeviceInfo  string `json:"device_info"`
}
