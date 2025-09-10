package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"mymod/new_zxwl/model/param"
	"mymod/new_zxwl/model/sqlModel"
)

type UniversityRepository struct {
	db *gorm.DB
}

func NewUniversityRepository(db *gorm.DB) *UniversityRepository {
	return &UniversityRepository{db: db}
}

// GetDetailByID 获取院校完整详细信息（包含所有关联表数据）
func (r *UniversityRepository) GetDetailByID(id int) (*sqlModel.UniversityDetailResponse, error) {
	var response sqlModel.UniversityDetailResponse

	// 查询基本信息
	detail, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	response.UniversitiesDetail = *detail

	// 查询学院与系信息
	var colleges []sqlModel.UniversitiesCollegesDepartments
	if err := r.db.Where("school_id = ?", id).Find(&colleges).Error; err == nil {
		response.CollegesDepartments = colleges
	}

	// 查询学科评估信息
	var rankings []sqlModel.UniversitiesDisciplineRankings
	if err := r.db.Where("school_id = ?", id).Find(&rankings).Error; err == nil {
		response.DisciplineRankings = rankings
	}

	// 查询双一流学科信息
	var dualClass []sqlModel.UniversitiesDualClassSubjects
	if err := r.db.Where("school_id = ?", id).Find(&dualClass).Error; err == nil {
		response.DualClassSubjects = dualClass
	}

	// 查询视频信息
	var videos []sqlModel.UniversityVideo
	if err := r.db.Table("universities_video").Where("school_id = ? and title = 'PC端宣传视频'", id).Order("create_time DESC").Find(&videos).Error; err == nil {
		response.Videos = videos
	}

	// 查询专业信息
	var specials []sqlModel.UniversitiesSpecialPrograms
	if err := r.db.Where("school_id = ?", id).Order("special_name").Find(&specials).Error; err == nil {
		response.SpecialPrograms = specials
	}

	// 查询录取分数线信息
	var admissions []sqlModel.AdmissionUniversities
	if err := r.db.Where("school_id = ?", id).Order("year DESC, province_id").Find(&admissions).Error; err == nil {
		response.AdmissionScores = admissions
	}

	return &response, nil
}

// GetByID 根据ID获取院校基础信息
func (r *UniversityRepository) GetByID(id int) (*sqlModel.UniversitiesDetail, error) {
	var university sqlModel.UniversitiesDetail
	err := r.db.Table("universities_detail").Where("id = ?", id).First(&university).Error
	if err != nil {
		return nil, err
	}
	return &university, nil
}

// GetCollegesBySchoolID 获取院校的学院与系信息
func (r *UniversityRepository) GetCollegesBySchoolID(schoolID int) ([]sqlModel.UniversitiesCollegesDepartments, error) {
	var colleges []sqlModel.UniversitiesCollegesDepartments
	err := r.db.Where("school_id = ?", schoolID).Find(&colleges).Error
	return colleges, err
}

// GetDisciplineRankingsBySchoolID 获取院校的学科评估信息
func (r *UniversityRepository) GetDisciplineRankingsBySchoolID(schoolID int) ([]sqlModel.UniversitiesDisciplineRankings, error) {
	var rankings []sqlModel.UniversitiesDisciplineRankings
	err := r.db.Where("school_id = ?", schoolID).Find(&rankings).Error
	return rankings, err
}

// GetDualClassSubjectsBySchoolID 获取院校的双一流学科信息
func (r *UniversityRepository) GetDualClassSubjectsBySchoolID(schoolID int) ([]sqlModel.UniversitiesDualClassSubjects, error) {
	var subjects []sqlModel.UniversitiesDualClassSubjects
	err := r.db.Where("school_id = ?", schoolID).Find(&subjects).Error
	return subjects, err
}

// GetVideosBySchoolID 获取院校的视频信息
func (r *UniversityRepository) GetVideosBySchoolID(schoolID int) ([]sqlModel.UniversityVideo, error) {
	var videos []sqlModel.UniversityVideo
	err := r.db.Where("school_id = ?", schoolID).Order("create_time DESC").Find(&videos).Error
	return videos, err
}

// GetSpecialProgramsBySchoolID 获取院校的专业信息
func (r *UniversityRepository) GetSpecialProgramsBySchoolID(schoolID int) ([]sqlModel.UniversitiesSpecialPrograms, error) {
	var programs []sqlModel.UniversitiesSpecialPrograms
	err := r.db.Where("school_id = ?", schoolID).Order("special_name").Find(&programs).Error
	return programs, err
}

// GetAdmissionScoresBySchoolID 获取院校的录取分数线信息
func (r *UniversityRepository) GetAdmissionScoresBySchoolID(schoolID int) ([]sqlModel.AdmissionUniversities, error) {
	var scores []sqlModel.AdmissionUniversities
	err := r.db.Where("school_id = ?", schoolID).Order("year DESC, province_id").Find(&scores).Error
	return scores, err
}

// GetSpecialProgramsWithFilter 根据条件筛选院校专业信息
func (r *UniversityRepository) GetSpecialProgramsWithFilter(schoolID int, filters map[string]interface{}) ([]sqlModel.UniversitiesSpecialPrograms, error) {
	query := r.db.Where("school_id = ?", schoolID)

	for key, value := range filters {
		switch key {
		case "level_name":
			query = query.Where("level_name = ?", value)
		case "nation_feature":
			query = query.Where("nation_feature = ?", value)
		case "nation_first_class":
			query = query.Where("nation_first_class = ?", value)
		case "ruanke_level":
			query = query.Where("ruanke_level = ?", value)
		case "year":
			query = query.Where("year = ?", value)
		}
	}

	var programs []sqlModel.UniversitiesSpecialPrograms
	err := query.Order("special_name").Find(&programs).Error
	return programs, err
}

// GetAdmissionScoresByProvince 获取院校在指定省份的录取分数线
func (r *UniversityRepository) GetAdmissionScoresByProvince(schoolID, provinceID int) ([]sqlModel.AdmissionUniversities, error) {
	var scores []sqlModel.AdmissionUniversities
	err := r.db.Where("school_id = ? AND province_id = ?", schoolID, provinceID).
		Order("year DESC").
		Find(&scores).Error
	return scores, err
}

// GetByIDs 批量获取院校基础信息
func (r *UniversityRepository) GetByIDs(ids []int) ([]sqlModel.UniversitiesDetail, error) {
	var universities []sqlModel.UniversitiesDetail
	err := r.db.Where("id IN ?", ids).Find(&universities).Error
	return universities, err
}

// GetByProvince 根据省份获取院校列表
func (r *UniversityRepository) GetByProvince(provinceID int) ([]sqlModel.UniversitiesDetail, error) {
	var universities []sqlModel.UniversitiesDetail
	err := r.db.Where("province_id = ?", provinceID).Find(&universities).Error
	return universities, err
}

// GetBySpecialAndScore 根据专业和分数筛选院校
func (r *UniversityRepository) GetBySpecialAndScore(specialID int, minScore uint16) ([]sqlModel.UniversitiesDetail, error) {
	var universities []sqlModel.UniversitiesDetail
	err := r.db.Joins("JOIN universities_special_programs ON universities_detail.id = universities_special_programs.school_id").
		Where("universities_special_programs.special_id = ? AND universities_detail.id IN "+
			"(SELECT school_id FROM admission_special WHERE min_score <= ?)", specialID, minScore).
		Find(&universities).Error
	return universities, err
}

// QuerySimpleColleges querySimpleColleges 查询精简的学校信息
func (r *UniversityRepository) QuerySimpleColleges(params *param.CollegeQueryParams) ([]param.UniversitiesSimple, int, error) {
	query := r.db.Table("universities_detail u").
		Select(`u.id, u.name, u.motto, u.province_id, p.name as province_name, 
			   u.postcode, u.belong, u.nature_name, u.type_name, u.level_name, 
			   u.create_date, u.f211, u.f985, u.logo_url, u.dual_class_name`).
		Joins("LEFT JOIN common_provinces p ON u.province_id = p.id")

	if params.Name != "" {
		query.Where("u.name LIKE ?", "%"+params.Name+"%")
	}
	if params.ProvinceID > 0 {
		query.Where("u.province_id = ?", params.ProvinceID)
	}
	if params.TypeName != "" {
		query.Where("u.type_name = ?", params.TypeName)
	}
	if params.LevelName != "" {
		query.Where("u.level_name = ?", params.LevelName)
	}
	if params.Is211 > 0 {
		query.Where("u.f211 = ?", params.Is211)
	}
	if params.Is985 > 0 {
		query.Where("u.f985 = ?", params.Is985)
	}
	if params.DualClass != "" {
		query.Where("u.dual_class_name = ?", params.DualClass)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "u.id ASC"
	if params.OrderBy != "" {
		order = "u." + params.OrderBy
		if params.OrderDesc {
			order += " DESC"
		}
	}
	query.Order(order).Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)

	var universities []param.UniversitiesSimple
	if err := query.Find(&universities).Error; err != nil {
		return nil, 0, err
	}

	return universities, int(total), nil
}

// GetIDByName 根据学校名称获取学校ID
func (r *UniversityRepository) GetIDByName(name string) (int, error) {
	var university sqlModel.UniversitiesDetail
	err := r.db.Table("universities_detail").Where("name = ?", name).First(&university).Error
	if err != nil {
		return 0, err
	}
	return university.ID, nil
}

// GetByName 根据学校名称获取学校信息
func (r *UniversityRepository) GetByName(name string) ([]sqlModel.UniversitiesDetail, error) {
	var universities []sqlModel.UniversitiesDetail
	err := r.db.Table("universities_detail").Where("name LIKE ?", "%"+name+"%").Find(&universities).Error
	if err != nil {
		return nil, err
	}
	if len(universities) == 0 {
		return nil, fmt.Errorf("未找到相关学校")
	}
	return universities, err
}
