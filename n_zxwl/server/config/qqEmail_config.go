package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"
	"strings"
)

// EmailConfig Config 结构体定义
type EmailConfig struct {
	QQEmail    string `json:"qq_email"`
	QQPassword string `json:"qq_password"` // 实际上是授权码
	SMTPHost   string `json:"smtp_host"`
	SMTPPort   int    `json:"smtp_port"`
}

// 确保目录存在
func ensureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0700) // 只允许所有者读写执行
}

// 加密函数
func encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建加密块失败: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM模式失败: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("生成随机数失败: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return []byte(encoded), nil
}

// 解密函数
func decrypt(key, encryptedText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建解密块失败: %v", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(encryptedText))
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM模式失败: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return nil, fmt.Errorf("密文太短")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// 生成加密密钥（32字节）
func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256需要32字节密钥
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("生成随机密钥失败: %v", err)
	}
	return key, nil
}

// 保存密钥到文件
func saveKeyToFile(key []byte, filename string) error {
	// 确保目录存在
	if err := ensureDir(filepath.Dir(filename)); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)
	return ioutil.WriteFile(filename, []byte(encodedKey), 0600) // 只允许所有者读写
}

// 从文件加载密钥
func loadKeyFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取密钥文件失败: %v", err)
	}

	// 移除可能的换行符
	keyStr := strings.TrimSpace(string(data))
	fmt.Println(base64.StdEncoding.DecodeString(keyStr))
	return base64.StdEncoding.DecodeString(keyStr)
}

// CreateEncryptedConfig 创建加密配置文件
func CreateEncryptedConfig(email, password, keyFile, configFile string) error {
	// 生成或加载密钥
	var key []byte
	var err error

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		// 密钥文件不存在，生成新密钥
		fmt.Printf("生成新的加密密钥到: %s\n", keyFile)
		key, err = generateRandomKey()
		if err != nil {
			return err
		}
		if err := saveKeyToFile(key, keyFile); err != nil {
			return fmt.Errorf("保存密钥失败: %v", err)
		}
		fmt.Printf("密钥已保存到: %s\n", keyFile)
	} else {
		// 加载现有密钥
		key, err = loadKeyFromFile(keyFile)
		if err != nil {
			return fmt.Errorf("加载密钥失败: %v", err)
		}
		fmt.Printf("使用现有密钥: %s\n", keyFile)
	}

	// 创建配置对象
	config := EmailConfig{
		QQEmail:    email,
		QQPassword: password,
		SMTPHost:   "smtp.qq.com",
		SMTPPort:   587,
	}

	// 序列化为JSON
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 加密配置
	encryptedConfig, err := encrypt(key, configJSON)
	if err != nil {
		return fmt.Errorf("加密配置失败: %v", err)
	}

	// 确保配置目录存在
	if err := ensureDir(filepath.Dir(configFile)); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 保存加密配置
	if err := ioutil.WriteFile(configFile, encryptedConfig, 0600); err != nil {
		return fmt.Errorf("保存加密配置失败: %v", err)
	}

	fmt.Printf("加密配置文件已创建: %s\n", configFile)
	fmt.Printf("QQ邮箱: %s\n", email)
	return nil
}

// GetEmailConfig 获取解密后的邮箱配置(核心)
func GetEmailConfig(keyFile, configFile string) (*EmailConfig, error) {
	// 加载密钥
	key, err := loadKeyFromFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("加载密钥失败: %v", err)
	}

	// 读取加密配置
	encryptedData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("读取加密配置失败: %v", err)
	}

	// 解密配置
	decrypted, err := decrypt(key, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("解密配置失败: %v", err)
	}

	// 解析JSON配置
	var config EmailConfig
	if err := json.Unmarshal(decrypted, &config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %v", err)
	}

	return &config, nil
}

// Validate 验证配置有效性
func (c *EmailConfig) Validate() error {
	if c.QQEmail == "" {
		return fmt.Errorf("QQ邮箱不能为空")
	}
	if !strings.Contains(c.QQEmail, "@qq.com") {
		return fmt.Errorf("请输入有效的QQ邮箱地址")
	}
	if c.QQPassword == "" {
		return fmt.Errorf("QQ邮箱授权码不能为空")
	}
	if c.SMTPHost == "" {
		return fmt.Errorf("SMTP主机不能为空")
	}
	if c.SMTPPort <= 0 {
		return fmt.Errorf("SMTP端口必须大于0")
	}
	return nil
}

// Print 打印配置信息（隐藏敏感信息）
func (c *EmailConfig) Print() {
	fmt.Printf("QQ邮箱: %s\n", c.QQEmail)
	fmt.Printf("SMTP服务器: %s:%d\n", c.SMTPHost, c.SMTPPort)
	fmt.Printf("授权码: %s\n", strings.Repeat("*", len(c.QQPassword)))
}
