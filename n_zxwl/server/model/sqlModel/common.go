package sqlModel

// CommonProvince 省份信息表
type CommonProvince struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"unique;not null" json:"name"`
	ShortName string `gorm:"not null" json:"short_name"`
	Code      string `gorm:"unique;not null" json:"code"`
	Pinyin    string `gorm:"not null" json:"pinyin"`
}

// CommonCity 城市信息表
type CommonCity struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	CityName   string `gorm:"not null" json:"city_name"`
	ProvinceID int    `gorm:"not null" json:"province_id"`
	CityLevel  int    `json:"city_level"`
}

// CommonDistrict 区县信息表
type CommonDistrict struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	DistrictName string `gorm:"not null" json:"district_name"`
	CityID       int    `gorm:"not null" json:"city_id"`
	DistrictType int    `json:"district_type"`
	IsUrban      bool   `json:"is_urban"`
}
