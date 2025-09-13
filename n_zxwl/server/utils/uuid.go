// file: utils/uuid.go
package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为后备方案
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// GenerateFileName 生成带UUID的文件名
func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)
	uuid := GenerateUUID()

	// 如果原始文件名有值，保留部分信息
	if baseName != "" {
		// 只保留文件名（不含路径）
		baseName = filepath.Base(baseName)
		// 限制长度，避免过长
		if len(baseName) > 20 {
			baseName = baseName[:20]
		}
		return fmt.Sprintf("%s_%s%s", baseName, uuid, ext)
	}

	return fmt.Sprintf("%s%s", uuid, ext)
}

// IsImageFile 检查是否是图片文件
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}
	for _, imageExt := range imageExts {
		if ext == imageExt {
			return true
		}
	}
	return false
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}
