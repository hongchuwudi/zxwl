package schoolApi

import (
	"encoding/json"
	"mymod/config"
	"mymod/repositories/schoolRepo"
	"net/http"
)

// GetSchoolIDByNameHandler 根据专业名称获取专业ID
func GetSchoolIDByNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"code": 1001}

	if r.Method != http.MethodGet {
		response["code"] = 405
		response["msg"] = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 从查询参数获取专业名称
	name := r.URL.Query().Get("name")
	if name == "" {
		response["msg"] = "缺少name参数"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	schoolRepo := schoolRepo.NewUniversityRepository(config.GetDB())

	// 调用repository方法
	id, err := schoolRepo.GetIDByName(name)
	if err != nil {
		response["msg"] = "未找到该学校"
		json.NewEncoder(w).Encode(response)
		return
	}

	response["code"] = 0
	response["data"] = map[string]interface{}{
		"id":   id,
		"name": name,
	}
	json.NewEncoder(w).Encode(response)
}
