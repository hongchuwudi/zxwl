// api/score_handler.go
package scoreApi

import (
	"encoding/json"
	"log"
	"mymod/config"
	"mymod/model/param"
	scoreRepo "mymod/repositories/scoreRepo"
	"mymod/service"
	"net/http"
	"strconv"
)

type ScoreHandler struct {
	scoreService *service.ScoreService
}

func NewScoreHandler() *ScoreHandler {
	scoreRepo := scoreRepo.NewScoreRepository(config.GetDB())
	scoreService := service.NewScoreService(scoreRepo)
	return &ScoreHandler{scoreService: scoreService}
}

// GetScoreDataHandler 获取分数数据
func (h *ScoreHandler) GetScoreDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodGet {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析查询参数
	queryParams := r.URL.Query()
	provinceID, _ := strconv.Atoi(queryParams.Get("province_id"))
	typeID, _ := strconv.Atoi(queryParams.Get("type_id"))
	year, _ := strconv.Atoi(queryParams.Get("year"))
	var batchID int
	if provinceID == 11 || provinceID == 12 || provinceID == 13 {
		batchID, _ = strconv.Atoi(queryParams.Get("batch"))
	} else {
		batchID = 3
	}

	if provinceID == 0 || typeID == 0 || year == 0 {
		response.Error = http.StatusBadRequest
		response.Message = "省份ID、考试类型和年份不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取分数数据
	scoreData, err := h.scoreService.GetScoreData(provinceID, typeID, year, batchID)
	if err != nil {
		log.Printf("获取分数数据失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "获取分数数据失败: " + err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Error = 0
	response.Message = "获取成功"
	response.Data = scoreData
	json.NewEncoder(w).Encode(response)
}

// GetAvailableYearsHandler 获取可用年份
func (h *ScoreHandler) GetAvailableYearsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodGet {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析查询参数
	queryParams := r.URL.Query()
	provinceID, _ := strconv.Atoi(queryParams.Get("province_id"))
	typeID, _ := strconv.Atoi(queryParams.Get("type_id"))

	if provinceID == 0 || typeID == 0 {
		response.Error = http.StatusBadRequest
		response.Message = "省份ID和考试类型不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取可用年份
	years, err := h.scoreService.GetAvailableYears(provinceID, typeID)
	if err != nil {
		log.Printf("获取可用年份失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "获取可用年份失败: " + err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Error = 0
	response.Message = "获取成功"
	response.Data = years
	json.NewEncoder(w).Encode(response)
}
