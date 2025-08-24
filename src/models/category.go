package models

import (
	"ledger/src/database/mysql"
	"time"
)

// OrmCategory 分类表
type OrmCategory struct {
	Kid        int64  `gorm:"column:kid;primaryKey"` // 主键ID
	InstanceId string `gorm:"column:instance_id"`    // 实例ID
	Name       string `gorm:"column:name"`           // 分类名称
	Direction  int8   `gorm:"column:direction"`      // 1=收入，2=支出
	Pid        int64  `gorm:"column:pid"`            // 父类的ID
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

// FindCategoryByUser 按类型查询用户的所有分类
func FindCategoryByUser(instanceId string, Direction int8) ([]*OrmCategory, error) {
	var category []*OrmCategory
	r := mysql.OrmDB.Model(&OrmCategory{}).
		Where("instance_id is NULL OR instance_id = ?", instanceId).
		Where("direction = ?", Direction).
		Find(&category)
	return category, r.Error
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
