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
	AppId   string     `json:"app_id" mapstructure:"app_id"`
	ApiHash string     `json:"api_hash" mapstructure:"api_hash"`
	Token   string     `json:"token" mapstructure:"token"`
	Start   string     `json:"start" mapstructure:"start"`
	Menus   [][]string `json:"menus" mapstructure:"menus"`
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
