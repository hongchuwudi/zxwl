package userApi

import (
	"encoding/json"
	"fmt"
	"log"
	"mymod/model/param"
	"mymod/model/param/userParam"
	"net/http"
)

// LoginHandler 用户登录
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	var req userParam.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Login == "" || req.Password == "" {
		fmt.Println(req.Login, req.Password)
		response.Error = http.StatusBadRequest
		response.Message = "登录名和密码不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	loginResp, err := h.userService.Login(req)
	if err != nil {
		log.Printf("登录失败: %v", err)
		response.Error = http.StatusUnauthorized
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	// 不返回密码等敏感信息
	loginResp.User.PasswordHash = ""

	response.Error = 0
	response.Code = 0
	response.Message = "登录成功"
	response.Data = loginResp
	json.NewEncoder(w).Encode(response)
}
