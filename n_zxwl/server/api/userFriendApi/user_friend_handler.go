// api/friend_api.go
package api

import (
	"encoding/json"
	"mymod/config"
	"mymod/repositories/userFriendsRepo"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"mymod/service"
)

type FriendHandler struct {
	friendService *service.FriendService
}

func NewFriendHandler() *FriendHandler {
	friendRepo := userFriendsRepo.NewFriendRepository(config.GetDB())
	friendServices := service.NewFriendService(friendRepo)
	return &FriendHandler{friendService: friendServices}
}

// SendFriendRequest 发送好友请求
func (h *FriendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	var request struct {
		Salutation *string `json:"salutation"`
		Message    *string `json:"request_message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"code": 400, "message": "无效的请求数据"}`, http.StatusBadRequest)
		return
	}

	err, friResId := h.friendService.SendFriendRequest(userID, friendID, request.Salutation, request.Message)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "发送请求失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "好友请求发送成功",
		"id":      friResId,
	})
}

// AcceptFriendRequest 接受好友请求
func (h *FriendHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	requestIDStr := vars["requestID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的请求ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.friendService.AcceptFriendRequest(requestID, userID); err != nil {
		http.Error(w, `{"code": 500, "message": "接受请求失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "好友请求已接受",
	})
}

// RejectFriendRequest 拒绝好友请求
func (h *FriendHandler) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	requestIDStr := vars["requestID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的请求ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.friendService.RejectFriendRequest(requestID, userID); err != nil {
		http.Error(w, `{"code": 500, "message": "拒绝请求失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "好友请求已接受",
	})
}

// GetUserFriends 获取用户好友列表
func (h *FriendHandler) GetUserFriends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	friends, err := h.friendService.GetUserFriends(userID)
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

// GetFriendRequestsToMe 获取发给我的好友请求
func (h *FriendHandler) GetFriendRequestsToMe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	requests, err := h.friendService.GetFriendRequestsToMe(userID)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "获取好友请求失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    requests,
	})
}

// SetFriendNickname 设置好友昵称
func (h *FriendHandler) SetFriendNickname(w http.ResponseWriter, r *http.Request) {
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

	var request struct {
		Nickname string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"code": 400, "message": "无效的请求数据"}`, http.StatusBadRequest)
		return
	}

	if err := h.friendService.SetFriendNickname(userID, friendID, request.Nickname); err != nil {
		http.Error(w, `{"code": 500, "message": "设置昵称失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "昵称设置成功",
	})
}

// DeleteFriend 删除好友
func (h *FriendHandler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
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

	if err := h.friendService.DeleteFriend(userID, friendID); err != nil {
		http.Error(w, `{"code": 500, "message": "删除好友失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "好友删除成功",
	})
}
