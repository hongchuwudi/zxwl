package newsApi

import (
	"encoding/json"
	"log"
	"mymod/config"
	"mymod/model/param"
	"mymod/model/sqlModel"
	"mymod/repositories/newsRepo"
	"net/http"
	"strconv"
	"time"
)

// UpdateNewsHandlers 更新资讯
func UpdateNewsHandlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPut {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求参数 - 包含id在请求体中
	var req struct {
		ID           uint64    `json:"id"`
		ProvinceID   string    `json:"province_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		Keywords     string    `json:"keywords"`
		Content      string    `json:"content"`
		VideoDetail  string    `json:"video_detail"`
		VideoType    string    `json:"video_type"`
		VideoImg     string    `json:"video_img"`
		FromSource   string    `json:"from_source"`
		NewsNum      string    `json:"news_num"`
		IsPush       int       `json:"is_push"`
		IsTop        int       `json:"is_top"`
		StyleType    string    `json:"style_type"`
		StyleURL     string    `json:"style_url"`
		Highlight    string    `json:"highlight"`
		PublishTime  time.Time `json:"publish_time"`
		CardSchoolID string    `json:"card_school_id"`
		CardLiveID   string    `json:"card_live_id"`
		NewsFrom     string    `json:"news_from"`
		ClassName    string    `json:"class_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证id参数
	if req.ID == 0 {
		response.Error = http.StatusBadRequest
		response.Message = "id参数不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	newsRepo := newsRepo.NewNewsRepository(config.GetDB())

	// 检查资讯是否存在
	existingNews, err := newsRepo.GetByNewsID(int64(req.ID))
	if err != nil || existingNews == nil {
		log.Printf("资讯不存在: ID=%d, 错误: %v", req.ID, err)
		response.Error = http.StatusNotFound
		response.Message = "资讯不存在"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 更新资讯对象
	updateData := &sqlModel.NewsInfo{
		ID:           req.ID,
		ProvinceID:   req.ProvinceID,
		Title:        req.Title,
		Description:  req.Description,
		Keywords:     req.Keywords,
		Content:      req.Content,
		VideoDetail:  req.VideoDetail,
		VideoType:    req.VideoType,
		VideoImg:     req.VideoImg,
		FromSource:   req.FromSource,
		NewsNum:      req.NewsNum,
		IsPush:       req.IsPush,
		IsTop:        req.IsTop,
		StyleType:    req.StyleType,
		StyleURL:     req.StyleURL,
		PublishTime:  req.PublishTime,
		CardSchoolID: req.CardSchoolID,
		CardLiveID:   req.CardLiveID,
		ClassName:    req.ClassName,
		UpdateTime:   time.Now(),
	}

	// 更新到数据库
	err = newsRepo.Update(updateData)
	if err != nil {
		log.Printf("更新资讯失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "更新失败: " + err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "更新成功"
	json.NewEncoder(w).Encode(response)
}

// UpdateNewsContentHandlers 更新资讯内容
func UpdateNewsContentHandlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPatch {
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

	// 解析请求参数
	var req struct {
		Content  string `json:"content"`
		Keywords string `json:"keywords"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	newsRepo := newsRepo.NewNewsRepository(config.GetDB())

	// 更新资讯内容
	err = newsRepo.UpdateContent(id, req.Content, req.Keywords)
	if err != nil {
		log.Printf("更新资讯内容失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "更新内容失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "更新内容成功"
	json.NewEncoder(w).Encode(response)
}
