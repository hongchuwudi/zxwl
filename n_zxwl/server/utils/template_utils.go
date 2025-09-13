package utils

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

// EmailTemplateData 邮件模板数据
type EmailTemplateData struct {
	VerificationCode string
}

// LoadEmailTemplate 加载邮件模板
func LoadEmailTemplate(templateName string, data EmailTemplateData) (string, error) {
	// 获取项目根目录
	proPath, err := GetProjectRoot()
	if err != nil {
		return "", err
	}

	// 构建模板文件路径
	templatePath := filepath.Join(proPath, "new_zxwl", "static", templateName)

	// 读取模板文件
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	// 解析模板
	tmpl, err := template.New("email").Parse(string(content))
	if err != nil {
		return "", err
	}

	// 应用数据到模板
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
