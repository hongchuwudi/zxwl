package MysqlMgr

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	chat "mymod/ChatMgr"
	"mymod/Const"
	"mymod/new_zxwl/config"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB // 学长的go原生sql数据库连接

var validate *validator.Validate

// ChatMessage 结构体，用于存储聊天消息。
type ChatMessage struct {
	SchoolID int    `json:"school_id"`
	Email    string `json:"email"`
	Content  string `json:"content"`
}

// 功能：初始化验证器，注册自定义验证规则（如 datetime 格式验证）。
func init() {
	validate = validator.New()
	err := validate.RegisterValidation("datetime", func(fl validator.FieldLevel) bool {
		_, err := time.Parse("2006-01-02 15:04:05", fl.Field().String())
		return err == nil
	})
	if err != nil {
		return
	}
}

// InitDB 功能：初始化数据库连接，连接到 MySQL 数据库。
func InitDB() {
	var err error
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	} else {
		log.Println("数据库连接成功")
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}
}

// stringPtr 功能: 将sql.NullString转换为*string
func stringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// GetUnreadCount 获取未读消息数
func GetUnreadCount(email string, schoolID int) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM chat_messages 
		WHERE school_id = ? 
		AND created_at > (
			SELECT COALESCE(last_visit, '1970-01-01') 
			FROM user_visits 
			WHERE email = ? 
			AND school_id = ?
		)`
	err := Db.QueryRow(query, schoolID, email, schoolID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

// UnreadHandler 获取未读消息数
func UnreadHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	schoolID, _ := strconv.Atoi(r.URL.Query().Get("school_id"))

	count, err := GetUnreadCount(email, schoolID)
	if err != nil {
		http.Error(w, jsonResponse(500, "数据库查询失败", nil), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":  0,
		"count": count,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// BatchUnreadHandler 批量获取未读消息数
func BatchUnreadHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email     string `json:"email"`
		SchoolIDs []int  `json:"school_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, jsonResponse(400, "无效的请求参数", nil), http.StatusBadRequest)
		return
	}

	if len(request.SchoolIDs) == 0 {
		http.Error(w, jsonResponse(400, "school_ids不能为空", nil), http.StatusBadRequest)
		return
	}
	log.Printf("进入")
	placeholders := strings.Repeat("?,", len(request.SchoolIDs)-1) + "?"
	query := fmt.Sprintf(`
		SELECT school_id, COUNT(*) 
		FROM chat_messages cm
		WHERE cm.school_id IN (%s)
		AND cm.created_at > (
			SELECT COALESCE(uv.last_visit, '1970-01-01') 
			FROM user_visits uv
			WHERE uv.email = ? 
			AND uv.school_id = cm.school_id
		)
		GROUP BY cm.school_id`, placeholders)

	args := make([]interface{}, len(request.SchoolIDs)+1)
	for i, id := range request.SchoolIDs {
		args[i] = id
	}
	args[len(request.SchoolIDs)] = request.Email

	rows, err := Db.Query(query, args...)
	if err != nil {
		http.Error(w, jsonResponse(500, "数据库查询失败", nil), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var schoolID, count int
		if err := rows.Scan(&schoolID, &count); err != nil {
			continue
		}
		result[schoolID] = count
	}

	for _, id := range request.SchoolIDs {
		if _, exists := result[id]; !exists {
			result[id] = 0
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": result,
	})
}

// MarkReadHandler 标记已读
func MarkReadHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		SchoolID int    `json:"school_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, jsonResponse(400, "无效的请求参数", nil), http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO user_visits (email, school_id, last_visit)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE last_visit = VALUES(last_visit)`

	_, err := Db.Exec(query, request.Email, request.SchoolID, time.Now().UTC())
	if err != nil {
		http.Error(w, jsonResponse(500, "更新访问时间失败", nil), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "标记已读成功",
	})
}

// SendHandler 发送消息
func SendHandler(w http.ResponseWriter, r *http.Request) {
	var msg ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, jsonResponse(400, "无效的请求参数", nil), http.StatusBadRequest)
		return
	}

	var name, picture string
	err := Db.QueryRow(`
        SELECT COALESCE(NULLIF(name, ''), '匿名用户'), 
               COALESCE(NULLIF(picture, ''), 'https://example.com/default-avatar.png')
        FROM profile WHERE email = ?`, msg.Email).Scan(&name, &picture)
	if err != nil {
		log.Printf("获取用户信息失败，使用默认值: %v", err)
		name = "匿名用户"
		picture = "https://example.com/default-avatar.png"
	}

	tx, err := Db.Begin()
	if err != nil {
		http.Error(w, jsonResponse(500, "事务启动失败", nil), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(
		"INSERT INTO chat_messages (school_id, email, content) VALUES (?, ?, ?)",
		msg.SchoolID, msg.Email, msg.Content,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, jsonResponse(500, "消息发送失败", nil), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(
		`INSERT INTO user_visits (email, school_id, last_visit)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE last_visit = VALUES(last_visit)`,
		msg.Email, msg.SchoolID, time.Now().UTC(),
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, jsonResponse(500, "更新访问时间失败", nil), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, jsonResponse(500, "事务提交失败", nil), http.StatusInternalServerError)
		return
	}

	chat.Broadcast <- chat.Message{
		SchoolID:  msg.SchoolID,
		Email:     msg.Email,
		Content:   msg.Content,
		CreatedAt: time.Now(),
		Name:      name,
		Picture:   picture,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "消息发送成功",
	})
}

// jsonResponse 创建JSON响应
func jsonResponse(code int, message string, data interface{}) string {
	response := map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	}
	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

// SchoolProfileHandler 获取学校简介
func SchoolProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"code": 1001}

	if r.Method != http.MethodPost {
		response["code"] = 405
		json.NewEncoder(w).Encode(response)
		return
	}

	var req Const.SchoolProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response["msg"] = "无效的请求参数"
		json.NewEncoder(w).Encode(response)
		return
	}

	// rows, err := Db.Query(`
	// 	SELECT
	// 		school_id, introduce, hightitle,
	// 		gong_major, li_major, jing_major,
	// 		guan_major, yi_major, fa_major,
	// 		jiao_major, yishu_major, picture
	// 	FROM sch_profile
	// 	WHERE hightitle = ?`, req.Hightitle)
	//学校专业暂时不添加，只需要补上学校
	rows, err := Db.Query(`
		SELECT 
			id, content, name,logo_url
		FROM universities_detail 
		WHERE name = ?`, req.Hightitle)
	//school_id
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		response["msg"] = "服务器内部错误"
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var profiles []Const.SchoolProfileResponse
	for rows.Next() {
		var (
			schoolID  int
			introduce sql.NullString
			hightitle string
			picture   sql.NullString
		)

		err := rows.Scan(
			&schoolID, &introduce, &hightitle, &picture,
		)

		if err != nil {
			log.Printf("行扫描失败: %v", err)
			continue
		}

		profile := Const.SchoolProfileResponse{
			SchoolLogo: fmt.Sprintf("https://static-data.gaokao.cn/upload/logo/%d.jpg", schoolID),
			Hightitle:  hightitle,
			Introduce:  stringPtr(introduce),
			SchoolPic:  stringPtr(picture),
		}

		profiles = append(profiles, profile)
	}

	if len(profiles) == 0 {
		response["msg"] = "未找到相关院校信息"
		json.NewEncoder(w).Encode(response)
		return
	}

	response["code"] = 0
	response["data"] = profiles
	json.NewEncoder(w).Encode(response)
}

// ProfessionalInsertHandler 专业信息插入
func ProfessionalInsertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.ProfessionalResponse{Code: 1001}

	if r.Method != http.MethodPost {
		response.Msg = "仅支持POST方法"
		json.NewEncoder(w).Encode(response)
		return
	}

	var req Const.ProfessionalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Msg = "请求解析失败"
		log.Printf("JSON解析错误: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	tx, err := Db.Begin()
	if err != nil {
		log.Printf("事务开启失败: %v", err)
		response.Msg = "服务器内部错误"
		json.NewEncoder(w).Encode(response)
		return
	}
	defer tx.Rollback()

	successCount := 0
	for _, item := range req.Items {
		// 处理男女比例
		boyRate := item.BoyRate
		girlRate := item.GirlRate
		if boyRate == "0" && girlRate == "0" {
			boyRate = "60"
			girlRate = "40"
		}

		// 处理salaryavg（增强版）
		salaryavg := processNumericField(item.Salaryavg, 111364)

		// 处理fivesalaryavg（增强版）
		fivesalaryavg := processNumericField(item.FiveSalaryavg, 12365)

		// 调试日志
		log.Printf("正在插入数据 - 专业: %s", item.Name)
		log.Printf("处理后的 salaryavg: %d", salaryavg)
		log.Printf("处理后的 fivesalaryavg: %d", fivesalaryavg)

		_, err := tx.Exec(`
            INSERT INTO professional (
                name, level1_name, salaryavg, limit_year,
                fivesalaryavg, level2_name, boy_rate,
                girl_rate, level3_name
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            ON DUPLICATE KEY UPDATE
                level1_name = VALUES(level1_name),
                salaryavg = VALUES(salaryavg),
                limit_year = VALUES(limit_year),
                fivesalaryavg = VALUES(fivesalaryavg),
                level2_name = VALUES(level2_name),
                boy_rate = VALUES(boy_rate),
                girl_rate = VALUES(girl_rate),
                level3_name = VALUES(level3_name)`,
			item.Name,
			item.Level1Name,
			salaryavg,
			item.LimitYear,
			fivesalaryavg,
			item.Level2Name,
			boyRate,
			girlRate,
			item.Level3Name,
		)

		if err != nil {
			log.Printf("插入失败: %v 原始数据: %+v", err, item)
			continue
		}
		successCount++
	}

	if err := tx.Commit(); err != nil {
		log.Printf("事务提交失败: %v", err)
		response.Msg = "数据提交失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0
	response.Msg = "操作完成"
	response.Total = successCount
	json.NewEncoder(w).Encode(response)
}

// 新增辅助函数处理数字字段
func processNumericField(value interface{}, defaultValue int) int {
	switch v := value.(type) {
	case string:
		// 处理带逗号的数字字符串（如："14,200"）
		cleaned := strings.ReplaceAll(v, ",", "")
		if cleaned == "" {
			return defaultValue
		}
		if result, err := strconv.Atoi(cleaned); err == nil {
			if result == 0 {
				return defaultValue
			}
			return result
		}
	case int:
		if v == 0 {
			return defaultValue
		}
		return v
	case float64: // 处理JSON可能解析为float的情况
		if v == 0 {
			return defaultValue
		}
		return int(v)
	}
	return defaultValue
}

// ProfessionalQueryHandler 专业信息查询
func ProfessionalQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.ProfessionalQueryResponse{Code: 1001}

	// 只允许POST方法
	if r.Method != http.MethodPost {
		response.Code = http.StatusMethodNotAllowed
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求参数
	var req Const.ProfessionalQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		response.Code = http.StatusBadRequest
		json.NewEncoder(w).Encode(response)
		return
	}

	// 执行数据库查询
	rows, err := Db.Query(`
        SELECT 
            name, level1_name, salaryavg,
            limit_year, fivesalaryavg, level2_name,
            boy_rate, girl_rate, level3_name
        FROM professional 
        WHERE level3_name = ?`,
		req.Level3Name)

	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		response.Code = http.StatusInternalServerError
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var results []Const.ProfessionalItemResponse
	for rows.Next() {
		var item Const.ProfessionalItemResponse
		err := rows.Scan(
			&item.Name,
			&item.Level1Name,
			&item.Salaryavg,
			&item.LimitYear,
			&item.FiveSalaryavg,
			&item.Level2Name,
			&item.BoyRate,
			&item.GirlRate,
			&item.Level3Name,
		)

		if err != nil {
			log.Printf("数据解析失败: %v", err)
			continue
		}
		results = append(results, item)
	}

	// 检查遍历结果时的错误
	if err = rows.Err(); err != nil {
		log.Printf("结果集遍历错误: %v", err)
		response.Code = http.StatusInternalServerError
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造成功响应
	response.Code = 0
	response.Data = results
	json.NewEncoder(w).Encode(response)
}

// respondError 志愿推荐响应错误
func respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    code,
		"message": message,
	})
}

// ProfessionalUpsertHandler 专业 upsert
func ProfessionalUpsertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.UpsertResponse{Code: 1001, Message: "操作失败"}

	if r.Method != http.MethodPost {
		response.Message = "仅支持POST请求"
		json.NewEncoder(w).Encode(response)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Message = "读取请求体失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	var requests []Const.MoniRequest
	if err := json.Unmarshal(body, &requests); err != nil {
		response.Message = "解析JSON数据失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	for _, req := range requests {
		// 转换空字符串为NULL
		major1 := toNullString(req.Major1)
		major2 := toNullString(req.Major2)
		major3 := toNullString(req.Major3)
		major4 := toNullString(req.Major4)
		major5 := toNullString(req.Major5)
		major6 := toNullString(req.Major6)

		// 检查记录是否存在
		var exists bool
		err := Db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM moni WHERE email = ? AND benke = ? AND yitianzhuanye = ?)",
			req.Email, req.Benke, req.Yitianzhuanye,
		).Scan(&exists)

		if err != nil {
			log.Printf("检查记录存在性失败: %v", err)
			continue // 或返回错误，根据需求决定
		}

		if exists {
			// 更新记录
			query := `UPDATE moni SET 
                college = ?, major1 = ?, major2 = ?, major3 = ?, 
                major4 = ?, major5 = ?, major6 = ?
                WHERE email = ? AND benke = ? AND yitianzhuanye = ?`
			_, err = Db.Exec(query,
				req.College,
				major1, major2, major3,
				major4, major5, major6,
				req.Email, req.Benke, req.Yitianzhuanye,
			)
		} else {
			// 插入新记录
			query := `INSERT INTO moni 
                (college, benke, yitianzhuanye, major1, major2, major3, major4, major5, major6, email)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
			_, err = Db.Exec(query,
				req.College,
				req.Benke, req.Yitianzhuanye,
				major1, major2, major3,
				major4, major5, major6,
				req.Email,
			)
		}

		if err != nil {
			log.Printf("操作失败: %v", err)
			continue
		}
	}

	response.Code = 0
	response.Message = "操作成功"
	json.NewEncoder(w).Encode(response)
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func FetchVolunteerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := Const.VolunteerFetchResponse{Status: 400}
	log.Printf("进入")
	// 验证请求方法
	if r.Method != http.MethodPost {
		resp.Message = "仅支持POST请求"
		json.NewEncoder(w).Encode(resp)
		return
	}

	// 解析请求体
	var req Const.VolunteerFetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Message = "请求格式错误"
		json.NewEncoder(w).Encode(resp)
		return
	}

	// 参数校验
	if req.UserEmail == "" {
		resp.Message = "邮箱参数缺失"
		json.NewEncoder(w).Encode(resp)
		return
	}

	// 执行数据库查询
	queryStr := `
        SELECT 
            college, 
            yitianzhuanye,
            COALESCE(major1, ''),
            COALESCE(major2, ''),
            COALESCE(major3, ''),
            COALESCE(major4, ''),
            COALESCE(major5, ''),
            COALESCE(major6, '')
        FROM moni
        WHERE email = ? AND benke = ?
        ORDER BY yitianzhuanye ASC`

	rows, err := Db.Query(queryStr, req.UserEmail, req.BatchType)
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		resp.Message = "数据检索失败"
		json.NewEncoder(w).Encode(resp)
		return
	}
	defer rows.Close()

	// 处理查询结果
	var records []Const.VolunteerItem
	for rows.Next() {
		var item Const.VolunteerItem
		if err := rows.Scan(
			&item.Institution,
			&item.Sequence,
			&item.FirstMajor,
			&item.SecondMajor,
			&item.ThirdMajor,
			&item.FourthMajor,
			&item.FifthMajor,
			&item.SixthMajor,
		); err != nil {
			log.Printf("数据解析错误: %v", err)
			continue
		}
		records = append(records, item)
	}

	// 处理空结果
	if len(records) == 0 {
		resp.Status = 404
		resp.Message = "未找到相关记录"
	} else {
		resp.Status = 0
		resp.Records = records
	}

	json.NewEncoder(w).Encode(resp)
}

func FamilyAddHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.FamilyAddResponse{Code: 1001}

	// 只处理POST请求
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req Const.FamilyAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(response)
		return
	}

	// 插入数据库
	_, err := Db.Exec("INSERT INTO familyshare (familyemail, myemail) VALUES (?, ?)",
		req.FamilyEmail, req.MyEmail)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			response.Code = 1002
		}
		log.Printf("插入家庭成员失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0
	json.NewEncoder(w).Encode(response)
}

func FamilyFindHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.FamilyFindResponse{
		Code:   1001,
		Emails: []string{},
	}

	// 只处理POST请求
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req Const.FamilyFindRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(response)
		return
	}

	// 查询数据库
	rows, err := Db.Query("SELECT familyemail FROM familyshare WHERE myemail = ?", req.MyEmail)
	if err != nil {
		log.Printf("查询家庭成员失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	// 收集结果
	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}
		emails = append(emails, email)
	}

	// 检查遍历过程中的错误
	if err = rows.Err(); err != nil {
		log.Printf("行遍历错误: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0
	response.Emails = emails
	json.NewEncoder(w).Encode(response)
}

func FamilyRemoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.FamilyAddResponse{Code: 1001}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	var req Const.FamilyRemoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("请求解析失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 执行删除操作
	result, err := Db.Exec(
		"DELETE FROM familyshare WHERE myemail = ? AND familyemail = ?",
		req.MyEmail,
		req.FamilyEmail,
	)

	if err != nil {
		log.Printf("数据库删除失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查实际删除的行数
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.Code = 1003 // 记录不存在
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0
	json.NewEncoder(w).Encode(response)
}

func PolicyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.PolicyResponse{Code: 1001} // 默认错误状态码

	// 仅处理GET请求
	if r.Method != http.MethodGet {
		response.Code = http.StatusMethodNotAllowed
		json.NewEncoder(w).Encode(response)
		return
	}

	// 查询数据库
	rows, err := Db.Query("SELECT id, title, foreword FROM policy")
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var policies []Const.PolicyItem
	for rows.Next() {
		var (
			id      int
			title   sql.NullString
			content sql.NullString
		)

		// 扫描字段
		if err := rows.Scan(&id, &title, &content); err != nil {
			log.Printf("数据扫描失败: %v", err)
			continue
		}

		// 处理可能为NULL的字段
		item := Const.PolicyItem{ID: id}
		if title.Valid {
			item.Title = &title.String
		}
		if content.Valid {
			item.Content = &content.String
		}

		policies = append(policies, item)
	}

	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		log.Printf("遍历数据行失败: %v", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造成功响应
	response.Code = 0
	response.Data = policies
	json.NewEncoder(w).Encode(response)
}

// PolicySearchHandler 政策查询
func PolicySearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.PolicySearchResponse{Code: 1001, Msg: "操作失败"}

	// 只处理POST请求
	if r.Method != http.MethodPost {
		response.Msg = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req Const.PolicySearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON解析失败: %v", err)
		response.Msg = "无效的请求格式"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证标题是否为空
	if strings.TrimSpace(req.Title) == "" {
		response.Msg = "标题不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 执行数据库查询
	rows, err := Db.Query(
		"SELECT id, title, contenet FROM policy WHERE title = ?",
		req.Title,
	)
	if err != nil {
		log.Printf("数据库查询错误: %v", err)
		response.Msg = "服务器内部错误"
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var results []Const.PolicyItem
	for rows.Next() {
		var (
			id      int
			title   sql.NullString
			content sql.NullString
		)

		if err := rows.Scan(&id, &title, &content); err != nil {
			log.Printf("数据扫描失败: %v", err)
			continue
		}

		item := Const.PolicyItem{ID: id}
		if title.Valid {
			item.Title = &title.String
		}
		if content.Valid {
			item.Content = &content.String
		}
		results = append(results, item)
	}

	// 处理查询结果
	if len(results) == 0 {
		response.Msg = "未找到相关数据"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 返回成功响应
	response.Code = 0
	response.Msg = ""
	response.Data = results
	json.NewEncoder(w).Encode(response)
}

// 辅助函数处理NULL值
func convertNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// 处理整数
func convertNullInt(ni sql.NullInt64) *int {
	if ni.Valid {
		v := int(ni.Int64)
		return &v
	}
	return nil
}

// 处理年份
func convertNullYear(ni sql.NullInt64) *int {
	if ni.Valid {
		// 处理年份特殊逻辑
		year := int(ni.Int64)
		if year < 1900 || year > 2100 {
			return nil
		}
		return &year
	}
	return nil
}

// PolicyCreateHandler 创建政策
func PolicyCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.PolicyCreateResponse{Code: 1001}

	// 只处理POST请求
	if r.Method != http.MethodPost {
		response.Msg = "无效的请求方法"
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req Const.PolicyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON解析失败: %v", err)
		response.Msg = "无效的请求格式"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 基础验证（示例：标题必填）
	if req.Title == nil || strings.TrimSpace(*req.Title) == "" {
		response.Msg = "标题不能为空"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 执行数据库插入
	result, err := Db.Exec(
		"INSERT INTO policy (title, contenet, foreword) VALUES (?, ?, ?)",
		req.Title,
		req.Contenet,
		req.Foreword,
	)
	if err != nil {
		log.Printf("数据库插入失败: %v", err)
		response.Msg = "创建政策失败"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取插入ID
	insertedID, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入ID失败: %v", err)
		// 不返回错误，因为数据已经插入成功
	}

	// 返回成功响应
	response.Code = 0
	response.Msg = "创建成功"
	response.ID = insertedID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// LogInsertHandler  插入日志
func LogInsertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.LogResponse{Code: 1001}

	// 验证请求方法
	if r.Method != http.MethodPost {
		response.Msg = "仅支持POST请求"
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req Const.LogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Msg = "无效的JSON格式"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证数据有效性
	if err := validate.Struct(req); err != nil {
		response.Msg = fmt.Sprintf("验证失败: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 转换日期格式
	parsedTime, err := time.Parse("2006-01-02 15:04:05", req.Date)
	if err != nil {
		response.Msg = "日期格式无效，请使用YYYY-MM-DD HH:MM:SS格式"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 执行数据库插入
	result, err := Db.Exec(
		"INSERT INTO logtextuser (email, date, operation) VALUES (?, ?, ?)",
		req.Email,
		parsedTime,
		req.Operation,
	)
	if err != nil {
		log.Printf("日志插入失败: %v", err)
		response.Msg = "服务器内部错误"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取插入ID
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取插入ID失败: %v", err)
	}

	// 返回成功响应
	response.Code = 0
	response.Msg = "日志记录成功"
	response.ID = lastID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// LogRetrieveHandler 日志查询处理函数
func LogRetrieveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Const.LogListResponse{Code: 1001}

	// 只处理GET请求
	if r.Method != http.MethodGet {
		response.Msg = "Method not allowed"
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取查询参数
	queryParams := r.URL.Query()
	emailFilter := queryParams.Get("email")
	operationFilter := queryParams.Get("operation")
	startDate := queryParams.Get("start_date")
	endDate := queryParams.Get("end_date")

	// 构建基础SQL
	sqlQuery := "SELECT email, date, operation FROM logtextuser WHERE 1=1"
	var args []interface{}

	// 添加过滤条件
	if emailFilter != "" {
		sqlQuery += " AND email = ?"
		args = append(args, emailFilter)
	}
	if operationFilter != "" {
		sqlQuery += " AND operation LIKE ?"
		args = append(args, "%"+operationFilter+"%")
	}
	if startDate != "" {
		sqlQuery += " AND date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		sqlQuery += " AND date <= ?"
		args = append(args, endDate)
	}

	// 添加排序
	sqlQuery += " ORDER BY date DESC"

	// 执行查询
	rows, err := Db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("数据库查询失败: %v", err)
		response.Msg = "数据获取失败"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var logs []Const.LogEntry
	for rows.Next() {
		var entry Const.LogEntry
		var rawTime []uint8 // 处理不同数据库驱动的时间类型

		if err := rows.Scan(&entry.Email, &rawTime, &entry.Operation); err != nil {
			log.Printf("数据解析失败: %v", err)
			continue
		}

		// 转换时间格式
		if t, err := time.Parse("2006-01-02 15:04:05", string(rawTime)); err == nil {
			entry.Date = t.Format("2006-01-02 15:04:05")
		} else {
			entry.Date = string(rawTime) // 原始格式输出
		}

		logs = append(logs, entry)
	}

	// 检查遍历错误
	if err = rows.Err(); err != nil {
		log.Printf("数据遍历失败: %v", err)
		response.Msg = "数据解析错误"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 返回结果
	response.Code = 0
	response.Data = logs
	response.Total = len(logs)
	json.NewEncoder(w).Encode(response)
}
