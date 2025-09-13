package userApi

import (
	"encoding/json"
	"log"
	"mymod/model/param"
	"mymod/model/param/userParam"
	"net/http"
	"strconv"
)

// GetUsersHandler 分页查询用户
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodGet {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	var req userParam.UserQueryRequest
	// 从查询参数解析
	query := r.URL.Query()
	req.Page, _ = strconv.Atoi(query.Get("page"))
	req.Size, _ = strconv.Atoi(query.Get("size"))
	req.Username = query.Get("username")
	req.Email = query.Get("email")
	req.StartTime = query.Get("startTime")
	req.EndTime = query.Get("endTime")
	req.Gender, _ = strconv.Atoi(query.Get("gender"))
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	result, err := h.userService.GetUsersWithPagination(req)
	if err != nil {
		log.Printf("查询用户列表失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Error = 0
	response.Message = "查询成功"
	response.Data = result
	json.NewEncoder(w).Encode(response)
}
