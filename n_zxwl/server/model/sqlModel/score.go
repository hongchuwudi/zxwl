package sqlModel

// ScoreExamInfo 一分一段
type ScoreExamInfo struct {
	ID         uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ExamYear   int     `gorm:"not null" json:"exam_year"`
	ExamName   string  `gorm:"not null" json:"exam_name"`
	ProvinceID int     `gorm:"not null" json:"province_id"`
	TypeID     int     `gorm:"not null" json:"type_id"`
	BatchID    int     `gorm:"not null" json:"batch_id"`
	FullScore  float64 `gorm:"not null" json:"full_score"`
	TotalNum   int     `gorm:"not null" json:"total_num"`
}

// ScoreSection 一分一段主表
type ScoreSection struct {
	ID            uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ExamYear      int     `gorm:"not null" json:"exam_year"`
	BatchName     string  `gorm:"not null" json:"batch_name"`
	ExamID        int     `gorm:"not null" json:"exam_id"`
	ScoreRange    string  `gorm:"not null" json:"score_range"`
	MinScore      float64 `gorm:"not null" json:"min_score"`
	MaxScore      float64 `gorm:"not null" json:"max_score"`
	RankRange     string  `gorm:"not null" json:"rank_range"`
	RankStart     int     `gorm:"not null" json:"rank_start"`
	RankEnd       int     `gorm:"not null" json:"rank_end"`
	TotalStudents int     `gorm:"not null" json:"total_students"`
}
