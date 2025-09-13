package recomService

import (
	"fmt"
	"mymod/model/sqlModel"
	"strings"
)

// PrintRequestPriInfo 打印请求信息
func PrintRequestPriInfo(testInput sqlModel.UserPriInput) {

	fmt.Println("=== 测试用户输入详细信息 ===")
	fmt.Printf("高考年份(2025,2026....):               %d\n", testInput.Year)
	fmt.Printf("参考年份(2024,2023....):               %d\n", testInput.GoalYear)
	fmt.Printf("生源地省份名称(北京,上海....):           %s\n", testInput.ProvinceName)
	fmt.Printf("意向上学省份(北京,上海....):             %s\n", testInput.GoalProvinceName)
	fmt.Printf("高考类型(文科,理科....):                %s\n", testInput.TypeName)
	fmt.Printf("高考分数(565.0....):                   %.1f\n", testInput.Score)
	fmt.Printf("位次(1234....):                       %s\n", testInput.Rank)
	fmt.Printf("兴趣爱好(我喜欢搞单片机...):             %s\n", testInput.Interest)
	fmt.Printf("家庭希望(家长希望我可以从事服务类行业):    %s\n", testInput.FamilyPref)
	fmt.Printf("期望毕业后一年的工资(25000):             %.1f\n\n", testInput.Salary)

}

func PrintRequestInfo(input sqlModel.UserInput) {

	fmt.Printf("=== 处理后的用户输入信息 ===\n")
	fmt.Printf("高考年份:              %d\n", input.Year)
	fmt.Printf("参考年份:              %d\n", input.GoalYear)
	fmt.Printf("高考省份ID:            %d\n", input.ProvinceID)
	fmt.Printf("目标省份ID:            %d\n", input.GoalProvinceID)
	fmt.Printf("高考类型ID:            %d\n", input.TypeID)
	fmt.Printf("一分一段表中是否分本专:  %d", input.IsBenKe)
	switch input.IsBenKe {
	case 1:
		fmt.Printf(" (本科)\n")
	case 2:
		fmt.Printf(" (专科)\n")
	case 3:
		fmt.Printf(" (不分本专)\n")
	default:
		fmt.Printf(" (未知)\n")
	}
	join1 := strings.Join(input.Level2, ", ")
	join2 := strings.Join(input.Level3, ", ")
	if join1 == "" || join2 == "" {
		fmt.Printf("专业大类:             -无-\n")
		fmt.Printf("专业小类:             -无-\n")
	} else {
		fmt.Printf("专业大类:             %v\n", join1)
		fmt.Printf("专业小类:             %v\n", join2)
	}

	fmt.Printf("期望工资:              %.1f\n\n", input.Salary)
	fmt.Printf("同位次分数例程-调试信息详情:\n%s\n", input.DebugInfo)
}

// PrintRecSchool 打印推荐学校信息
func PrintRecSchool(rec sqlModel.Recommendation) {
	// 打印大学推荐结果
	uNumber := rec.UniversityRecommendations

	if len(rec.UniversityRecommendations) > 0 {
		fmt.Printf("\n🏫 === 大学推荐 (%d所) ===\n", len(uNumber))
		// 按录取概率分组显示
		var chong, wen, bao []sqlModel.UniversityRecommendation
		for i := 0; i < len(uNumber); i++ {
			if rec.UniversityRecommendations[i].AdmissionProbability <= 0.4 {
				chong = append(chong, rec.UniversityRecommendations[i])
			} else if rec.UniversityRecommendations[i].AdmissionProbability > 0.5 &&
				rec.UniversityRecommendations[i].AdmissionProbability <= 0.9 {
				wen = append(wen, rec.UniversityRecommendations[i])
			} else {
				bao = append(bao, rec.UniversityRecommendations[i])
			}
		}
		if len(chong) > 0 {
			fmt.Printf("\n🚀 【冲刺院校】\n")
			for i, u := range chong {
				fmt.Printf("%d. %s (%.1f%%) - %.1f分\n", i+1, u.SchoolName, u.AdmissionProbability*100, u.MinScore)
			}
		}
		if len(wen) > 0 {
			fmt.Printf("\n✅ 【稳妥院校】\n")
			for i, u := range wen {
				fmt.Printf("%d. %s (%.1f%%) - %.1f分\n", i+1, u.SchoolName, u.AdmissionProbability*100, u.MinScore)
			}
		}
		if len(bao) > 0 {
			fmt.Printf("\n🛡️  【保底院校】\n")
			for i, u := range bao {
				fmt.Printf("%d. %s (%.1f%%) - %.1f分\n", i+1, u.SchoolName, u.AdmissionProbability*100, u.MinScore)
			}
		}
	} else {
		fmt.Println("\n❌ 无大学推荐结果")
	}
}

// PrintRecSpecial 打印推荐专业信息
func PrintRecSpecial(rec sqlModel.Recommendation) {
	// 打印专业推荐结果
	if len(rec.SpecialRecommendations) > 0 {
		fmt.Printf("\n🎓 === 专业推荐 (%d个) ===\n", len(rec.SpecialRecommendations))

		// 定义列宽
		const (
			schoolWidth   = 25
			specialWidth  = 30
			probWidth     = 12
			scoreWidth    = 25
			rankWidth     = 20
			salaryWidth   = 15
			categoryWidth = 30
			batchWidth    = 25
			subjectWidth  = 20
		)

		// 打印表头
		fmt.Printf("%-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s\n",
			schoolWidth, "🏫 学校名称",
			specialWidth, "📚 专业名称",
			probWidth, "🎯 录取概率",
			scoreWidth, "📊 分数范围",
			rankWidth, "📈 排名范围",
			salaryWidth, "💰 平均薪资",
			categoryWidth, "🏷️ 专业分类",
			batchWidth, "📝 批次/招生类型",
			subjectWidth, "📚 选科要求")

		// 打印分隔线
		fmt.Println(strings.Repeat("-", schoolWidth+specialWidth+probWidth+scoreWidth+rankWidth+salaryWidth+categoryWidth+batchWidth+subjectWidth+8))

		for _, r := range rec.SpecialRecommendations {
			// 格式化学校名称和专业名称
			schoolName := truncateString(r.SchoolName, schoolWidth)
			specialName := truncateString(r.SpecialName, specialWidth)

			// 格式化录取概率
			prob := fmt.Sprintf("%.1f%%", r.AdmissionProbability*100)

			// 格式化分数
			var scoreInfo string
			if r.ScoreMaxScore == 0 || r.ScoreAvgScore == 0 {
				scoreInfo = fmt.Sprintf("最低: %.1f分", r.ScoreMinScore)
			} else {
				scoreInfo = fmt.Sprintf("%.1f-%.1f分(均%.1f)", r.ScoreMinScore, r.ScoreMaxScore, r.ScoreAvgScore)
			}

			// 格式化排名
			var rankInfo string
			if r.ScoreMaxRank == 0 {
				rankInfo = fmt.Sprintf("最低: %d名", r.ScoreMinRank)
			} else {
				rankInfo = fmt.Sprintf("%d-%d名", r.ScoreMinRank, r.ScoreMaxRank)
			}

			// 格式化薪资
			var salaryInfo string
			if r.SpecialAvgSalary == 0 {
				salaryInfo = "暂无数据"
			} else {
				salaryInfo = fmt.Sprintf("%.1f元/年", r.SpecialAvgSalary)
			}

			// 格式化专业分类
			category := truncateString(fmt.Sprintf("%s/%s/%s", r.SpecialLevel1Name, r.SpecialLevel2Name, r.SpecialLevel3Name), categoryWidth)

			// 格式化批次和招生类型
			var batchInfo string
			if r.ZslxName == "" {
				batchInfo = truncateString(fmt.Sprintf("%s | 招生类型: 暂无", r.BatchName), batchWidth)
			} else {
				batchInfo = truncateString(fmt.Sprintf("%s | %s", r.BatchName, r.ZslxName), batchWidth)
			}

			// 格式化选科要求
			subjectReq := truncateString(r.SpecialSubjectRequirements, subjectWidth)
			if subjectReq == "" {
				subjectReq = "暂无"
			}

			// 打印行数据
			fmt.Printf("%-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s\n",
				schoolWidth, schoolName,
				specialWidth, specialName,
				probWidth, prob,
				scoreWidth, scoreInfo,
				rankWidth, rankInfo,
				salaryWidth, salaryInfo,
				categoryWidth, category,
				batchWidth, batchInfo,
				subjectWidth, subjectReq)
		}
	} else {
		fmt.Println("\n❌ 无专业推荐结果")
	}
}

// PrintRecAny 打印推荐分析结果
func PrintRecAny(rec sqlModel.Recommendation) {
	fmt.Printf("\n=== 推荐结果分析 ===\n")
	fmt.Printf("📊 %s\n", rec.AdmissionAnalysis)
	fmt.Printf("🎯 %s\n", rec.InterestAnalysis)

	fmt.Printf("\n✨ === 推荐总结 ===\n")
	fmt.Printf("共推荐 %d 个专业和 %d 所大学\n",
		len(rec.SpecialRecommendations),
		len(rec.UniversityRecommendations))
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
