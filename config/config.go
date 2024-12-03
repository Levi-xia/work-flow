package config

import (
	"fmt"
	"log"
	"os"
	"workflow/internal/common"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = &Configuration{}

type Config struct {
	viper *viper.Viper
}

type Configuration struct {
	Server        ServerConfig        `mapstructure:"server" json:"server" yaml:"server"`
	Jwt           JwtConfig           `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	DB            MySqlConfig         `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	MiniProgram   MiniProgramConfig   `mapstructure:"mini_program" json:"mini_program" yaml:"mini_program"`
	ElasticSearch ElasticSearchConfig `mapstructure:"elasticsearch" json:"elasticsearch" yaml:"elasticsearch"`
	Redis         RedisConfig         `mapstructure:"redis" json:"redis" yaml:"redis"`
	Logger        LoggerConfig        `mapstructure:"logger" json:"logger" yaml:"logger"`
	Qiniu         QiniuConfig         `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
}

// / 服务配置
type ServerConfig struct {
	Port int `mapstructure:"port" json:"port" yaml:"port"`
}

// / 数据库配置
type MySqlConfig struct {
	Driver              string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                int    `mapstructure:"port" json:"port" yaml:"port"`
	Database            string `mapstructure:"database" json:"database" yaml:"database"`
	Username            string `mapstructure:"username" json:"username" yaml:"username"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Charset             string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

// / jwt 配置
type JwtConfig struct {
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	JwtTtl int64  `mapstructure:"jwt_ttl" json:"jwt_ttl" yaml:"jwt_ttl"` // token 有效期（秒）
}

// / 小程序
type MiniProgramConfig struct {
	AppId     string `mapstructure:"app_id" json:"app_id" yaml:"app_id"`
	AppSecret string `mapstructure:"app_secret" json:"app_secret" yaml:"app_secret"`
}

// ES
type ElasticSearchConfig struct {
	Addresses []string `mapstructure:"addresses" json:"addresses" yaml:"addresses"`
}

// Redis
type RedisConfig struct {
	Host             string `mapstructure:"host" json:"host" yaml:"host"`
	Port             string `mapstructure:"port" json:"port" yaml:"port"`
	MaxIdle          int    `mapstructure:"max_idle" json:"max_idle" yaml:"max_idle"`
	MaxActive        int    `mapstructure:"max_active" json:"max_active" yaml:"max_active"`
	IdleTimeout      int    `mapstructure:"idle_timeout" json:"idle_timeout" yaml:"idle_timeout"`
	Password         string `mapstructure:"password" json:"password" yaml:"password"`
	DialReadTimeout  int    `mapstructure:"dial_read_timeout" json:"dial_read_timeout" yaml:"dial_read_timeout"`
	DialWriteTimeout int    `mapstructure:"dial_write_timeout" json:"dial_write_timeout" yaml:"dial_write_timeout"`
}

// Logger
type LoggerConfig struct {
	DebugFileName string `mapstructure:"debugFileName" json:"debugFileName" yaml:"debugFileName"`
	InfoFileName  string `mapstructure:"infoFileName" json:"infoFileName" yaml:"infoFileName"`
	WarnFileName  string `mapstructure:"warnFileName" json:"warnFileName" yaml:"warnFileName"`
	ErrorFileName string `mapstructure:"errorFileName" json:"errorFileName" yaml:"errorFileName"`
	MaxSize       int    `mapstructure:"maxSize" json:"maxSize" yaml:"maxSize"`
	MaxAge        int    `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`
	MaxBackups    int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
}

// Qiniu
type QiniuConfig struct {
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	AccessKey string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Expire    uint64 `mapstructure:"expire" json:"expire" yaml:"expire"`
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
}

func InitConfig() *Config {
	config := &Config{
		viper: viper.New(),
	}
	workDir, _ := os.Getwd()
	config.viper.SetConfigName(fmt.Sprintf("app-%s", common.ENV))
	config.viper.AddConfigPath(workDir + "/config")
	config.viper.SetConfigType("yaml")

	if err := config.viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %v\n", err)
	}
	config.viper.WatchConfig()
	config.viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("config file changed:", in.Name)
		if err := config.viper.Unmarshal(Conf); err != nil {
			log.Println("Unmarshal config failed, err:", err)
		}
	})
	if err := config.viper.Unmarshal(Conf); err != nil {
		log.Fatalf("Unmarshal config failed, err:%v\n", err)
	}
	return config
}
