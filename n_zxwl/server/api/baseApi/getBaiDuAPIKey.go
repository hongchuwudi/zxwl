package baseApi

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"mymod/model/param"
	"mymod/utils"
	"net/http"
	"path/filepath"
	"strings"
)

// 加密密钥（32字节）
var encryptionKey = []byte("your-32-byte-encryption-key-here")

// GetAPIKeyHandler 功能：获取API密钥
func GetAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1, Message: "获取API密钥失败"}

	if r.Method != http.MethodGet {
		response.Message = "只支持GET请求"
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取项目根目录
	proPath, pathErr := utils.GetProjectRoot()
	if pathErr != nil {
		response.Message = "获取项目路径失败"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构建API密钥文件路径
	apiKeyPath := filepath.Dir(proPath) + "\\env\\baidu_api_ak.txt"

	// 读取API密钥文件
	apiKeyBytes, err := ioutil.ReadFile(apiKeyPath)
	if err != nil {
		log.Printf("读取API密钥文件失败: %v", err)
		response.Message = "读取API密钥文件失败"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	apiKey := strings.TrimSpace(string(apiKeyBytes))
	if apiKey == "" {
		response.Message = "API密钥为空"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 加密API密钥
	encryptedKey, err := EncryptAPIKey(apiKey)
	if err != nil {
		log.Printf("加密API密钥失败: %v", err)
		response.Message = "加密失败"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Error = 0
	response.Message = "成功获取API密钥"
	response.Data = encryptedKey
	json.NewEncoder(w).Encode(response)
}

// EncryptAPIKey 加密API密钥 - 最简单的方式
func EncryptAPIKey(apiKey string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	// 简单的PKCS7填充
	paddedData := pkcs7Pad([]byte(apiKey), aes.BlockSize)
	ciphertext := make([]byte, len(paddedData))

	// 使用ECB模式（最简单）
	for i := 0; i < len(paddedData); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], paddedData[i:i+aes.BlockSize])
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// PKCS7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7去除填充
func pkcs7Unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
