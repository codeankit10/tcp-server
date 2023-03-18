package config

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerHost            string `envconfig:"SERVER_HOST"`
	ServerPort            int    `envconfig:"Server_Port"`
	CacheHost             string `envconfig:"Cache_Host"`
	CachePort             int    `envconfig:"Cache_Port"`
	HashcashZerosCount    int
	HashcashDuration      int64
	HashcashMaxIterations int
}

func Load(path string) (*Config, error) {
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err

	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return &config, err

	}
	err = envconfig.Process("", &config)
	return &config, err

}
