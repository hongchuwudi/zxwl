// file: api/upload_handler.go
package baseApi

import (
	"encoding/json"
	"log"
	"mymod/model/param"
	"mymod/service"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		uploadService: service.NewUploadService(),
	}
}

// UploadAvatarHandler 上传头像
func (h *UploadHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "上传失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析multipart表单
	err := r.ParseMultipartForm(10 << 20) // 10MB限制
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "解析表单数据失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取文件
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "获取文件失败"
		json.NewEncoder(w).Encode(response)
		return
	}
	defer file.Close()

	// 上传图片
	fileURL, err := h.uploadService.UploadImage(fileHeader, "avatars")
	if err != nil {
		log.Printf("上传头像失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Error = 0
	response.Message = "上传成功"
	response.Data = map[string]string{
		"url":  fileURL,
		"name": fileHeader.Filename,
		"type": filepath.Ext(fileHeader.Filename),
		"size": string(fileHeader.Size),
	}
	json.NewEncoder(w).Encode(response)
}

// UploadFileHandler 通用文件上传
func (h *UploadHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "上传失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 从URL参数获取文件夹类型
	vars := mux.Vars(r)
	folder := vars["folder"]
	if folder == "" {
		folder = "files"
	}

	err := r.ParseMultipartForm(20 << 20) // 20MB限制
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "解析表单数据失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "获取文件失败"
		json.NewEncoder(w).Encode(response)
		return
	}
	defer file.Close()

	// 上传文件
	fileURL, err := h.uploadService.UploadFile(fileHeader, folder)
	if err != nil {
		log.Printf("上传文件失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Error = 0
	response.Message = "上传成功"
	response.Data = map[string]interface{}{
		"url":    fileURL,
		"name":   fileHeader.Filename,
		"type":   filepath.Ext(fileHeader.Filename),
		"size":   fileHeader.Size,
		"folder": folder,
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteFileHandler 删除文件
func (h *UploadHandler) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "删除失败"}

	if r.Method != http.MethodDelete {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.URL == "" {
		response.Error = http.StatusBadRequest
		response.Message = "文件URL不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	err := h.uploadService.DeleteFile(req.URL)
	if err != nil {
		log.Printf("删除文件失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Error = 0
	response.Message = "删除成功"
	json.NewEncoder(w).Encode(response)
}
