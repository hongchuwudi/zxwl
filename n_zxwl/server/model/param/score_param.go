// score_response.go
package param

import "mymod/model/sqlModel"

// ScoreDataResponse 分数数据响应
type ScoreDataResponse struct {
	ExamInfo  *sqlModel.ScoreExamInfo `json:"exam_info"`
	Sections  []*ScoreSectionDetail   `json:"sections"`
	ScoreStat *ScoreStatistic         `json:"score_stat,omitempty"`
}

// ScoreSectionDetail 分数段详情
type ScoreSectionDetail struct {
	ScoreRange    string  `json:"score_range"`
	MinScore      float64 `json:"min_score"`
	MaxScore      float64 `json:"max_score"`
	RankRange     string  `json:"rank_range"`
	RankStart     int     `json:"rank_start"`
	RankEnd       int     `json:"rank_end"`
	TotalStudents int     `json:"total_students"`
	Percentage    float64 `json:"percentage"` // 该分段人数占比
}

// ScoreStatistic 分数统计信息
type ScoreStatistic struct {
	AverageScore float64 `json:"average_score"` // 平均分（估算）
	MaxScore     float64 `json:"max_score"`     // 最高分
	MinScore     float64 `json:"min_score"`     // 最低分
	MedianScore  float64 `json:"median_score"`  // 中位数分数
}

// ScoreQueryParams 分数查询参数
type ScoreQueryParams struct {
	ProvinceID int `json:"province_id"`
	TypeID     int `json:"type_id"`
	Year       int `json:"year"`
}
