package schoolApi

import (
	"encoding/json"
	"fmt"
	"mymod/config"
	"mymod/model/param"
	"mymod/repositories/schoolRepo"
	"net/http"
	"strconv"
)

// CollegesHandlers 学校分页条件查询处理器
func CollegesHandlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	repo := schoolRepo.NewUniversityRepository(config.GetDB())
	params := parseQueryParams(r.URL.Query())
	universities, total, err := repo.QuerySimpleColleges(params)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Database query failed: %v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(param.CollegeListResponse{
		Total:        total,
		Page:         params.Page,
		PageSize:     params.PageSize,
		TotalPages:   (total + params.PageSize - 1) / params.PageSize,
		Universities: universities,
	})
}

// parseQueryParams 解析查询参数
func parseQueryParams(query map[string][]string) *param.CollegeQueryParams {
	params := &param.CollegeQueryParams{Page: 1, PageSize: 20}

	if vals := query["page"]; len(vals) > 0 {
		if page, _ := strconv.Atoi(vals[0]); page > 0 {
			params.Page = page
		}
	}
	if vals := query["page_size"]; len(vals) > 0 {
		if pageSize, _ := strconv.Atoi(vals[0]); pageSize > 0 {
			if pageSize > 100 {
				pageSize = 100
			}
			params.PageSize = pageSize
		}
	}

	params.Name = getVal(query, "name")
	params.ProvinceID, _ = strconv.Atoi(getVal(query, "province_id"))
	params.TypeName = getVal(query, "type_name")
	params.LevelName = getVal(query, "level_name")
	params.Is211, _ = strconv.Atoi(getVal(query, "is_211"))
	params.Is985, _ = strconv.Atoi(getVal(query, "is_985"))
	params.DualClass = getVal(query, "dual_class")
	params.OrderBy = getVal(query, "order_by")
	params.OrderDesc = getVal(query, "order_desc") == "true" || getVal(query, "order_desc") == "1"

	return params
}

// getVal 辅助函数 获取查询参数的值
func getVal(query map[string][]string, key string) string {
	if vals := query[key]; len(vals) > 0 {
		return vals[0]
	}
	return ""
}
