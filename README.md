# MTA-STS Exporter

**WIP** Prometheus metrics not implemented yet!

Configuration

- `REPORTS_PORT` (default: 8080)
- `METRICS_PORT` (default: 8081)
- `METRICS_PATH` (default: /metrics)
- `MAX_BODY_SIZE` (default: 1 MiB)
- `MAX_JSON_SIZE` (default: 5 MiB)
- `SAVE_REPORTS` (default: true)
- `SAVE_REPORTS_PATH` (default: /tmp/reports)
- `COLLECT_GO_STATS` (default: false)

Example

    cat test/example.json | gzip | curl -X POST -v --data-binary @- localhost:8080
    curl localhost:8081/metrics
