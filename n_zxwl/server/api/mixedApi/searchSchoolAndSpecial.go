package mixedApi

import (
	"encoding/json"
	"mymod/config"
	"mymod/model/param"
	repositories2 "mymod/repositories/schoolRepo"
	specialRepo2 "mymod/repositories/specialRepo"
	"net/http"
)

// SearchSchoolAndSpecial 根据关键词来查询专业和学校
func SearchSchoolAndSpecial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"code": 1001}

	// 异常请求方法
	if r.Method != http.MethodGet {
		response["code"] = 405
		response["msg"] = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取查询参数
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		response["code"] = 400
		response["msg"] = "关键词不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取数据库连接
	db := config.GetDB()

	// 创建repository实例
	specialRepo := specialRepo2.NewSpecialRepository(db)
	universityRepo := repositories2.NewUniversityRepository(db)

	// 查询专业信息（返回数组）
	specials, err := specialRepo.GetByName(keyword)
	if err != nil {
		// 如果没找到专业，不返回错误，继续查询学校
		specials = nil
	}

	// 查询学校信息（返回数组）
	universities, err := universityRepo.GetByName(keyword)
	if err != nil {
		// 如果没找到学校，不返回错误
		universities = nil
	}

	// 构建响应数据
	var professionalItems []param.ProfessionalItemResponse
	var schoolProfiles []param.UniversitiesSimple

	// 如果有找到专业，转换为响应格式
	if specials != nil && len(specials) > 0 {
		for _, special := range specials {
			professionalItem := param.ProfessionalItemResponse{
				ID:             special.ID,
				Name:           special.Name,
				Level1Name:     special.Level1Name,
				Level2Name:     special.Level2Name,
				Level3Name:     special.Level3Name,
				Salaryavg:      special.AvgSalary,
				LimitYear:      special.LimitYear,
				Code:           special.Code,
				Degree:         special.Degree,
				EmploymentRate: special.EmploymentRate,
				TopIndustry:    special.TopIndustry,
				TopPosition:    special.TopPosition,
			}
			professionalItems = append(professionalItems, professionalItem)
		}
	}

	// 如果有找到学校，转换为响应格式
	if universities != nil && len(universities) > 0 {
		for _, university := range universities {
			schoolProfile := param.UniversitiesSimple{
				ID:            university.ID,
				Name:          university.Name,
				Motto:         university.Motto,
				ProvinceID:    university.ProvinceID,
				Postcode:      university.Postcode,
				Belong:        university.Belong,
				NatureName:    university.NatureName,
				TypeName:      university.TypeName,
				LevelName:     university.LevelName,
				CreateDate:    university.CreateDate,
				F211:          university.F211,
				F985:          university.F985,
				LogoURL:       university.LogoURL,
				DualClassName: university.DualClassName,
			}
			schoolProfiles = append(schoolProfiles, schoolProfile)
		}
	}

	// 构建最终响应
	result := param.SchAndSpeResponse{
		ProfessionalItemResArr: professionalItems,
		SchoolProfileRes:       schoolProfiles,
	}

	response["code"] = 200
	response["msg"] = "成功"
	response["data"] = result

	// 返回响应
	json.NewEncoder(w).Encode(response)
}
