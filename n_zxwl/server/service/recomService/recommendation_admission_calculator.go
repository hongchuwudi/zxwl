// admission_calculator 计算录取概率
package recomService

import (
	"fmt"
	"math"
	"mymod/model/sqlModel"
	"sort"
)

// CalculateProbability 计算录取概率（sigmoid函数）
func CalculateProbability(userScore, targetScore, gap float64) float64 {
	// 分数差异
	scoreDiff := userScore - targetScore

	// 使用sigmoid函数计算概率
	// 公式: P = 1 / (1 + e^(-(x + gap)/gap))
	// 其中 x = userScore - targetScore
	z := (scoreDiff + gap) / gap
	probability := 1.0 / (1.0 + math.Exp(-z))

	// 限制概率在0.05到0.98之间，避免极端值
	return math.Max(0.01, math.Min(0.99, probability))
}

// CalculateSpecialAdmissionProbability 计算专业录取概率
func CalculateSpecialAdmissionProbability(userScore, minScore, scoreGap float64) float64 {
	return CalculateProbability(userScore, minScore, scoreGap)
}

// CalculateUniversityAdmissionProbability 计算学校录取概率
func CalculateUniversityAdmissionProbability(userScore, minScore, scoreGap float64) float64 {
	// 学校录取更保守，使用0.89倍的scoreGap
	return CalculateProbability(userScore, minScore, scoreGap*0.62)
}

// ClassifyAndDistribute 分类并分配冲稳保区间
func ClassifyAndDistribute(schools []sqlModel.UniversityRecommendation, userScore, scoreGap float64) ([]sqlModel.UniversityRecommendation, []sqlModel.UniversityRecommendation, []sqlModel.UniversityRecommendation) {
	// 1. 计算每个学校的录取概率
	for i := range schools {
		schools[i].AdmissionProbability = CalculateUniversityAdmissionProbability(
			userScore, schools[i].MinScore, scoreGap,
		)
		fmt.Println(schools[i].SchoolName, "-录取概率:", schools[i].AdmissionProbability)
	}

	// 2. 按录取概率排序（从低到高）
	sort.Slice(schools, func(i, j int) bool {
		return schools[i].AdmissionProbability < schools[j].AdmissionProbability
	})

	// 3. 根据录取概率分类
	var chongSchools, wenSchools, baoSchools []sqlModel.UniversityRecommendation

	for _, school := range schools {
		if school.AdmissionProbability < 0.4 && len(chongSchools) < 8 {
			chongSchools = append(chongSchools, school)
		} else if school.AdmissionProbability >= 0.5 && school.AdmissionProbability < 0.9 && len(wenSchools) < 8 {
			wenSchools = append(wenSchools, school)
		} else if school.AdmissionProbability >= 0.9 && len(baoSchools) < 8 {
			baoSchools = append(baoSchools, school)
		}

		// 如果所有区间都已满，停止处理
		if len(chongSchools) >= 8 && len(wenSchools) >= 8 && len(baoSchools) >= 8 {
			break
		}
	}

	return chongSchools, wenSchools, baoSchools
}

// ProcessSpecialRecommendations 处理专业推荐结果
func ProcessSpecialRecommendations(specials []sqlModel.SpecialRecommendation, userScore, scoreGap float64) []sqlModel.SpecialRecommendation {
	// 1. 计算每个专业的录取概率
	for i := range specials {
		specials[i].AdmissionProbability = CalculateSpecialAdmissionProbability(
			userScore, specials[i].ScoreMinScore, scoreGap,
		)
	}

	// 2. 按录取概率从高到低排序
	sort.Slice(specials, func(i, j int) bool {
		if specials[i].AdmissionProbability == specials[j].AdmissionProbability {
			// 概率相同时保持原顺序
			return i < j
		}
		return specials[i].AdmissionProbability > specials[j].AdmissionProbability
	})

	// 3. 取前50个专业
	if len(specials) > 50 {
		return specials[:50]
	}
	return specials
}
