package oldManyMgrApi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	oldChatApi2 "mymod/api/oldChatApi"
	"mymod/config"
	"mymod/model/oldModel"
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
        FROM s WHERE email = ?`, msg.Email).Scan(&name, &picture)
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

	oldChatApi2.Broadcast <- oldChatApi2.Message{
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

// ProfessionalUpsertHandler 志愿 upsert
func ProfessionalUpsertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := oldModel.UpsertResponse{Code: 1001, Message: "操作失败"}

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

	var requests []oldModel.MoniRequest
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
	resp := oldModel.VolunteerFetchResponse{Status: 400}
	log.Printf("进入")
	// 验证请求方法
	if r.Method != http.MethodPost {
		resp.Message = "仅支持POST请求"
		json.NewEncoder(w).Encode(resp)
		return
	}

	// 解析请求体
	var req oldModel.VolunteerFetchRequest
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
	var records []oldModel.VolunteerItem
	for rows.Next() {
		var item oldModel.VolunteerItem
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
	response := oldModel.FamilyAddResponse{Code: 1001}

	// 只处理POST请求
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req oldModel.FamilyAddRequest
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
	response := oldModel.FamilyFindResponse{
		Code:   1001,
		Emails: []string{},
	}

	// 只处理POST请求
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req oldModel.FamilyFindRequest
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
	response := oldModel.FamilyAddResponse{Code: 1001}

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(response)
		return
	}

	var req oldModel.FamilyRemoveRequest
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

// LogInsertHandler  插入日志
func LogInsertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := oldModel.LogResponse{Code: 1001}

	// 验证请求方法
	if r.Method != http.MethodPost {
		response.Msg = "仅支持POST请求"
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req oldModel.LogRequest
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
		"INSERT INTO log_user_do (email, date, operation) VALUES (?, ?, ?)",
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
	response := oldModel.LogListResponse{Code: 1001}

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
	sqlQuery := "SELECT email, date, operation FROM log_user_do WHERE 1=1"
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

	var logs []oldModel.LogEntry
	for rows.Next() {
		var entry oldModel.LogEntry
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
