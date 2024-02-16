package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Ports struct {
	Reports string
	Metrics string
}

type Limits struct {
	MaxBodySize int64
	MaxJsonSize int64
}

type Path struct {
	Reports string
	Metrics string
}

type ReportConfig struct {
	Save     bool
	SavePath string
	Limits   Limits
	Path     Path
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

func createRegistry(collectGoStats bool) *prometheus.Registry {
	registry := prometheus.NewRegistry()

	/*
		counter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "docker_events",
			Help: "Number of docker container events",
		}, []string{"type", "action", "scope", "from", "name", "namespace", "service_name", "node_id", "service_id"})
	*/

	if collectGoStats {
		registry.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}

	return registry
}

func handleReport(config ReportConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		limitedBody := io.LimitReader(r.Body, config.Limits.MaxBodySize)
		defer r.Body.Close()

		gzipReader, err := gzip.NewReader(limitedBody)
		if err != nil {
			slog.Warn("Gzip error", "remote", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		limitedJson := io.LimitReader(gzipReader, config.Limits.MaxJsonSize)
		defer gzipReader.Close()

		if config.Save {
			err := os.MkdirAll(config.SavePath, os.ModePerm)
			if err != nil {
				slog.Error("Failed to create directory", "path", config.SavePath, "error", err)
				return
			}

			filename := time.Now().Format(time.RFC3339Nano) + ".json"
			target := filepath.Join(config.SavePath, filename)
			slog.Info("Saving report", "target", target)

			out, err := os.Create(target)
			if err != nil {
				slog.Error("Failed to create file", "target", target, "error", err)
				return
			}
			defer out.Close()

			/// NOTE: It seems `TeeReader` writes complete json even when parsing fails later.
			limitedJson = io.TeeReader(limitedJson, out)
		}

		report, err := parseReport(limitedJson)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				slog.Warn("Request too large", "remote", r.RemoteAddr)
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}

			slog.Warn("Report error", "remote", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("DONE", "report", report)
	}
}

func main() {
	ports := Ports{
		Reports: getEnv("REPORTS_PORT", "8080"),
		Metrics: getEnv("METRICS_PORT", "8081"),
	}

	reportConfig := ReportConfig{
		Save:     getEnvBool("SAVE_REPORTS", true),
		SavePath: getEnv("SAVE_REPORTS_PATH", "/tmp/reports"),
		Path: Path{
			Reports: getEnv("REPORTS_PATH", "/"),
			Metrics: getEnv("METRICS_PATH", "/metrics"),
		},
		Limits: Limits{
			MaxBodySize: getEnvInt64("MAX_BODY_SIZE", 1*1024*1024),
			MaxJsonSize: getEnvInt64("MAX_JSON_SIZE", 5*1024*1024),
		},
	}

	collectGoStats := getEnvBool("COLLECT_GO_STATS", false)
	registry := createRegistry(collectGoStats)
	metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})

	metricsHttp := http.NewServeMux()
	metricsHttp.Handle(reportConfig.Path.Metrics, metricsHandler)

	go func() {
		slog.Info("Serving metrics", "port", ports.Metrics, "path", reportConfig.Path.Metrics)
		log.Fatal(http.ListenAndServe(":"+ports.Metrics, metricsHttp))
	}()

	http.HandleFunc(reportConfig.Path.Reports, handleReport(reportConfig))
	slog.Info("Listening for reports", "port", ports.Reports, "path", reportConfig.Path.Reports)
	log.Fatal(http.ListenAndServe(":"+ports.Reports, nil))
}
