package models

import "time"

// OrmCategory 分类表
type OrmCategory struct {
	Kid        int    `gorm:"column:kid;primaryKey"` // 主键ID
	Title      string `gorm:"column:title"`          // 分类名称
	Pid        int    `gorm:"column:pid"`            // 父类的ID,1代表大类
	Type       string `gorm:"column:type"`           // 收入1，支出2
	Sort       uint8  `gorm:"column:sort"`           // 数字越小越靠前
	UpdateTime int64  `gorm:"column:update_time"`    // 记录时间
	CreateTime int64  `gorm:"column:create_time"`    // 更新时间
}

func (OrmCategory) TableName() string {
	return SetTableName("category")
}

// FindCategoryByType 查询收入或支出的分类
func FindCategoryByType(t string) ([]*OrmCategory, error) {
	var category []*OrmCategory
	err := db.Model(&OrmCategory{}).Where("type = ?", t).Find(&category).Error
	return category, err
}

// FindCategoryAll 查询所有分类
func FindCategoryAll(pageSize, page int, count *int64) ([]*OrmCategory, error) {
	var category []*OrmCategory
	err := db.Model(&OrmCategory{}).Count(count).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&category).Error
	return category, err
}

// CreateCategory 创建分类
func CreateCategory(category []*OrmCategory) error {
	return db.Model(&OrmCategory{}).Create(&category).Error
}

// DeleteCategory 删除分类
func DeleteCategory(ids []int) (int64, error) {
	r := db.Model(&OrmCategory{}).Where("kid IN (?)", ids).Delete(&OrmCategory{})
	return r.RowsAffected, r.Error
}

// UpdateCategory 更新分类
func UpdateCategory(title string, pid uint32) error {
	r := db.Model(&OrmCategory{}).Updates(map[string]interface{}{
		"title":       title,
		"pid":         pid,
		"update_time": time.Now().UnixMilli(),
	})
	return r.Error
}
