// service/user_choose_service.go
package service

import (
	"fmt"
	"mymod/model/sqlModel"
	"mymod/repositories/userChooseRepo"
)

type UserChooseService struct {
	repo userChooseRepo.UserChooseRepository
}

func NewUserChooseService(repo userChooseRepo.UserChooseRepository) UserChooseService {
	return UserChooseService{repo: repo}
}

// GetUserChoices 获取
func (s *UserChooseService) GetUserChoices(userID int64) ([]sqlModel.UserChooseResponse, error) {
	return s.repo.FindUserChoicesWithDetails(userID)
}

// CreateUserChoice 创建
func (s *UserChooseService) CreateUserChoice(choice *sqlModel.UserChoose) error {
	if !choice.IsValid() {
		return fmt.Errorf("无效的志愿数据")
	}
	return s.repo.CreateUserChoice(choice)
}

// DeleteUserChoice 删除
func (s *UserChooseService) DeleteUserChoice(id int64, userID int64) error {
	// 先检查是否存在
	_, err := s.repo.GetUserChoiceByID(id, userID)
	if err != nil {
		return fmt.Errorf("志愿不存在或无权操作")
	}
	return s.repo.DeleteUserChoice(id, userID)
}

// DeleteAllUserChoices 删除所有
func (s *UserChooseService) DeleteAllUserChoices(userID int64) error {
	return s.repo.DeleteAllUserChoices(userID)
}

// DeleteUserChoicesByIDs 批量删除
func (s *UserChooseService) DeleteUserChoicesByIDs(ids []int64, userID int64) error {
	return s.repo.DeleteUserChoicesByIDs(ids, userID)
}
