package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"ledger/src/cfgtypes"
	"ledger/src/slog"
	"os"
)

// FlagArgs represents the command-line arguments for the application.
type FlagArgs struct {
	CfgPath      string
	PrintVersion bool
	Plain        string // 接收命令行字符串，用于加密
}

// NewFlagArgs creates a new FlagArgs object and parses command line flags.
func NewFlagArgs() *FlagArgs {
	fa := &FlagArgs{}
	flag.StringVar(&fa.CfgPath, "c", "ledger.yaml", "Configuration file path.")
	flag.BoolVar(&fa.PrintVersion, "version", false, "Print version information and quit.")
	flag.StringVar(&fa.Plain, "encrypt", "", "Encrypted string.")
	flag.Parse()
	return fa
}

// configPath is a global variable that stores a pointer to the configuration file path.
var configPath *string
var AppCfg *cfgtypes.Config

// Initializer function is used to initialize the application's configuration.
func Initializer() {
	fa := NewFlagArgs()
	if fa.PrintVersion { // 显示版本
		versions, _ := newVersions(Version, GoVersion, GitCommit)
		versions.Print(versions)
	}
	if fa.Plain != "" { // 加密命令行字符串
		encryption(fa.Plain)
	}
	configPath = &fa.CfgPath
}

// InitConfig 初始化配置
func InitConfig() *cfgtypes.Config {
	klog := slog.FromContext(nil)
	klog.Info("Read configuration file: ", *configPath)

	configData, err := os.ReadFile(*configPath)
	if err != nil {
		klog.Error("Failed to read the configuration file: ", err)
		os.Exit(1)
	}

	var cfg cfgtypes.Config
	// 解析配置文件
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		klog.Errorf("Unmarshal configuration file error: %v", err)
		os.Exit(1)
	}

	decryCfgCipher(&cfg.Ledger.Database.Passwd)      // 解密数据库密码
	decryCfgCipher(&cfg.Ledger.Uias.AccessKeySecret) // 解密uias的sk
	AppCfg = &cfg
	return &cfg
}
