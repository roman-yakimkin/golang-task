package configmanager

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	BindAddr       string `yaml:"bind_addr"`
	SecretString   string `yaml:"secret_string"`
	JsonPathTask   string `yaml:"json_path_task"`
	MongoDBConnStr string `yaml:"mongo_db_connection_string"`
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
