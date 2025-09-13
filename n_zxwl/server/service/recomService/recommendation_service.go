package recomService

import (
	"fmt"
	"gorm.io/gorm"
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/repositories/admissionRepo"
	"mymod/repositories/commonRepo"
	"mymod/repositories/recommendRepo"
	"mymod/repositories/schoolRepo"
	"mymod/repositories/scoreRepo"
	"mymod/repositories/specialRepo"
	"regexp"
	"strconv"
	"strings"
)

type RecommendationRepository struct {
	db                       *gorm.DB
	admissionRepo            *admissionRepo.AdmissionRepository
	universityRepo           *schoolRepo.UniversityRepository
	specialRepo              *specialRepo.SpecialRepository
	recommendationRepo       *recommendRepo.RecommendationRepo
	classifyService          *ClassificationService
	admissionAnalysisService *AdmissionAnalysisService
}

func NewRecommendationRepository(db *gorm.DB) *RecommendationRepository {
	// 加载配置
	cfg := config.LoadConfig()

	return &RecommendationRepository{
		db:                       db,
		admissionRepo:            admissionRepo.NewAdmissionRepository(db),
		universityRepo:           schoolRepo.NewUniversityRepository(db),
		specialRepo:              specialRepo.NewSpecialRepository(db),
		recommendationRepo:       recommendRepo.NewRecommendationRepo(db),
		classifyService:          NewClassificationService(cfg), // 使用配置创建分类服务
		admissionAnalysisService: NewAdmissionAnalysisService(cfg),
	}
}

// GetRecommendations 综合推荐算法
func (r *RecommendationRepository) GetRecommendations(priInput sqlModel.UserPriInput) (sqlModel.UserInput, sqlModel.Recommendation, error) {
	var input sqlModel.UserInput
	var recommendation sqlModel.Recommendation
	// 打印用户原始数据
	PrintRequestPriInfo(priInput)

	// TODO 1.用户数据预处理
	// TODO 1.1 处理用户的兴趣爱好 调用qwen语言处理模型来匹配专业分类标签表,挑选出合适的分类组
	interestsText := priInput.Interest
	if interestsText != "" {
		fmt.Printf("正在分析用户兴趣爱好: %s\n", interestsText)

		classificationResult, err := r.classifyService.ClassifyInterests(interestsText)
		if err != nil {
			return input, recommendation, fmt.Errorf("兴趣爱好分类失败: %v", err)
		}

		fmt.Printf("AI分类结果 - 置信度: %.2f\n", classificationResult.Confidence)
		fmt.Printf("推荐大类: %v\n", classificationResult.MajorCategories)
		fmt.Printf("推荐小类: %v\n", classificationResult.MinorCategories)

		// 将AI分类结果保存到处理后的用户输入中
		input.Level2 = classificationResult.MajorCategories // 保持为 []string
		input.Level3 = classificationResult.MinorCategories // 保持为 []string

		// 将AI的推荐结果保存到返回结果中去
		recommendation.InterestAnalysis = classificationResult.AnalyzeReasons

	} else {
		fmt.Println("用户未提供兴趣爱好信息，将基于分数和排名进行推荐")
		recommendation.InterestAnalysis = "用户未提供兴趣爱好信息"
	}

	// TODO 1.2 处理用户的生源地高考类型(文科,理科,体育文,体育理,历史类,物理类)......
	input.Year = priInput.Year

	if priInput.TypeName != "" {
		typeID, err := admissionRepo.NewAdmissionRepository(r.db).GetTypeIDByName(priInput.TypeName)
		if err != nil {
			return input, recommendation, fmt.Errorf("高考类型ID获取失败: %v", err)
		}
		input.TypeID = typeID
	} else {
		input.TypeID = 0
	}

	// TODO 1.3 转换目标省份和生源省份的name为省份ID
	var err1, err2 error
	// 先查生源省份ID,这个是必查的
	input.ProvinceID, err2 = commonRepo.NewProvinceRepository(r.db).GetProvinceIDByName(priInput.ProvinceName)
	if priInput.GoalProvinceName == "" {
		input.GoalProvinceID = 0
	} else {
		input.GoalProvinceID, err1 = commonRepo.NewProvinceRepository(r.db).GetProvinceIDByName(priInput.GoalProvinceName)
	}
	if err1 != nil {
		return input, recommendation, fmt.Errorf("目的省份ID获取失败: %v, %v", err1, err2)
	}
	if err2 != nil {
		return input, recommendation, fmt.Errorf("源省份ID获取失败: %v, %v", err1, err2)
	}

	// TODO 1.4 处理用户的是否为专科或者本科(专门针对北上津,他们的分数位次参考系有两套)
	// TODO 1.4 对应admission_batches中对应is_benke这个字段,本科报本,专,科特殊批,专科只能报专,特殊批次
	// TODO 1.4.1 如果用户是高贵的北上津爷,在分数转化前就要确定他们的批次了,因为他们的一分一段有两套表,
	// TODO 1.4.2 其他地方不用管
	if input.ProvinceID == 11 || input.ProvinceID == 12 || input.ProvinceID == 31 {
		if priInput.BatchName == "本科" {
			input.IsBenKe = 1
		} else if priInput.BatchName == "专科" {
			input.IsBenKe = 2
		} else {
			//	默认参考本科位次系
			input.IsBenKe = 1
		}
	} else {
		input.IsBenKe = 3
	}

	// TODO 1.5 处理用户的分数/位次,这里调用高见写的例程(逻辑如下)
	// TODO 1.5.1 处理用户的分数排名
	// TODO 1.5.2 根据排名将用户当前的分数等比转化为指定年份分数来处理(默认2024,如果后续有其他年份数据再做调整)
	// TODO 1.5.3 新高考物理类对应理科,新高考历史类对应物理类
	if priInput.Score == 0 {
		return input, recommendation, fmt.Errorf("用户分数不能为0")
	}
	if priInput.GoalYear == 0 {
		input.GoalYear = 2024
	} else {
		input.GoalYear = priInput.GoalYear

	}

	// 调用存储过程
	equivalentResult, err := scoreRepo.NewScoreRepository(r.db).GetEquivalentScore(
		int(priInput.Year),  // currentYear
		input.ProvinceID,    // provinceID
		input.TypeID,        // categoryID
		priInput.Score,      // currentScore
		int(input.GoalYear), // targetYear
		input.IsBenKe,       // batchID (需要根据batch名称获取)
	)

	if err != nil {
		return input, recommendation, fmt.Errorf("等效分数转换失败: %v", err)
	}

	// 将存储过程返回的结果填充到input中
	input.CurrRankStart = equivalentResult.CurrRankStart
	input.CurrRankEnd = equivalentResult.CurrRankEnd
	input.EquivalentScoreStart = equivalentResult.EquivalentScoreStart
	input.EquivalentScoreEnd = equivalentResult.EquivalentScoreEnd
	input.DebugInfo = equivalentResult.DebugInfo

	recommendation.DebugInfo = equivalentResult.DebugInfo

	// 提取debug-info步骤1的type类型
	if re := regexp.MustCompile(`步骤1: 考试类型转化: \d+ --> (\d+)`); re.MatchString(equivalentResult.DebugInfo) {
		matches := re.FindStringSubmatch(equivalentResult.DebugInfo)
		if len(matches) > 1 {
			if newTypeID, err := strconv.Atoi(matches[1]); err == nil {
				input.TypeID = newTypeID
			}
		}
	}

	// 提取debug-info步骤5的内容
	if re := regexp.MustCompile(`步骤5: 获取当前年份考试分数是否过线: (.+?)\n`); re.MatchString(equivalentResult.DebugInfo) {
		matches := re.FindStringSubmatch(equivalentResult.DebugInfo)
		if len(matches) > 1 {
			batchType := strings.TrimSpace(matches[1])
			switch {
			case strings.Contains(batchType, "本科"):
				input.IsBenKe = 1
			case strings.Contains(batchType, "专科"):
				input.IsBenKe = 2
			default:
				input.IsBenKe = 3
			}
		}
	}

	PrintRequestInfo(input)

	// TODO 2.根据各种条件对admission_special表进行专业筛选,以专业为先选专业
	// TODO 2.0 筛选条件:
	// TODO -----------------------------条件---------------------说明--------可选?--------
	// TODO			1 				GoalYear       				参考年份 	硬性条件
	// TODO 		2 				TypeId         				高考类型ID 	硬性条件
	// TODO 		3 				ProvinceID     				高考省份ID 	硬性条件
	// TODO 		4 			EquivalentScoreStart-End 		同位分数		硬性条件
	// TODO 		5 			CurrRankStart-End 				当前排名		硬性条件
	// TODO         6	            scoreGap                    冲刺分数		硬性条件
	// TODO 		6 				IsBenKe  					本专筛选		可选条件
	// TODO 		7 			 GoalProvinceID 				目标省份ID	可选条件
	// TODO 		8 			 Level2- Level3 				专业分类		可选条件
	// TODO 		9 				Salary  				  每年期望薪资	可选条件
	scoreGap := 18.0
	specialResults, err := r.recommendationRepo.GetSpecialRecommendations(input, scoreGap)
	if err != nil {
		return input, recommendation, fmt.Errorf("专业推荐失败: %v", err)
	}

	// TODO 3.1 根据最低分对admission_universities表筛选可捡漏的大学(调用学校推荐存储过程 ),以学校优先来选取学校
	// TODO 3.0 筛选条件:
	// TODO -----------------------------条件---------------------说明--------可选?--------
	// TODO			1 				GoalYear       				参考年份 	硬性条件
	// TODO 		2 				TypeId         				高考类型ID 	硬性条件
	// TODO 		3 				ProvinceID     				高考省份ID 	硬性条件
	// TODO 		4 			EquivalentScoreStart-End 		同位分数		硬性条件
	// TODO         5	            scoreGap                    冲刺分数		硬性条件
	// TODO 		6 			 GoalProvinceID 				目标省份ID	可选条件
	universityResults, err := r.recommendationRepo.GetUniversityRecommendations(input, scoreGap)

	if err != nil {
		return input, recommendation, fmt.Errorf("学校推荐失败: %v", err)
	}

	// TODO 4.计算录取率并二次筛序专业学校
	// TODO 4.1 计算方法:看项目里的计算方法markdown文件
	// TODO 注意:如果录取率相同,依然按照原来的顺序排序
	// 计算用户平均分数
	userAvgScore := (input.EquivalentScoreStart + input.EquivalentScoreEnd) / 2

	// TODO 4.2.1 录取专业:根据考生分数与专业最低分作差换,得到录取率
	// TODO 4.2.2 根据录取率排序 取专业前50个
	// 处理专业推荐结果
	processedSpecials := ProcessSpecialRecommendations(specialResults, userAvgScore, scoreGap)
	// 将处理结果保存到返回结果中
	recommendation.SpecialRecommendations = processedSpecials
	PrintRecSpecial(recommendation)

	// TODO 4.3.1 录取学校:根据考生分数与学校最低分作差换,得到录取率
	// TODO 4.3.1 根据录取率排序 将学校数组分为三个区间(冲稳保),在三个区间每个区间都选取一定数量(学校区间数量配比:5:8:3)
	// 处理学校推荐结果并分区间
	chongSchools, wenSchools, baoSchools := ClassifyAndDistribute(universityResults, userAvgScore, scoreGap)
	// 合并所有学校（按冲稳保顺序）
	allSchools := append(append(chongSchools, wenSchools...), baoSchools...)
	// 将处理结果保存到返回结果中
	recommendation.UniversityRecommendations = allSchools
	PrintRecSchool(recommendation)

	//// TODO 5.结合AI综合性分析
	//// TODO 说明 : 结合AI,考生分数,高考人数,高考地区,高考政策,当下环境等综合因素给出用户百字短言
	if priInput.FamilyPref != "" {
		admissionAnalysis, err := r.admissionAnalysisService.GenerateAdmissionAnalysis(priInput, input)
		if err != nil {
			return input, recommendation, fmt.Errorf("最终结果AI分析失败: %v", err)
		}
		// 将处理结果保存到返回结果中
		recommendation.AdmissionAnalysis = admissionAnalysis.Analysis
	} else {
		fmt.Println("用户没有家庭偏好,将不再使用AI分析")
		recommendation.AdmissionAnalysis = "用户没有家庭偏好,暂时不使用AI分析"
	}

	// 输出调试信息
	PrintRecAny(recommendation)

	// TODO 6.输出结果
	return input, recommendation, nil
}
