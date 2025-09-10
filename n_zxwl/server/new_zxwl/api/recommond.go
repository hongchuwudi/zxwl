package api

import (
	"encoding/json"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/model/sqlModel"
	"mymod/new_zxwl/service"
	"net/http"
)

// RecommendationHandlers 志愿推荐HTTP处理函数
func RecommendationHandlers(w http.ResponseWriter, r *http.Request) {
	// CORS设置
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求
	var req sqlModel.UserPriInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	// 基本验证
	if req.Year == 0 || req.ProvinceName == "" || req.Score == 0 {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 获取推荐结果
	db := config.GetDB()
	repo := service.NewRecommendationRepository(db)
	_, recommendation, err := repo.GetRecommendations(req)
	if err != nil {
		http.Error(w, "推荐失败", http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": recommendation,
	})
}
