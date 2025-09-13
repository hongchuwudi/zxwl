package schoolApi

import (
	"encoding/json"
	"log"
	"mymod/config"
	"mymod/repositories/schoolRepo"
	"net/http"
	"strconv"
	"strings"
)

// GetSchoolProfileByIDHandler 根据学校ID获取学校简介（GET）
func GetSchoolProfileByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"code": 1001}

	if r.Method != http.MethodGet {
		response["code"] = 405
		response["msg"] = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 从URL路径中提取school_id
	// 路径格式应该是 /schools/profiles/{school_id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response["msg"] = "无效的URL路径"
		json.NewEncoder(w).Encode(response)
		return
	}

	schoolIDStr := pathParts[3] // /schools/profiles/123 中的 123

	if schoolIDStr == "" {
		response["msg"] = "缺少school_id参数"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换school_id为整数
	schoolID, err := strconv.Atoi(schoolIDStr)
	if err != nil {
		response["msg"] = "无效的school_id参数"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化数据库连接和repository
	universityRepo := schoolRepo.NewUniversityRepository(config.GetDB())

	// 获取学校详细信息
	universityDetailResponse, err := universityRepo.GetDetailByID(schoolID)
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		response["msg"] = "服务器内部错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	if universityDetailResponse == nil {
		response["msg"] = "未找到相关院校信息"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构建响应数据
	profile := map[string]interface{}{
		"universities_detail":  universityDetailResponse.UniversitiesDetail,
		"colleges_departments": universityDetailResponse.CollegesDepartments,
		"discipline_rankings":  universityDetailResponse.DisciplineRankings,
		"dual_class_subjects":  universityDetailResponse.DualClassSubjects,
		"videos":               universityDetailResponse.Videos,
		"special_programs":     universityDetailResponse.SpecialPrograms,
		"admission_scores":     universityDetailResponse.AdmissionScores,
	}

	response["code"] = 0
	response["data"] = profile
	json.NewEncoder(w).Encode(response)
}
