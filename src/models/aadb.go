package models

import (
	"gorm.io/gorm"
)

var db *gorm.DB

const (
	TableNameInstance    = "ledger_instance"
	TableNameCategory    = "ledger_category"
	TableNameCaterela    = "ledger_caterela"
	TableNameTransaction = "ledger_transaction"
	TableNameDetail      = "ledger_detail"
	TableNameTag         = "ledger_tag"
)

// SetConnectPool 设置models层db连接
func SetConnectPool(_db *gorm.DB) {
	db = _db
}
