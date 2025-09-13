package param

// CollegeQueryParams 学校查询参数
type CollegeQueryParams struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	Name        string `json:"name"`
	ProvinceID  int    `json:"province_id"`
	CityID      int    `json:"city_id"`
	TypeName    string `json:"type_name"`
	LevelName   string `json:"level_name"`
	Is211       int    `json:"is_211"`
	Is985       int    `json:"is_985"`
	DualClass   string `json:"dual_class"`
	MinRuanke   int    `json:"min_ruanke"`
	MaxRuanke   int    `json:"max_ruanke"`
	SpecialID   int    `json:"special_id"`
	SpecialName string `json:"special_name"`
	OrderBy     string `json:"order_by"`
	OrderDesc   bool   `json:"order_desc"`
}

// CollegeListResponse 学校列表响应
type CollegeListResponse struct {
	Total        int                  `json:"total"`
	Page         int                  `json:"page"`
	PageSize     int                  `json:"page_size"`
	TotalPages   int                  `json:"total_pages"`
	Universities []UniversitiesSimple `json:"universities"`
}

// UniversitiesSimple 大学简要信息表
type UniversitiesSimple struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"not null" json:"name"`
	Motto         string `json:"motto"`
	ProvinceID    int    `json:"province_id"`
	ProvinceName  string `json:"province_name"`
	Postcode      string `json:"postcode"`
	Belong        string `json:"belong"`
	NatureName    string `json:"nature_name"`
	TypeName      string `json:"type_name"`
	LevelName     string `json:"level_name"`
	CreateDate    string `json:"create_date"`
	F211          int    `json:"f211"`
	F985          int    `json:"f985"`
	LogoURL       string `json:"logo_url"`
	DualClassName string `json:"dual_class_name"`
}

// ProfessionalItemResponse 专业返回字段结构体
type ProfessionalItemResponse struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Level1Name     string   `json:"level1_name"`
	Level2Name     string   `json:"level2_name"`
	Level3Name     string   `json:"level3_name"`
	Salaryavg      int      `json:"salaryavg"`
	LimitYear      string   `json:"limit_year"`
	FiveSalaryavg  int      `json:"fivesalaryavg"`
	BoyRate        string   `json:"boy_rate"`
	GirlRate       string   `json:"girl_rate"`
	Code           string   `json:"code,omitempty"`
	Degree         string   `json:"degree,omitempty"`
	EmploymentRate *float64 `json:"employment_rate,omitempty"`
	TopIndustry    string   `json:"top_industry,omitempty"`
	TopPosition    string   `json:"top_position,omitempty"`
}

// SchAndSpeResponse 搜索学校与专业返回字段结构体
type SchAndSpeResponse struct {
	ProfessionalItemResArr []ProfessionalItemResponse `json:"professional_item_res_arr"`
	SchoolProfileRes       []UniversitiesSimple       `json:"school_profile_res"`
}
