package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mymod/new_zxwl/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	QwenAPIKey string
}

func LoadConfig() *Config {
	config := &Config{DBHost: "", DBPort: "", DBUser: "", DBPassword: "", DBName: "", QwenAPIKey: ""}

	// 读取数据库配置文件
	if proPath, err := utils.GetProjectRoot(); err == nil {
		dbConfigFile := filepath.Join(filepath.Dir(proPath), "env", "database_config.txt")
		if content, err := os.ReadFile(dbConfigFile); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue // 跳过空行和注释
				}

				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])

					switch key {
					case "DBHost":
						config.DBHost = value
					case "DBPort":
						config.DBPort = value
					case "DBUser":
						config.DBUser = value
					case "DBPassword":
						config.DBPassword = value
					case "DBName":
						config.DBName = value
					}
				}
			}
		} else {
			log.Printf("Warning: 未找到文件: %v", err)
		}
	}

	// 读取API Key文件
	if proPath, err := utils.GetProjectRoot(); err == nil {
		keyFile := filepath.Join(filepath.Dir(proPath), "env", "qwen_api_ak.txt")
		if content, err := os.ReadFile(keyFile); err == nil {
			config.QwenAPIKey = strings.TrimSpace(string(content))
		}
	}

	return config
}

// GetDB 获取数据库连接单例
func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		cfg := LoadConfig()
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // 启用SQL日志
		})
		if err != nil {
			log.Fatal("数据库连接失败:", err)
		}
		log.Println("数据库连接成功")
	})
	return db
}
