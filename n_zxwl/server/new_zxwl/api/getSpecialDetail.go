package api

import (
	"encoding/json"
	"log"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/repositories"
	"net/http"
	"strconv"
	"strings"
)

// GetSpecialProfileByIDHandler 根据专业ID获取专业详细信息（GET）
func GetSpecialProfileByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"code": 1001}

	if r.Method != http.MethodGet {
		response["code"] = 405
		response["msg"] = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 从URL路径中提取special_id
	// 路径格式: /specials/profiles/{special_id}
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 3 || parts[0] != "specials" || parts[1] != "profiles" {
		response["msg"] = "无效的URL路径"
		json.NewEncoder(w).Encode(response)
		return
	}

	specialIDStr := parts[2]
	if specialIDStr == "" {
		response["msg"] = "缺少special_id参数"
		json.NewEncoder(w).Encode(response)
		return
	}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response["msg"] = "无效的URL路径"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换special_id为整数
	specialID, err := strconv.Atoi(specialIDStr)
	if err != nil {
		response["msg"] = "无效的special_id参数"
		json.NewEncoder(w).Encode(response)
		return
	}
	repository := repositories.NewSpecialRepository(config.GetDB())
	specialDetailResponse, err := repository.GetDetailByID(specialID)

	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		response["msg"] = "服务器内部错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构建响应数据
	profile := map[string]interface{}{
		"special_detail":          specialDetailResponse.SpecialDetail,
		"special_contents":        specialDetailResponse.SpecialContents,
		"employment_rates":        specialDetailResponse.EmploymentRates,
		"famous_schools":          specialDetailResponse.FamousSchools,
		"videos":                  specialDetailResponse.Videos,
		"impression_tags":         specialDetailResponse.ImpressionTags,
		"job_distributions":       specialDetailResponse.JobDistributions,
		"salary_data":             specialDetailResponse.SalaryData,
		"university_special_info": specialDetailResponse.UniversitySpecialInfo,
	}

	response["code"] = 0
	response["data"] = profile
	json.NewEncoder(w).Encode(response)
}
