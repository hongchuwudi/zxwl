package main

import (
	"bufio"
	"fmt"
	"log"
	"mymod/config"
	"os"
	"strings"
)

func maifan() {
	// 如果要改发送方qq邮箱,密码,先删掉.env文件,再运行程序

	keyFile := "email_key.txt"
	configFile := "email_config.txt"
	// 控制台输入
	fmt.Println("QQ邮箱&授权码加密工具")
	createNewConfig(keyFile, configFile)
}

func createNewConfig(keyFile, configFile string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔐 创建加密邮箱配置(一定要先创建文件)")
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
	if err := config.CreateEncryptedConfig(email, password, keyFile, configFile); err != nil {
		log.Fatalf("创建配置失败: %v", err)
	}

	fmt.Println("\n🎉 配置创建成功！")
	fmt.Printf("⚠️  密钥文件: %s\n", keyFile)
	fmt.Printf("⚠️  配置文件: %s\n", configFile)
	fmt.Println("⚠️  请妥善保管这些文件！并且只有放到指定位置下才能生效")
	fmt.Println("⚠️  不要将这些文件提交到版本控制系统！")
}
