package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//app config
type Config struct {
	AppName       string         `json:"app_name"`
	AppModel      string         `json:"app_model"`
	AppHost       string         `json:"app_host"`
	AppPort       int            `json:"app_port"`
	Database      DatabaseConfig `json:"database"`
	RedisConfig   RedisConfig    `json:"redis_config"`
	WSConfig      WSConfig       `json:"ws_config"`
	K8sConfig     K8sConfig      `json:"k8s"`
	AuthingConfig AuthingConfig  `json:"authing"`
	JwtConfig     JwtConfig      `json:"jwt"`
	Statistic     Statistic      `json:"statistic"`
}

type K8sConfig struct {
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
	FfileType string `json:"ffileType"`
}

//sql config
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	DBUser   string `json:"db_user"`
	Password string `json:"password"`
	DBHost   string `json:"db_host"`
	DBPort   string `json:"db_port"`
	DbName   string `json:"db_name"`
	Chartset string `json:"charset"`
	ShowSql  bool   `json:"show_sql"`
}

//Redis config
type RedisConfig struct {
	Addr     string `json:"addr"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

//websocket config
type WSConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	CheckOrigin bool   `json:"check_origin"`
}

//Authing Config
type AuthingConfig struct {
	UserPoolID  string `json:"userPoolID"`
	Secret      string `json:"secret"`
	AppID       string `json:"appID"`
	AppSecret   string `json:"appSecret"`
	RedirectURI string `json:"redirect_uri"`
}

//Jwt Jwt
type JwtConfig struct {
	Expire int    `json:"expire"`
	JwtKey string `json:"jwtKey"`
}

//Statistic function
type Statistic struct {
	Dir           string `json:"dir"`
	LogFile       string `json:"log_file"`
	LogFileSize   int64  `json:"log_file_size"`
	LogFileSuffix string `json:"log_file_suffix"`
}

//MessageQueue
type MessageQueue struct {
	KafkaServer string `json:"kafka_server"`
}

func InitConfig(path string) error {
	//app.json must be set right folder

	if dir, err := os.Getwd(); err == nil {
		if path == "" {
			dir = dir + "/conf/app.json"
		} else {
			dir = path
		}

		err = parseConfig(dir)
		if err != nil {
			return err
		}
	}

	if os.Getenv("GIN_MODE") != "" {
		cfg.AppModel = os.Getenv("GIN_MODE")
	}

	if os.Getenv("APP_PORT") != "" {
		cfg.AppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	}

	if os.Getenv("DB_USER") != "" {
		cfg.Database.DBUser = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PSWD") != "" {
		cfg.Database.Password = os.Getenv("DB_PSWD")
	}
	if os.Getenv("DB_HOST") != "" {
		cfg.Database.DBHost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_NAME") != "" {
		cfg.Database.DbName = os.Getenv("DB_NAME")
	}
	if os.Getenv("REDIS_ADDR") != "" {
		cfg.RedisConfig.Addr = os.Getenv("REDIS_ADDR")
	}
	if os.Getenv("REDIS_DB") != "" {
		cfg.RedisConfig.Db, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	}
	if os.Getenv("REDIS_PSWD") != "" {
		cfg.RedisConfig.Password = os.Getenv("REDIS_PSWD")
	}
	if os.Getenv("WS_HOST") != "" {
		cfg.WSConfig.Host = os.Getenv("WS_HOST")
	}
	if os.Getenv("WS_PORT") != "" {
		cfg.WSConfig.Port, _ = strconv.Atoi(os.Getenv("WS_PORT"))
	}

	if os.Getenv("AUTHING_APP_ID") != "" {
		cfg.AuthingConfig.AppID = os.Getenv("AUTHING_APP_ID")
	}
	if os.Getenv("AUTHING_APP_SECRET") != "" {
		cfg.AuthingConfig.AppSecret = os.Getenv("AUTHING_APP_SECRET")
	}
	if os.Getenv("AUTHING_SECRET") != "" {
		cfg.AuthingConfig.Secret = os.Getenv("AUTHING_SECRET")
	}
	if os.Getenv("AUTHING_USER_POOL_ID") != "" {
		cfg.AuthingConfig.UserPoolID = os.Getenv("AUTHING_USER_POOL_ID")
	}
	if os.Getenv("JWT_KEY") != "" {
		cfg.JwtConfig.JwtKey = os.Getenv("JWT_KEY")
	}

	if GetConfig().AppModel == "dev" || GetConfig().AppModel == "debug" {
		Log.SetLevel(logrus.DebugLevel)
		GetConfig().AppModel = gin.DebugMode
	} else {
		Log.SetLevel(logrus.InfoLevel)
		GetConfig().AppModel = gin.ReleaseMode
	}
	gin.SetMode(GetConfig().AppModel)
	return nil
}

//external
func GetConfig() *Config {
	return cfg
}

//internal
var cfg *Config = nil

func parseConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("read config file failed, please check path .  app exit now .")
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&cfg); err != nil {
		err = fmt.Errorf("load app.json failed, app must exit and error:%s", err)

		return err
	}
	return nil
}
