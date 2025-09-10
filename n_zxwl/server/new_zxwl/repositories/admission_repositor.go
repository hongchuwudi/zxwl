package repositories

import (
	"errors"
	"gorm.io/gorm"
	"mymod/new_zxwl/model/sqlModel"
	"strings"
)

// AdmissionRepository 数据库操作
type AdmissionRepository struct {
	db *gorm.DB
}

// NewAdmissionRepository 构造函数
func NewAdmissionRepository(db *gorm.DB) *AdmissionRepository {
	return &AdmissionRepository{db: db}
}

// GetTypeIDByName 根据高考类型名称查询类型ID
func (r *AdmissionRepository) GetTypeIDByName(name string) (int, error) {
	var admissionType sqlModel.AdmissionType

	// 精确匹配优先
	err := r.db.Where("name like CONCAT('%',?,'%') ", name).First(&admissionType).Error
	if err == nil {
		return admissionType.ID, nil
	}

	// 如果精确匹配没找到，尝试模糊匹配
	var admissionTypes []sqlModel.AdmissionType
	err = r.db.Where("name LIKE ?", "%"+name+"%").Find(&admissionTypes).Error
	if err != nil {
		return 0, err
	}

	// 检查匹配结果
	if len(admissionTypes) == 0 {
		return 0, errors.New("未找到匹配的高考类型")
	}

	if len(admissionTypes) > 1 {
		// 如果有多个匹配，可以返回第一个或者报错
		// 这里返回第一个匹配的
		return admissionTypes[0].ID, nil
	}

	return admissionTypes[0].ID, nil
}

// GetTypeByName 根据高考类型名称查询类型信息
func (r *AdmissionRepository) GetTypeByName(name string) (*sqlModel.AdmissionType, error) {
	var admissionType sqlModel.AdmissionType

	// 精确匹配优先
	err := r.db.Where("name = ?", name).First(&admissionType).Error
	if err == nil {
		return &admissionType, nil
	}

	// 模糊匹配
	err = r.db.Where("name LIKE ?", "%"+name+"%").First(&admissionType).Error
	if err != nil {
		return nil, err
	}

	return &admissionType, nil
}

// GetAllTypes 获取所有高考类型
func (r *AdmissionRepository) GetAllTypes() ([]sqlModel.AdmissionType, error) {
	var types []sqlModel.AdmissionType
	err := r.db.Find(&types).Error
	return types, err
}

// FilterAdmissionData 调用存储过程筛选录取数据
func (r *AdmissionRepository) FilterAdmissionData(input sqlModel.UserInput) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 将数组转换为逗号分隔的字符串
	level2Str := strings.Join(input.Level2, ",")
	level3Str := strings.Join(input.Level3, ",")

	// 调用存储过程
	err := r.db.Raw("CALL sp_FilterAdmissionData(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		input.GoalYear,
		input.TypeID,
		input.ProvinceID,
		input.EquivalentScoreStart,
		input.EquivalentScoreEnd,
		input.CurrRankStart,
		input.CurrRankEnd,
		input.IsBenKe,
		input.GoalProvinceID,
		level2Str,
		level3Str,
		input.Salary,
	).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
