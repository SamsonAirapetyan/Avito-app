package config

import (
	logger2 "SergeyProject/pkg/logger"
	"github.com/spf13/viper"
	"os"
)

// Config для принятия данных об Бд и сервере
type Config struct {
	PostgresDB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
		SSLmode  string `yaml:"sslmode"`
		MaxConns string `yaml:"maxconns"`
	}
	Server struct {
		BindAddr     string `yaml:"bindAddr"`
		ReadTimeout  string `yaml:"readTimeout"`
		WriteTimeout string `yaml:"writeTimeout"`
		IdleTimeout  string `yaml:"idleTimeout"`
	}
}

// ConfigViper Считать конфигурацию с файла yaml
func ConfigViper() *viper.Viper {
	logger := logger2.GetLogger()
	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		logger.Error("Unable read config file", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Config loaded successfully")

	return v
}

// ParseConfig Переношу все данные с файла в структуру и передаю её
func ParseConfig(v *viper.Viper) *Config {
	logger := logger2.GetLogger()

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		logger.Error("Unable to parse the configuration file.")
	}
	logger.Info("Configuratin file parsed successfully.")

	return cfg
}
