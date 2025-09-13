// Package userFriendsApi api/user_friends_api.go
package userFriendsApi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/repositories/userFriendsRepo"
	"mymod/service"
	"net/http"
	"strconv"
)

type UserFriendsHandler struct {
	service *service.UserFriendsService
}

// NewUserFriendsHandler 使用默认配置创建handler
func NewUserFriendsHandler() *UserFriendsHandler {
	repo := userFriendsRepo.NewUserFriendsRepository(config.GetDB())
	serv := service.NewUserFriendsService(repo)
	return &UserFriendsHandler{
		service: &serv,
	}
}

// AddFriend 添加好友
func (h *UserFriendsHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	var request sqlModel.AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"code": 400, "message": "无效的请求数据"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.AddFriend(userID, &request); err != nil {
		http.Error(w, `{"code": 500, "message": "添加好友失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    201,
		"message": "好友请求已发送",
	})
}

// AcceptFriend 同意好友请求
func (h *UserFriendsHandler) AcceptFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	friendIDStr := vars["friendID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的好友ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.AcceptFriend(userID, friendID); err != nil {
		http.Error(w, `{"code": 500, "message": "同意好友请求失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "好友请求已同意",
	})
}

// RejectFriend 拒绝好友请求
func (h *UserFriendsHandler) RejectFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	friendIDStr := vars["friendID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的好友ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.RejectFriend(userID, friendID); err != nil {
		http.Error(w, `{"code": 500, "message": "拒绝好友请求失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "好友请求已拒绝",
	})
}

// DeleteFriend 删除好友
func (h *UserFriendsHandler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	friendIDStr := vars["friendID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的好友ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteFriend(userID, friendID); err != nil {
		http.Error(w, `{"code": 500, "message": "删除好友失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "好友已删除",
	})
}

// GetFriends 获取好友列表
func (h *UserFriendsHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friends, err := h.service.GetFriends(userID)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取好友列表失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    friends,
	})
}

// GetPendingRequests 获取待处理的好友请求
func (h *UserFriendsHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	requests, err := h.service.GetPendingRequests(userID)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取待处理请求失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    requests,
	})
}

// GetFriendCount 获取好友数量
func (h *UserFriendsHandler) GetFriendCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	count, err := h.service.GetFriendCount(userID)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取好友数量失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    count,
	})
}
