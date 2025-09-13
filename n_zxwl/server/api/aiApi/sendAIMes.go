// ai_handler.go
package aiApi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mymod/config"
	"mymod/model/param"
	"net/http"
	"strings"
)

// AIHandlers 处理AI请求（支持SSE流式传输）
func AIHandlers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var aiReq param.AIRequest
	if err := json.NewDecoder(r.Body).Decode(&aiReq); err != nil {
		http.Error(w, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 设置流式传输
	aiReq.Stream = true

	// 调用通义千问API并转发SSE流
	err := callQwenAPIStream(w, aiReq)
	if err != nil {
		log.Printf("调用AI API失败: %v", err)
		// 发送错误消息给前端
		fmt.Fprintf(w, "data: %s\n\n", `{"error": "AI服务调用失败"}`)
		w.(http.Flusher).Flush()
		return
	}
}

// callQwenAPIStream 调用通义千问API并转发SSE流
func callQwenAPIStream(w http.ResponseWriter, req param.AIRequest) error {
	// 从配置或环境变量获取API密钥
	apiKey := config.LoadConfig().QwenAPIKey
	apiURL := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": req.Model,
		"input": map[string]interface{}{
			"messages": req.Messages,
		},
		"parameters": map[string]interface{}{
			"stream":        true,
			"temperature":   req.Temperature,
			"max_tokens":    req.MaxTokens,
			"result_format": "text",
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("X-DashScope-SSE", "enable")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API请求失败: %s", string(body))
	}

	// 转发SSE流到前端
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	buffer := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			// 直接转发数据
			w.Write(buffer[:n])
			flusher.Flush()
		}
		if err != nil {
			if err == io.EOF {
				// 发送结束标记
				fmt.Fprintf(w, "data: [DONE]\n\n")
				flusher.Flush()
				break
			}
			return err
		}
	}

	return nil
}
