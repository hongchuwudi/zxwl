package main

import (
	"fmt"
	"log"
	"mymod/config"
	"mymod/model/sqlModel"
	"mymod/repositories/specialRepo"
)

func test(id int) (*sqlModel.SpecialDetailResponse, error) {
	// 数据库连接配置
	db := config.GetDB()

	// 创建 Repository 实例
	repo := specialRepo.NewSpecialRepository(db)

	// 测试 GetDetailByID 函数（假设专业ID为1）
	detail, err := repo.GetDetailByID(id)
	if err != nil {
		log.Fatal("获取专业详情失败:", err)
		return nil, err
	}

	// 打印基本信息
	fmt.Printf("专业名称: %s\n", detail.SpecialDetail.Name)
	fmt.Printf("专业代码: %s\n", detail.SpecialDetail.Code)
	fmt.Printf("学历层次: %s\n", detail.SpecialDetail.Degree)
	fmt.Printf("男女比例: %s\n", detail.SpecialDetail.GenderRatio)
	// 平均薪资
	fmt.Printf("平均薪资: %d\n\n", detail.SpecialDetail.AvgSalary)

	// 打印关联数据数量（可选）
	fmt.Printf("大文本内容数量: %d\n", len(detail.SpecialContents))
	fmt.Printf("就业率记录数量: %d\n", len(detail.EmploymentRates))
	fmt.Printf("名校示例数量: %d\n", len(detail.FamousSchools))
	fmt.Printf("视频数量: %d\n", len(detail.Videos))
	fmt.Printf("印象标签数量: %d\n", len(detail.ImpressionTags))
	fmt.Printf("就业分布数量: %d\n", len(detail.JobDistributions))
	fmt.Printf("薪资数据数量: %d\n", len(detail.SalaryData))
	fmt.Printf("开设院校数量: %d\n", len(detail.UniversitySpecialInfo))
	return detail, nil
}
