// services/score_service.go
package service

import (
	"mymod/model/param"
	"mymod/model/sqlModel"
	scoreRepo "mymod/repositories/scoreRepo"
	"sort"
)

type ScoreService struct {
	scoreRepo *scoreRepo.ScoreRepository
}

func NewScoreService(scoreRepo *scoreRepo.ScoreRepository) *ScoreService {
	return &ScoreService{scoreRepo: scoreRepo}
}

// GetScoreData 获取分数数据
func (s *ScoreService) GetScoreData(provinceID, typeID, year, batchID int) (*param.ScoreDataResponse, error) {
	// 获取考试信息
	examInfo, err := s.scoreRepo.GetExamInfo(provinceID, typeID, year, batchID)
	if err != nil {
		return nil, err
	}

	// 获取分数段数据
	sections, err := s.scoreRepo.GetScoreSections(int64(examInfo.ID), year)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	sectionDetails := s.convertToSectionDetails(sections, examInfo.TotalNum)

	// 计算统计信息
	statistic := s.calculateStatistics(sectionDetails, examInfo.TotalNum)

	response := &param.ScoreDataResponse{
		ExamInfo:  examInfo,
		Sections:  sectionDetails,
		ScoreStat: statistic,
	}

	return response, nil
}

// convertToSectionDetails 转换分数段数据到响应格式
func (s *ScoreService) convertToSectionDetails(sections []*sqlModel.ScoreSection, totalNum int) []*param.ScoreSectionDetail {
	sectionDetails := make([]*param.ScoreSectionDetail, 0, len(sections))

	for _, section := range sections {
		percentage := 0.0
		if totalNum > 0 {
			percentage = float64(section.TotalStudents) / float64(totalNum) * 100
		}

		sectionDetails = append(sectionDetails, &param.ScoreSectionDetail{
			ScoreRange:    section.ScoreRange,
			MinScore:      section.MinScore,
			MaxScore:      section.MaxScore,
			RankRange:     section.RankRange,
			RankStart:     section.RankStart,
			RankEnd:       section.RankEnd,
			TotalStudents: section.TotalStudents,
			Percentage:    percentage,
		})
	}

	return sectionDetails
}

// calculateStatistics 计算分数统计信息
func (s *ScoreService) calculateStatistics(sections []*param.ScoreSectionDetail, totalStudents int) *param.ScoreStatistic {
	if len(sections) == 0 {
		return nil
	}

	stat := &param.ScoreStatistic{
		MaxScore: sections[0].MaxScore,               // 第一个是最高分段
		MinScore: sections[len(sections)-1].MinScore, // 最后一个是最低分段
	}

	// 计算平均分和中位数
	totalScore := 0.0
	allScores := make([]float64, 0)

	for _, section := range sections {
		// 使用分段中间值作为该分段所有学生的分数估算
		midScore := (section.MinScore + section.MaxScore) / 2
		totalScore += midScore * float64(section.TotalStudents)

		// 为计算中位数，添加每个学生的分数（估算）
		for i := 0; i < section.TotalStudents; i++ {
			allScores = append(allScores, midScore)
		}
	}

	if totalStudents > 0 {
		stat.AverageScore = totalScore / float64(totalStudents)
	}

	// 计算中位数
	if len(allScores) > 0 {
		sort.Float64s(allScores)
		if len(allScores)%2 == 0 {
			stat.MedianScore = (allScores[len(allScores)/2-1] + allScores[len(allScores)/2]) / 2
		} else {
			stat.MedianScore = allScores[len(allScores)/2]
		}
	}

	return stat
}

// GetAvailableYears 获取可用的年份列表
func (s *ScoreService) GetAvailableYears(provinceID, typeID int) ([]int, error) {
	return s.scoreRepo.GetAvailableYears(provinceID, typeID)
}
