package config

import (
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

type Configs struct {
	Rest    RestConfigs
	Store   StoreConfigs
	Indexer IndexerConfigs
}

type RestConfigs struct {
	Port string
}

type StoreConfigs struct {
	Path string
}

type IndexerConfigs struct {
	Address  string
	Host     string
	Token    string
	Endpoint string
	Limiter  struct {
		Rate     int
		Duration time.Duration
	}
	Workers  int
	Timeout  time.Duration
	ChanSize int `mapstructure:"chan_size"`
}

func Init() *Configs {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	viper.SetConfigName("config")
	viper.AddConfigPath(basePath)
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	var configs Configs
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading configs: %s", err)
	}

	err := viper.Unmarshal(&configs)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &configs
}
