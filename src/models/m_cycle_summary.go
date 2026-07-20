package models

import "github.com/shopspring/decimal"

type CycleSummary struct {
	Income  decimal.Decimal `gorm:"column:income"`  // 总收入（正数）
	Expense decimal.Decimal `gorm:"column:expense"` // 总支出（正数）
}

// CycleSummary 周期汇总
func (d *DB) CycleSummary(from, to int64) (*CycleSummary, error) {
	var res CycleSummary

	db := d.WithInstance().Table(TableNameTransaction)
	db = db.Where("occ_at >= ? AND occ_at < ?", from, to)
	err := db.Select(`
		IFNULL(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS income,
		IFNULL(ABS(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END)), 0) AS expense
	`).Scan(&res).Error

	return &res, err
}
