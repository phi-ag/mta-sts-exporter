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
	PolicyRequestsTotal     prometheus.Counter
	ReportRequestsTotal     prometheus.Counter
	ReportErrorsTotal       *prometheus.CounterVec
	SuccessfulSessionsTotal *prometheus.CounterVec
	FailureSessionsTotal    *prometheus.CounterVec
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
		ReportErrorsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "report_errors_total",
			Help:      "Total number of report errors.",
		}, []string{"cause"}),
		SuccessfulSessionsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "successful_sessions_total",
			Help:      "Total number of successful sessions.",
		}, []string{"organization"}),
		FailureSessionsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "failure_sessions_total",
			Help:      "Total number of failure sessions.",
		}, []string{"organization"}),
	}
}

func createRegistry(config Config, metrics Metrics) *prometheus.Registry {
	registry := prometheus.NewRegistry()

	registry.MustRegister(
		metrics.PolicyRequestsTotal,
		metrics.ReportRequestsTotal,
		metrics.ReportErrorsTotal,
		metrics.SuccessfulSessionsTotal,
		metrics.FailureSessionsTotal)

	if config.Metrics.Collectors.Go {
		registry.MustRegister(
			collectors.NewGoCollector(),
		)
	}

	if config.Metrics.Collectors.Process {
		registry.MustRegister(
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
		_, err := fmt.Fprint(w, response)
		if err != nil {
			slog.Warn("Policy error", "remote", r.RemoteAddr, "error", err)
		}
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
		defer func() {
			err := r.Body.Close()
			if err != nil {
				slog.Warn("Failed to close report body", "remote", r.RemoteAddr, "error", err)
			}
		}()

		gzipReader, err := gzip.NewReader(bodyReader)
		if err != nil {
			slog.Warn("Gzip error", "remote", r.RemoteAddr, "error", err)
			metrics.ReportErrorsTotal.With(prometheus.Labels{"cause": "gzip"}).Inc()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			err := gzipReader.Close()
			if err != nil {
				slog.Warn("Failed to close report gzip reader", "remote", r.RemoteAddr, "error", err)
			}
		}()

		jsonReader := io.LimitReader(gzipReader, config.Reports.Max.Json)

		if config.Reports.Save.Enabled {
			saveReader, file, err := saveReport(config, jsonReader)
			if err != nil {
				slog.Warn("Save failed")
				metrics.ReportErrorsTotal.With(prometheus.Labels{"cause": "save"}).Inc()
			} else {
				defer func() {
					err := file.Close()
					if err != nil {
						slog.Warn("Failed to close report file", "remote", r.RemoteAddr, "error", err)
					}
				}()
				jsonReader = saveReader
			}
		}

		report, err := parseReport(jsonReader)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				slog.Warn("Request too large", "remote", r.RemoteAddr)
				metrics.ReportErrorsTotal.With(prometheus.Labels{"cause": "request_too_large"}).Inc()
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}

			slog.Warn("Report error", "remote", r.RemoteAddr, "error", err)
			metrics.ReportErrorsTotal.With(prometheus.Labels{"cause": "parse"}).Inc()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, policy := range report.Policies {
			if policy.Summary.TotalSuccessfulSessionCount > 0 {
				labels := prometheus.Labels{"organization": report.OrganizationName}
				value := float64(policy.Summary.TotalSuccessfulSessionCount)
				metrics.SuccessfulSessionsTotal.With(labels).Add(value)
			}

			if policy.Summary.TotalFailureSessionCount > 0 {
				labels := prometheus.Labels{"organization": report.OrganizationName}
				value := float64(policy.Summary.TotalFailureSessionCount)
				metrics.FailureSessionsTotal.With(labels).Add(value)
			}

			if len(policy.FailureDetails) > 0 {
				slog.Warn("Policy contains failure details", "details", policy.FailureDetails)
			}
		}
	}
}

func start(config Config) {
	slog.Info("Config", "config", config)

	metrics := createMetrics()
	registry := createRegistry(config, metrics)

	var handlerOpts promhttp.HandlerOpts
	if config.Metrics.Collectors.Exporter {
		handlerOpts = promhttp.HandlerOpts{Registry: registry}
	} else {
		handlerOpts = promhttp.HandlerOpts{}
	}

	metricsHandler := promhttp.HandlerFor(registry, handlerOpts)

	metricsHttp := http.NewServeMux()
	metricsHttp.Handle(config.Metrics.Path, metricsHandler)

	publicHttp := http.NewServeMux()
	publicHttp.HandleFunc("/healthz", func(http.ResponseWriter, *http.Request) {})

	if config.Metrics.Enabled {
		go func() {
			slog.Info("Serving metrics", "port", config.Metrics.Port, "path", config.Metrics.Path)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Metrics.Port), metricsHttp))
		}()
	}

	if config.Reports.Enabled {
		slog.Info("Listening for reports", "port", config.Port, "path", config.Reports.Path)
		publicHttp.HandleFunc(config.Reports.Path, handleReport(config, metrics))
	}

	if config.Policy.Enabled {
		slog.Info("Serving policy", "port", config.Port, "path", config.Policy.Path)
		publicHttp.HandleFunc(config.Policy.Path, handlePolicy(config, metrics))
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), publicHttp))
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
		start(config)
	}
}
