// Package service file: service/upload_service.go
package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"mymod/config"
	"mymod/utils"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type UploadService struct {
	bucket *oss.Bucket
}

func NewUploadService() *UploadService {
	return &UploadService{
		bucket: config.GetOSSBucket(),
	}
}

// UploadFile 上传文件到OSS
func (s *UploadService) UploadFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// 验证文件大小（限制为10MB）
	if fileHeader.Size > 10<<20 {
		return "", errors.New("文件大小不能超过10MB")
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 生成文件名
	fileName := utils.GenerateFileName(fileHeader.Filename)

	// 构建OSS对象路径
	objectKey := fmt.Sprintf("%s/%s/%s", folder, time.Now().Format("2006/01/02"), fileName)

	// 上传到OSS
	err = s.bucket.PutObject(objectKey, bytes.NewReader(fileBytes))
	if err != nil {
		return "", fmt.Errorf("上传到OSS失败: %v", err)
	}

	// 生成访问URL
	ossConfig := config.GetOSSConfig()
	fileURL := fmt.Sprintf("%s/%s", ossConfig.Domain, objectKey)

	return fileURL, nil
}

// UploadImage 上传图片文件（额外验证图片类型）
func (s *UploadService) UploadImage(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// 验证文件类型
	if !utils.IsImageFile(fileHeader.Filename) {
		return "", errors.New("只支持图片文件格式: jpg, jpeg, png, gif, bmp, webp, svg")
	}

	return s.UploadFile(fileHeader, folder)
}

// DeleteFile 删除OSS上的文件
func (s *UploadService) DeleteFile(fileURL string) error {
	ossConfig := config.GetOSSConfig()

	// 从完整的URL中提取object key
	objectKey := strings.TrimPrefix(fileURL, ossConfig.Domain+"/")
	if objectKey == fileURL {
		return errors.New("无效的文件URL")
	}

	err := s.bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// GetFileInfo 获取文件信息
func (s *UploadService) GetFileInfo(fileURL string) (map[string]interface{}, error) {
	ossConfig := config.GetOSSConfig()
	objectKey := strings.TrimPrefix(fileURL, ossConfig.Domain+"/")

	meta, err := s.bucket.GetObjectMeta(objectKey)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	info := map[string]interface{}{
		"size":         meta.Get("Content-Length"),
		"contentType":  meta.Get("Content-Type"),
		"lastModified": meta.Get("Last-Modified"),
		"etag":         meta.Get("ETag"),
	}

	return info, nil
}
