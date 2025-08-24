package models

import (
	"github.com/shopspring/decimal"
	"ledger/src/database/mysql"
)

// OrmTransaction 收支明细表
type OrmTransaction struct {
	Kid        int64           `gorm:"column:kid;primaryKey"` // 主键ID
	InstanceId string          `gorm:"column:instance_id"`    // 实例ID
	CategoryId int64           `gorm:"column:category_id"`    // 分类ID
	Amount     decimal.Decimal `gorm:"column:amount"`         // 金额,0.0000
	OccTime    string          `gorm:"column:occ_time"`       // 账目发生时间(occurrence)
	Remark     string          `gorm:"column:remark"`         // 详细备注
	UpdateTime int64           `gorm:"column:update_time"`    // 记录时间
	CreateTime int64           `gorm:"column:create_time"`    // 更新时间
}

func (OrmTransaction) TableName() string {
	return SetTableName("transaction")
}

// InstallTransaction 新增
func InstallTransaction(detail []*OrmTransaction) error {
	return mysql.OrmDB.Create(&detail).Error
}

// FindTransactionForUser 按用户实例和id查询
func FindTransactionForUser(instanceId string, transactionId int64) ([]*OrmTransaction, error) {
	var detail []*OrmTransaction
	r := mysql.OrmDB.Model(&OrmTransaction{}).
		Where("instance_id = ? AND kid = ?", instanceId, transactionId).Find(&detail)
	return detail, r.Error
}

func DeleteTransaction(instanceId string, transactionId int64) error {
	r := mysql.OrmDB.Model(&OrmTransaction{}).
		Where("instance_id = ? AND kid = ?", instanceId, transactionId).Delete(&OrmTransaction{})
	return r.Error
}

func UpdateTransaction(detail OrmTransaction) error {
	return db.Model(&OrmTransaction{}).Updates(map[string]interface{}{
		"category_id": detail.CategoryId,
		"occ_time":    detail.OccTime,
		"amount":      detail.Amount,
		"remark":      detail.Remark,
		"update_time": detail.UpdateTime,
	}).Error
}

func FindTransactionAll(instanceId string, limit, offset int, count *int64) ([]*OrmTransaction, error) {
	var detail []*OrmTransaction
	err := mysql.OrmDB.Model(&OrmTransaction{}).
		Where("instance_id = ?", instanceId).
		Count(count).Order("occ_time DESC").
		Limit(limit).Offset(offset).
		Find(&detail).Error

	return detail, err
}
