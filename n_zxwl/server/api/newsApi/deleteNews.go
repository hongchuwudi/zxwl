package newsApi

import (
	"encoding/json"
	"log"
	"mymod/config"
	"mymod/model/param"
	"mymod/repositories/newsRepo"
	"net/http"
	"strconv"
)

// DeleteNewsHandlers 删除资讯
func DeleteNewsHandlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodDelete {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 从查询参数获取id
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response.Error = http.StatusBadRequest
		response.Message = "id参数不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "id格式错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	newsRepo := newsRepo.NewNewsRepository(config.GetDB())

	// 检查资讯是否存在
	existingNews, err := newsRepo.GetByNewsID(int64(id))
	if err != nil || existingNews == nil {
		response.Error = http.StatusNotFound
		response.Message = "资讯不存在"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 删除资讯
	err = newsRepo.Delete(id)
	if err != nil {
		log.Printf("删除资讯失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "删除失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "删除成功"
	json.NewEncoder(w).Encode(response)
}
