package api

import (
	"encoding/json"
	"log"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/model/param"
	"mymod/new_zxwl/repositories"
	"net/http"
	"strconv"
)

// GetNewsByIDHandlers 按NewsID查询资讯
func GetNewsByIDHandlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodGet {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 从URL路径参数获取newsID
	newsID := r.URL.Query().Get("id")
	if newsID == "" {
		response.Error = http.StatusBadRequest
		response.Message = "newsId参数不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}
	id, _ := strconv.ParseInt(newsID, 10, 64) // 将newsID转换成long

	// 初始化repository
	newsRepo := repositories.NewNewsRepository(config.GetDB())

	// 查询资讯详情
	news, err := newsRepo.GetByNewsID(id)
	if err != nil {
		log.Printf("查询资讯详情失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "查询失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	if news == nil {
		response.Error = http.StatusNotFound
		response.Message = "资讯不存在"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造详情响应
	detailResponse := map[string]interface{}{
		"id":              news.ID,
		"title":           news.Title,
		"description":     news.Description,
		"keywords":        news.Keywords,
		"content":         news.Content,
		"from_source":     news.FromSource,
		"class_name":      news.ClassName,
		"publish_time":    news.PublishTime,
		"create_time":     news.CreateTime,
		"update_time":     news.UpdateTime,
		"style_url":       news.StyleURL,
		"video_detail":    news.VideoDetail,
		"video_img":       news.VideoImg,
		"video_type":      news.VideoType,
		"news_num":        news.NewsNum,
		"is_push":         news.IsPush,
		"is_top":          news.IsTop,
		"province_id":     news.ProvinceID,
		"card_school_id":  news.CardSchoolID,
		"card_live_id":    news.CardLiveID,
		"add_time":        news.AddTime,
		"publisher_email": news.PublisherEmail,
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "成功"
	response.Data = detailResponse
	json.NewEncoder(w).Encode(response)
}

// GetNewsCountHandler 获取资讯统计数量
func GetNewsCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求参数
	var req struct {
		NewsID uint64 `json:"news_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证参数
	if req.NewsID == 0 {
		response.Error = http.StatusBadRequest
		response.Message = "缺少资讯ID参数"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	newsRepo := repositories.NewNewsRepository(config.GetDB())

	// 获取数量
	counts, err := newsRepo.GetCount(req.NewsID)
	if err != nil {
		log.Printf("获取资讯统计失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "获取失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "获取成功"
	response.Data = map[string]interface{}{
		"news_id": req.NewsID,
		"counts":  counts,
	}
	json.NewEncoder(w).Encode(response)
}
