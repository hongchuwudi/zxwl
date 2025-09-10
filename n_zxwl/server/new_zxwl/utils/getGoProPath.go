package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetProjectRoot 获取项目根目录的绝对路径
// 通过查找包含go.mod文件的目录来确定项目根目录
func GetProjectRoot() (string, error) {
	// 首先尝试通过运行时调用栈信息获取路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrNotExist
	}

	// 从当前文件开始向上查找包含go.mod的目录
	dir := filepath.Dir(filename)
	for {
		// 检查当前目录是否包含go.mod文件
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// 到达文件系统根目录，停止查找
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	// 如果没有找到go.mod，尝试其他方法
	return getFallbackRootPath()
}

// getFallbackRootPath 当无法通过go.mod找到项目根目录时的备用方案
func getFallbackRootPath() (string, error) {
	// 方法1: 通过可执行文件路径
	if exePath, err := os.Executable(); err == nil {
		if resolvedPath, err := filepath.EvalSymlinks(exePath); err == nil {
			dir := filepath.Dir(resolvedPath)
			// 如果是临时目录（如go run运行），尝试其他方法
			if !isTempDir(dir) {
				return dir, nil
			}
		}
	}

	// 方法2: 通过当前工作目录
	if wd, err := os.Getwd(); err == nil {
		return wd, nil
	}

	return "", os.ErrNotExist
}

// isTempDir 检查目录是否为临时目录（如go run运行时）
func isTempDir(dir string) bool {
	tempDir := os.TempDir()
	return strings.Contains(dir, tempDir) ||
		strings.Contains(dir, "/var/folders/") || // macOS临时目录
		strings.Contains(dir, "/tmp/") // Linux临时目录
}

// GetModuleName 获取当前项目的模块名（从go.mod中）
func GetModuleName() (string, error) {
	rootDir, err := GetProjectRoot()
	if err != nil {
		return "", err
	}

	goModPath := filepath.Join(rootDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	// 解析go.mod文件的第一行获取模块名
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(line[len("module "):]), nil
		}
	}

	return "", os.ErrNotExist
}
