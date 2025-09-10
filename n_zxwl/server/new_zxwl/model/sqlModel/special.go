package sqlModel

import (
	"time"
)

// SpecialDetail 专业基本信息表
type SpecialDetail struct {
	ID                  int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code                string    `gorm:"not null;column:code" json:"code"`
	Name                string    `gorm:"not null;column:name" json:"name"`
	Degree              string    `gorm:"column:degree" json:"degree"`
	Direction           string    `gorm:"column:direction" json:"direction"`
	Type                string    `gorm:"column:type" json:"type"`
	TypeDetail          string    `gorm:"column:type_detail" json:"type_detail"`
	LimitYear           string    `gorm:"column:limit_year" json:"limit_year"`
	Level1              int       `gorm:"column:level1" json:"level1"`
	Level2              int       `gorm:"column:level2" json:"level2"`
	Level3              int       `gorm:"column:level3" json:"level3"`
	Level1Name          string    `gorm:"column:level1_name" json:"level1_name"`
	Level2Name          string    `gorm:"column:level2_name" json:"level2_name"`
	Level3Name          string    `gorm:"column:level3_name" json:"level3_name"`
	EmploymentRate      *float64  `gorm:"column:employment_rate" json:"employment_rate"`
	AvgSalary           int       `gorm:"column:avg_salary" json:"avg_salary"`
	TopIndustry         string    `gorm:"column:top_industry" json:"top_industry"`
	TopPosition         string    `gorm:"column:top_position" json:"top_position"`
	TopArea             string    `gorm:"column:top_area" json:"top_area"`
	MonthlyViews        string    `gorm:"column:monthly_views" json:"monthly_views"`
	TotalViews          string    `gorm:"column:total_views" json:"total_views"`
	SubjectRequirements string    `gorm:"column:subject_requirements" json:"subject_requirements"`
	GenderRatio         string    `gorm:"column:gender_ratio" json:"gender_ratio"`
	Celebrities         string    `gorm:"column:celebrities" json:"celebrities"`
	Courses             string    `gorm:"column:courses" json:"courses"`
	Content             string    `gorm:"column:content" json:"content"`
	CareerProspects     string    `gorm:"column:career_prospects" json:"career_prospects"`
	Description         string    `gorm:"column:description" json:"description"`
	UpdateTime          time.Time `gorm:"column:update_time" json:"update_time"`
	CreateTime          time.Time `gorm:"column:create_time" json:"create_time"`
}

// SpecialContent 专业大文本内容表
type SpecialContent struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SpecialID   int       `gorm:"not null" json:"special_id"`
	ContentType int       `gorm:"not null" json:"content_type"` // 1-专业介绍 2-就业方向 3-专业描述 4-主要课程
	Content     string    `json:"content"`
	CreateTime  time.Time `json:"create_time"`
}

// SpecialEmploymentRate 专业就业率表
type SpecialEmploymentRate struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SpecialID  string    `gorm:"not null" json:"special_id"`
	Year       string    `gorm:"not null" json:"year"`
	Rate       string    `gorm:"not null" json:"rate"`
	CreateTime time.Time `json:"create_time"`
}

// SpecialFamousSchool 专业名校示例表
type SpecialFamousSchool struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SpecialID  string    `gorm:"not null" json:"special_id"`
	SchoolName string    `gorm:"not null" json:"school_name"`
	CreateTime time.Time `json:"create_time"`
}

// SpecialVideo 专业视频表
type SpecialVideo struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SpecialID       int       `gorm:"not null" json:"special_id"`
	Title           string    `gorm:"not null" json:"title"`
	CoverImage      string    `json:"cover_image"`
	VideoURL        string    `gorm:"not null" json:"video_url"`
	UpdateTime      time.Time `json:"update_time"`
	SchoolID        string    `json:"school_id"`
	SchoolSpecialID string    `json:"school_special_id"`
	Ranks           int       `json:"ranks"`
	URLType         int       `json:"url_type"`
	CreateTime      time.Time `json:"create_time"`
}

// SpecialImpressionTag 专业印象标签表
type SpecialImpressionTag struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SpecialID  string    `gorm:"not null" json:"special_id"`
	ImageURL   string    `gorm:"not null" json:"image_url"`
	Keyword    string    `gorm:"not null" json:"keyword"`
	CreateTime time.Time `json:"create_time"`
}

// SpecialJobDistribution 专业就业分布表
type SpecialJobDistribution struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SpecialID        string    `gorm:"not null" json:"special_id"`
	DistributionType int       `gorm:"not null" json:"distribution_type"`
	Name             string    `json:"name"`
	Position         string    `json:"position"`
	JobDescription   string    `json:"job_description"`
	Rate             string    `gorm:"not null" json:"rate"`
	Sort             int       `json:"sort"`
	CreateTime       time.Time `json:"create_time"`
}

// SpecialSalaryData 专业薪资数据表
type SpecialSalaryData struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SpecialID   string    `gorm:"not null" json:"special_id"`
	SalaryType  int       `gorm:"not null" json:"salary_type"`
	SalaryYear  int       `gorm:"not null" json:"salary_year"`
	SalaryValue int       `gorm:"not null" json:"salary_value"`
	CreateTime  time.Time `json:"create_time"`
}

// SpecialDetailResponse 专业详细信息返回实体
type SpecialDetailResponse struct {
	SpecialDetail         SpecialDetail                 `json:"detail"`                  // 专业基本信息
	SpecialContents       []SpecialContent              `json:"contents"`                // 专业大文本内容
	EmploymentRates       []SpecialEmploymentRate       `json:"employment_rates"`        // 专业就业率信息
	FamousSchools         []SpecialFamousSchool         `json:"famous_schools"`          // 专业名校示例
	Videos                []SpecialVideo                `json:"videos"`                  // 专业视频信息
	ImpressionTags        []SpecialImpressionTag        `json:"impression_tags"`         // 专业印象标签
	JobDistributions      []SpecialJobDistribution      `json:"job_distributions"`       // 专业就业分布
	SalaryData            []SpecialSalaryData           `json:"salary_data"`             // 专业薪资数据
	UniversitySpecialInfo []UniversitiesSpecialPrograms `json:"university_special_info"` // 开设该专业的院校信息
}
