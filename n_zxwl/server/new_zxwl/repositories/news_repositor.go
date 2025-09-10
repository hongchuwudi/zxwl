package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"mymod/new_zxwl/model/sqlModel"
)

type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

// Create 新增资讯
func (r *NewsRepository) Create(news *sqlModel.NewsInfo) error {
	err := r.db.Create(news).Error
	if err != nil {
		return fmt.Errorf("新增资讯失败: %v", err)
	}
	return nil
}

// Update 更新资讯
func (r *NewsRepository) Update(news *sqlModel.NewsInfo) error {
	err := r.db.Model(&sqlModel.NewsInfo{}).Where("id = ?", news.ID).Updates(news).Error
	if err != nil {
		return fmt.Errorf("更新资讯失败: %v", err)
	}
	return nil
}

// Delete 删除资讯
func (r *NewsRepository) Delete(id uint64) error {
	err := r.db.Where("id = ?", id).Delete(&sqlModel.NewsInfo{}).Error
	if err != nil {
		return fmt.Errorf("删除资讯失败: %v", err)
	}
	return nil
}

// GetByNewsID 按ID查询资讯
func (r *NewsRepository) GetByNewsID(newsID int64) (*sqlModel.NewsInfo, error) {
	var news sqlModel.NewsInfo
	err := r.db.Table("news_info").Where("id = ?", newsID).First(&news).Error
	if err != nil {
		return nil, fmt.Errorf("查询资讯失败: %v", err)
	}
	return &news, nil
}

// FindByCondition 分页条件查询资讯
func (r *NewsRepository) FindByCondition(condition map[string]interface{}, page, pageSize int) ([]sqlModel.NewsInfo, int64, error) {
	var newsList []sqlModel.NewsInfo
	var total int64

	// 构建查询条件并指定表名
	query := r.db.Table("news_info") // 替换为实际的表名

	// 添加条件 - 逐个处理，字符串类型使用模糊查询
	for key, value := range condition {
		if value != "" && value != nil {
			// 对字符串类型的字段使用模糊查询
			switch key {
			case "title", "keywords", "from_source":
				query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
			default:
				query = query.Where(key+" = ?", value)
			}
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("publish_time DESC").Find(&newsList).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询资讯列表失败: %v", err)
	}

	return newsList, total, nil
}

// BatchCreate 批量新增资讯
func (r *NewsRepository) BatchCreate(newsList []sqlModel.NewsInfo) error {
	err := r.db.Create(&newsList).Error
	if err != nil {
		return fmt.Errorf("批量新增资讯失败: %v", err)
	}
	return nil
}

// UpdateContent 更新资讯内容
func (r *NewsRepository) UpdateContent(id uint64, content, keywords string) error {
	updateData := map[string]interface{}{
		"content":  content,
		"keywords": keywords,
	}

	err := r.db.Model(&sqlModel.NewsInfo{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return fmt.Errorf("更新资讯内容失败: %v", err)
	}
	return nil
}

// UpdateCount 更新资讯的点赞、收藏、分享数量
// field: 要更新的字段名 ("like_count", "favorite_count", "share_count")
// action: 操作类型 ("increment" 增加, "decrement" 减少)
func (r *NewsRepository) UpdateCount(newsID uint64, field string, action string) error {
	if newsID == 0 {
		return fmt.Errorf("资讯ID不能为空")
	}

	// 验证字段名
	validFields := map[string]bool{
		"like_count":     true,
		"favorite_count": true,
		"share_count":    true,
	}

	if !validFields[field] {
		return fmt.Errorf("无效的字段名: %s", field)
	}

	// 构建更新表达式
	var updateValue interface{}
	if action == "increment" {
		updateValue = gorm.Expr(field + " + 1")
	} else if action == "decrement" {
		updateValue = gorm.Expr("GREATEST(" + field + " - 1, 0)")
	} else {
		return fmt.Errorf("无效的操作类型: %s", action)
	}

	// 使用正确的Updates语法
	result := r.db.Model(&sqlModel.NewsInfo{}).
		Where("id = ?", newsID).
		Update(field, updateValue)

	if result.Error != nil {
		return fmt.Errorf("更新%s失败: %v", field, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("资讯不存在或未更新任何记录, ID: %d", newsID)
	}

	return nil
}

// GetCount 获取资讯的统计数量
func (r *NewsRepository) GetCount(newsID uint64) (map[string]int, error) {
	var news sqlModel.NewsInfo
	err := r.db.Select("like_count, favorite_count, share_count, comment_count").
		Where("id = ?", newsID).
		First(&news).Error

	if err != nil {
		return nil, fmt.Errorf("获取资讯统计失败: %v", err)
	}

	counts := map[string]int{
		"like_count":     news.LikeCount,
		"favorite_count": news.FavoriteCount,
		"share_count":    news.ShareCount,
		"comment_count":  news.CommentCount,
	}

	return counts, nil
}
