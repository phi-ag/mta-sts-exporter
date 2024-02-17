package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

type Config struct {
	Log struct {
		Json bool
	}
	Reports Reports
	Metrics Metrics
}

type DateRange struct {
	StartDatetime time.Time `json:"start-datetime"`
	EndDatetime   time.Time `json:"end-datetime"`
}

type Policy struct {
	PolicyType   string   `json:"policy-type"`
	PolicyString []string `json:"policy-string"`
	PolicyDomain string   `json:"policy-domain"`
	MxHost       string   `json:"mx-host"`
}

type Summary struct {
	TotalSuccessfulSessionCount int64 `json:"total-successful-session-count"`
	TotalFailureSessionCount    int64 `json:"total-failure-session-count"`
}

type FailureDetail struct {
	ResultType            string `json:"result-type"`
	SendingMtaIp          string `json:"sending-mta-ip"`
	ReceivingIp           string `json:"receiving-ip"`
	ReceivingMxHostname   string `json:"receiving-mx-hostname"`
	FailedSessionCount    int64  `json:"failed-session-count"`
	FailureReasonCode     string `json:"failure-reason-code"`
	AdditionalInformation string `json:"additional-information"`
}

type PolicyItem struct {
	Policy         Policy          `json:"policy"`
	Summary        Summary         `json:"summary"`
	FailureDetails []FailureDetail `json:"failure-details"`
}

type Report struct {
	OrganizationName string       `json:"organization-name"`
	DateRange        DateRange    `json:"date-range"`
	ContactInfo      string       `json:"contact-info"`
	ReportId         string       `json:"report-id"`
	Policies         []PolicyItem `json:"policies"`
}

func parseReport(reader io.Reader) (Report, error) {
	report := &Report{}
	err := json.NewDecoder(reader).Decode(report)
	return *report, err
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return result
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return result
	}
	return fallback
}

func createRegistry(config Config) *prometheus.Registry {
	registry := prometheus.NewRegistry()

	/*
		counter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "docker_events",
			Help: "Number of docker container events",
		}, []string{"type", "action", "scope", "from", "name", "namespace", "service_name", "node_id", "service_id"})
	*/

	if config.Metrics.Go {
		registry.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}

	return registry
}

func save(config Config, reader io.Reader) (io.Reader, *os.File, error) {
	err := os.MkdirAll(config.Reports.SavePath, os.ModePerm)
	if err != nil {
		slog.Error("Failed to create directory", "path", config.Reports.SavePath, "error", err)
		return reader, nil, err
	}

	filename := time.Now().Format(time.RFC3339Nano) + ".json"
	target := filepath.Join(config.Reports.SavePath, filename)
	slog.Info("Saving report", "target", target)

	file, err := os.Create(target)
	if err != nil {
		slog.Error("Failed to create file", "target", target, "error", err)
		return reader, nil, err
	}

	/// NOTE: It seems `TeeReader` writes the complete stream even when parsing fails later.
	/// This is probably only true for small payloads.
	return io.TeeReader(reader, file), file, nil
}

func handleReport(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		limitedBody := io.LimitReader(r.Body, config.Reports.MaxBodySize)
		defer r.Body.Close()

		gzipReader, err := gzip.NewReader(limitedBody)
		if err != nil {
			slog.Warn("Gzip error", "remote", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()

		limitedJson := io.LimitReader(gzipReader, config.Reports.MaxJsonSize)

		if config.Reports.Save {
			saveReader, file, err := save(config, limitedJson)
			if err != nil {
				slog.Warn("Save failed")
			} else {
				defer file.Close()
				limitedJson = saveReader
			}
		}

		report, err := parseReport(limitedJson)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				slog.Warn("Request too large", "remote", r.RemoteAddr)
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}

			slog.Warn("Report error", "remote", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		slog.Info("DONE", "report", report)
	}
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
	viper.SetDefault("Reports.Path", "/")
	viper.SetDefault("Reports.MaxBodySize", 1*1024*1024)
	viper.SetDefault("Reports.MaxJsonSize", 5*1024*1024)
	viper.SetDefault("Reports.Save", true)
	viper.SetDefault("Reports.SavePath", "/tmp/reports")
	viper.SetDefault("Metrics.Port", 8081)
	viper.SetDefault("Metrics.Path", "/metrics")
	viper.SetDefault("Metrics.Go", false)

	if _, err := os.Stat(filepath.Join(configPath, configName)); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Failed to read config file", "error", err)
		}
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		slog.Warn("Failed to unmarshal config", "error", err)
	}

	return config
}

func main() {
	config := createConfig()

	if config.Log.Json {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	}

	slog.Info("Config", "config", config)

	registry := createRegistry(config)
	metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})

	metricsHttp := http.NewServeMux()
	metricsHttp.Handle(config.Metrics.Path, metricsHandler)

	go func() {
		slog.Info("Serving metrics", "port", config.Metrics.Port, "path", config.Metrics.Path)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Metrics.Port), metricsHttp))
	}()

	http.HandleFunc(config.Reports.Path, handleReport(config))
	slog.Info("Listening for reports", "port", config.Reports.Port, "path", config.Reports.Path)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Reports.Port), nil))
}
