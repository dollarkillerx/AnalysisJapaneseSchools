package conf

import cfg "github.com/dollarkillerx/common/pkg/config"

var conf *configuration

type configuration struct {
	PostgresConfig cfg.PostgresConfiguration
	ListenAddr     string
}

func InitConfig(configName string, configPaths []string) error {
	return cfg.InitConfiguration(configName, configPaths, &conf)
}

func GetConfig() *configuration {
	return conf
}
