# MTA-STS Exporter

**WIP** Prometheus metrics not implemented yet!

Configuration

- `CONFIG_PATH` (default: /etc/mta-sts-exporter/config.yaml)
- `LOG_JSON` (default: true)
- `REPORTS_PORT` (default: 8080)
- `METRICS_PORT` (default: 8081)
- `REPORTS_PATH` (default: /)
- `METRICS_PATH` (default: /metrics)
- `REPORTS_MAXBODYSIZE` (default: 1 MiB)
- `REPORTS_MAXJSONSIZE` (default: 5 MiB)
- `REPORTS_SAVE` (default: true)
- `REPORTS_SAVEPATH` (default: /tmp/reports)
- `METRICS_GO` (default: false)

Example

    docker run phiag/mta-sts-exporter
    cat test/example.json | gzip | curl -X POST -v --data-binary @- localhost:8080
    curl localhost:8081/metrics
