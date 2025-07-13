package mysql

import (
	"ledger/src/cfgtypes"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

func gormLogLev(lev string) logger.LogLevel {
	lowerLev := strings.ToLower(lev)
	switch lowerLev {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	case "fatal":
		return logger.Error
	default:
		return logger.Warn
	}
}

func NewGormLogger(cfg *cfgtypes.Ledger) logger.Interface {
	return logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,    // 慢查询阈值
			LogLevel:                  gormLogLev(cfg.Log.Level), // 日志级别
			IgnoreRecordNotFoundError: true,                      // 忽略记录未找到错误
			Colorful:                  false,                     // 禁用颜色
		},
	)
}
