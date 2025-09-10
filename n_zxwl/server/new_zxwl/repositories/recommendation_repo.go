// filename: recommendation_repo.go 调用推荐存储过程的repo
package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/model/sqlModel"
	"strings"
)

type RecommendationRepo struct {
	db *gorm.DB
}

func NewRecommendationRepo(db *gorm.DB) *RecommendationRepo {
	return &RecommendationRepo{db: db}
}

// GetDBConfig 从配置获取数据库连接信息
func GetDBConfig() *config.Config {
	return config.LoadConfig()
}

// GetSpecialRecommendations 调用专业推荐存储过程
func (r *RecommendationRepo) GetSpecialRecommendations(input sqlModel.UserInput, scoreGap float64) ([]sqlModel.SpecialRecommendation, error) {
	var results []sqlModel.SpecialRecommendation

	// 转换分类数组为逗号分隔的字符串
	level2Str := ""
	level3Str := ""

	if input.Level2 != nil {
		level2Str = strings.Join(input.Level2, ",")
	}
	if input.Level3 != nil {
		level3Str = strings.Join(input.Level3, ",")
	}

	// 调试信息
	fmt.Printf("调用专业推荐存储过程参数:\n")
	fmt.Printf("GoalYear: %d\n", input.GoalYear)
	fmt.Printf("TypeID: %d\n", input.TypeID)
	fmt.Printf("ProvinceID: %d\n", input.ProvinceID)
	fmt.Printf("EquivalentScore: %.2f-%.2f\n", input.EquivalentScoreStart, input.EquivalentScoreEnd)
	fmt.Printf("Rank: %d-%d\n", input.CurrRankStart, input.CurrRankEnd)
	fmt.Printf("Level2: %s\n", level2Str)
	fmt.Printf("Level3: %s\n", level3Str)

	err := r.db.Raw(`
		CALL sp_filter(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		int(input.GoalYear),        // p_goal_year
		input.TypeID,               // p_type_id
		input.ProvinceID,           // p_province_id
		input.EquivalentScoreStart, // p_equivalent_score_start
		input.EquivalentScoreEnd,   // p_equivalent_score_end
		input.CurrRankStart,        // p_curr_rank_start
		input.CurrRankEnd,          // p_curr_rank_end
		scoreGap,                   // score_gap
		input.IsBenKe,              // p_is_benke
		input.GoalProvinceID,       // p_goal_province_id
		level2Str,                  // p_level2
		level3Str,                  // p_level3
		input.Salary,               // p_salary
	).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("调用专业推荐存储过程失败: %v", err)
	}

	fmt.Printf("专业推荐结果数量: %d\n", len(results))
	return results, nil
}

// GetUniversityRecommendations 调用学校推荐存储过程
func (r *RecommendationRepo) GetUniversityRecommendations(input sqlModel.UserInput, scoreGap float64) ([]sqlModel.UniversityRecommendation, error) {
	var results []sqlModel.UniversityRecommendation

	// 转换分类数组为逗号分隔的字符串
	level2Str := ""
	level3Str := ""

	if input.Level2 != nil {
		level2Str = strings.Join(input.Level2, ",")
	}
	if input.Level3 != nil {
		level3Str = strings.Join(input.Level3, ",")
	}

	// 调试信息
	fmt.Printf("调用学校推荐存储过程参数:\n")
	fmt.Printf("GoalYear: %d\n", input.GoalYear)
	fmt.Printf("TypeID: %d\n", input.TypeID)
	fmt.Printf("ProvinceID: %d\n", input.ProvinceID)
	fmt.Printf("EquivalentScore: %.2f-%.2f\n", input.EquivalentScoreStart, input.EquivalentScoreEnd)

	err := r.db.Raw(`
		CALL Campus_filter(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		int(input.GoalYear),        // p_goal_year
		input.TypeID,               // p_type_id
		input.ProvinceID,           // p_province_id
		input.EquivalentScoreStart, // p_equivalent_score_start
		input.EquivalentScoreEnd,   // p_equivalent_score_end
		input.CurrRankStart,        // p_curr_rank_start
		input.CurrRankEnd,          // p_curr_rank_end
		scoreGap,                   // p_score_gap
		input.IsBenKe,              // p_is_benke
		input.GoalProvinceID,       // p_goal_province_id
		level2Str,                  // p_level2
		level3Str,                  // p_level3
		input.Salary,               // p_salary
	).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("调用学校推荐存储过程失败: %v", err)
	}

	fmt.Printf("学校推荐结果数量: %d\n", len(results))
	return results, nil
}
