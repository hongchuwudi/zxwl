// file: config/oss_config.go
package config

import (
	"log"
	"mymod/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	Domain          string // CDN域名或OSS外网域名
}

var (
	ossConfig *OSSConfig
	ossClient *oss.Client
	ossBucket *oss.Bucket
	ossOnce   sync.Once
)

// LoadOSSConfig 从配置文件加载OSS配置
func LoadOSSConfig() *OSSConfig {
	config := &OSSConfig{
		Endpoint:        "",
		AccessKeyID:     "",
		AccessKeySecret: "",
		BucketName:      "",
		Domain:          "",
	}

	// 读取OSS配置文件
	if proPath, err := utils.GetProjectRoot(); err == nil {
		ossConfigFile := filepath.Join(filepath.Dir(proPath), "env", "oss_config.txt")
		if content, err := os.ReadFile(ossConfigFile); err == nil {
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
					case "OSSEndpoint":
						config.Endpoint = value
					case "OSSAccessKeyID":
						config.AccessKeyID = value
					case "OSSAccessKeySecret":
						config.AccessKeySecret = value
					case "OSSBucketName":
						config.BucketName = value
					case "OSSDomain":
						config.Domain = value
					}
				}
			}
		} else {
			log.Printf("Warning: 未找到OSS配置文件: %v", err)
		}
	}

	return config
}

// InitOSS 初始化OSS连接
func InitOSS() {
	ossOnce.Do(func() {
		ossConfig = LoadOSSConfig()

		// 检查必要的配置项
		if ossConfig.Endpoint == "" || ossConfig.AccessKeyID == "" ||
			ossConfig.AccessKeySecret == "" || ossConfig.BucketName == "" {
			log.Fatal("OSS配置不完整，请检查env/oss_config.txt文件")
		}

		// 创建OSS客户端
		client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyID, ossConfig.AccessKeySecret)
		if err != nil {
			log.Fatalf("OSS初始化失败: %v", err)
		}

		// 获取Bucket
		bucket, err := client.Bucket(ossConfig.BucketName)
		if err != nil {
			log.Fatalf("获取OSS Bucket失败: %v", err)
		}

		ossClient = client
		ossBucket = bucket
		log.Println("OSS初始化成功")
	})
}

func GetOSSBucket() *oss.Bucket {
	if ossBucket == nil {
		InitOSS()
	}
	return ossBucket
}

func GetOSSConfig() *OSSConfig {
	if ossConfig == nil {
		InitOSS()
	}
	return ossConfig
}
