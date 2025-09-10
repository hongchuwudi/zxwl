package api

import (
	"encoding/json"
	"log"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/model/param"
	"mymod/new_zxwl/repositories"
	"net/http"
)

// NewsQueryHandlers 资讯分页查询
func NewsQueryHandlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求参数
	var req param.NewsQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 50
	}

	// 初始化repository
	newsRepo := repositories.NewNewsRepository(config.GetDB())

	// 构建筛选条件
	filters := make(map[string]interface{})
	if req.Title != "" {
		filters["title"] = req.Title
	}
	if req.Keywords != "" {
		filters["keywords"] = req.Keywords
	}
	if req.FromSource != "" {
		filters["from_source"] = req.FromSource
	}

	// 查询资讯列表
	newsList, total, err := newsRepo.FindByCondition(filters, req.Page, req.PageSize)
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "查询失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换为VO
	var results []param.NewsVO
	for _, news := range newsList {
		item := param.NewsVO{
			ID:          news.ID,
			Title:       news.Title,
			Description: news.Description,
			Keywords:    news.Keywords,
			FromSource:  news.FromSource,
			ClassName:   news.ClassName,
			PublishTime: news.PublishTime,
			CreateTime:  news.CreateTime,
			StyleURL:    news.StyleURL,
		}
		results = append(results, item)
	}

	// 构造分页响应数据
	paginationData := map[string]interface{}{
		"list":      results,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
		"pages":     (total + int64(req.PageSize) - 1) / int64(req.PageSize),
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "成功"
	response.Data = paginationData
	json.NewEncoder(w).Encode(response)
}
