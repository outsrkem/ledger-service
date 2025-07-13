package models

import (
	"ledger/src/database/mysql"
)

var db = mysql.OrmDB

const TablesPrefix = "ledger_" // 表前缀

func SetTableName(name string) string {
	return TablesPrefix + name
}
