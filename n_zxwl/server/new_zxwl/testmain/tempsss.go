package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 定义API响应结构
type ApiResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID         string `json:"id"`
		SchoolID   string `json:"school_id"`
		Type       string `json:"type"`
		Content    string `json:"content"`
		Status     string `json:"status"`
		UpdateTime string `json:"update_time"`
		TypeName   string `json:"type_name"`
	} `json:"data"`
	MD5 string `json:"md5"`
}

// 数据库配置
const (
	dbDriver   = "mysql"
	dbUser     = "root"           // 替换为你的数据库用户名
	dbPassword = "hongchu"        // 替换为你的数据库密码
	dbName     = "zxwl"           // 替换为你的数据库名
	dbHost     = "localhost:3306" // 替换为你的数据库地址
)

// 184 897 1139 2199 2239 2380 3448 3688 3692
// API基础URL
const baseURL = "https://static-data.gaokao.cn/www/2.0/school/%s/detail/69000.json?a=www.gaokao.cn"

func mains() {
	// 连接数据库
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}
	fmt.Println("数据库连接成功")

	// 获取所有大学的ID
	universities, err := getAllUniversityIDs(db)
	if err != nil {
		log.Fatal("获取大学ID列表失败:", err)
	}

	fmt.Printf("共找到 %d 所大学\n", len(universities))

	// 遍历所有大学，更新简介内容
	for i, uni := range universities {
		fmt.Printf("正在处理第 %d/%d 所大学: %s (ID: %d)\n", i+1, len(universities), uni.Name, uni.ID)

		// 从API获取简介内容
		content, err := fetchUniversityContent(uni.ID)
		if err != nil {
			log.Printf("获取大学 %d 的简介失败: %v\n", uni.ID, err)
			continue
		}

		// 更新数据库
		err = updateUniversityContent(db, uni.ID, content)
		if err != nil {
			log.Printf("更新大学 %d 的简介失败: %v\n", uni.ID, err)
			continue
		}

		fmt.Printf("成功更新大学 %s 的简介\n", uni.Name)

		// 添加延迟，避免请求过于频繁
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("所有大学简介更新完成")
}

// 大学结构体
type University struct {
	ID   int
	Name string
}

// 获取所有大学的ID和名称
func getAllUniversityIDs(db *sql.DB) ([]University, error) {
	query := "SELECT id, name FROM universities_detail WHERE status = '1' ORDER BY id"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var universities []University
	for rows.Next() {
		var uni University
		err := rows.Scan(&uni.ID, &uni.Name)
		if err != nil {
			return nil, err
		}
		universities = append(universities, uni)
	}

	return universities, nil
}

// 从API获取大学简介内容
func fetchUniversityContent(universityID int) (string, error) {
	url := fmt.Sprintf(baseURL, fmt.Sprintf("%d", universityID))

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回非200状态码: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("解析JSON失败: %v", err)
	}

	if apiResponse.Code != "0000" {
		return "", fmt.Errorf("API返回错误: %s - %s", apiResponse.Code, apiResponse.Message)
	}

	return apiResponse.Data.Content, nil
}

// 更新数据库中的大学简介内容
func updateUniversityContent(db *sql.DB, universityID int, content string) error {
	query := "UPDATE universities_detail SET content = ?, last_updated = NOW() WHERE id = ?"
	_, err := db.Exec(query, content, universityID)
	return err
}
