# MTA-STS Exporter

[![GitHub release](https://img.shields.io/github/release/phi-ag/mta-sts-exporter.svg?logo=github&style=flat-square)](https://github.com/phi-ag/mta-sts-exporter/releases/latest)
[![Docker pulls](https://img.shields.io/docker/pulls/phiag/mta-sts-exporter.svg?logo=docker&style=flat-square)](https://hub.docker.com/r/phiag/mta-sts-exporter/tags)

**WIP** Prometheus metrics not implemented yet!

Configuration

- `CONFIG_PATH` (default: /etc/mta-sts-exporter/config.yaml)
- `PORT` (default: 8080)
- `LOG_JSON` (default: true)
- `POLICY_ENABLED` (default: true)
- `POLICY_PATH` (default: /.well-known/mta-sts.txt)
- `POLICY_VERSION` (default: STSv1)
- `POLICY_MODE` (default: enforce)
- `POLICY_MX` (default: example.com)
- `POLICY_MAXAGE` (default: 86400)
- `REPORTS_PATH` (default: /report)
- `REPORTS_MAXBODYSIZE` (default: 1 MiB)
- `REPORTS_MAXJSONSIZE` (default: 5 MiB)
- `REPORTS_SAVE` (default: true)
- `REPORTS_SAVEPATH` (default: /tmp/reports)
- `METRICS_PORT` (default: 8081)
- `METRICS_PATH` (default: /metrics)
- `METRICS_GO` (default: false)

Save reports

    mkdir reports
    chown 65532:65532 reports
    docker run -it --rm -p 8080:8080 -p 8081:8081 -v ${PWD}/reports:/tmp/reports phiag/mta-sts-exporter:latest

Post examples

    docker run -it --rm -p 8080:8080 -p 8081:8081 phiag/mta-sts-exporter:latest

    cat examples/rfc.json | gzip | curl -X POST -v --data-binary @- localhost:8080/report
    cat examples/google.json | gzip | curl -X POST -v --data-binary @- localhost:8080/report
    cat examples/microsoft.json | gzip | curl -X POST -v --data-binary @- localhost:8080/report

    curl localhost:8081/metrics
    curl http://localhost:8080/.well-known/mta-sts.txt

## References

- [RFC 8460: SMTP TLS Reporting](https://www.rfc-editor.org/rfc/rfc8460.html)
- [RFC 8461: SMTP MTA Strict Transport Security (MTA-STS)](https://www.rfc-editor.org/rfc/rfc8461.html)
