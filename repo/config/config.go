package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Db             DbConfig                 `yaml:"db"`
	QueryEndpoints []string                 `yaml:"query_endpoints"`
	Accounts       map[string]AccountConfig `yaml:"accounts"`
}

type AccountConfig struct {
	PrivateKey string `yaml:"private_key"`
	Name       string `yaml:"name"`
}

type DbConfig struct {
	Name string
	Pass string
}

func LoadConfig() (r *Config) {
	data, err := os.ReadFile("config/config.yml")
	if err != nil {
		panic(err)
		return
	}

	r = &Config{}
	err = yaml.Unmarshal(data, r)
	if err != nil {
		panic(err)
	}
	return
}
