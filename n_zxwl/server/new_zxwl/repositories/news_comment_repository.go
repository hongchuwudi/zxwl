package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"mymod/new_zxwl/model/sqlModel"
)

type NewsCommentRepository struct {
	db *gorm.DB
}

func NewNewsCommentRepository(db *gorm.DB) *NewsCommentRepository {
	return &NewsCommentRepository{db: db}
}

// CreateComment 创建评论
func (r *NewsCommentRepository) CreateComment(comment *sqlModel.NewsComment) error {
	err := r.db.Create(comment).Error
	if err != nil {
		return fmt.Errorf("创建评论失败: %v", err)
	}
	return nil
}

// GetCommentsByNewsID 根据资讯ID获取评论列表
func (r *NewsCommentRepository) GetCommentsByNewsID(newsID uint64, page, pageSize int) ([]sqlModel.NewsComment, int64, error) {
	var comments []sqlModel.NewsComment
	var total int64

	// 获取总数
	err := r.db.Model(&sqlModel.NewsComment{}).Where("news_id = ? AND parent_id = 0", newsID).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询评论总数失败: %v", err)
	}

	// 分页查询顶级评论
	offset := (page - 1) * pageSize
	err = r.db.Where("news_id = ? AND parent_id = 0", newsID).
		Order("comment_time DESC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询评论列表失败: %v", err)
	}

	return comments, total, nil
}

// GetRepliesByCommentID 获取评论的回复
func (r *NewsCommentRepository) GetRepliesByCommentID(commentID uint64) ([]sqlModel.NewsComment, error) {
	var replies []sqlModel.NewsComment
	err := r.db.Where("parent_id = ?", commentID).
		Order("comment_time ASC").
		Find(&replies).Error
	if err != nil {
		return nil, fmt.Errorf("查询回复失败: %v", err)
	}
	return replies, nil
}

// GetCommentByID 根据ID获取评论
func (r *NewsCommentRepository) GetCommentByID(id uint64) (*sqlModel.NewsComment, error) {
	var comment sqlModel.NewsComment
	err := r.db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, fmt.Errorf("查询评论失败: %v", err)
	}
	return &comment, nil
}

// UpdateComment 更新评论
func (r *NewsCommentRepository) UpdateComment(comment *sqlModel.NewsComment) error {
	err := r.db.Model(&sqlModel.NewsComment{}).Where("id = ?", comment.ID).Updates(comment).Error
	if err != nil {
		return fmt.Errorf("更新评论失败: %v", err)
	}
	return nil
}

// DeleteComment 删除评论
func (r *NewsCommentRepository) DeleteComment(id uint64) error {
	err := r.db.Where("id = ?", id).Delete(&sqlModel.NewsComment{}).Error
	if err != nil {
		return fmt.Errorf("删除评论失败: %v", err)
	}
	return nil
}

// IncrementLikeCount 增加点赞数
func (r *NewsCommentRepository) IncrementLikeCount(id uint64) error {
	err := r.db.Model(&sqlModel.NewsComment{}).Where("id = ?", id).
		Update("like_count", gorm.Expr("like_count + 1")).Error
	if err != nil {
		return fmt.Errorf("增加点赞数失败: %v", err)
	}
	return nil
}

// IncrementReplyCount 增加回复数
func (r *NewsCommentRepository) IncrementReplyCount(id uint64) error {
	err := r.db.Model(&sqlModel.NewsComment{}).Where("id = ?", id).
		Update("reply_count", gorm.Expr("reply_count + 1")).Error
	if err != nil {
		return fmt.Errorf("增加回复数失败: %v", err)
	}
	return nil
}

// GetCommentCountByNewsID 获取资讯的评论总数
func (r *NewsCommentRepository) GetCommentCountByNewsID(newsID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&sqlModel.NewsComment{}).Where("news_id = ?", newsID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("查询评论总数失败: %v", err)
	}
	return count, nil
}
