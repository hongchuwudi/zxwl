package sqlModel

import (
	"time"
)

// UniversitiesDetail 大学基本信息表
type UniversitiesDetail struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	Motto          string    `json:"motto"`
	Address        string    `json:"address"`
	ProvinceID     int       `json:"province_id"`
	CityID         int       `json:"city_id"`
	CityName       string    `json:"city_name"`
	CountyID       int       `json:"county_id"`
	TownName       string    `json:"town_name"`
	Postcode       string    `json:"postcode"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Site           string    `json:"site"`
	SchoolSite     string    `json:"school_site"`
	Belong         string    `json:"belong"`
	NatureName     string    `json:"nature_name"`
	TypeName       string    `json:"type_name"`
	LevelName      string    `json:"level_name"`
	CreateDate     string    `json:"create_date"`
	Area           float64   `json:"area"`
	NumAcademician int       `json:"num_academician"`
	NumDoctor      int       `json:"num_doctor"`
	NumMaster      int       `json:"num_master"`
	NumLab         int       `json:"num_lab"`
	NumLibrary     string    `json:"num_library"`
	NumSubject     int       `json:"num_subject"`
	Content        string    `json:"content"`
	F211           int       `json:"f211"`
	F985           int       `json:"f985"`
	LogoURL        string    `json:"logo_url"`
	DualClassName  string    `json:"dual_class_name"`
	QsWorld        string    `json:"qs_world"`
	UsRank         string    `json:"us_rank"`
	RuankeRank     string    `json:"ruanke_rank"`
	XyhRank        string    `json:"xyh_rank"`
	Status         string    `json:"status"`
	AddTime        time.Time `json:"add_time"`
	LastUpdated    time.Time `json:"last_updated"`
}

// UniversitiesCollegesDepartments 学院与系表
type UniversitiesCollegesDepartments struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID       int    `json:"school_id"`
	CampusName     string `json:"campus_name"`
	DepartmentID   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
}

// UniversitiesDisciplineRankings 学科评估表
type UniversitiesDisciplineRankings struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID  int    `json:"school_id"`
	RankLevel string `json:"rank_level"`
	Count     int    `json:"count"`
}

// UniversitiesDualClassSubjects 双一流学科表
type UniversitiesDualClassSubjects struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID    int    `json:"school_id"`
	SubjectName string `json:"subject_name"`
	SubjectID   int    `json:"subject_id"`
}

// UniversityVideo 院校视频信息表
type UniversityVideo struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID   int       `gorm:"not null" json:"school_id"`
	VideoType  int       `json:"video_type"`
	Title      string    `json:"title"`
	URL        string    `gorm:"not null" json:"url"`
	ImgURL     string    `json:"img_url"`
	URLType    int       `json:"url_type"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// UniversitiesSpecialPrograms 学校专业表
type UniversitiesSpecialPrograms struct {
	ID               int    `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID         int    `json:"school_id"`
	SpecialID        int    `json:"special_id"`
	SpecialName      string `json:"special_name"`
	LevelName        string `json:"level_name"`
	LimitYear        string `json:"limit_year"`
	NationFeature    string `json:"nation_feature"`     // 是否国家级特色专业（1是，2否）
	NationFirstClass string `json:"nation_first_class"` // 是否国家一流专业（1是，2否）
	ProvinceFeature  string `json:"province_feature"`   // 是否省级特色专业（1是，2否）
	RuankeLevel      string `json:"ruanke_level"`       // 软科评级（如"A+"）
	RuankeRank       int    `json:"ruanke_rank"`        // 软科排名
	XuekeRankScore   string `json:"xueke_rank_score"`   // 学科评估结果（如"A+"）
	Year             string `json:"year"`               // 数据年份（如"2022"）
	IsImportant      string `json:"is_important"`       // 是否重点专业（1是，2否）
}

// UniversityDetailResponse 大学详细信息返回实体
type UniversityDetailResponse struct {
	UniversitiesDetail  UniversitiesDetail                `json:"detail"`               // 学校基本信息
	CollegesDepartments []UniversitiesCollegesDepartments `json:"colleges_departments"` // 学院与系信息
	DisciplineRankings  []UniversitiesDisciplineRankings  `json:"discipline_rankings"`  // 学科评估信息
	DualClassSubjects   []UniversitiesDualClassSubjects   `json:"dual_class_subjects"`  // 双一流学科信息
	Videos              []UniversityVideo                 `json:"videos"`               // 视频信息
	SpecialPrograms     []UniversitiesSpecialPrograms     `json:"special_programs"`     // 专业信息
	AdmissionScores     []AdmissionUniversities           `json:"admission_scores"`     // 录取分数线信息
}
