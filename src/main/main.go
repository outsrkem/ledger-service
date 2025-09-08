package main

import (
	"ledger/src/config"
	"ledger/src/database/mysql"
	"ledger/src/database/sql"
	"ledger/src/models"
	"ledger/src/pkg/migrator"
	"ledger/src/route"
	"ledger/src/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func init() {
	config.Initializer()
}

// AutoMigrator performs automatic database migrations using embedded SQL scripts
func AutoMigrator(klog *logrus.Entry, db *gorm.DB) {
	// Initialize migrator with embedded scripts
	m := migrator.NewDefault(db,
		&migrator.Config{
			Logger:         migrator.NewKlogLogger(klog),
			TableName:      "ledger_migrations",
			ScriptProvider: migrator.NewEmbeddedScriptProvider("script", sql.SQLScript),
		})

	// Execute migration workflow
	if err := migrator.RunWithMigrator(m); err != nil {
		panic("Database migration failed: " + err.Error())
	}
}

// main is the entry point of the ledger application
func main() {
	cfg := config.InitConfig()
	app := cfg.Ledger.App
	slog.InitLogger(&cfg.Ledger.Log)
	klog := slog.FromContext(nil)

	// Initialize database connection
	db := mysql.InitDB(cfg.Ledger)
	models.SetConnectPool(db)

	// Execute critical database migrations
	AutoMigrator(klog, db)

	klog.Info("start server...")
	svc := server.Default(
		server.WithHostPorts(app.Bind),
		server.WithExitWaitTime(0*time.Second))

	route.Middleware(svc)
	route.AppRoute(svc)
	svc.Spin()
}
