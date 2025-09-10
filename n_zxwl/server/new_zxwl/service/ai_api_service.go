package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mymod/new_zxwl/config"
	"net/http"
	"strings"
)

type QwenService struct {
	apiKey string
}

// NewQwenService 创建一个新的 QwenService 实例
func NewQwenService(cfg *config.Config) *QwenService {
	return &QwenService{apiKey: cfg.QwenAPIKey}
}

// 公共的API调用方法
func (q *QwenService) callQwenAPI(url string, payload map[string]interface{}) ([]byte, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+q.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetEmbedding 获取文本的嵌入向量
func (q *QwenService) GetEmbedding(text string) ([]float64, error) {
	url := "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding"

	payload := map[string]interface{}{
		"model": "text-embedding-v1",
		"input": map[string]string{
			"text": text,
		},
	}

	body, err := q.callQwenAPI(url, payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Output struct {
			Embedding []float64 `json:"embedding"`
		} `json:"output"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Output.Embedding, nil
}

// CalculateSimilarity 计算两个文本的相似度
func (q *QwenService) CalculateSimilarity(text1, text2 string) (float64, error) {
	url := "https://dashscope.aliyuncs.com/api/v1/services/similarity/calculation"

	payload := map[string]interface{}{
		"model": "text-similarity-v1",
		"input": map[string]string{
			"text1": text1,
			"text2": text2,
		},
	}

	body, err := q.callQwenAPI(url, payload)
	if err != nil {
		return 0, err
	}

	var result struct {
		Output struct {
			Similarity float64 `json:"similarity"`
		} `json:"output"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	return result.Output.Similarity, nil
}

// ExtractNumber 从文本中提取特定字段的数值
func (q *QwenService) ExtractNumber(text, targetField string) (float64, error) {
	url := "https://dashscope.aliyuncs.com/api/v1/services/number-extraction/extract"

	payload := map[string]interface{}{
		"model": "number-extraction-v1",
		"input": map[string]string{
			"text":        text,
			"targetField": targetField,
		},
	}

	body, err := q.callQwenAPI(url, payload)
	if err != nil {
		return 0, err
	}

	var result struct {
		Output struct {
			Number float64 `json:"number"`
		} `json:"output"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	return result.Output.Number, nil
}

// TextGeneration 通用文本生成方法
func (q *QwenService) TextGeneration(prompt string, resultFormat string) (string, error) {
	url := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

	payload := map[string]interface{}{
		"model": "qwen-max",
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": prompt,
				},
			},
		},
		"parameters": map[string]string{
			"result_format": resultFormat,
		},
	}

	body, err := q.callQwenAPI(url, payload)
	if err != nil {
		return "", err
	}

	var result struct {
		Output struct {
			Text string `json:"text"`
		} `json:"output"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Output.Text, nil
}

// ExtractJSONFromText 从文本中提取JSON
func (q *QwenService) ExtractJSONFromText(text string, v interface{}) error {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}") + 1

	if start == -1 || end == -1 {
		return fmt.Errorf("未找到有效的JSON响应")
	}

	jsonStr := text[start:end]
	return json.Unmarshal([]byte(jsonStr), v)
}
