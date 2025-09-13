// Package userRepo file: verify_code_repository.go
package userRepo

import (
	"gorm.io/gorm"
	"mymod/model/param/userParam"
	sqlModel2 "mymod/model/sqlModel"
	"time"
)

type VerifyCodeRepository struct {
	db *gorm.DB
}

func NewVerifyCodeRepository(db *gorm.DB) *VerifyCodeRepository {
	return &VerifyCodeRepository{db: db}
}

// CreateVerifyCode 创建验证码
func (r *VerifyCodeRepository) CreateVerifyCode(code *sqlModel2.UserVerifyCode) error {
	return r.db.Create(code).Error
}

// GetUsersWithPagination 分页查询用户
func (r *UserRepository) GetUsersWithPagination(req userParam.UserQueryRequest) ([]*sqlModel2.UserProfile, int64, error) {
	var users []*sqlModel2.UserProfile
	var total int64

	db := r.db.Table("user_profile")

	// 添加查询条件
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		db = db.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.StartTime != "" {
		startTime, _ := time.Parse("2006-01-02", req.StartTime)
		db = db.Where("created_at >= ?", startTime)
	}
	if req.EndTime != "" {
		endTime, _ := time.Parse("2006-01-02", req.EndTime)
		db = db.Where("created_at <= ?", endTime.AddDate(0, 0, 1))
	}
	if req.Gender == 3 || req.Gender == 1 || req.Gender == 2 {
		db = db.Where("gender = ?", req.Gender)
	}
	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	err := db.Offset(offset).Limit(req.Size).
		Order("created_at DESC").
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetLatestVerifyCode 获取指定邮箱的最新验证码
func (r *VerifyCodeRepository) GetLatestVerifyCode(email string) (*sqlModel2.UserVerifyCode, error) {
	var code sqlModel2.UserVerifyCode
	err := r.db.Where("email = ?", email).
		Order("created_at DESC").
		First(&code).Error
	if err != nil {
		return nil, err
	}
	return &code, nil
}

// DeleteExpiredCodes 删除过期验证码
func (r *VerifyCodeRepository) DeleteExpiredCodes() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&sqlModel2.UserVerifyCode{}).Error
}
