package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"mymod/new_zxwl/model/sqlModel"
)

type SpecialRepository struct {
	db *gorm.DB
}

func NewSpecialRepository(db *gorm.DB) *SpecialRepository {
	return &SpecialRepository{db: db}
}

// GetDetailByID 获取专业完整详细信息（包含所有关联表数据）
func (r *SpecialRepository) GetDetailByID(id int) (*sqlModel.SpecialDetailResponse, error) {
	var response sqlModel.SpecialDetailResponse

	// 查询基本信息
	detail, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	response.SpecialDetail = *detail

	// 查询专业大文本内容
	var contents []sqlModel.SpecialContent
	if err := r.db.Table("special_content").Where("special_id = ?", id).Find(&contents).Error; err == nil {
		response.SpecialContents = contents
	}

	// 查询专业就业率信息
	var employmentRates []sqlModel.SpecialEmploymentRate
	if err := r.db.Table("special_employment_rate").Where("special_id = ?", id).Order("year DESC").Find(&employmentRates).Error; err == nil {
		response.EmploymentRates = employmentRates
	}

	// 查询专业名校示例
	var famousSchools []sqlModel.SpecialFamousSchool
	if err := r.db.Table("special_famous_school").Where("special_id = ?", id).Order("school_name").Find(&famousSchools).Error; err == nil {
		response.FamousSchools = famousSchools
	}

	// 查询专业视频信息
	var videos []sqlModel.SpecialVideo
	if err := r.db.Table("special_video").Where("special_id = ?", id).Order("create_time DESC").Find(&videos).Error; err == nil {
		response.Videos = videos
	}

	// 查询专业印象标签
	var impressionTags []sqlModel.SpecialImpressionTag
	if err := r.db.Table("special_impression_tag").Where("special_id = ?", id).Order("keyword").Find(&impressionTags).Error; err == nil {
		response.ImpressionTags = impressionTags
	}

	// 查询专业就业分布
	var jobDistributions []sqlModel.SpecialJobDistribution
	if err := r.db.Table("special_job_distribution").Where("special_id = ?", id).Order("sort").Find(&jobDistributions).Error; err == nil {
		response.JobDistributions = jobDistributions
	}

	// 查询专业薪资数据
	var salaryData []sqlModel.SpecialSalaryData
	if err := r.db.Table("special_salary_data").Where("special_id = ? and salary_value <> 0", id).Order("salary_year DESC, salary_type").Find(&salaryData).Error; err == nil {
		response.SalaryData = salaryData
	}

	// 查询开设该专业的院校信息
	var universitySpecialInfo []sqlModel.UniversitiesSpecialPrograms
	if err := r.db.Table("universities_special_programs").Where("special_id = ?", id).Order("school_id").Find(&universitySpecialInfo).Error; err == nil {
		response.UniversitySpecialInfo = universitySpecialInfo
	}

	return &response, nil
}

// GetByID 根据ID获取专业基础信息
func (r *SpecialRepository) GetByID(id int) (*sqlModel.SpecialDetail, error) {
	var special sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("id = ?", id).First(&special).Error
	if err != nil {
		return nil, err
	}
	return &special, nil
}

// GetByCode 根据专业代码获取专业信息
func (r *SpecialRepository) GetByCode(code string) (*sqlModel.SpecialDetail, error) {
	var special sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("code = ?", code).First(&special).Error
	if err != nil {
		return nil, err
	}
	return &special, nil
}

// GetByName 根据专业名称获取专业信息
func (r *SpecialRepository) GetByName(name string) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("name LIKE ?", "%"+name+"%").Find(&specials).Error
	if err != nil {
		return nil, err
	}
	if len(specials) == 0 {
		return nil, fmt.Errorf("未找到相关专业")
	}
	return specials, nil
}

// GetIDByName 根据专业名称获取专业ID
func (r *SpecialRepository) GetIDByName(name string) (int, error) {
	var id int
	err := r.db.Raw("SELECT id FROM special_detail WHERE name = ? LIMIT 1", name).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("未找到专业: %s", name)
	}
	return id, nil
}

// GetContentsBySpecialID 获取专业的大文本内容
func (r *SpecialRepository) GetContentsBySpecialID(specialID int) ([]sqlModel.SpecialContent, error) {
	var contents []sqlModel.SpecialContent
	err := r.db.Table("special_content").Where("special_id = ?", specialID).Find(&contents).Error
	return contents, err
}

// GetEmploymentRatesBySpecialID 获取专业的就业率信息
func (r *SpecialRepository) GetEmploymentRatesBySpecialID(specialID string) ([]sqlModel.SpecialEmploymentRate, error) {
	var rates []sqlModel.SpecialEmploymentRate
	err := r.db.Table("special_employment_rate").Where("special_id = ?", specialID).Order("year DESC").Find(&rates).Error
	return rates, err
}

// GetFamousSchoolsBySpecialID 获取专业的名校示例
func (r *SpecialRepository) GetFamousSchoolsBySpecialID(specialID string) ([]sqlModel.SpecialFamousSchool, error) {
	var schools []sqlModel.SpecialFamousSchool
	err := r.db.Table("special_famous_school").Where("special_id = ?", specialID).Order("school_name").Find(&schools).Error
	return schools, err
}

// GetVideosBySpecialID 获取专业的视频信息
func (r *SpecialRepository) GetVideosBySpecialID(specialID int) ([]sqlModel.SpecialVideo, error) {
	var videos []sqlModel.SpecialVideo
	err := r.db.Table("special_video").Where("special_id = ?", specialID).Order("create_time DESC").Find(&videos).Error
	return videos, err
}

// GetImpressionTagsBySpecialID 获取专业的印象标签
func (r *SpecialRepository) GetImpressionTagsBySpecialID(specialID string) ([]sqlModel.SpecialImpressionTag, error) {
	var tags []sqlModel.SpecialImpressionTag
	err := r.db.Table("special_impression_tag").Where("special_id = ?", specialID).Order("keyword").Find(&tags).Error
	return tags, err
}

// GetJobDistributionsBySpecialID 获取专业的就业分布
func (r *SpecialRepository) GetJobDistributionsBySpecialID(specialID string) ([]sqlModel.SpecialJobDistribution, error) {
	var distributions []sqlModel.SpecialJobDistribution
	err := r.db.Table("special_job_distribution").Where("special_id = ?", specialID).Order("sort").Find(&distributions).Error
	return distributions, err
}

// GetSalaryDataBySpecialID 获取专业的薪资数据
func (r *SpecialRepository) GetSalaryDataBySpecialID(specialID string) ([]sqlModel.SpecialSalaryData, error) {
	var salaryData []sqlModel.SpecialSalaryData
	err := r.db.Table("special_salary_data").Where("special_id = ?", specialID).Order("salary_year DESC, salary_type").Find(&salaryData).Error
	return salaryData, err
}

// GetUniversitySpecialInfoBySpecialID 获取开设该专业的院校信息
func (r *SpecialRepository) GetUniversitySpecialInfoBySpecialID(specialID int) ([]sqlModel.UniversitiesSpecialPrograms, error) {
	var programs []sqlModel.UniversitiesSpecialPrograms
	err := r.db.Table("universities_special_programs").Where("special_id = ?", specialID).Order("school_id").Find(&programs).Error
	return programs, err
}

// GetByInterest 根据兴趣关键词获取相关专业
func (r *SpecialRepository) GetByInterest(keyword string) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("name LIKE ? OR direction LIKE ? OR career_prospects LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Find(&specials).Error
	return specials, err
}

// GetHotSpecials 获取热门专业
func (r *SpecialRepository) GetHotSpecials(limit int) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Order("total_views DESC").Limit(limit).Find(&specials).Error
	return specials, err
}

// GetByLevel 根据专业层次获取专业列表
func (r *SpecialRepository) GetByLevel(level string) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("degree = ?", level).Find(&specials).Error
	return specials, err
}

// GetByType 根据专业类型获取专业列表
func (r *SpecialRepository) GetByType(specialType string) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("type = ?", specialType).Find(&specials).Error
	return specials, err
}

// GetByIDs 批量获取专业基础信息
func (r *SpecialRepository) GetByIDs(ids []int) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("id IN ?", ids).Find(&specials).Error
	return specials, err
}

// GetBySalaryRange 根据薪资范围获取专业列表
func (r *SpecialRepository) GetBySalaryRange(minSalary, maxSalary int) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("avg_salary BETWEEN ? AND ?", minSalary, maxSalary).Find(&specials).Error
	return specials, err
}

// GetByEmploymentRate 根据就业率获取专业列表
func (r *SpecialRepository) GetByEmploymentRate(minRate float64) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").Where("employment_rate >= ?", minRate).Find(&specials).Error
	return specials, err
}

// GetAllSpecialDetails 获取所有专业的基本信息
func (r *SpecialRepository) GetAllSpecialDetails() ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail
	err := r.db.Table("special_detail").
		Select("id", "code", "name", "degree", "direction", "type", "type_detail",
			"limit_year", "level1", "level2", "level3", "employment_rate",
			"avg_salary", "top_industry", "top_position", "top_area",
			"monthly_views", "total_views", "subject_requirements",
			"gender_ratio", "celebrities", "courses", "content",
			"career_prospects", "description", "update_time", "create_time").
		Order("name ASC").
		Find(&specials).Error
	if err != nil {
		return nil, err
	}
	return specials, nil
}

// GetAllSpecialDetailsWithPagination 分页获取所有专业的基本信息（支持筛选）
func (r *SpecialRepository) GetAllSpecialDetailsWithPagination(page, pageSize int, filters map[string]interface{}) ([]sqlModel.SpecialDetail, int64, error) {
	var specials []sqlModel.SpecialDetail
	var total int64

	// 构建基础查询
	query := r.db.Table("special_detail")

	// 应用筛选条件
	for key, value := range filters {
		switch key {
		case "degree":
			query = query.Where("degree = ?", value)
		case "type":
			query = query.Where("type = ?", value)
		case "level3":
			query = query.Where("level3 = ?", value)
		case "min_salary":
			query = query.Where("avg_salary >= ?", value)
		case "max_salary":
			query = query.Where("avg_salary <= ?", value)
		case "min_employment_rate":
			query = query.Where("employment_rate >= ?", value)
		case "keyword":
			query = query.Where("name LIKE ? OR direction LIKE ? OR career_prospects LIKE ?",
				"%"+value.(string)+"%", "%"+value.(string)+"%", "%"+value.(string)+"%")
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.
		Select("id", "code", "name", "degree", "direction", "type", "type_detail",
			"limit_year", "level1", "level2", "level3", "employment_rate",
			"avg_salary", "top_industry", "top_position", "top_area",
			"monthly_views", "total_views", "subject_requirements",
			"gender_ratio", "celebrities", "courses", "content",
			"career_prospects", "description", "update_time", "create_time").
		Order("name ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&specials).Error
	if err != nil {
		return nil, 0, err
	}

	return specials, total, nil
}

// GetAllSpecialDetailsWithFilter 根据条件筛选获取专业基本信息
func (r *SpecialRepository) GetAllSpecialDetailsWithFilter(filters map[string]interface{}) ([]sqlModel.SpecialDetail, error) {
	query := r.db.Table("special_detail")

	for key, value := range filters {
		switch key {
		case "degree":
			query = query.Where("degree = ?", value)
		case "type":
			query = query.Where("type = ?", value)
		case "level1":
			query = query.Where("level1 = ?", value)
		case "level2":
			query = query.Where("level2 = ?", value)
		case "level3":
			query = query.Where("level3 = ?", value)
		case "min_salary":
			query = query.Where("avg_salary >= ?", value)
		case "max_salary":
			query = query.Where("avg_salary <= ?", value)
		case "min_employment_rate":
			query = query.Where("employment_rate >= ?", value)
		case "keyword":
			query = query.Where("name LIKE ? OR direction LIKE ? OR career_prospects LIKE ?",
				"%"+value.(string)+"%", "%"+value.(string)+"%", "%"+value.(string)+"%")
		}
	}

	var specials []sqlModel.SpecialDetail
	err := query.
		Select("id", "code", "name", "degree", "direction", "type", "type_detail",
			"limit_year", "level1", "level2", "level3", "employment_rate",
			"avg_salary", "top_industry", "top_position", "top_area",
			"monthly_views", "total_views", "subject_requirements",
			"gender_ratio", "celebrities", "courses", "content",
			"career_prospects", "description", "update_time", "create_time").
		Order("name ASC").
		Find(&specials).Error

	return specials, err
}

// GetProfessionalList 获取专业列表（支持筛选）
func (r *SpecialRepository) GetProfessionalList(filters map[string]interface{}) ([]sqlModel.SpecialDetail, error) {
	var specials []sqlModel.SpecialDetail

	query := r.db.Table("special_detail")

	// 应用筛选条件
	for key, value := range filters {
		switch key {
		case "level3_name":
			query = query.Where("level3_name = ?", value)
		case "type":
			query = query.Where("type = ?", value)
		case "type_detail":
			query = query.Where("type_detail = ?", value)
		case "degree":
			query = query.Where("degree = ?", value)
		case "keyword":
			query = query.Where("name LIKE ? OR direction LIKE ?",
				"%"+value.(string)+"%", "%"+value.(string)+"%")
		case "min_salary":
			query = query.Where("avg_salary >= ?", value)
		case "max_salary":
			query = query.Where("avg_salary <= ?", value)
		case "min_employment_rate":
			query = query.Where("employment_rate >= ?", value)
		}
	}

	// 调试：打印最终执行的SQL
	sql := query.
		Select("id", "code", "name", "degree", "limit_year",
			"level1_name", "level2_name", "level3_name",
			"employment_rate", "avg_salary", "top_industry",
			"top_position", "gender_ratio").
		Order("name ASC").
		ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.Find(&specials)
		})
	fmt.Printf("专业大全页面执行的SQL: %s\n", sql)

	err := query.
		Select("id", "code", "name", "degree", "limit_year",
			"level1_name", "level2_name", "level3_name",
			"employment_rate", "avg_salary", "top_industry",
			"top_position", "gender_ratio").
		Order("name ASC").
		Find(&specials).Error

	if err != nil {
		return nil, err
	}

	return specials, nil
}

// GetFiveYearSalary 获取专业5年薪资数据
func (r *SpecialRepository) GetFiveYearSalary(specialID string) (int, error) {
	var salaryData sqlModel.SpecialSalaryData
	err := r.db.Table("special_salary_data").
		Where("special_id = ? AND salary_type = ? AND salary_year = ?",
			specialID, 1, 4). // 1-专业薪资, 4-毕业5年
		First(&salaryData).Error

	if err != nil {
		return 0, err
	}
	return salaryData.SalaryValue, nil
}
