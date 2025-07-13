package cfgtypes

// Config yaml配置结构体
type Config struct {
	Ledger *Ledger `yaml:"ledger"`
}

type Ledger struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
	Ats      Ats      `yaml:"ats"`
	Uias     Uias     `yaml:"uias"`
	Log      Log      `yaml:"log"`
}

type App struct {
	Bind string `yaml:"bind"`
}

type Database struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Name   string `yaml:"name"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
}

type Ats struct {
	Endpoint      string `yaml:"endpoint"`
	SkipTlsVerify bool   `yaml:"skipTlsVerify"`
}

type Uias struct {
	Endpoint        string `yaml:"endpoint"`
	SkipTlsVerify   bool   `yaml:"skipTlsVerify"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
}

type Log struct {
	Level  string `yaml:"level"`
	Output Output `yaml:"output"`
}

type Output struct {
	File   File   `yaml:"file"`
	Stdout string `yaml:"stdout"`
}

type File struct {
	Name       string `yaml:"name"`
	MaxSize    int    `yaml:"maxsize"`
	MaxBackups int    `yaml:"maxbackups"`
	MaxAge     int    `yaml:"maxage"`
	Compress   bool   `yaml:"compress"`
}
