package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

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

		bodyReader := io.LimitReader(r.Body, config.Reports.MaxBodySize)
		defer r.Body.Close()

		gzipReader, err := gzip.NewReader(bodyReader)
		if err != nil {
			slog.Warn("Gzip error", "remote", r.RemoteAddr, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()

		jsonReader := io.LimitReader(gzipReader, config.Reports.MaxJsonSize)

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
