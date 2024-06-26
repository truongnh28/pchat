package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type AppConfig struct {
	Server              ServerConfig          `mapstructure:"api"`
	Logger              LoggerConfig          `mapstructure:"logger"`
	Authentication      *AuthenticationConfig `mapstructure:"authentication"`
	ChatAppDatabase     *DatabaseConfig       `mapstructure:"chatAppDatabase"`
	ChatMessageDatabase *DatabaseConfig       `mapstructure:"chatMessageDatabase"`
	Env                 string                `mapstructure:"env"`
	IsProduction        bool                  `mapstructure:"isProduction"`
	Redis               *RedisConfig          `mapstructure:"redis"`
	Mail                *MailConfig           `mapstructure:"mail"`
	Cloudinary          *CloudinaryConfig     `mapstructure:"cloudinary"`
}

type ServerConfig struct {
	Port  string `mapstructure:"port"  default:"8080"`
	Debug bool   `mapstructure:"debug" default:"false"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type AuthenticationConfig struct {
	SecretKey    string `mapstructure:"secret"`
	ExpiredTime  int64  `mapstructure:"expiredTime"`
	CookieName   string `mapstructure:"cookieName"`
	CookiePath   string `mapstructure:"cookiePath"`
	CookieSecure bool   `mapstructure:"cookieSecure"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"databaseName"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

type MailConfig struct {
	MailSender string `mapstructure:"mailSender"`
	Password   string `mapstructure:"password"`
	SmtpHost   string `mapstructure:"smtpHost"`
	SmtpPort   int    `mapstructure:"smtpPort"`
}

type CloudinaryConfig struct {
	Name      string `mapstructure:"cloudName"`
	APIKey    string `mapstructure:"apiKey"`
	APISecret string `mapstructure:"apiSecret"`
	Folder    string `mapstructure:"folder"`
}

var (
	_, b, _, _        = runtime.Caller(0)
	basePath          = filepath.Dir(b) //get absolute directory of current file
	defaultConfigFile = basePath + "/local.yaml"
	v                 = viper.New()
	appConfig         AppConfig
)

func init() {
	Load()
}

func Load() {
	var configFile string
	if configFile = os.Getenv("CONFIG_PATH"); len(configFile) == 0 {
		configFile = defaultConfigFile
	}

	if err := loadConfigFile(configFile); err != nil {
		panic(err)
	}

	if err := scanConfigFile(&appConfig); err != nil {
		panic(err)
	}

}

func loadConfigFile(configFile string) error {
	configFileName := filepath.Base(configFile)
	configFilePath := filepath.Dir(configFile)

	v.AddConfigPath(configFilePath)
	v.SetConfigName(strings.TrimSuffix(configFileName, filepath.Ext(configFileName)))
	v.AutomaticEnv()

	return v.ReadInConfig()
}

func scanConfigFile(config any) error {
	return v.Unmarshal(&config)
}

func GetAppConfig() *AppConfig {
	return &appConfig
}
