package config

import (
	"github.com/spf13/viper"
)

// Config Service configs

type Config struct {
	Grpc            Grpc            `json:"grpc" yaml:"grpc"`
	Http            Http            `json:"http" yaml:"http"`
	ServiceEndpoint ServiceEndpoint `json:"serviceEndpoint" yaml:"serviceEndpoint"`
	Database        Database        `json:"database" yaml:"database"`
	QtumDataSource  QtumDataSource  `json:"net_type" yaml:"qtumDataSourceZ"`
}

// NewConfig Initial app's configs
func NewConfig(cfg string) *Config {

	if cfg == "" {
		panic("load configs file failed.configs file can not be empty.")
	}

	viper.SetConfigFile(cfg)

	// Read configs file
	if err := viper.ReadInConfig(); err != nil {
		panic("read configs failed.[ERROR]=>" + err.Error())
	}
	conf := &Config{}
	// Assign the overloaded configuration to the global
	if err := viper.Unmarshal(conf); err != nil {
		panic("assign configs failed.[ERROR]=>" + err.Error())
	}

	return conf

}
