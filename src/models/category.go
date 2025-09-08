package models

import (
	"time"

	"gorm.io/gorm"
)

// OrmCategory 分类表
type OrmCategory struct {
	Kid        int64  `gorm:"column:kid;primaryKey"` // 主键ID
	InstanceId string `gorm:"column:instance_id"`    // 实例ID
	Name       string `gorm:"column:name"`           // 分类名称
	Direction  int8   `gorm:"column:direction"`      // 1=收入，2=支出
	Layer      int8   `gorm:"column:layer"`          // 1=大类，2=子类
	CreateTime int64  `gorm:"column:create_time"`    // 创建时间
}

func (OrmCategory) TableName() string {
	return TableNameCategory
}

// OrmCaterela 分类关联表
type OrmCaterela struct {
	Kid        int64 `gorm:"column:kid;primaryKey"` // 主键ID
	ParentId   int64 `gorm:"column:parent_id"`      // 父ID
	ChildId    int64 `gorm:"column:child_id"`       // 子项ID
	CreateTime int64 `gorm:"column:create_time"`    // 创建时间
}

func (OrmCaterela) TableName() string {
	return TableNameCaterela
}

// CategoryWithRelResult 自定义结果结构体：接收关联查询后的字段（用rel.parent_id作为pid）
type CategoryWithRelResult struct {
	Kid        int64  `gorm:"column:kid" json:"kid"`                 // 对应ledger_category.kid
	InstanceId string `gorm:"column:instance_id" json:"instance_id"` // 对应ledger_category.instance_id
	Name       string `gorm:"column:name" json:"name"`               // 对应ledger_category.name
	Direction  int8   `gorm:"column:direction" json:"direction"`     // 对应ledger_category.direction
	Pid        int64  `gorm:"column:pid" json:"pid"`                 // 对应ledger_caterela.parent_id（别名pid）
	CreateTime int64  `gorm:"column:create_time" json:"create_time"` // 对应ledger_category.create_time
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
	r := db.Model(&OrmCategory{}).
		Where("instance_id is NULL OR instance_id = ?", instanceId).
		Where("direction = ?", Direction).
		Find(&category)
	return category, r.Error
}

// GetCategoryWithRel 关联查询分类（修复表别名问题）
func GetCategoryWithRel(instanceId string, direction int8) ([]*CategoryWithRelResult, error) {
	var results []*CategoryWithRelResult

	err := db.Table(TableNameCategory+" AS c").
		Select(`c.kid, c.instance_id, c.name, c.direction, rel.parent_id AS pid, c.create_time`).
		Joins("LEFT JOIN "+TableNameCaterela+" AS rel ON c.kid = rel.child_id").
		Where("c.instance_id IS NULL OR c.instance_id = ?", instanceId).
		Where("c.direction = ?", direction).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// CreateCategory 创建分类
func CreateCategory(category *OrmCategory, rela *OrmCaterela) (int64, error) {
	var kid int64
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(category).Error; err != nil {
			return err
		}
		kid = category.Kid
		if rela != nil {
			rela.ChildId = kid
			if err := tx.Create(&rela).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return kid, err
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
