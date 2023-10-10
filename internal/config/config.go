package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig ServerConfig `json:"server" mapstructure:"server"`
	Redis        Redis        `json:"redis" mapstructure:"redis"`
	Proxy        Proxy        `json:"proxy" mapstructure:"proxy"`
	Bot          Bot          `json:"bot" mapstructure:"bot"`
	Index        Index        `json:"index" mapstructure:"index"`
	Gorm         Gorm         `json:"gorm" mapstructure:"gorm"`
	PrintConfig  bool         `json:"print_config" mapstructure:"print_config"`
}

var (
	C    = new(Config)
	once sync.Once
)

// Load config file (toml/json/yaml)
func MustLoad(path string) {
	once.Do(func() {
		viper.SetConfigFile(path)
		viper.ReadInConfig()
		err := viper.Unmarshal(&C)
		if err != nil {
			logrus.Warn(err)
		}
	})
}

func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

type ServerConfig struct {
	Host         string   `json:"host" mapstructure:"host"`
	GrpcProtocol string   `json:"grpc_protocol" mapstructure:"grpc_protocol"`
	GrpcPort     string   `json:"grpc_port" mapstructure:"grpc_port"`
	HttpPort     string   `json:"http_port" mapstructure:"http_port"`
	Cert         string   `json:"cert" mapstructure:"cert"`
	Key          string   `json:"key" mapstructure:"key"`
	CAName       string   `json:"ca_name" mapstructure:"ca_name"`
	AuthKey      string   `json:"auth_key" mapstructure:"auth_key"`
	Tokens       []string `json:"tokens" mapstructure:"tokens"`
}

type Redis struct {
	Addr      string `json:"addr" mapstructure:"addr"`
	Password  string `json:"password" mapstructure:"password"`
	DB        int    `json:"db" mapstructure:"db"`
	KeyPrefix string `json:"key_prefix" mapstructure:"key_prefix"`
}

type Bot struct {
	AppId    string     `json:"app_id" mapstructure:"app_id"`
	ApiHash  string     `json:"api_hash" mapstructure:"api_hash"`
	Token    string     `json:"token" mapstructure:"token"`         // 机器人token
	Menus    [][]string `json:"menus" mapstructure:"menus"`         // 菜单配置
	UseCache bool       `json:"use_cache" mapstructure:"use_cache"` // 使用缓存
	Manager  BotManager `json:"manager" mapstructure:"manager"`     // 管理信息
}

type Index struct {
	Languages []string `json:"languages" mapstructure:"languages"` // 启动的语言
	PageSize  int64    `json:"page_size" mapstructure:"page_size"` // 搜索单页大小
	MaxPage   int64    `json:"max_page" mapstructure:"max_page"`   // 最大页码
	Detail    bool     `json:"detail" mapstructure:"detail"`       // 详情模式(直接返回详情信息)
	Order     string   `json:"order" mapstructure:"order"`         // 排序
	Tags      []string `json:"tags" mapstructure:"tags"`           // 数据标记, 建立索引
}

type Proxy struct {
	Protocal string `json:"protocal" mapstructure:"protocal"`
	Ip       string `json:"ip" mapstructure:"ip"`
	Port     int    `json:"port" mapstructure:"port"`
}

type Gorm struct {
	Debug             bool   `json:"debug" mapstructure:"debug"`
	DBType            string `json:"db_type" mapstructure:"db_type"`
	DbName            string `json:"db_name" mapstructure:"db_name"`
	MaxLifetime       int    `json:"max_lifetime" mapstructure:"max_lifetime"`
	MaxOpenConns      int    `json:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns      int    `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	TablePrefix       string `json:"table_prefix" mapstructure:"table_prefix"`
	EnableAutoMigrate bool   `json:"enable_auto_migrate" mapstructure:"enable_auto_migrate"`
	Dsn               string `json:"dsn" mapstructure:"dsn"`
}

type BotManager struct {
	UserId      int64  `json:"user_id" mapstructure:"user_id"`
	Username    string `json:"username" mapstructure:"username"`
	DeleteDelay int    `json:"delete_delay" mapstructure:"delete_delay"`
}
