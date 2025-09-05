package models

import (
	"ledger/src/database/mysql"
)

var db = mysql.OrmDB

const (
	TableNameInstance    = "ledger_instance"
	TableNameCategory    = "ledger_category"
	TableNameTransaction = "ledger_transaction"
	TableNameDetail      = "ledger_detail"
	TableNameTag         = "ledger_tag"
)
