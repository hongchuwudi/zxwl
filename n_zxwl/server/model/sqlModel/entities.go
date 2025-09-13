package sqlModel

// TODO 记录数据库中没有的实体类型

// UserPriInput 用户初次输入信息
type UserPriInput struct {
	Year             uint    `json:"year"`               // 高考年份
	GoalYear         uint    `json:"goal_year"`          // 参考年份（通常是当前年份）
	Score            float64 `json:"score"`              // 用户输入的分数
	ProvinceName     string  `json:"province_name"`      // 用户输入的高考省份名称
	TypeName         string  `json:"type_name"`          // 用户输入的高考类型名称(文科/理科)
	BatchName        string  `json:"batch_name"`         // 高考批次类型(本科,专科两个选择)
	Rank             string  `json:"rank"`               // 用户输入的排名
	GoalProvinceName string  `json:"goal_province_name"` // 用户输入的目标省份名称
	Interest         string  `json:"interest"`           // 用户输入的兴趣爱好
	FamilyPref       string  `json:"family_pref"`        // 用户输入的家庭偏好
	Salary           float64 `json:"salary"`             // 用户输入的希望工资
}

// UserInput 处理后的用户输入信息
type UserInput struct {
	Year                 uint     `json:"year"`                   // 高考年份
	GoalYear             uint     `json:"goal_year"`              // 参考年份
	TypeID               int      `json:"type_id"`                // 高考类型(文科/理科)
	IsBenKe              int      `json:"is_ben_ke"`              // 一分一段表中是否为本科?(1为本科2为专科3不分本专,默认为3)
	CurrRankStart        int64    `json:"curr_rank_start"`        // 当前排名起始
	CurrRankEnd          int64    `json:"curr_rank_end"`          // 当前排名结束
	EquivalentScoreStart float64  `json:"equivalent_score_start"` // 目标年份等效分数起始
	EquivalentScoreEnd   float64  `json:"equivalent_score_end"`   // 目标年份等效分数结束
	DebugInfo            string   `json:"debug_info"`             // 调试信息
	ProvinceID           int      `json:"province_id"`            // 高考省份ID
	GoalProvinceID       int      `json:"goal_provice_id"`        // 目标省份ID
	Level2               []string `json:"level1_arr"`             // 专业大类
	Level3               []string `json:"level2_arr"`             // 专业小类
	Salary               float64  `json:"salary"`                 // 用户输入的期望工资
}

// EquivalentScoreResult 同位分转化返回结果(同位分存储过程)
type EquivalentScoreResult struct {
	CurrRankStart        int64   `gorm:"column:curr_rank_start"`
	CurrRankEnd          int64   `gorm:"column:curr_rank_end"`
	EquivalentScoreStart float64 `gorm:"column:equivalent_score_start"`
	EquivalentScoreEnd   float64 `gorm:"column:equivalent_score_end"`
	DebugInfo            string  `gorm:"column:debug_info"`
}

// SpecialRecommendation 专业推荐结果（对应sp_filter存储过程返回的字段）
type SpecialRecommendation struct {
	SchoolID                   int64   `json:"school_id"`                    // 学校表ID
	SchoolName                 string  `json:"school_name"`                  // 高校名称
	SchoolAddress              string  `json:"school_address"`               // 高校地址
	SpecialID                  int64   `json:"special_id"`                   // 专业表ID
	SpecialName                string  `json:"special_name"`                 // 专业名称
	SpecialLevel1Name          string  `json:"special_level1_name"`          // 专业一级分类名称
	SpecialLevel2Name          string  `json:"special_level2_name"`          // 专业二级分类名称
	SpecialLevel3Name          string  `json:"special_level3_name"`          // 专业三级分类名称
	SpecialAvgSalary           float64 `json:"special_avg_salary"`           // 专业平均薪资
	SpecialInfo                string  `json:"special_info"`                 // 专业补充信息
	AdmissionSpecialRemark     string  `json:"admission_special_remark"`     // 招生信息补充信息
	SpecialCode                string  `json:"special_code"`                 // 专业代码
	SpecialSubjectRequirements string  `json:"special_subject_requirements"` // 选科建议
	SpecialKeywords            string  `json:"special_keywords"`             // 专业关键词（合并多个标签）
	AdmissionID                int64   `json:"admission_id"`                 // 录取信息表ID
	ScoreMinRank               int64   `json:"score_min_rank"`               // 最低排名
	ScoreMaxRank               int64   `json:"score_max_rank"`               // 最高排名
	ScoreMinScore              float64 `json:"score_min_score"`              // 最低分数
	ScoreMaxScore              float64 `json:"score_max_score"`              // 最高分数
	ScoreAvgScore              float64 `json:"score_avg_score"`              // 平均分数
	AdmissionCount             int64   `json:"admission_count"`              // 录取人数
	CommonProvinceName         string  `json:"common_province_name"`         // 生源地省份名称
	ZslxName                   string  `json:"zslx_name"`                    // 招生类型名称
	BatchName                  string  `json:"batch_name"`                   // 批次名称
	AdmissionProbability       float64 `json:"admission_probability"`        // 录取概率
}

// UniversityRecommendation 学校推荐结果（对应Campus_filter存储过程返回的字段）
type UniversityRecommendation struct {
	UniversityID         int64   `json:"university_id"`         // 学校ID
	SchoolName           string  `json:"school_name"`           // 高校名称
	SchoolAddress        string  `json:"school_address"`        // 高校地址
	RuankeRank           string  `json:"ruanke_rank"`           // 软科排名
	XyhRank              string  `json:"xyh_rank"`              // 校友会排名
	ProvinceName         string  `json:"province_name"`         // 省份名称
	TypeName             string  `json:"type_name"`             // 类型名称
	MinScore             float64 `json:"min_score"`             // 最低分数
	HasRkRank            int     `json:"has_rk_rank"`           // 是否有rk排名
	HasXyhRank           int     `json:"has_xyh_rank"`          // 是否有xyh排名
	AdmissionProbability float64 `json:"admission_probability"` // 录取概率
}

// Recommendation 推荐结果
type Recommendation struct {
	SpecialRecommendations    []SpecialRecommendation    `json:"special_recommendations"`    // 推荐的专业数组（来自sp_filter存储过程）
	UniversityRecommendations []UniversityRecommendation `json:"university_recommendations"` // 推荐的学校数组（来自Campus_filter存储过程）
	AdmissionAnalysis         string                     `json:"admission_analysis"`         // 高考分析结果
	InterestAnalysis          string                     `json:"interest_analysis"`          // 个人兴趣爱好结果
	DebugInfo                 string                     `gorm:"column:debug_info"`          // 调试信息
}
