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
