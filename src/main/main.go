package main

import (
	"ledger/src/config"
	"ledger/src/database/mysql"
	"ledger/src/database/sql"
	"ledger/src/pkg/migrator"
	"ledger/src/route"
	"ledger/src/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/sirupsen/logrus"
)

func init() {
	config.Initializer()
}

// AutoMigrator performs automatic database migrations using embedded SQL scripts
func AutoMigrator(klog *logrus.Entry, tableName string) {
	// Initialize migrator with embedded scripts
	m := migrator.NewDefault(mysql.OrmDB,
		&migrator.Config{
			Logger:         migrator.NewKlogLogger(klog),
			TableName:      tableName,
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
	mysql.InitDB(cfg.Ledger)
	// Execute critical database migrations
	AutoMigrator(klog, "ledger_migrations")

	klog.Info("start server...")
	svc := server.Default(
		server.WithHostPorts(app.Bind),
		server.WithExitWaitTime(0*time.Second))

	route.Middleware(svc)
	route.AppRoute(svc)
	svc.Spin()
}
