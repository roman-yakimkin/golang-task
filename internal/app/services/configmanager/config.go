package configmanager

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	BindAddr     string `yaml:"bind_addr"`
	GRPCBindAddr string `yaml:"grpc_bind_addr"`
	SecretString string `yaml:"secret_string"`
	JsonPathTask string `yaml:"json_path_task"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}

func (cm *Config) Init(configPath string) error {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, cm)
	return err
}
