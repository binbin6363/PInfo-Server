package config

import (
	"os"

	"log"

	"gopkg.in/yaml.v2"
)

type ServerInfo struct {
	Name          string `yaml:"name"`
	Listen        string `yaml:"listen"`
	Timeout       int    `yaml:"timeout"`
	Secret        string `yaml:"secret"`
	TokenExpire   int    `yaml:"token_expire"`
	DataCenterId  int64  `yaml:"data_center_id"`
	WorkerId      int64  `yaml:"worker_id"`
	DebugReqRsp   bool   `yaml:"debug_req_rsp"`
	ResourceRoot  string `yaml:"resource_root"`
	RemoteUrlRoot string `yaml:"remote_url_root"`
	Mode          string `yaml:"mode"` // debug, release, test
}

// DBInfo db信息
type DBInfo struct {
	// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	Dsn          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxLifeTime  int    `yaml:"max_life_time"` // 单位秒
}

// ConnInfo 接入服务信息
type ConnInfo struct {
	Addr    string `yaml:"addr"`
	Timeout int    `yaml:"timeout"`
}

type ServerCfg struct {
	ServerInfo *ServerInfo `yaml:"server"`
	DBInfo     *DBInfo     `yaml:"db"`
	ConnInfo   *ConnInfo   `yaml:"conn"`
	CosInfo    *CosInfo    `yaml:"cos"`
	LogInfo    *LogInfo    `yaml:"log"`
}

type CosInfo struct {
	SecretID       string `yaml:"secret_id"`
	SecretKey      string `yaml:"secret_key"`
	Domain         string `yaml:"domain"`
	Region         string `yaml:"region"`
	DisableSSL     bool   `yaml:"disable_ssl"`
	ForcePathStyle bool   `yaml:"force_path_style"`
	AvatarBucket   string `yaml:"avatar_bucket"`
	Expire         int64  `yaml:"expire"` // 单位小时
}

type LogInfo struct {
	Path       string `yaml:"path"`
	Level      int    `yaml:"level"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackUps int    `yaml:"max_backups"`
}

// 配置实例
var cfg = &ServerCfg{}

// AppConfig 获取配置单例
func AppConfig() *ServerCfg {
	return cfg
}

// Init 初始化配置
func Init(file string) {
	configFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("load conf fail, path:%s, err:%v", file, err)
	}

	if err = yaml.Unmarshal(configFile, cfg); err != nil {
		log.Fatalf("Unmarshal conf fail, err:%v", err)
	}

	log.Printf("load conf ok, path:%s, conf:%v", file, string(configFile))
}
