package userApi

import (
	"encoding/json"
	"log"
	"mymod/model/param"
	"mymod/model/param/userParam"
	"net/http"
)

// RegisterHandler 用户注册
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	var req userParam.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证必要字段
	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.Error = http.StatusBadRequest
		response.Message = "用户名、邮箱和密码不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		log.Printf("注册失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	// 不返回密码等敏感信息
	user.PasswordHash = ""
	user.AuthToken = ""

	response.Error = 0
	response.Message = "注册成功"
	response.Data = user
	json.NewEncoder(w).Encode(response)
}
