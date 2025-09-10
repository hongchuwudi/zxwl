package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	chat "mymod/ChatMgr"
	"mymod/MysqlMgr"
	"net/http/httptest"
)

func maan() {
	MysqlMgr.InitDB()
	chat.InitDB()
	// 创建测试请求
	reqBody := map[string]string{"email": "2772167017@qq.com"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/verify-code", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	// 调用被测试的函数
	MysqlMgr.GetVerifyCodeHandler(w, req)

	// 检查响应
	resp := w.Result()
	var response struct {
		Error int    `json:"error"`
		Msg   string `json:"msg"`
	}
	json.NewDecoder(resp.Body).Decode(&response)

	fmt.Printf("响应状态: %d\n", resp.StatusCode)
	fmt.Printf("错误码: %d\n", response.Error)
	fmt.Printf("消息: %s\n", response.Msg)
}
