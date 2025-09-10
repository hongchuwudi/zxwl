package Const

import (
	"mymod/new_zxwl/model/sqlModel"
	"time"
)

type AuthRequest struct {
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
}

// 验证码请求结构体
type VerifyCodeRequest struct {
	Email string `json:"email"`
}

type RegisterRequest struct {
	UserName        string `json:"user"`
	Email           string `json:"email"`
	Password        string `json:"passwd"`
	ConfirmPassword string `json:"confirm"`
	VerifyCode      string `json:"varifycode"`
}

type VerifyCodeResponse struct {
	Error int    `json:"code"`
	Msg   string `json:"msg,omitempty"`
	Token string `json:"token,omitempty"`
}

// 注册验证响应
type RegisterResponse struct {
	Code int    `json:"error"`
	Msg  string `json:"msg"`
}

type ProfileRequest struct {
	Email string `json:"email"`
}

type AuthResponse struct {
	Code int `json:"code"`
}

type Profile struct {
	Name     *string `json:"name"`
	Email    string  `json:"email"`
	Sex      *int    `json:"sex"`
	Graduate *int    `json:"graduate"`
	Address  *string `json:"address"`
	Picture  *string `json:"picture"`
}

type ProfileResponse struct {
	Code int     `json:"code"`
	Data Profile `json:"data,omitempty"`
}

type ProfileUpdateRequest struct {
	Email    string  `json:"email"`
	Name     *string `json:"name,omitempty"`
	Sex      *int    `json:"sex,omitempty"`
	Graduate *int    `json:"graduate,omitempty"`
	Address  *string `json:"address,omitempty"`
	Picture  *string `json:"picture,omitempty"`
}

type ProfileUpdateResponse struct {
	Code int `json:"code"`
}

type CollegeResponse struct {
	CityName      *string `json:"city_name"`
	DualClassName *string `json:"dual_class_name"`
	F985          *int    `json:"f985"`
	F211          *int    `json:"f211"`
	HighTitle     *string `json:"hightitle"`
	LevelName     *string `json:"level_name"`
	NatureName    *string `json:"nature_name"`
	TypeName      *string `json:"type_name"`
	// IsBiaoShi     int     `json:"isbiaoshi"`
	SchoolLogo string `json:"school_logo"`
}

type CollegesResponse struct {
	Code int               `json:"code"`
	Data []CollegeResponse `json:"data"`
}

type ChatRequest struct {
	CollegeID int    `json:"college_id"`
	UserEmail string `json:"user_email"`
	Content   string `json:"content"`
}

type ChatMessage struct {
	ID        int       `json:"id"`
	CollegeID int       `json:"college_id"`
	UserEmail string    `json:"user_email"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatResponse struct {
	Code     int           `json:"code"`
	Message  ChatMessage   `json:"message,omitempty"`
	Messages []ChatMessage `json:"messages,omitempty"`
}

type SchoolProfileRequest struct {
	Hightitle string `json:"hightitle"`
}

type SchoolProfileResponse struct {
	SchoolLogo string  `json:"school_logo"`
	Introduce  *string `json:"introduce,omitempty"`
	Hightitle  string  `json:"hightitle"`
	SchoolPic  *string `json:"school_pic,omitempty"`
}

type ProfessionalItem struct {
	Name          string      `json:"name"`
	Level1Name    string      `json:"level1_name"`
	Salaryavg     interface{} `json:"salaryavg"` // 支持string和int类型
	LimitYear     string      `json:"limit_year"`
	FiveSalaryavg interface{} `json:"fivesalaryavg"` // 支持string和int类型
	Level2Name    string      `json:"level2_name"`
	BoyRate       string      `json:"boy_rate"`
	GirlRate      string      `json:"girl_rate"`
	Level3Name    string      `json:"level3_name"`
}

type ProfessionalRequest struct {
	Items []ProfessionalItem `json:"item"`
}

type ProfessionalResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg,omitempty"`
	Total int    `json:"total,omitempty"`
}

// 推荐请求结构体
type RecommendationRequest struct {
	Origin           string `json:"origin"`            // 生源地
	TargetLocation   string `json:"target_location"`   // 目标大学位置
	Score            int    `json:"score"`             // 高考分数
	Rank             int    `json:"rank"`              // 高考排名
	PreferredSubject string `json:"preferred_subject"` // 喜欢学科
}

// 推荐响应结构体
type RecommendationResponse struct {
	Code int                `json:"code"`
	Data RecommendationData `json:"data"`
}

type RecommendationData struct {
	Rush   []University `json:"rush"`
	Stable []University `json:"stable"`
	Safe   []University `json:"safe"`
}

type University struct {
	Name        string   `json:"name"`               // 大学名称
	Probability float64  `json:"probability"`        // 录取概率
	Majors      []string `json:"recommended_majors"` // 推荐专业
}

type MoniRequest struct {
	College       string `json:"college"`
	Benke         int    `json:"benke"`
	Yitianzhuanye int    `json:"yitianzhuanye"`
	Major1        string `json:"major1"`
	Major2        string `json:"major2"`
	Major3        string `json:"major3"`
	Major4        string `json:"major4"`
	Major5        string `json:"major5"`
	Major6        string `json:"major6"`
	Email         string `json:"email"`
}

type UpsertResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// 志愿查询请求
type VolunteerFetchRequest struct {
	UserEmail string `json:"user_email"`
	BatchType int    `json:"batch_type"`
}

// 志愿数据项
type VolunteerItem struct {
	Institution string `json:"institution"`  // 对应原college
	Sequence    int    `json:"sequence"`     // 对应原yitianzhuanye
	FirstMajor  string `json:"first_major"`  // 对应原major1
	SecondMajor string `json:"second_major"` // 对应原major2
	ThirdMajor  string `json:"third_major"`  // 对应原major3
	FourthMajor string `json:"fourth_major"` // 对应原major4
	FifthMajor  string `json:"fifth_major"`  // 对应原major5
	SixthMajor  string `json:"sixth_major"`  // 对应原major6
}

// 志愿查询响应
type VolunteerFetchResponse struct {
	Status  int             `json:"status"` // 0=成功
	Message string          `json:"message,omitempty"`
	Records []VolunteerItem `json:"records,omitempty"`
}

// Const包新增结构体
type FamilyAddRequest struct {
	MyEmail     string `json:"myemail"`
	FamilyEmail string `json:"familyemail"`
}

type FamilyAddResponse struct {
	Code int `json:"code"`
}

type FamilyFindRequest struct {
	MyEmail string `json:"myemail"`
}

type FamilyFindResponse struct {
	Code   int      `json:"code"`
	Emails []string `json:"emails"`
}

type FamilyRemoveRequest struct {
	MyEmail     string `json:"myemail"`
	FamilyEmail string `json:"familyemail"`
}

type PolicyItem struct {
	ID      int     `json:"id"`
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

type PolicyResponse struct {
	Code int          `json:"code"`
	Data []PolicyItem `json:"data,omitempty"`
}

// 请求结构体
type PolicySearchRequest struct {
	Title string `json:"title"`
}

// 响应结构体
type PolicySearchResponse struct {
	Code int          `json:"code"`
	Data []PolicyItem `json:"data,omitempty"`
	Msg  string       `json:"msg,omitempty"`
}

type ProfileResponsed struct {
	Address  *string `json:"address,omitempty"`
	Graduate *int    `json:"graduate,omitempty"` // 年份用int类型
	Picture  *string `json:"picture,omitempty"`
	Sex      *int    `json:"sex,omitempty"`
	Email    *string `json:"email,omitempty"`
	Name     *string `json:"name,omitempty"`
}

type ProfileListResponse struct {
	Code int                `json:"code"`
	Data []ProfileResponsed `json:"data"`
	Msg  string             `json:"msg,omitempty"`
}

// 添加在Const包中
type ProfileDeleteRequest struct {
	Email string `json:"email"`
}

type ProfileDeleteResponse struct {
	Code int `json:"code"`
}

// 请求结构体
type PolicyCreateRequest struct {
	Title    *string `json:"title,omitempty"`
	Contenet *string `json:"content,omitempty"` // 注意字段拼写与数据库一致
	Foreword *string `json:"foreword,omitempty"`
}

// 响应结构体
type PolicyCreateResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	ID   int64  `json:"id,omitempty"` // 返回插入的ID
}

type LogRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Date      string `json:"date" validate:"required,datetime=2006-01-02 15:04:05"`
	Operation string `json:"operation" validate:"required,min=1,max=255"`
}

type LogResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	ID   int64  `json:"id,omitempty"`
}

type LogEntry struct {
	Email     string `json:"email"`
	Date      string `json:"date"`
	Operation string `json:"operation"`
}

type LogListResponse struct {
	Code  int        `json:"code"`
	Data  []LogEntry `json:"data"`
	Msg   string     `json:"msg,omitempty"`
	Total int        `json:"total,omitempty"`
}

// 查询请求结构体
type ProfessionalQueryRequest struct {
	Level3Name        string  `json:"level3_name"`         // 三级分类名称
	Type              string  `json:"type"`                // 学科门类
	TypeDetail        string  `json:"type_detail"`         // 专业类
	Degree            string  `json:"degree"`              // 学位
	Keyword           string  `json:"keyword"`             // 关键词搜索
	MinSalary         int     `json:"min_salary"`          // 最低薪资
	MaxSalary         int     `json:"max_salary"`          // 最高薪资
	MinEmploymentRate float64 `json:"min_employment_rate"` // 最低就业率
}

// 返回字段结构体
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

// 创建一个专门用于列表查询的结构体
type SpecialList struct {
	ID             int      `json:"id"`
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	Degree         string   `json:"degree"`
	Direction      string   `json:"direction"`
	LimitYear      string   `json:"limit_year"`
	Level1Name     string   `json:"level1_name"`
	Level2Name     string   `json:"level2_name"`
	Level3Name     string   `json:"level3_name"`
	EmploymentRate *float64 `json:"employment_rate"`
	AvgSalary      *int     `json:"avg_salary"` // 改为指针类型
	TopIndustry    string   `json:"top_industry"`
	TopPosition    string   `json:"top_position"`
	GenderRatio    *string  `json:"gender_ratio"` // 改为指针类型
}

// 全部专业查询请求
type ProfessionalQueryResponse struct {
	Code  int                        `json:"code"`
	Msg   string                     `json:"msg,omitempty"`
	Data  []ProfessionalItemResponse `json:"data,omitempty"`
	Total int64                      `json:"total,omitempty"`
}

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
	Total        int                           `json:"total"`
	Page         int                           `json:"page"`
	PageSize     int                           `json:"page_size"`
	TotalPages   int                           `json:"total_pages"`
	Universities []sqlModel.UniversitiesDetail `json:"universities"`
}
