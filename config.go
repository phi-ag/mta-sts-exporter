package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Policy struct {
	Enabled bool
	Path    string
	Version string
	Mode    string
	Mx      []string
	MaxAge  int64
}

type ReportsMax struct {
	Body int64
	Json int64
}

type ReportsSave struct {
	Enabled bool
	Path    string
}

type Reports struct {
	Enabled bool
	Path    string
	Max     ReportsMax
	Save    ReportsSave
}

type ConfigMetrics struct {
	Enabled    bool
	Port       uint16
	Path       string
	Collectors struct {
		Go       bool
		Process  bool
		Exporter bool
	}
}

type Config struct {
	Port uint16
	Log  struct {
		Json bool
	}
	Policy  Policy
	Reports Reports
	Metrics ConfigMetrics
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
	viper.SetDefault("Port", 8080)
	viper.SetDefault("Policy.Enabled", true)
	viper.SetDefault("Policy.Path", "/.well-known/mta-sts.txt")
	viper.SetDefault("Policy.Version", "STSv1")
	viper.SetDefault("Policy.Mode", "enforce")
	viper.SetDefault("Policy.Mx", []string{"mx1.example.com", "mx2.example.com"})
	viper.SetDefault("Policy.MaxAge", "86400")
	viper.SetDefault("Reports.Enabled", true)
	viper.SetDefault("Reports.Path", "/report")
	viper.SetDefault("Reports.Max.Body", 1*1024*1024)
	viper.SetDefault("Reports.Max.Json", 5*1024*1024)
	viper.SetDefault("Reports.Save.Enabled", false)
	viper.SetDefault("Reports.Save.Path", "/tmp/reports")
	viper.SetDefault("Metrics.Enabled", true)
	viper.SetDefault("Metrics.Port", 8081)
	viper.SetDefault("Metrics.Path", "/metrics")
	viper.SetDefault("Metrics.Collectors.Go", false)
	viper.SetDefault("Metrics.Collectors.Process", false)
	viper.SetDefault("Metrics.Collectors.Exporter", false)

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
