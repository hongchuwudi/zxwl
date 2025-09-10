// main.go
package main

import (
	"bufio"
	"fmt"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/service"
	"os"
	"strings"
)

func main1() {
	// 设置配置
	cfg := config.LoadConfig()

	// 创建分类服务
	classificationService := service.NewClassificationService(cfg)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🎯 兴趣爱好专业分类测试工具")
	fmt.Println("📝 输入你的兴趣爱好，AI会帮你匹配最合适的专业分类")
	fmt.Println("⏎  直接按回车退出")
	fmt.Println("")

	for {
		fmt.Print("➡️  请输入兴趣爱好: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("👋 再见！")
			break
		}

		fmt.Printf("🔍 正在分析: \"%s\"...\n", input)

		// 调用分类服务
		result, err := classificationService.ClassifyInterests(input)
		if err != nil {
			fmt.Printf("❌ 错误: %v\n\n", err)
			continue
		}

		// 显示漂亮的结果
		fmt.Println("\n📊 分析结果:")
		fmt.Printf("  置信度: %.0f%%\n", result.Confidence*100)

		fmt.Println("\n🎓 推荐大类专业:")
		for i, category := range result.MajorCategories {
			fmt.Printf("  %d. %s\n", i+1, category)
		}

		fmt.Println("\n📚 推荐小类专业:")
		for i, category := range result.MinorCategories {
			fmt.Printf("  %d. %s\n", i+1, category)
		}

		fmt.Println("\n" + strings.Repeat("─", 60))
		fmt.Println()
	}
}
