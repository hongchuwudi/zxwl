package commonRepo

import (
	"gorm.io/gorm"
	"mymod/model/sqlModel"
)

// ProvinceRepository 省份数据库操作
type ProvinceRepository struct {
	db *gorm.DB
}

// NewProvinceRepository 构造函数
func NewProvinceRepository(db *gorm.DB) *ProvinceRepository {
	return &ProvinceRepository{db: db}
}

// GetProvinceIDByName 根据省份名称模糊查询省份ID
func (r *ProvinceRepository) GetProvinceIDByName(name string) (int, error) {
	var provinces []sqlModel.CommonProvince

	// 使用GORM的Where方法进行模糊查询
	err := r.db.
		Table("common_provinces").
		Where("name LIKE ?", "%"+name+"%").
		Find(&provinces).Error

	if err != nil {
		return provinces[0].ID, err
	}

	return provinces[0].ID, nil
}
