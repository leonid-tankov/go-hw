package config

import (
	"os"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf  `yaml:"logger"`
	HTTP     HTTPConfig  `yaml:"http"`
	Storage  StorageConf `yaml:"storage,omitempty"`
	Postgres PostgresConf
}

type StorageConf struct {
	Type string `yaml:"type"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type PostgresConf struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	if err = viper.ReadConfig(file); err != nil {
		panic(err)
	}
	conf := &Config{}
	if err = viper.Unmarshal(conf); err != nil {
		panic(err)
	}
	viper.SetEnvPrefix("postgres")
	viper.AutomaticEnv()
	postgresConf := &PostgresConf{
		Host:     viper.GetString("host"),
		Port:     viper.GetString("port"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
		Database: viper.GetString("database"),
	}
	conf.Postgres = *postgresConf

	return *conf
}
