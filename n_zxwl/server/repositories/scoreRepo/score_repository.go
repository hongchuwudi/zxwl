package scoreRepo

import (
	"fmt"
	"gorm.io/gorm"
	"mymod/model/sqlModel"
)

type ScoreRepository struct {
	db *gorm.DB
}

func NewScoreRepository(db *gorm.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

// GetEquivalentScore 调用存储过程获取等效分数
func (r *ScoreRepository) GetEquivalentScore(currentYear, provinceID, categoryID int, currentScore float64, targetYear int, batchID int) (*sqlModel.EquivalentScoreResult, error) {
	var result sqlModel.EquivalentScoreResult
	// 第一步：调用存储过程，设置用户变量
	err := r.db.Exec(`
        CALL get_equivalent_score(?, ?, ?, ?, ?, ?, 
            @curr_rank_start, 
            @curr_rank_end, 
            @equivalent_score_start, 
            @equivalent_score_end, 
            @debug_info);
    `, currentYear, provinceID, categoryID, currentScore, targetYear, batchID).Error
	if err != nil {
		return nil, fmt.Errorf("调用等效分数存储过程失败: %v", err)
	}

	// 第二步：单独查询用户变量
	err = r.db.Raw(`
        SELECT 
            @curr_rank_start as curr_rank_start, 
            @curr_rank_end as curr_rank_end,
            @equivalent_score_start as equivalent_score_start,
            @equivalent_score_end as equivalent_score_end,
            @debug_info as debug_info;
    `).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("获取等效分数结果失败: %v", err)
	}

	return &result, nil
}

// GetScoreSectionsByParams 直接通过参数查询分数段数据
func (r *ScoreRepository) GetScoreSectionsByParams(provinceID, typeID, year int) ([]*sqlModel.ScoreSection, error) {
	var sections []*sqlModel.ScoreSection

	err := r.db.Table("score_section as ss").
		Joins("LEFT JOIN score_exam_info as sei ON ss.exam_id = sei.id").
		Where("sei.province_id = ? AND sei.type_id = ? AND sei.exam_year = ?", provinceID, typeID, year).
		Order("ss.min_score DESC").
		Find(&sections).Error

	if err != nil {
		return nil, err
	}
	return sections, nil
}

// GetExamInfo 获取考试基本信息
func (r *ScoreRepository) GetExamInfo(provinceID, typeID, year, batchID int) (*sqlModel.ScoreExamInfo, error) {
	var examInfo sqlModel.ScoreExamInfo
	err := r.db.Where("province_id = ? AND type_id = ? AND exam_year = ? AND batch_id = ?", provinceID, typeID, year, batchID).
		First(&examInfo).Error
	if err != nil {
		return nil, err
	}
	return &examInfo, nil
}

// GetScoreSections 获取分数段数据
func (r *ScoreRepository) GetScoreSections(examID int64, year int) ([]*sqlModel.ScoreSection, error) {
	var sections []*sqlModel.ScoreSection
	err := r.db.Where("exam_id = ? AND exam_year = ?", examID, year).
		Order("min_score DESC").
		Find(&sections).Error
	if err != nil {
		return nil, err
	}
	return sections, nil
}

// GetAvailableYears 获取可用的年份列表
func (r *ScoreRepository) GetAvailableYears(provinceID, typeID int) ([]int, error) {
	var years []int
	err := r.db.Model(&sqlModel.ScoreExamInfo{}).
		Where("province_id = ? AND type_id = ?", provinceID, typeID).
		Distinct("exam_year").
		Pluck("exam_year", &years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}

// GetTotalStudentsByExam 根据考试ID获取总人数
func (r *ScoreRepository) GetTotalStudentsByExam(examID int) (int, error) {
	var examInfo sqlModel.ScoreExamInfo
	err := r.db.Select("total_num").Where("id = ?", examID).First(&examInfo).Error
	if err != nil {
		return 0, err
	}
	return examInfo.TotalNum, nil
}
