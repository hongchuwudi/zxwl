// Package service file: services/user_service.go
package service

import (
	"errors"
	userParam2 "mymod/model/param/userParam"
	"mymod/model/sqlModel"
	userRepo2 "mymod/repositories/userRepo"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo       *userRepo2.UserRepository
	verifyCodeRepo *userRepo2.VerifyCodeRepository
}

func NewUserService(userRepo *userRepo2.UserRepository, verifyCodeRepo *userRepo2.VerifyCodeRepository) *UserService {
	return &UserService{
		userRepo:       userRepo,
		verifyCodeRepo: verifyCodeRepo,
	}
}

// Login 用户登录
func (s *UserService) Login(req userParam2.LoginRequest) (*userParam2.LoginResponse, error) {
	var user *sqlModel.UserProfile
	var err error

	// 判断是邮箱登录还是用户名登录
	if contains(req.Login, "@") {
		user, err = s.userRepo.GetUserByEmail(req.Login)
	} else {
		user, err = s.userRepo.GetUserByUsername(req.Login)
	}

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成token（这里简化处理，实际应该使用JWT等）
	token := generateToken(user.ID)
	user.DeviceInfo = req.DeviceInfo

	// 更新最后在线时间
	now := time.Now()
	user.LastOnlineTime = &now
	user.IsOnline = true
	s.userRepo.UpdateUser(user)

	return &userParam2.LoginResponse{
		User:        user,
		AccessToken: token,
	}, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID int) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	return s.userRepo.DeleteUser(userID)
}

// GetUserByID 获取用户信息
func (s *UserService) GetUserByID(userID int) (*sqlModel.UserProfile, error) {
	return s.userRepo.GetUserByID(userID)
}

// GetUserByEmail 获取用户信息
func (s *UserService) GetUserByEmail(email string) (*sqlModel.UserProfile, error) {
	return s.userRepo.GetUserByEmail(email)
}

// Register 用户注册（带验证码验证）
func (s *UserService) Register(req userParam2.RegisterRequest) (*sqlModel.UserProfile, error) {
	// 验证验证码
	if req.VerifyCode == "" {
		return nil, errors.New("验证码不能为空")
	}

	latestCode, err := s.verifyCodeRepo.GetLatestVerifyCode(req.Email)
	if err != nil {
		return nil, errors.New("验证码获取失败")
	}

	if latestCode.Code != req.VerifyCode {
		return nil, errors.New("验证码错误")
	}

	if time.Now().After(latestCode.ExpiresAt) {
		return nil, errors.New("验证码已过期")
	}

	// 检查邮箱是否已存在
	emailExists, err := s.userRepo.CheckEmailExists(req.Email, 0)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("邮箱已被注册")
	}

	// 检查用户名是否已存在
	usernameExists, err := s.userRepo.CheckUsernameExists(req.Username, 0)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.New("用户名已被使用")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	//fmt.Println("user: ", req.DeviceInfo)

	// 创建用户
	user := &sqlModel.UserProfile{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeviceInfo:   req.DeviceInfo,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID int, req sqlModel.UserProfile) (*sqlModel.UserProfile, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 检查邮箱是否被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		emailExists, err := s.userRepo.CheckEmailExists(req.Email, userID)
		if err != nil {
			return nil, err
		}
		if emailExists {
			return nil, errors.New("邮箱已被其他用户使用")
		}
		user.Email = req.Email
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Gender != 0 {
		user.Gender = req.Gender
	}
	if req.BirthYear != 0 {
		user.BirthYear = req.BirthYear
	}
	if req.Location != "" {
		user.Location = req.Location
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	user.UpdatedAt = time.Now()

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 通过邮件验证码更改密码
func (s *UserService) ChangePassword(req userParam2.ChangePasswordRequest) error {
	// 验证必要字段
	if req.Email == "" || req.VerifyCode == "" || req.NewPassword == "" {
		return errors.New("邮箱、验证码和新密码不能为空")
	}

	// 验证验证码（与注册逻辑保持一致）
	latestCode, err := s.verifyCodeRepo.GetLatestVerifyCode(req.Email)
	if err != nil {
		return errors.New("验证码获取失败")
	}

	if latestCode == nil || latestCode.Code != req.VerifyCode {
		return errors.New("验证码错误")
	}

	if time.Now().After(latestCode.ExpiresAt) {
		return errors.New("验证码已过期")
	}

	// 获取用户信息
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.UpdateUser(user)
}

// GetUsersWithPagination 分页查询用户
func (s *UserService) GetUsersWithPagination(req userParam2.UserQueryRequest) (*userParam2.PaginationResponse, error) {
	users, total, err := s.userRepo.GetUsersWithPagination(req)
	if err != nil {
		return nil, err
	}

	// 移除敏感信息
	for i := range users {
		users[i].PasswordHash = ""
		users[i].AuthToken = ""
	}

	return &userParam2.PaginationResponse{
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
		Data:  users,
	}, nil
}

// 辅助函数
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func generateToken(userID int) string {
	// 这里应该使用JWT生成token，简化处理
	return "token_" + string(rune(userID)) + "_" + time.Now().Format("20060102150405")
}
