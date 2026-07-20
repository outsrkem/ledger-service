package models

import "github.com/shopspring/decimal"

type DayStat struct {
	Day     string          `gorm:"column:day"`
	Income  decimal.Decimal `gorm:"column:income"`
	Expense decimal.Decimal `gorm:"column:expense"`
}

func (d *DB) DayTrend(from, to int64) ([]*DayStat, error) {
	var list []*DayStat

	err := d.WithInstance().
		Table("ledger_transaction").
		Where("occ_at >= ? AND occ_at < ?", from, to).
		Select(`
			FROM_UNIXTIME(FLOOR(occ_at / 1000 / 86400) * 86400, '%Y-%m-%d') AS day,
			IFNULL(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS income,
			IFNULL(ABS(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END)), 0) AS expense
		`).
		Group("day").
		Order("day ASC").
		Scan(&list).Error

	return list, err
}

// MonthTrend 按月统计收支，入参毫秒起止戳，返回格式 Day="2026-03"
func (d *DB) MonthTrend(from, to int64) ([]*DayStat, error) {
	var list []*DayStat
	err := d.WithInstance().
		Table("ledger_transaction").
		Where("occ_at >= ? AND occ_at < ?", from, to).
		Select(`
			FROM_UNIXTIME(occ_at / 1000, '%Y-%m') AS day,
			IFNULL(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS income,
			IFNULL(ABS(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END)), 0) AS expense
		`).
		Group("day").
		Order("day ASC").
		Scan(&list).Error
	return list, err
}

// YearTrend 按年汇总收支，参数毫秒起止时间戳
func (d *DB) YearTrend(from, to int64) ([]*DayStat, error) {
	var list []*DayStat
	err := d.WithInstance().
		Table("ledger_transaction").
		Where("occ_at >= ? AND occ_at < ?", from, to).
		Select(`
			FROM_UNIXTIME(occ_at / 1000, '%Y') AS day,
			IFNULL(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS income,
			IFNULL(ABS(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END)), 0) AS expense
		`).
		Group("day").
		Order("day ASC").
		Scan(&list).Error
	return list, err
}
