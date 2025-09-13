// Package userRepo file: user_repository.go
package userRepo

import (
	"errors"
	"gorm.io/gorm"
	"mymod/model/sqlModel"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser 创建用户
func (r *UserRepository) CreateUser(user *sqlModel.UserProfile) error {
	return r.db.Create(user).Error
}

// GetUserByEmail 根据邮箱获取用户
func (r *UserRepository) GetUserByEmail(email string) (*sqlModel.UserProfile, error) {
	var user sqlModel.UserProfile
	err := r.db.Table("user_profile").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (r *UserRepository) GetUserByUsername(username string) (*sqlModel.UserProfile, error) {
	var user sqlModel.UserProfile
	err := r.db.Table("user_profile").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (r *UserRepository) GetUserByID(id int) (*sqlModel.UserProfile, error) {
	var user sqlModel.UserProfile
	err := r.db.Table("user_profile").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (r *UserRepository) UpdateUser(user *sqlModel.UserProfile) error {
	return r.db.Table("user_profile").Save(user).Error
}

// DeleteUser 删除用户
func (r *UserRepository) DeleteUser(id int) error {
	return r.db.Table("user_profile").Where("id = ?", id).Delete(&sqlModel.UserProfile{}).Error
}

// CheckEmailExists 检查邮箱是否存在
func (r *UserRepository) CheckEmailExists(email string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Table("user_profile").Where("email = ?", email)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// CheckUsernameExists 检查用户名是否存在
func (r *UserRepository) CheckUsernameExists(username string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Table("user_profile").Model(&sqlModel.UserProfile{}).Where("username = ?", username)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
