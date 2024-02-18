package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	PolicyRequestsTotal prometheus.Counter
	ReportRequestsTotal prometheus.Counter
	ReportRequestsValid prometheus.Counter
	ReportGzipError     prometheus.Counter
	ReportSaveError     prometheus.Counter
	ReportTooLarge      prometheus.Counter
	ReportError         prometheus.Counter
}

func createMetrics() Metrics {
	namespace := "mta_sts"

	return Metrics{
		PolicyRequestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "policy_requests_total",
			Help:      "Total number of policy requests.",
		}),
		ReportRequestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "report_requests_total",
			Help:      "Total number of report requests.",
		}),
		ReportRequestsValid: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "report_requests_valid",
			Help:      "Total number of valid report requests.",
		}),
		ReportGzipError: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "report_gzip_error",
			Help:      "Total number of report gzip errors.",
		}),
		ReportSaveError: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "report_save_error",
			Help:      "Total number of report save errors.",
		}),
		ReportTooLarge: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "report_too_large",
			Help:      "Total number of too large report requests.",
		}),
		ReportError: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "report_error",
			Help:      "Total number of report errors.",
		}),
	}
}

func createRegistry(config Config, metrics Metrics) *prometheus.Registry {
	registry := prometheus.NewRegistry()

	registry.MustRegister(
		metrics.PolicyRequestsTotal,
		metrics.ReportRequestsTotal,
		metrics.ReportRequestsValid,
		metrics.ReportGzipError,
		metrics.ReportSaveError,
		metrics.ReportTooLarge,
		metrics.ReportError)

	if config.Metrics.Go {
		registry.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}

	return registry
}

func healthCheck(config Config) int {
	res, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", config.Port))
	if err != nil || res.StatusCode != http.StatusOK {
		slog.Error("Healthcheck failed", "error", err, "statusCode", res.StatusCode)
		return 1
	}
	return 0
}

func policyResponse(policy Policy) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("version: %s\n", policy.Version))
	sb.WriteString(fmt.Sprintf("mode: %s\n", policy.Mode))

	for _, mx := range policy.Mx {
		sb.WriteString(fmt.Sprintf("mx: %s\n", mx))
	}

	sb.WriteString(fmt.Sprintf("max_age: %d\n", policy.MaxAge))
	return sb.String()
}

func handlePolicy(config Config, metrics Metrics) http.HandlerFunc {
	response := policyResponse(config.Policy)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		metrics.PolicyRequestsTotal.Inc()
		fmt.Fprint(w, response)
	}
}

func handleReport(config Config, metrics Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		metrics.ReportRequestsTotal.Inc()

		bodyReader := io.LimitReader(r.Body, config.Reports.Max.Body)
		defer r.Body.Close()

		gzipReader, err := gzip.NewReader(bodyReader)
		if err != nil {
			slog.Warn("Gzip error", "remote", r.RemoteAddr, "error", err)
			metrics.ReportGzipError.Inc()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()

		jsonReader := io.LimitReader(gzipReader, config.Reports.Max.Json)

		if config.Reports.Save.Enabled {
			saveReader, file, err := saveReport(config, jsonReader)
			if err != nil {
				slog.Warn("Save failed")
				metrics.ReportSaveError.Inc()
			} else {
				defer file.Close()
				jsonReader = saveReader
			}
		}

		report, err := parseReport(jsonReader)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				slog.Warn("Request too large", "remote", r.RemoteAddr)
				metrics.ReportTooLarge.Inc()
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}

			slog.Warn("Report error", "remote", r.RemoteAddr, "error", err)
			metrics.ReportError.Inc()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metrics.ReportRequestsValid.Inc()
		slog.Info("DONE", "report", report)
	}
}

func main() {
	config := createConfig()

	if config.Log.Json {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	}

	healthCheckFlag := flag.Bool("health", false, "run health check")
	flag.Parse()

	if *healthCheckFlag {
		os.Exit(healthCheck(config))
	} else {
		slog.Info("Config", "config", config)

		metrics := createMetrics()
		registry := createRegistry(config, metrics)
		metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})

		metricsHttp := http.NewServeMux()
		metricsHttp.Handle(config.Metrics.Path, metricsHandler)

		http.HandleFunc("/healthz", func(http.ResponseWriter, *http.Request) {})

		if config.Metrics.Enabled {
			go func() {
				slog.Info("Serving metrics", "port", config.Metrics.Port, "path", config.Metrics.Path)
				log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Metrics.Port), metricsHttp))
			}()
		}

		if config.Reports.Enabled {
			slog.Info("Listening for reports", "port", config.Port, "path", config.Reports.Path)
			http.HandleFunc(config.Reports.Path, handleReport(config, metrics))
		}

		if config.Policy.Enabled {
			slog.Info("Serving policy", "port", config.Port, "path", config.Policy.Path)
			http.HandleFunc(config.Policy.Path, handlePolicy(config, metrics))
		}

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
	}
}
