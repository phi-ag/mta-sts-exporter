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

func handleReport(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		bodyReader := io.LimitReader(r.Body, config.Reports.Max.Body)
		defer r.Body.Close()

		gzipReader, err := gzip.NewReader(bodyReader)
		if err != nil {
			slog.Warn("Gzip error", "remote", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()

		jsonReader := io.LimitReader(gzipReader, config.Reports.Max.Json)

		if config.Reports.Save {
			saveReader, file, err := saveReport(config, jsonReader)
			if err != nil {
				slog.Warn("Save failed")
			} else {
				defer file.Close()
				jsonReader = saveReader
			}
		}

		report, err := parseReport(jsonReader)
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

func handlePolicy(config Config) http.HandlerFunc {
	response := policyResponse(config.Policy)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprint(w, response)
	}
}

func healthCheck(config Config) int {
	res, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", config.Port))
	if err != nil || res.StatusCode != http.StatusOK {
		slog.Error("Healthcheck failed", "error", err, "statusCode", res.StatusCode)
		return 1
	}
	return 0
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

		registry := createRegistry(config)
		metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})

		metricsHttp := http.NewServeMux()
		metricsHttp.Handle(config.Metrics.Path, metricsHandler)

		go func() {
			slog.Info("Serving metrics", "port", config.Metrics.Port, "path", config.Metrics.Path)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Metrics.Port), metricsHttp))
		}()

		http.HandleFunc(config.Reports.Path, handleReport(config))
		http.HandleFunc("/healthz", func(http.ResponseWriter, *http.Request) {})

		if config.Policy.Enabled {
			http.HandleFunc(config.Policy.Path, handlePolicy(config))
		}

		slog.Info("Listening for reports", "port", config.Port, "path", config.Reports.Path)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
	}
}
