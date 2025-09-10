package api

import (
	"encoding/json"
	"fmt"
	"log"
	"mymod/Const"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/repositories"
	"net/http"
	"strconv"
	"strings"
)

// ProfessionalQueryHandlers 专业信息查询
func ProfessionalQueryHandlers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.ProfessionalQueryResponse{Code: 1001}

	if r.Method != http.MethodPost {
		response.Code = http.StatusMethodNotAllowed
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求参数
	var req Const.ProfessionalQueryRequest
	fmt.Printf("%+v\n", req)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		response.Code = http.StatusBadRequest

		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	specialRepo := repositories.NewSpecialRepository(config.GetDB())

	// 构建筛选条件
	filters := make(map[string]interface{})
	if req.Level3Name != "" {
		filters["level3_name"] = req.Level3Name
	}
	if req.Type != "" {
		filters["type"] = req.Type
	}
	if req.TypeDetail != "" {
		filters["type_detail"] = req.TypeDetail
	}
	if req.Degree != "" {
		filters["degree"] = req.Degree
	}
	if req.Keyword != "" {
		filters["keyword"] = req.Keyword
	}
	if req.MinSalary > 0 {
		filters["min_salary"] = req.MinSalary
	}
	if req.MaxSalary > 0 {
		filters["max_salary"] = req.MaxSalary
	}
	if req.MinEmploymentRate > 0 {
		filters["min_employment_rate"] = req.MinEmploymentRate
	}

	// 查询专业列表
	specials, err := specialRepo.GetProfessionalList(filters)
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		response.Code = http.StatusInternalServerError
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
		return
	}

	// 转换为响应格式
	var results []Const.ProfessionalItemResponse
	for _, special := range specials {

		// 解析男女比例
		boyRate, girlRate := parseGenderRatio(special.GenderRatio)

		// 获取5年薪资
		fiveYearSalary, _ := specialRepo.GetFiveYearSalary(strconv.Itoa(special.ID))
		if fiveYearSalary == 0 {
			fiveYearSalary = special.AvgSalary // 如果没有5年数据，使用平均薪资
			special.AvgSalary = fiveYearSalary
		}

		item := Const.ProfessionalItemResponse{
			ID:             special.ID,
			Name:           special.Name,
			Level1Name:     special.Level1Name,
			Salaryavg:      special.AvgSalary,
			LimitYear:      special.LimitYear,
			FiveSalaryavg:  fiveYearSalary,
			Level2Name:     special.Level2Name,
			BoyRate:        boyRate,
			GirlRate:       girlRate,
			Level3Name:     special.Level3Name,
			Code:           special.Code,
			Degree:         special.Degree,
			EmploymentRate: special.EmploymentRate,
			TopIndustry:    special.TopIndustry,
			TopPosition:    special.TopPosition,
		}
		results = append(results, item)
	}

	// 构造成功响应
	response.Code = 0
	response.Data = results
	json.NewEncoder(w).Encode(response)
}

// parseGenderRatio 解析男女比例字符串
func parseGenderRatio(ratio string) (string, string) {
	parts := strings.Split(ratio, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]

	}
	//return "50.00", "50.00"
	return "60.00", "40.00"
}
