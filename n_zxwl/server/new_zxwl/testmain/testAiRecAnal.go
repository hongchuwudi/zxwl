package main

import (
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/model/sqlModel"
	"mymod/new_zxwl/service"
)

func mains1() {
	analysisService := service.NewAdmissionAnalysisService(config.LoadConfig())
	var priI sqlModel.UserPriInput
	var I sqlModel.UserInput
	priI = sqlModel.UserPriInput{
		Year:             2025,
		GoalYear:         2024,
		ProvinceName:     "北京",    // 测试省份
		GoalProvinceName: "北京",    // 测试目标省份
		TypeName:         "综合",    // 测试高考类型
		Score:            577.0,   // 测试分数
		Rank:             "",      // 测试排名
		Interest:         "想要吃炸鸡", // 测试兴趣爱好
		FamilyPref:       "家长让我经济独立",
	}
	I = sqlModel.UserInput{
		Year:                 2025,
		GoalYear:             2024,
		TypeID:               3,
		IsBenKe:              1,
		CurrRankStart:        17005,
		CurrRankEnd:          17266,
		EquivalentScoreStart: 579.0,
		EquivalentScoreEnd:   580.0,
		DebugInfo:            "步骤0: 是否为特殊城市: 是\n步骤1: 考试类型转化: 3 --> 3\n步骤2: 考试id: 284, 考试名称: 2025年-北京省-高考-综合（本科）\n步骤3: 获取当前分数对应的排名区间: 17005 - 17264\n步骤4: 获取当前年份总人数: 65434\n步骤5: 获取当前年份考试分数是否过线: 本科批\n步骤6: 目标考试id: 230, 考试名称: 2024年-北京省-高考-综合（本科）\n步骤7: 获取目标年份总人数: 55936\n步骤8: 当前排名在目标年份中的等效排名: 14537 - 14759\n步骤9和10: 当前分数在目标年份中的等效分数: 580.00 - 579.00\n步骤11: 确保分数顺序正确: 579.00 - 580.00",
		ProvinceID:           11,
		GoalProvinceID:       11,
		Level2:               []string{"食品药品与粮食大类", "农林牧渔大类", "轻工纺织大类", "生物与化工大类", "财经商贸大类"},
		Level3:               []string{"食品类", "药品与医疗器械类", "粮食类", "农业类", "畜牧业类", "渔业类", "轻化工类", "纺织服装类", "包装类", "经济贸易类"},
		Salary:               0,
	}

	analysis, err := analysisService.GenerateAdmissionAnalysis(priI, I)
	if err != nil {
		panic(err)
	} else {
		println(analysis.Analysis)
	}
}
