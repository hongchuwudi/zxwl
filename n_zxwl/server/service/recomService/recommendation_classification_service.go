// service/classification_service.go 处理用户的兴趣爱好
package recomService

import (
	"fmt"
	"mymod/config"
	"mymod/service"
	"strings"
)

type ClassificationService struct {
	qwenService *service.QwenService
}

type ClassificationRequest struct {
	InterestsText string `json:"interests_text"`
}

type ClassificationResponse struct {
	MajorCategories []string `json:"major_categories"` // 5个最匹配的大类
	MinorCategories []string `json:"minor_categories"` // 10个最匹配的小类
	Confidence      float64  `json:"confidence"`       // 整体匹配置信度
	AnalyzeReasons  string   `json:"analyze_reasons"`  // 分析原因
}

// 定义所有大类名称
var majorCategories = []string{
	"哲学", "经济学", "法学", "教育学", "文学", "历史学",
	"理学", "工学", "农学", "医学", "管理学", "艺术学",
	"农林牧渔大类", "资源环境与安全大类", "能源动力与材料大类",
	"土木建筑大类", "水利大类", "装备制造大类", "生物与化工大类",
	"轻工纺织大类", "食品药品与粮食大类", "交通运输大类", "电子与信息大类",
	"医药卫生大类", "财经商贸大类", "旅游大类", "文化艺术大类",
	"新闻传播大类", "教育与体育大类", "公安与司法大类", "公共管理与服务大类",
	"能源动力与材料大类", "文化艺术大类", "土木建筑大类", "装备制造大类",
	"财经商贸大类", "电子与信息大类", "医药卫生大类", "新闻传播大类",
	"农林牧渔大类", "资源环境与安全大类", "轻工纺织大类", "教育与体育大类",
	"食品药品与粮食大类", "交通运输大类", "公安与司法大类", "旅游大类",
	"公共管理与服务大类", "生物与化工大类", "水利大类",
}

// 定义所有小类名称
var minorCategories = []string{
	"哲学类", "经济学类", "财政学类", "金融学类",
	"经济与贸易类", "法学类", "政治学类", "社会学类", "民族学类", "马克思主义理论类", "公安学类",
	"教育学类", "体育学类", "中国语言文学类", "外国语言文学类", "新闻传播学类", "历史学类", "数学类",
	"物理学类", "化学类", "天文学类", "地理科学类", "大气科学类", "海洋科学类", "地球物理学类",
	"地质学类", "生物科学类", "心理学类", "统计学类", "力学类", "机械类", "仪器类", "材料类",
	"能源动力类", "电气类", "电子信息类", "自动化类", "计算机类", "土木类", "水利类",
	"测绘类", "化工与制药类", "地质类", "矿业类", "纺织类", "轻工类", "交通运输类",
	"海洋工程类", "航空航天类", "兵器类", "核工程类", "农业工程类", "林业工程类",
	"环境科学与工程类", "生物医学工程类", "食品科学与工程类", "建筑类", "安全科学与工程类",
	"生物工程类", "公安技术类", "植物生产类", "自然保护与环境生态类", "动物生产类",
	"动物医学类", "林学类", "水产类", "草学类", "基础医学类", "临床医学类", "口腔医学类",
	"公共卫生与预防医学类", "中医学类", "中西医结合类", "药学类", "中药学类", "法医学类",
	"医学技术类", "护理学类", "管理科学与工程类", "工商管理类", "农业经济管理类", "公共管理类",
	"图书情报与档案管理类", "物流管理与工程类", "工业工程类", "电子商务类", "旅游管理类",
	"艺术学理论类", "音乐与舞蹈学类", "戏剧与影视学类", "美术学类", "设计学类", "农业类", "林业类",
	"畜牧业类", "渔业类", "资源勘查类", "地质类", "测绘地理信息类", "石油与天然气类", "煤炭类", "金属与非金属矿类",
	"气象类", "环境保护类", "安全类", "电力技术类", "热能与发电工程类", "新能源发电工程类", "黑色金属材料类",
	"有色金属材料类", "非金属材料类", "建筑材料类", "建筑设计类", "城乡规划与管理类", "土建施工类", "建筑设备类",
	"建设工程管理类", "市政工程类", "房地产类", "水文水资源类", "水利工程与管理类", "水利水电设备类",
	"水土保持与水环境类", "机械设计制造类", "机电设备类", "自动化类", "轨道装备类", "船舶与海洋工程装备类",
	"航空装备类", "汽车制造类", "生物技术类", "化工技术类", "轻化工类", "包装类", "印刷类", "纺织服装类",
	"食品类", "药品与医疗器械类", "粮食类", "铁道运输类", "道路运输类", "水上运输类", "航空运输类", "管道运输类",
	"城市轨道交通类", "邮政类", "电子信息类", "计算机类", "通信类", "临床医学类", "护理类", "药学类", "医学技术类",
	"康复治疗类", "公共卫生与卫生管理类", "健康管理与促进类", "财政税务类", "金融类", "财务会计类", "统计类",
	"经济贸易类", "工商管理类", "电子商务类", "物流类", "旅游类", "餐饮类", "艺术设计类", "表演艺术类",
	"民族文化艺术类", "文化服务类", "新闻出版类", "广播影视类", "教育类", "语言类", "体育类", "公安管理类",
	"公安技术类", "侦查类", "法律实务类", "法律执行类", "司法技术类", "公共事业类", "公共管理类",
	"公共服务类", "中医药类", "安全防范类", "眼视光类", "集成电路类", "文秘类", "电力技术类",
	"热能与发电工程类", "新能源发电工程类", "艺术设计类", "黑色金属材料类", "有色金属材料类",
	"非金属材料类", "建筑设计类", "建筑材料类", "城乡规划与管理类", "土建施工类", "建筑设备类",
	"建设工程管理类", "市政工程类", "房地产类", "机械设计制造类", "机电设备类", "自动化类", "轨道装备类",
	"船舶与海洋工程装备类", "航空装备类", "汽车制造类", "电子信息类", "计算机类", "通信类", "集成电路类",
	"表演艺术类", "护理类", "药学类", "中医药类", "医学技术类", "康复治疗类", "文化服务类", "公共卫生与卫生管理类",
	"健康管理与促进类", "眼视光类", "财政税务类", "金融类", "财务会计类", "经济贸易类", "工商管理类", "电子商务类",
	"物流类", "新闻出版类", "农业类", "林业类", "畜牧业类", "渔业类", "资源勘查类", "轻化工类", "广播影视类", "地质类",
	"包装类", "测绘地理信息类", "印刷类", "石油与天然气类", "纺织服装类", "煤炭类", "教育类", "气象类", "环境保护类",
	"食品类", "安全类", "药品与医疗器械类", "语言类", "粮食类", "铁道运输类", "体育类", "道路运输类", "水上运输类",
	"航空运输类", "城市轨道交通类", "邮政类", "公安技术类", "侦查类", "法律实务类", "法律执行类", "司法技术类",
	"安全防范类", "旅游类", "公共事业类", "餐饮类", "公共管理类", "公共服务类", "生物技术类", "化工技术类",
	"水文水资源类", "水利工程与管理类", "水利水电设备类", "水土保持与水环境类", "交叉工程类", "文秘类",
}

func NewClassificationService(cfg *config.Config) *ClassificationService {
	return &ClassificationService{
		qwenService: service.NewQwenService(cfg),
	}
}

// ClassifyInterests 使用AI API对兴趣爱好文本进行分类
func (cs *ClassificationService) ClassifyInterests(interestsText string) (*ClassificationResponse, error) {
	if interestsText == "" {
		return nil, fmt.Errorf("兴趣爱好文本不能为空")
	}

	// 构建AI提示词
	prompt := cs.buildClassificationPrompt(interestsText)

	// 调用AI API
	result, err := cs.callAIClassification(prompt)
	if err != nil {
		return nil, fmt.Errorf("AI分类失败: %v", err)
	}

	return result, nil
}

// 构建AI提示词
func (cs *ClassificationService) buildClassificationPrompt(interestsText string) string {
	// 构建详细的提示词
	return fmt.Sprintf(`请根据用户的兴趣爱好文本，分析并匹配最合适的专业分类。

用户兴趣爱好: "%s"

可选的大类分类:
%s

可选的小类分类:
%s

请严格按照以下要求返回结果:
1. 返回最匹配的5个大类名称
2. 返回最匹配的10个小类名称
3. 返回整体匹配置信度(0-1之间的小数)
4. 返回的分析原因(也就是分析过程)一定要具体(250字左右)
5. 返回格式必须是严格的json:
{
    "major_categories": ["大类1", "大类2", "大类3", "大类4", "大类5"],
    "minor_categories": ["小类1", "小类2", "小类3", "小类4", "小类5", "小类6", "小类7", "小类8", "小类9", "小类10"],
    "analyze_reasons" : "这里是分析的思考过程",
	"confidence": 0.85
}

请确保返回的分类名称完全来自上面提供的列表，不要创造新的分类名称。`,
		interestsText,
		strings.Join(majorCategories, "、"),
		strings.Join(minorCategories, "、"))
}

// 调用AI API
func (cs *ClassificationService) callAIClassification(prompt string) (*ClassificationResponse, error) {
	// 使用新的TextGeneration方法
	aiText, err := cs.qwenService.TextGeneration(prompt, "text")
	if err != nil {
		return nil, err
	}

	// 使用新的ExtractJSONFromText方法
	var result ClassificationResponse
	if err := cs.qwenService.ExtractJSONFromText(aiText, &result); err != nil {
		return nil, err
	}

	// 验证结果
	if len(result.MajorCategories) != 5 {
		return nil, fmt.Errorf("返回的大类数量不正确")
	}
	if len(result.MinorCategories) != 10 {
		return nil, fmt.Errorf("返回的小类数量不正确")
	}
	if result.Confidence < 0 || result.Confidence > 1 {
		return nil, fmt.Errorf("置信度必须在0-1之间")
	}

	return &result, nil
}
