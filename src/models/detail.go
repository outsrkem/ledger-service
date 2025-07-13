package models

import "github.com/shopspring/decimal"

// OrmDetail 收支明细表
type OrmDetail struct {
	Kid        int64           `gorm:"column:kid;primaryKey"` // 主键ID
	CategoryId int64           `gorm:"column:category_id"`    // 分类ID
	OccTime    int64           `gorm:"column:occ_time"`       // 账目发生时间(occurrence)
	Amount     decimal.Decimal `gorm:"column:amount"`         // 金额
	Total      int64           `gorm:"column:total"`          // 是否计入本月收支，1：计入；0：不计入
	Direction  int64           `gorm:"column:direction"`      // 金钱流向，1：支出；2：收入
	Status     int64           `gorm:"column:status"`         // 状态，1：封存，不可编辑；2：可编辑
	Remark     string          `gorm:"column:remark"`         // 详细备注
	UpdateTime int64           `gorm:"column:update_time"`    // 记录时间
	CreateTime int64           `gorm:"column:create_time"`    // 更新时间
}

func (OrmDetail) TableName() string {
	return SetTableName("detail")
}

// InstallDetail 新增
func InstallDetail(detail []*OrmDetail) error {
	return db.Create(detail).Error
}
func DeleteDetail(detail []*OrmDetail) error {
	return db.Delete(detail).Error
}
func UpdateDetail(detail OrmDetail) error {
	return db.Model(&OrmDetail{}).Updates(map[string]interface{}{
		"category_id": detail.CategoryId,
		"occ_time":    detail.OccTime,
		"amount":      detail.Amount,
		"total":       detail.Total,
		"direction":   detail.Direction,
		"remark":      detail.Remark,
		"update_time": detail.UpdateTime,
	}).Error
}
