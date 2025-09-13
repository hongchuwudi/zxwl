package newsApi

import (
	"encoding/json"
	"log"
	"mymod/config"
	"mymod/model/param"
	"mymod/model/sqlModel"
	"mymod/repositories/newsRepo"
	"net/http"
	"time"
)

// CreateNewsHandlers 创建资讯
func CreateNewsHandlers(w http.ResponseWriter, r *http.Request) {
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
		ProvinceID     string    `json:"province_id"`
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		Keywords       string    `json:"keywords"`
		Content        string    `json:"content"`
		VideoDetail    string    `json:"video_detail"`
		VideoType      string    `json:"video_type"`
		VideoImg       string    `json:"video_img"`
		FromSource     string    `json:"from_source"`
		NewsNum        string    `json:"news_num"`
		IsPush         int       `json:"is_push"`
		IsTop          int       `json:"is_top"`
		StyleType      string    `json:"style_type"`
		StyleURL       string    `json:"style_url"`
		Highlight      string    `json:"highlight"`
		PublishTime    time.Time `json:"publish_time"`
		CardSchoolID   string    `json:"card_school_id"`
		CardLiveID     string    `json:"card_live_id"`
		NewsFrom       string    `json:"news_from"`
		ClassName      string    `json:"class_name"`
		PublisherEmail string    `json:"publisher_email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证必填字段
	if req.Title == "" {
		response.Error = http.StatusBadRequest
		response.Message = "标题为必填字段"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	newsRepo := newsRepo.NewNewsRepository(config.GetDB())

	// 创建资讯对象
	now := time.Now()
	news := &sqlModel.NewsInfo{
		ProvinceID:     req.ProvinceID,
		Title:          req.Title,
		Description:    req.Description,
		Keywords:       req.Keywords,
		Content:        req.Content,
		VideoDetail:    req.VideoDetail,
		VideoType:      req.VideoType,
		VideoImg:       req.VideoImg,
		FromSource:     req.FromSource,
		NewsNum:        req.NewsNum,
		IsPush:         req.IsPush,
		IsTop:          req.IsTop,
		StyleType:      req.StyleType,
		StyleURL:       req.StyleURL,
		PublishTime:    now,
		CardSchoolID:   req.CardSchoolID,
		CardLiveID:     req.CardLiveID,
		ClassName:      req.ClassName,
		AddTime:        now,
		CreateTime:     now,
		UpdateTime:     now,
		PublisherEmail: req.PublisherEmail,
	}

	// 保存到数据库
	err := newsRepo.Create(news)
	if err != nil {
		log.Printf("创建资讯失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "创建失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "创建成功"
	response.Data = map[string]interface{}{
		"id": news.ID,
	}
	json.NewEncoder(w).Encode(response)
}

// UpdateNewsCountHandler 更新资讯统计数量
func UpdateNewsCountHandler(w http.ResponseWriter, r *http.Request) {
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
		Field  string `json:"field"`
		Action string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证参数
	if req.NewsID == 0 || req.Field == "" || req.Action == "" {
		response.Error = http.StatusBadRequest
		response.Message = "缺少必要参数"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	newsRepo := newsRepo.NewNewsRepository(config.GetDB())

	// 更新数量
	err := newsRepo.UpdateCount(req.NewsID, req.Field, req.Action)
	if err != nil {
		log.Printf("更新资讯统计失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "数据更新失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取更新后的数量
	counts, err := newsRepo.GetCount(req.NewsID)
	if err != nil {
		log.Printf("获取更新后数量失败: %v", err)
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "更新成功"
	response.Data = map[string]interface{}{
		"news_id": req.NewsID,
		"field":   req.Field,
		"action":  req.Action,
		"counts":  counts,
	}
	json.NewEncoder(w).Encode(response)
}
