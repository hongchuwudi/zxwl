package oldManyMgrApi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	config2 "mymod/config"
	"mymod/model/oldModel"
	utils2 "mymod/utils"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strconv"
	"time"
)

// GetVerifyCodeHandler 功能：发送验证码到指定邮箱。
func GetVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := oldModel.VerifyCodeResponse{Error: 1001}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	var req oldModel.VerifyCodeRequest
	// 从请求体中解析JSON数据，并将其存储到req结构体中，所以req.Email可以用于获取前端传来的邮箱地址
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs := json.NewEncoder(w).Encode(response)
		if errs != nil {
			return
		}
		return
	}

	// 获取数据库连接
	db := config2.GetDB()

	// 生成6位数字验证码
	rh := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", rh.Intn(1000000))
	expiresAt := time.Now().Add(15 * time.Minute)

	// 存储验证码到数据库（需要先创建verify_codes表）
	err := db.Exec(
		"INSERT INTO user_verify_codes (email, code, expires_at,created_at) VALUES (?, ?, ?, now()) ON DUPLICATE KEY UPDATE code= ?, expires_at= ?",
		req.Email, code, expiresAt, code, expiresAt,
	).Error
	if err != nil {
		log.Printf("存储验证码失败: %v", err)
		errs := json.NewEncoder(w).Encode(response)
		if errs != nil {
			return
		}
		return
	}

	// 发送验证码到邮箱
	// 发送验证码到邮箱
	to := []string{req.Email}

	// 获取邮箱配置
	proPath, pathErr := utils2.GetProjectRoot()
	if pathErr != nil {
		return
	}
	keyFile := filepath.Dir(proPath) + "\\env\\email_key.txt"
	configFile := filepath.Dir(proPath) + "\\env\\email_config.txt"
	configs, err := config2.GetEmailConfig(keyFile, configFile)
	if err != nil {
		return
	}

	// 加载邮件模板
	templateData := utils2.EmailTemplateData{
		VerificationCode: code,
	}

	emailBody, err := utils2.LoadEmailTemplate("email_verify_code.html", templateData)
	if err != nil {
		log.Printf("加载邮件模板失败: %v", err)
		// 可以回退到简单文本
		emailBody = fmt.Sprintf("您的验证码是：%s，有效期15分钟", code)
	}

	auth := smtp.PlainAuth("", configs.QQEmail, configs.QQPassword, configs.SMTPHost)
	subject := "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte("验证码通知")) + "?="

	// 构建完整的邮件内容
	fullBody := fmt.Sprintf(`To: %s
From: "智选高考志愿服务平台" <%s>
Subject: %s
Content-Type: text/html; charset=UTF-8

%s`, req.Email, configs.QQEmail, subject, emailBody)

	err = smtp.SendMail(
		configs.SMTPHost+":"+strconv.Itoa(configs.SMTPPort),
		auth,
		configs.QQEmail,
		to,
		[]byte(fullBody),
	)
	if err != nil {
		log.Printf("发送邮件失败(short short response可以忽略): %v", err)
	}

	log.Printf("邮件服务器: %s:%s", configs.SMTPHost, configs.SMTPPort)
	log.Printf("发件人: %s", configs.QQEmail)
	log.Printf("收件人: %v", to)
	log.Printf("验证码%s已发送到 %s", code, req.Email)

	response.Error = 0
	response.Msg = "验证码已发送"
	json.NewEncoder(w).Encode(response)
	return
}
