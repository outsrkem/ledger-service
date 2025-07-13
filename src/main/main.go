package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"ledger/src/config"
	"ledger/src/database/mysql"
	"ledger/src/route"
	"ledger/src/slog"
	"time"
)

func init() {
	config.Initializer()
}

func main() {
	cfg := config.InitConfig()
	app := cfg.Ledger.App
	slog.InitLogger(&cfg.Ledger.Log)
	klog := slog.FromContext(nil)
	mysql.InitDB(cfg.Ledger)

	klog.Info("start server")
	svc := server.Default(server.WithHostPorts(app.Bind), server.WithExitWaitTime(0*time.Second))
	route.Middleware(svc)
	route.AppRoute(svc)
	svc.Spin()
}
