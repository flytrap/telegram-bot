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
	UseIndex bool       `json:"use_index" mapstructure:"use_index"` // 使用索引，否则纯数据库查询
	Manager  BotManager `json:"manager" mapstructure:"manager"`     // 管理信息
}

type Index struct {
	Name       string     `json:"name" mapstructure:"name"`               // 索引名称前缀
	Language   string     `json:"language" mapstructure:"language"`       // 语言设定
	Mode       string     `json:"mode" mapstructure:"mode"`               // 索引模式(min:小内存模式,normal,max:大内存模式)
	PageSize   int64      `json:"page_size" mapstructure:"page_size"`     // 搜索单页大小
	MaxPage    int64      `json:"max_page" mapstructure:"max_page"`       // 最大页码
	ItemMode   string     `json:"item_mode" mapstructure:"item_mode"`     // 条目模式(tg_link:点击直接跳转tg链接,private:隐私模式，显示隐私信息)
	Order      string     `json:"order" mapstructure:"order"`             // 排序
	DescWeight float64    `json:"desc_weight" mapstructure:"desc_weight"` // 描述信息权重
	NumFilter  []string   `json:"num_filter" mapstructure:"num_filter"`   // 数组类型需要筛选                             // 数组类型筛选
	Commend    BotCommand `json:"commend" mapstructure:"command"`         // 命令配置
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

// 机器人命令
type BotCommand struct {
	Category       string `json:"category" mapstructure:"category"`
	CategoryHelp   string `json:"category_help" mapstructure:"category_help"`
	CategoryTag    string `json:"category_tag" mapstructure:"category_tag"`
	CategorySearch string `json:"category_search" mapstructure:"category_search"`
}
