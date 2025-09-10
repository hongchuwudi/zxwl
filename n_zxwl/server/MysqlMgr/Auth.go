package MysqlMgr

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"mymod/Const"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/model/param"
	"mymod/new_zxwl/model/sqlModel"
	"mymod/new_zxwl/utils"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// AuthHandler 功能：用户登录认证，验证邮箱和密码是否匹配。
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.AuthResponse{Code: 1001} // 默认失败状态
	if r.Method != http.MethodPost {
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
		return
	}

	var req Const.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
		return
	}

	log.Printf("客户端发来的JSON数据: Email=%s, Passwd=%s\n", req.Email, req.Passwd)

	var userID int
	err = Db.QueryRow(
		"SELECT id FROM user WHERE email = ? AND pwd = ?",
		req.Email,
		req.Passwd,
	).Scan(&userID)

	if err == nil {
		response.Code = 0
	} else if err == sql.ErrNoRows {
		response.Code = 1001
	} else {
		log.Println("数据库查询错误:", err)
	}

	json.NewEncoder(w).Encode(response)
}

// GetVerifyCodeHandler 功能：发送验证码到指定邮箱。
func GetVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.VerifyCodeResponse{Error: 1001}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	var req Const.VerifyCodeRequest
	// 从请求体中解析JSON数据，并将其存储到req结构体中，所以req.Email可以用于获取前端传来的邮箱地址
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs := json.NewEncoder(w).Encode(response)
		if errs != nil {
			return
		}
		return
	}

	// 获取数据库连接
	db := config.GetDB()

	// 生成6位数字验证码
	rh := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", rh.Intn(1000000))
	expiresAt := time.Now().Add(15 * time.Minute)

	// 存储验证码到数据库（需要先创建verify_codes表）
	err := db.Exec(
		"INSERT INTO verify_codes (email, code, expires_at,created_at) VALUES (?, ?, ?, now()) ON DUPLICATE KEY UPDATE code= ?, expires_at= ?",
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
	proPath, pathErr := utils.GetProjectRoot()
	if pathErr != nil {
		return
	}
	keyFile := filepath.Dir(proPath) + "\\env\\email_key.txt"
	configFile := filepath.Dir(proPath) + "\\env\\email_config.txt"
	configs, err := config.GetEmailConfig(keyFile, configFile)
	if err != nil {
		return
	}

	// 加载邮件模板
	templateData := utils.EmailTemplateData{
		VerificationCode: code,
	}

	emailBody, err := utils.LoadEmailTemplate("email_verify_code.html", templateData)
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

// ChangePasswordHandler 功能:更改密码
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 只处理POST请求
	if r.Method != http.MethodPost {
		response := param.Response{
			Error:   1,
			Message: "只支持POST请求",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求数据
	var req sqlModel.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := param.Response{
			Error:   1,
			Message: "请求数据格式错误",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证必填字段
	if req.User == "" || req.Email == "" || req.Passwd == "" || req.VerifyCode == "" {
		response := param.Response{
			Error:   1,
			Message: "用户名、邮箱、密码和验证码不能为空",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 1. 验证验证码
	var dbCode string
	var expiresAt time.Time
	err := Db.QueryRow(
		"SELECT code, expires_at FROM verify_codes WHERE email = ? ORDER BY created_at DESC LIMIT 1",
		req.Email,
	).Scan(&dbCode, &expiresAt)

	log.Printf("正在从数据库中查询验证码")
	if err != nil || dbCode != req.VerifyCode || time.Now().After(expiresAt) {
		if err == sql.ErrNoRows {
			response := param.Response{
				Error:   1,
				Message: "验证码错误",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		} else if dbCode != req.VerifyCode {
			response := param.Response{
				Error:   1,
				Message: "验证码错误",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		} else if time.Now().After(expiresAt) {
			response := param.Response{
				Error:   1,
				Message: "验证码已过期",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		} else {
			log.Printf("验证码查询错误: %v", err)
			response := param.Response{
				Error:   1,
				Message: "验证码验证失败",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
		return
	}
	log.Printf("验证码验证通过")

	// 2. 验证用户身份
	var userID int
	var currentPassword string
	err = Db.QueryRow(
		"SELECT id, pwd FROM user WHERE email = ? AND name = ?",
		req.Email, req.User,
	).Scan(&userID, &currentPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			response := param.Response{
				Error:   1,
				Message: "用户不存在",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
		} else {
			log.Printf("用户查询错误: %v", err)
			response := param.Response{
				Error:   1,
				Message: "用户验证失败",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	// 4. 更新密码
	_, err = Db.Exec(
		"UPDATE user SET pwd = ? WHERE id = ?",
		req.Passwd, // 这里直接使用明文密码，建议使用加密后的密码
		userID,
	)
	if err != nil {
		log.Printf("密码更新失败: %v", err)
		response := param.Response{
			Error:   1,
			Message: "密码更新失败",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 5. 清除已使用的验证码
	_, err = Db.Exec("DELETE FROM verify_codes WHERE email = ?", req.Email)
	if err != nil {
		log.Printf("清除验证码失败: %v", err)
		// 这里不返回错误，因为主要操作已经成功
	}

	// 6. 返回成功响应
	response := param.Response{
		Error:   0,
		Message: "密码重置成功",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Printf("用户 %s 密码重置成功", req.User)
}

// UserRegisterHandler 功能：用户注册，创建新用户。
func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.RegisterResponse{Code: 1001}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	var req Const.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(response)
		return
	}

	// 1. 验证验证码
	var dbCode string
	var expiresAt time.Time
	err := Db.QueryRow(
		"SELECT code, expires_at FROM verify_codes WHERE email = ? ORDER BY created_at DESC LIMIT 1",
		req.Email,
	).Scan(&dbCode, &expiresAt)

	log.Printf("正在从数据库中查询验证码")
	if err != nil || dbCode != req.VerifyCode || time.Now().After(expiresAt) {
		if err == sql.ErrNoRows {
			response.Code = 1004
			response.Msg = "验证码错误"
		} else if dbCode != req.VerifyCode {
			response.Code = 1004
			response.Msg = "验证码错误"
		} else if time.Now().After(expiresAt) {
			response.Code = 1003
			response.Msg = "验证码已过期"
		} else {
			log.Printf("验证码查询错误: %v", err)
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Printf("验证码验证通过")
	// 2. 检查用户是否已存在
	var existingUserID int
	err = Db.QueryRow(
		"SELECT id FROM user WHERE email = ? OR name = ?",
		req.Email, req.UserName,
	).Scan(&existingUserID)

	if err == nil {
		response.Code = 1005
		response.Msg = "用户或电子邮件已存在"
		json.NewEncoder(w).Encode(response)
		return
	} else if err != sql.ErrNoRows {
		log.Printf("检查用户存在性错误: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Printf("用户不存在，准备创建新用户")
	// 4. 创建新用户
	_, err = Db.Exec(
		"INSERT INTO user (name, pwd,email) VALUES (?, ?, ?)",
		req.UserName, req.Password, req.Email,
	)

	if err != nil {
		log.Printf("创建用户失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 5. 注册成功，清除已使用的验证码
	_, err = Db.Exec("DELETE FROM verify_codes WHERE email = ?", req.Email)
	if err != nil {
		log.Printf("清除验证码失败: %v", err)
	}

	response.Code = 0
	response.Msg = "注册成功"
	json.NewEncoder(w).Encode(response)
}

// ProfileHandler 功能：获取用户个人信息。
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.ProfileResponse{Code: 1001}
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}
	var req Const.ProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Printf("客户端发来的JSON数据: Email=%s", req.Email)
	var (
		name     sql.NullString
		sex      sql.NullInt64
		graduate sql.NullInt64
		address  sql.NullString
		picture  sql.NullString
		profile  Const.Profile
	)
	err := Db.QueryRow(`
		SELECT name, email, sex, graduate, address, picture 
		FROM profile 
		WHERE email = ?`,
		req.Email,
	).Scan(&name, &profile.Email, &sex, &graduate, &address, &picture)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("邮箱不存在: %s", req.Email)
		} else {
			log.Printf("数据库错误: %v", err)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if name.Valid {
		profile.Name = &name.String
	}
	if sex.Valid {
		s := int(sex.Int64)
		profile.Sex = &s
	}
	if graduate.Valid {
		g := int(graduate.Int64)
		profile.Graduate = &g
	}
	if address.Valid {
		profile.Address = &address.String
	}
	if picture.Valid {
		profile.Picture = &picture.String
	}

	response.Code = 0
	response.Data = profile
	json.NewEncoder(w).Encode(response)
}

// UpdateProfileHandler 功能：更新用户个人信息。
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.ProfileUpdateResponse{Code: 1001}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	var req Const.ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Printf("客户端发来的JSON数据: Email=%s,G=%v", req.Email, req.Graduate)
	var setClauses []string
	var args []interface{}

	if req.Name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Sex != nil {
		setClauses = append(setClauses, "sex = ?")
		args = append(args, *req.Sex)
	}
	if req.Graduate != nil {
		setClauses = append(setClauses, "graduate = ?")
		args = append(args, *req.Graduate)
	}
	if req.Address != nil {
		setClauses = append(setClauses, "address = ?")
		args = append(args, *req.Address)
	}
	if req.Picture != nil {
		setClauses = append(setClauses, "picture = ?")
		args = append(args, *req.Picture)
	}
	if len(setClauses) == 0 {
		response.Code = 1001
		json.NewEncoder(w).Encode(response)
		return
	}

	args = append(args, req.Email)

	query := "UPDATE profile SET " + strings.Join(setClauses, ", ") + " WHERE email = ?"
	result, err := Db.Exec(query, args...)
	if err != nil {
		log.Printf("数据库更新失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取影响行数失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}
	if rowsAffected == 0 {
		log.Printf("邮箱不存在: %s", req.Email)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0
	json.NewEncoder(w).Encode(response)
}

// DeleteProfileHandler 删除用户
func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.ProfileDeleteResponse{Code: 1001} // 默认失败状态

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	var req Const.ProfileDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("删除请求: 邮箱=%s", req.Email)

	// 执行删除操作
	result, err := Db.Exec("DELETE FROM profile WHERE email = ?", req.Email)
	if err != nil {
		log.Printf("删除失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查是否实际删除了记录
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.Code = 1002 // 用户不存在
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0 // 成功
	json.NewEncoder(w).Encode(response)
}

// ProfileListHandler 个人信息列表
func ProfileListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.ProfileListResponse{Code: 1001}

	// 仅允许GET方法
	if r.Method != http.MethodGet {
		response.Msg = "Method not allowed"
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 执行查询
	rows, err := Db.Query(`
        SELECT address, graduate, picture, sex, email, name 
        FROM profile
    `)
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		response.Msg = "服务器内部错误"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var profiles []Const.ProfileResponsed
	for rows.Next() {
		var (
			address  sql.NullString
			graduate sql.NullInt64
			picture  sql.NullString
			sex      sql.NullInt64
			email    sql.NullString
			name     sql.NullString
		)

		if err := rows.Scan(
			&address,
			&graduate,
			&picture,
			&sex,
			&email,
			&name,
		); err != nil {
			log.Printf("数据扫描失败: %v", err)
			continue
		}

		profile := Const.ProfileResponsed{
			Address:  convertNullString(address),
			Graduate: convertNullYear(graduate),
			Picture:  convertNullString(picture),
			Sex:      convertNullInt(sex),
			Email:    convertNullString(email),
			Name:     convertNullString(name),
		}
		profiles = append(profiles, profile)
	}

	// 处理结果
	if len(profiles) == 0 {
		response.Msg = "无数据"
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0
	response.Data = profiles
	json.NewEncoder(w).Encode(response)
}
