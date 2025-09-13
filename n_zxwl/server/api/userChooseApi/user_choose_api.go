package userChooseApi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/repositories/userChooseRepo"
	"mymod/service"
	"net/http"
	"strconv"
)

type UserChooseHandler struct {
	service *service.UserChooseService
}

// NewUserChooseHandler 使用默认配置创建handler
func NewUserChooseHandler() *UserChooseHandler {
	repo := userChooseRepo.NewUserChooseRepository(config.GetDB())
	serv := service.NewUserChooseService(repo)
	return &UserChooseHandler{
		service: &serv,
	}
}

// GetUserChoices 获取用户所有志愿
func (h *UserChooseHandler) GetUserChoices(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取用户ID
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	// 查询数据
	choices, err := h.service.GetUserChoices(userID)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "查询失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    choices,
	})
}

// CreateUserChoice 创建用户志愿
func (h *UserChooseHandler) CreateUserChoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	var choice sqlModel.UserChoose
	if err := json.NewDecoder(r.Body).Decode(&choice); err != nil {
		http.Error(w, `{"code": 400, "message": "无效的请求数据"}`, http.StatusBadRequest)
		return
	}

	choice.UserID = userID

	if err := h.service.CreateUserChoice(&choice); err != nil {
		http.Error(w, `{"code": 500, "message": "创建失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    201,
		"message": "创建成功",
		"data":    choice,
	})
}

// DeleteUserChoice 删除用户志愿
func (h *UserChooseHandler) DeleteUserChoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	choiceIDStr := vars["choiceID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	choiceID, err := strconv.ParseInt(choiceIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的志愿ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUserChoice(choiceID, userID); err != nil {
		http.Error(w, `{"code": 500, "message": "删除失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "删除成功",
	})
}

// DeleteAllUserChoices 删除用户所有志愿
func (h *UserChooseHandler) DeleteAllUserChoices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"code": 400, "message": "无效的用户ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAllUserChoices(userID); err != nil {
		http.Error(w, `{"code": 500, "message": "删除失败: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "删除成功",
	})
}
