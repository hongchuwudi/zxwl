package main

import (
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/service/recomService"
)

func mainw() {

	// 创建推荐服务实例
	recommendationService := recomService.NewRecommendationRepository(config.GetDB())

	// 测试数据 - 模拟用户输入
	testInput := sqlModel.UserPriInput{
		Year:             2025,
		GoalYear:         2024,
		ProvinceName:     "陕西",  // 测试省份
		GoalProvinceName: "",    // 测试目标省份
		TypeName:         "物理类", // 测试高考类型
		Score:            599.0, // 测试分数
		Rank:             "",    // 测试排名
		Interest:         "",    // 测试兴趣爱好
		FamilyPref:       "",
		Salary:           0.0,
	}

	// 调用处理函数
	_, _, err := recommendationService.GetRecommendations(testInput)
	if err != nil {
		return
	}

}
