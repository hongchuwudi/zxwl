package api

import (
	"encoding/json"
	"log"
	"mymod/new_zxwl/config"
	"mymod/new_zxwl/model/param"
	"mymod/new_zxwl/model/sqlModel"
	"mymod/new_zxwl/repositories"
	"mymod/new_zxwl/utils"
	"net/http"
	_ "strconv"
	"time"
)

// GetCommentsByNewsIDHandler 根据文章ID获取该文章的所有评论
func GetCommentsByNewsIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析JSON请求体
	var request struct {
		NewsID uint64 `json:"news_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.NewsID == 0 {
		response.Error = http.StatusBadRequest
		response.Message = "news_id参数不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	commentRepo := repositories.NewNewsCommentRepository(config.GetDB())

	// 获取评论总数
	total, err := commentRepo.GetCommentCountByNewsID(request.NewsID)
	if err != nil {
		log.Printf("获取评论总数失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "获取评论失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取所有评论（不使用分页，使用repository的现有方法）
	comments, _, err := commentRepo.GetCommentsByNewsID(request.NewsID, 1, 1000) // 使用一个大数来获取所有
	if err != nil {
		log.Printf("查询评论列表失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "获取评论失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 获取每个评论的回复
	for i := range comments {
		replies, err := commentRepo.GetRepliesByCommentID(comments[i].ID)
		if err != nil {
			log.Printf("获取评论回复失败: %v", err)
			continue
		}

		comments[i].Replies = replies
	}

	// 构造成功响应
	responseData := map[string]interface{}{
		"comments": comments,
		"total":    total,
	}

	response.Error = 0
	response.Message = "获取成功"
	response.Data = responseData
	json.NewEncoder(w).Encode(response)
}

// AddCommentHandler 增加文章的评论
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析请求体
	var comment sqlModel.NewsComment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 验证必要字段
	if comment.NewsID == 0 || comment.CommenterName == "" || comment.CommentContent == "" {
		response.Error = http.StatusBadRequest
		response.Message = "缺少必要参数"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 设置默认值
	comment.CommentTime = time.Now()
	if comment.IsApproved == 0 {
		comment.IsApproved = 1 // 默认审核通过
	}

	// 初始化repository
	commentRepo := repositories.NewNewsCommentRepository(config.GetDB())

	// 创建评论
	err = commentRepo.CreateComment(&comment)
	if err != nil {
		log.Printf("创建评论失败: %v", err)
		response.Error = http.StatusInternalServerError
		response.Message = "评论发布失败"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 如果是回复评论，增加父评论的回复数
	if comment.ParentID > 0 {
		err = commentRepo.IncrementReplyCount(comment.ParentID)
		if err != nil {
			log.Printf("更新回复数失败: %v", err)
		}
	}

	// 构造成功响应
	response.Error = 0
	response.Message = "评论发布成功"
	response.Data = comment
	json.NewEncoder(w).Encode(response)
}

// ToggleCommentLikeHandler 点赞和取消点赞功能
func ToggleCommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := param.Response{Error: 1001, Message: "请求失败"}

	if r.Method != http.MethodPost {
		response.Error = http.StatusMethodNotAllowed
		response.Message = "方法不允许"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 解析JSON请求体
	var request struct {
		CommentID uint64 `json:"comment_id"`
		Action    string `json:"action"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error = http.StatusBadRequest
		response.Message = "请求参数错误"
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.CommentID == 0 || request.Action == "" {
		response.Error = http.StatusBadRequest
		response.Message = "comment_id和action参数不能为空"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 初始化repository
	commentRepo := repositories.NewNewsCommentRepository(config.GetDB())

	// 检查评论是否存在
	existingComment, err := commentRepo.GetCommentByID(request.CommentID)
	if err != nil || existingComment == nil {
		response.Error = http.StatusNotFound
		response.Message = "评论不存在"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 根据action执行点赞或取消点赞
	if request.Action == "like" {
		// 增加点赞数
		err = commentRepo.IncrementLikeCount(request.CommentID)
		if err != nil {
			log.Printf("点赞失败: %v", err)
			response.Error = http.StatusInternalServerError
			response.Message = "点赞失败"
			json.NewEncoder(w).Encode(response)
			return
		}
		response.Message = "点赞成功"
	} else if request.Action == "unlike" {
		// 减少点赞数
		err = commentRepo.UpdateComment(&sqlModel.NewsComment{
			ID:        request.CommentID,
			LikeCount: utils.Max(0, existingComment.LikeCount-1),
		})
		if err != nil {
			log.Printf("取消点赞失败: %v", err)
			response.Error = http.StatusInternalServerError
			response.Message = "取消点赞失败"
			json.NewEncoder(w).Encode(response)
			return
		}
		response.Message = "取消点赞成功"
	} else {
		response.Error = http.StatusBadRequest
		response.Message = "action参数错误，应为like或unlike"
		json.NewEncoder(w).Encode(response)
		return
	}

	// 构造成功响应
	response.Error = 0
	json.NewEncoder(w).Encode(response)
}
