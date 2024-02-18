package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Reports struct {
	Port        uint16
	Path        string
	MaxBodySize int64
	MaxJsonSize int64
	Save        bool
	SavePath    string
}

type Metrics struct {
	Port uint16
	Path string
	Go   bool
}

type Policy struct {
	Enabled bool
	Content string
}

type Config struct {
	Log struct {
		Json bool
	}
	Reports Reports
	Metrics Metrics
	Policy  Policy
}

func createConfig() Config {
	configPathFull := getEnv("CONFIG_PATH", "/etc/mta-sts-exporter/config.yaml")
	configPath := filepath.Dir(configPathFull)
	configName := filepath.Base(configPathFull)

	if filepath.Ext(configName) == "" {
		viper.SetConfigType("yaml")
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("Log.Json", true)
	viper.SetDefault("Reports.Port", 8080)
	viper.SetDefault("Reports.Path", "/report")
	viper.SetDefault("Reports.MaxBodySize", 1*1024*1024)
	viper.SetDefault("Reports.MaxJsonSize", 5*1024*1024)
	viper.SetDefault("Reports.Save", true)
	viper.SetDefault("Reports.SavePath", "/tmp/reports")
	viper.SetDefault("Metrics.Port", 8081)
	viper.SetDefault("Metrics.Path", "/metrics")
	viper.SetDefault("Metrics.Go", false)
	viper.SetDefault("Policy.Enabled", true)
	viper.SetDefault("Policy.Content", "version: STSv1\nmode: testing\nmx: example.com\nmax_age: 86400\n")

	if _, err := os.Stat(configPathFull); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("Failed to read config file:", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln("Failed to unmarshal config:", err)
	}

	return config
}
