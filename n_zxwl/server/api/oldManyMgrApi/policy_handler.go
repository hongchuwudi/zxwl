package oldManyMgrApi

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"mymod/model/oldModel"
	"net/http"
	"strconv"
	"strings"
)

// PolicySearchHandler 政策查询
func PolicySearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := oldModel.PolicySearchResponse{Code: 1001, Msg: "操作失败"}

	// 只处理POST请求
	if r.Method != http.MethodPost {
		response.Msg = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req oldModel.PolicySearchRequest
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

	var results []oldModel.PolicyItem
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

		item := oldModel.PolicyItem{ID: id}
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

// PolicyCreateHandler 创建政策
func PolicyCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := oldModel.PolicyCreateResponse{Code: 1001}

	// 只处理POST请求
	if r.Method != http.MethodPost {
		response.Msg = "无效的请求方法"
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var req oldModel.PolicyCreateRequest
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

// PolicyHandler 获取政策列表
func PolicyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := oldModel.PolicyResponse{Code: 1001} // 默认错误状态码

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

	var policies []oldModel.PolicyItem
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
		item := oldModel.PolicyItem{ID: id}
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

// PolicyDeleteHandler 删除政策
func PolicyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := oldModel.PolicySearchResponse{Code: 1001, Msg: "请求失败"}

	if r.Method != http.MethodDelete {
		response.Code = http.StatusMethodNotAllowed
		response.Msg = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 从URL路径参数获取政策ID
	vars := mux.Vars(r)
	policyIDStr := vars["id"]
	if policyIDStr == "" {
		response.Code = http.StatusBadRequest
		response.Msg = "政策ID不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	policyID, err := strconv.Atoi(policyIDStr)
	if err != nil || policyID <= 0 {
		response.Code = http.StatusBadRequest
		response.Msg = "政策ID格式错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 执行数据库删除操作
	result, err := Db.Exec("DELETE FROM policy WHERE id = ?", policyID)
	if err != nil {
		log.Printf("删除政策失败: %v", err)
		response.Code = http.StatusInternalServerError
		response.Msg = "服务器内部错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查是否成功删除
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取受影响行数失败: %v", err)
		response.Code = http.StatusInternalServerError
		response.Msg = "服务器内部错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	if rowsAffected == 0 {
		response.Code = http.StatusNotFound
		response.Msg = "未找到要删除的政策"
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Code = 0
	response.Msg = "删除成功"
	json.NewEncoder(w).Encode(response)
}
