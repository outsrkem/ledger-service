package models

import (
	"github.com/shopspring/decimal"
)

// OrmDetail 明细详情
type OrmDetail struct {
	Kid           int64           `gorm:"column:kid;primaryKey"` // 主键ID
	TransactionId int64           `gorm:"column:transaction_id"` // 交易ID
	Name          string          `gorm:"column:name"`           // 物品名称（如“牛奶”）
	Quantity      decimal.Decimal `gorm:"column:quantity"`       // 数量（支持小数）
	Unit          string          `json:"unit"`                  // 单位
	Price         decimal.Decimal `gorm:"column:price"`          // 单价
	Total         decimal.Decimal `gorm:"column:total"`          // 该项总价
	Remark        string          `gorm:"column:remark"`         // 详细备注
}

func (OrmDetail) TableName() string {
	return TableNameDetail
}

// GetDetail 查询明细
func GetDetail(billId int64) ([]*OrmDetail, error) {
	var detail []*OrmDetail
	err := db.Model(&OrmDetail{}).
		Where("transaction_id = ?", billId).
		Find(&detail).Error

	return detail, err
}
