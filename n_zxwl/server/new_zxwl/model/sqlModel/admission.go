package sqlModel

import "time"

// AdmissionSpecial 高校招生专业信息表
type AdmissionSpecial struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Year               int       `gorm:"not null" json:"year"`
	SchoolID           int       `gorm:"not null" json:"school_id"`
	ProvinceID         int       `gorm:"not null" json:"province_id"`
	SpecialID          int       `json:"special_id"`
	Type               int       `gorm:"not null" json:"type"`  // 高考类型ID
	Batch              int       `gorm:"not null" json:"batch"` // 批次类型ID
	Zslx               int       `json:"zslx"`                  // 招生类型ID
	MaxScore           int       `json:"max_score"`             // 最高分
	MinScore           int       `json:"min_score"`             // 最低分
	SpID               int       `gorm:"not null" json:"sp_id"` // 专业ID
	AverageScore       int       `json:"average_score"`         // 平均分
	MinSection         string    `json:"min_section"`           // 最低位次
	BatchDiff          int       `json:"batch_diff"`            // 批次线差
	SpecialInfo        string    `json:"special_info"`          // 专业补充信息
	AdmissionCount     int       `json:"admission_count"`       // 录取人数
	SubjectRequirement string    `json:"subject_requirement"`   // 选科要求文本
	SubjectCodes       string    `json:"subject_codes"`         // 选科科目代码
	MinRankRange       string    `json:"min_rank_range"`        // 最低位次范围
	RangeMaxRank       string    `json:"range_max_rank"`        // 最大位次范围
	ScoreRangeFlag     int       `json:"score_range_flag"`      // 是否分数段
	Remark             string    `json:"remark"`                // 备注
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// AdmissionUniversities 省份录取分数线表
type AdmissionUniversities struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolID   int    `json:"school_id"`
	ProvinceID int    `json:"province_id"`
	Year       int    `json:"year"`
	ScoreType  string `json:"score_type"`
	MinScore   int    `json:"min_score"`
}

// AdmissionBatch 录取批次表
type AdmissionBatch struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	BatchInfo string `json:"batch_info"`
}

// AdmissionType 高考类型表
type AdmissionType struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	TypeInfo string `json:"type_info"`
}

// AdmissionZhaoshengType 招生类型表
type AdmissionZhaoshengType struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	ZslxInfo string `json:"zslx_info"`
}
