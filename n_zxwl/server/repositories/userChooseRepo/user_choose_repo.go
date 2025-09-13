// Package userChooseRepo repository/user_choose_repo.go
package userChooseRepo

import (
	"gorm.io/gorm"
	"mymod/model/sqlModel"
)

type UserChooseRepository struct {
	db *gorm.DB
}

// NewUserChooseRepository 创建志愿仓库实例
func NewUserChooseRepository(db *gorm.DB) UserChooseRepository {
	return UserChooseRepository{db: db}
}

// FindUserChoicesWithDetails 获取用户志愿详情列表
func (r *UserChooseRepository) FindUserChoicesWithDetails(userID int64) ([]sqlModel.UserChooseResponse, error) {
	var choices []sqlModel.UserChooseResponse

	err := r.db.Table("user_choose").
		Select("id, user_id, school_name, major_name, priority, create_time").
		Where("user_id = ?", userID).
		Order("priority ASC").
		Scan(&choices).Error

	if err != nil {
		return nil, err
	}

	return choices, nil
}

// CreateUserChoice 创建用户志愿
func (r *UserChooseRepository) CreateUserChoice(choice *sqlModel.UserChoose) error {
	return r.db.Create(choice).Error
}

// DeleteUserChoice 删除用户志愿
func (r *UserChooseRepository) DeleteUserChoice(id int64, userID int64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&sqlModel.UserChoose{}).Error
}

// GetUserChoiceByID 根据ID获取用户志愿
func (r *UserChooseRepository) GetUserChoiceByID(id int64, userID int64) (*sqlModel.UserChoose, error) {
	var choice sqlModel.UserChoose
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&choice).Error
	if err != nil {
		return nil, err
	}
	return &choice, nil
}

// DeleteUserChoicesByIDs 通过ID列表批量删除用户志愿
func (r *UserChooseRepository) DeleteUserChoicesByIDs(ids []int64, userID int64) error {
	return r.db.Where("id IN ? AND user_id = ?", ids, userID).Delete(&sqlModel.UserChoose{}).Error
}

// DeleteAllUserChoices 删除用户的所有志愿
func (r *UserChooseRepository) DeleteAllUserChoices(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&sqlModel.UserChoose{}).Error
}
