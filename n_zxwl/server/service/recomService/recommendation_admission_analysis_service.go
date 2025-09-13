package recomService

import (
	"fmt"
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/service"
	"strings"
)

type AdmissionAnalysisService struct {
	qwenService *service.QwenService
}

type AdmissionAnalysisRequest struct {
	UserPriInput sqlModel.UserPriInput
	UserInput    sqlModel.UserInput
}

type AdmissionAnalysisResponse struct {
	Analysis string `json:"analysis"` // 300-400字的详细建议
}

func NewAdmissionAnalysisService(cfg *config.Config) *AdmissionAnalysisService {
	return &AdmissionAnalysisService{
		qwenService: service.NewQwenService(cfg),
	}
}

// GenerateAdmissionAnalysis 生成高考志愿填报分析建议
func (a *AdmissionAnalysisService) GenerateAdmissionAnalysis(userPriInput sqlModel.UserPriInput, userInput sqlModel.UserInput) (*AdmissionAnalysisResponse, error) {
	if userPriInput.Year == 0 || userPriInput.Score == 0 {
		return nil, fmt.Errorf("考生信息不完整")
	}

	// 构建AI提示词
	prompt := a.buildAnalysisPrompt(userPriInput, userInput)

	// 调用AI API
	result, err := a.callAIAnalysis(prompt)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %v", err)
	}

	return result, nil
}

// 构建AI提示词
func (a *AdmissionAnalysisService) buildAnalysisPrompt(userPriInput sqlModel.UserPriInput, userInput sqlModel.UserInput) string {
	return fmt.Sprintf(`请为%d年%s考生提供志愿填报分析建议，要求300-400字，内容具体实用,且必须返回完整内容。

考生基本信息：
- 所在省份：%s，目标省份：%s
- 高考类型：%s，报考批次：%s
- 高考分数：%.1f分，排名情况：%s
- 兴趣爱好：%s
- 家庭偏好：%s
- 期望薪资：%.0f元/月

分数分析数据：
- 参考%d年等效分数：%.1f-%.1f分
- 对应排名区间：%d-%d名

请结合%d年高考政策环境、就业市场趋势、经济发展方向等因素，从以下维度提供专业建议：
1. 专业选择方向（结合兴趣爱好和就业前景匹配度）
2. 院校选择策略（考虑地域、学校层次、专业实力）
3. 职业发展路径规划
4. 风险控制和备选方案

要求分析客观专业，建议具体可操作，避免空泛表述。`,
		userPriInput.Year, userPriInput.ProvinceName,
		userPriInput.ProvinceName, userPriInput.GoalProvinceName,
		userPriInput.TypeName, userPriInput.BatchName,
		userPriInput.Score, userPriInput.Rank,
		userPriInput.Interest,
		userPriInput.FamilyPref,
		userPriInput.Salary,
		userInput.GoalYear, userInput.EquivalentScoreStart, userInput.EquivalentScoreEnd,
		userInput.CurrRankStart, userInput.CurrRankEnd,
		userPriInput.Year)
}

// 调用AI API
func (a *AdmissionAnalysisService) callAIAnalysis(prompt string) (*AdmissionAnalysisResponse, error) {
	// 使用QwenService进行文本生成
	aiText, err := a.qwenService.TextGeneration(prompt, "text")
	if err != nil {
		return nil, err
	}

	// 清理和分析AI返回的文本
	analysis := a.cleanAnalysisText(aiText)

	return &AdmissionAnalysisResponse{
		Analysis: analysis,
	}, nil
}

// 获取基础分析（备用方案）
func (a *AdmissionAnalysisService) getBasicAnalysis(userPriInput sqlModel.UserPriInput, userInput sqlModel.UserInput) string {
	var analysis strings.Builder

	analysis.WriteString(fmt.Sprintf("基于%d年%s高考情况分析：\n\n", userPriInput.Year, userPriInput.ProvinceName))

	analysis.WriteString("📊 分数定位：")
	if userInput.EquivalentScoreStart > 600 {
		analysis.WriteString("你的分数属于高分段，可冲刺985/211院校。")
	} else if userInput.EquivalentScoreStart > 500 {
		analysis.WriteString("你的分数处于中上水平，适合选择一本院校优势专业。")
	} else {
		analysis.WriteString("你的分数建议重点关注二本院校和优质专科院校。")
	}
	analysis.WriteString("\n\n")

	analysis.WriteString("🎯 专业建议：")
	if strings.Contains(userPriInput.Interest, "技术") || strings.Contains(userPriInput.Interest, "工程") {
		analysis.WriteString("建议考虑工科类专业，如计算机、电子信息等，就业前景较好。")
	} else if strings.Contains(userPriInput.Interest, "文学") || strings.Contains(userPriInput.Interest, "艺术") {
		analysis.WriteString("可考虑文史类或艺术类专业，注重实践能力的培养。")
	} else {
		analysis.WriteString("建议选择与兴趣爱好相关的专业，保持学习热情。")
	}
	analysis.WriteString("\n\n")

	analysis.WriteString("🏫 院校选择：")
	if userPriInput.FamilyPref != "" {
		analysis.WriteString(fmt.Sprintf("结合家庭偏好「%s」，", userPriInput.FamilyPref))
	}
	analysis.WriteString("建议选择办学实力强、就业率高的院校。\n\n")

	analysis.WriteString("💡 温馨提示：志愿填报要形成梯度，既有冲刺院校也要有保底选择。")

	return analysis.String()
}

// 清理和分析AI返回的文本 - 改进版本
func (a *AdmissionAnalysisService) cleanAnalysisText(aiText string) string {
	// 移除可能的JSON格式标记
	cleaned := strings.ReplaceAll(aiText, "```json", "")
	cleaned = strings.ReplaceAll(cleaned, "```", "")
	cleaned = strings.TrimSpace(cleaned)

	// 如果文本过长，智能截断
	if len([]rune(cleaned)) > 400 {
		cleaned = a.smartTruncate(cleaned, 400)
	}

	return cleaned
}

// 智能截断，确保在完整的句子处结束
func (a *AdmissionAnalysisService) smartTruncate(text string, maxLength int) string {
	runeText := []rune(text)
	if len(runeText) <= maxLength {
		return text
	}

	// 优先在句号处截断
	for i := maxLength; i > maxLength-50 && i > 0; i-- {
		if i < len(runeText) && (runeText[i] == '。' || runeText[i] == '！' || runeText[i] == '？') {
			return string(runeText[:i+1])
		}
	}

	// 其次在逗号处截断
	for i := maxLength; i > maxLength-30 && i > 0; i-- {
		if i < len(runeText) && runeText[i] == '，' {
			return string(runeText[:i]) + "。"
		}
	}

	// 最后在合适的词语边界处截断
	return string(runeText[:maxLength])
}
