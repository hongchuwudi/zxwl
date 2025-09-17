// Package api file: api/user_handler.go
package userApi

import (
	"encoding/json"
	"log"
	"mymod/config"
	"mymod/model/param"
	userParam2 "mymod/model/param/userParam"
	"mymod/model/sqlModel"
	userRepo2 "mymod/repositories/userRepo"
	"mymod/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	userRepo := userRepo2.NewUserRepository(config.GetDB())
	verifyCodeRepo := userRepo2.NewVerifyCodeRepository(config.GetDB())
	userService := service.NewUserService(userRepo, verifyCodeRepo)
	return &UserHandler{userService: userService}
}

// DeleteUserHandler 删除用户
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodDelete {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 从URL参数获取用户ID
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	if userIDStr == "" {
		response.Error = http.StatusBadRequest
		response.Message = "用户ID不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "用户ID格式错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	err = h.userService.DeleteUser(userID)
	if err != nil {
		log.Printf("删除用户失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Error = 0
	response.Message = "删除成功"
	json.NewEncoder(w).Encode(response)
}

// GetUserHandler 根据ID获取用户信息
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodGet {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["id"]
	if userIDStr == "" {
		response.Error = http.StatusBadRequest
		response.Message = "用户ID不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "用户ID格式错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	if user == nil {
		response.Error = http.StatusNotFound
		response.Message = "用户不存在"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 不返回敏感信息
	user.PasswordHash = ""
	user.AuthToken = ""

	response.Error = 0
	response.Message = "获取成功"
	response.Data = user
	json.NewEncoder(w).Encode(response)
}

// GetUserHandlerByEmail 根据邮箱获取用户信息
func (h *UserHandler) GetUserHandlerByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodGet {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	vars := mux.Vars(r)
	email := vars["email"]
	if email == "" {
		response.Error = http.StatusBadRequest
		response.Message = "email不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	if user == nil {
		response.Error = http.StatusNotFound
		response.Message = "用户不存在"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 不返回敏感信息
	user.PasswordHash = ""
	user.AuthToken = ""

	response.Error = 0
	response.Message = "获取成功"
	response.Data = user
	json.NewEncoder(w).Encode(response)
}

// UpdateUserHandler 更新用户信息
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPut {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["id"]
	if userIDStr == "" {
		response.Error = http.StatusBadRequest
		response.Message = "用户ID不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "用户ID格式错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	var req sqlModel.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.userService.UpdateUser(userID, req)
	if err != nil {
		log.Printf("更新用户信息失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	// 不返回敏感信息
	user.PasswordHash = ""
	user.AuthToken = ""

	response.Error = 0
	response.Message = "更新成功"
	response.Data = user
	json.NewEncoder(w).Encode(response)
}

// ChangePasswordHandler 发送密码重置邮件
func (h *UserHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Code: 1001, Message: "请求失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	var req userParam2.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证必要字段
	if req.Email == "" || req.VerifyCode == "" || req.NewPassword == "" {
		response.Error = http.StatusBadRequest
		response.Message = "邮箱、验证码和新密码不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	err := h.userService.ChangePassword(req)
	if err != nil {
		log.Printf("更改密码失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0
	response.Message = "密码更改成功"
	json.NewEncoder(w).Encode(response)
}
