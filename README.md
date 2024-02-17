# MTA-STS Exporter

[![GitHub release](https://img.shields.io/github/release/phi-ag/mta-sts-exporter.svg?logo=github&style=flat-square)](https://github.com/phi-ag/mta-sts-exporter/releases/latest)
[![Docker pulls](https://img.shields.io/docker/pulls/phiag/mta-sts-exporter.svg?logo=docker&style=flat-square)](https://hub.docker.com/r/phiag/mta-sts-exporter)

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

Save reports

    mkdir reports
    chown 65532:65532 reports
    docker run -it --rm -v ${PWD}/reports:/tmp/reports phiag/mta-sts-exporter:latest

Post example

    docker run -it --rm phiag/mta-sts-exporter:latest
    cat test/example.json | gzip | curl -X POST -v --data-binary @- localhost:8080
    curl localhost:8081/metrics
