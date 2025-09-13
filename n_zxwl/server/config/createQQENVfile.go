package config

import (
	"bufio"
	"fmt"
	"log"
	"mymod/utils"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 如果要改发送方qq邮箱,密码,先删掉.env文件,再运行程序

	proPath, _ := utils.GetProjectRoot() // 获取项目根目录
	keyFile := filepath.Dir(proPath) + "\\env\\email_key.enc"
	configFile := filepath.Dir(proPath) + "\\env\\email_config.enc"
	// 控制台输入
	fmt.Println("邮箱配置文件创建工具")
	createNewConfig(keyFile, configFile)

	// 正常使用：读取配置
	useConfig(keyFile, configFile)
}

func createNewConfig(keyFile, configFile string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔐 创建加密邮箱配置(如果有先删除)")
	fmt.Println("====================")
	fmt.Printf("密钥文件: %s\n", keyFile)
	fmt.Printf("配置文件: %s\n", configFile)
	fmt.Println("")

	// 获取QQ邮箱
	fmt.Print("请输入QQ邮箱: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	// 获取授权码（注意：不是QQ密码！）
	fmt.Print("请输入QQ邮箱授权码: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// 确认信息
	fmt.Printf("\n即将创建配置:\n")
	fmt.Printf("QQ邮箱: %s\n", email)
	fmt.Printf("授权码: %s\n", strings.Repeat("*", len(password)))
	fmt.Printf("密钥文件: %s\n", keyFile)
	fmt.Printf("配置文件: %s\n", configFile)
	fmt.Print("\n确认创建？(y/N): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("操作取消")
		return
	}

	// 创建加密配置
	if err := CreateEncryptedConfig(email, password, keyFile, configFile); err != nil {
		log.Fatalf("创建配置失败: %v", err)
	}

	fmt.Println("\n🎉 配置创建成功！")
	fmt.Printf("⚠️  密钥文件: %s\n", keyFile)
	fmt.Printf("⚠️  配置文件: %s\n", configFile)
	fmt.Println("⚠️  请妥善保管这些文件！并且只有放到指定位置下才能生效")
	fmt.Println("⚠️  不要将这些文件提交到版本控制系统！")
}

func useConfig(keyFile, configFile string) {
	config, err := GetEmailConfig(keyFile, configFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("配置文件不存在，请先运行:\n")
			fmt.Printf("  go run main.go create\n")
			fmt.Printf("当前查找路径:\n")
			fmt.Printf("  密钥文件: %s\n", keyFile)
			fmt.Printf("  配置文件: %s\n", configFile)
			return
		}
		log.Fatalf("获取配置失败: %v", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 显示配置信息
	fmt.Println("✅ 成功加载加密配置:")
	config.Print()

	fmt.Println("\n📧 配置可用于发送邮件")
}
