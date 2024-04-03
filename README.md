# MTA-STS Exporter

[![GitHub Release](https://img.shields.io/github/v/release/phi-ag/mta-sts-exporter?style=for-the-badge)](https://github.com/phi-ag/mta-sts-exporter/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/phiag/mta-sts-exporter?style=for-the-badge)](https://hub.docker.com/r/phiag/mta-sts-exporter)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/phi-ag/mta-sts-exporter/check.yml?style=for-the-badge&label=Check)](https://github.com/phi-ag/mta-sts-exporter/actions)

## Configuration

Use environment variables or a configuration file (see [compose.yaml](compose.yaml))

- `CONFIG_PATH` (default: /etc/mta-sts-exporter/config.yaml)
- `PORT` (default: 8080)
- `LOG_JSON` (default: true)
- `POLICY_ENABLED` (default: true)
- `POLICY_PATH` (default: /.well-known/mta-sts.txt)
- `POLICY_VERSION` (default: STSv1)
- `POLICY_MODE` (default: enforce)
- `POLICY_MX` (default: example.com)
- `POLICY_MAXAGE` (default: 86400)
- `REPORTS_ENABLED` (default: true)
- `REPORTS_PATH` (default: /report)
- `REPORTS_MAX_BODY` (default: 1 MiB)
- `REPORTS_MAX_JSON` (default: 5 MiB)
- `REPORTS_SAVE_ENABLED` (default: false)
- `REPORTS_SAVE_PATH` (default: /tmp/reports)
- `METRICS_ENABLED` (default: true)
- `METRICS_PORT` (default: 8081)
- `METRICS_PATH` (default: /metrics)
- `METRICS_COLLECTORS_GO` (default: false)
- `METRICS_COLLECTORS_PROCESS` (default: false)
- `METRICS_COLLECTORS_EXPORTER` (default: false)

## DNS

    mta-sts    A    <IPv4>
    mta-sts    AAAA <IPv6>
    _mta-sts   TXT  "v=STSv1; id=20240101T010101;"
    _smtp._tls TXT  "v=TLSRPTv1;rua=https://mta-sts.example.com/report"

## Usage

Post examples

    docker run -it --rm -p 8080:8080 -p 8081:8081 phiag/mta-sts-exporter:latest

    cat examples/rfc.json | gzip | curl -X POST -v --data-binary @- localhost:8080/report
    cat examples/google.json | gzip | curl -X POST -v --data-binary @- localhost:8080/report
    cat examples/microsoft.json | gzip | curl -X POST -v --data-binary @- localhost:8080/report

    curl localhost:8080/.well-known/mta-sts.txt
    curl localhost:8081/metrics

Save reports

    mkdir reports
    chown 65532:65532 reports
    docker run -it --rm -p 8080:8080 -p 8081:8081 --env REPORTS_SAVE_ENABLED=true -v ${PWD}/reports:/tmp/reports phiag/mta-sts-exporter:latest

## Metrics

- `mta_sts_policy_requests_total` Total number of policy requests.
- `mta_sts_report_requests_total` Total number of report requests.
- `mta_sts_report_errors_total` Total number of report errors partitioned by `cause`.
- `mta_sts_successful_sessions_total` Total number of reported successful sessions partitioned by `organization`.
- `mta_sts_failure_sessions_total` Total number of reported failure sessions partitioned by `organization`.

## Development

Test

    go test

Lint

    golangci-lint run

Format

    go fmt

Run

    PORT=1234 METRICS_PORT=3333 go run .
    cat examples/rfc.json | gzip | curl -X POST -v --data-binary @- localhost:1234/report
    curl localhost:3333/metrics

## References

- [RFC 8460: SMTP TLS Reporting](https://www.rfc-editor.org/rfc/rfc8460.html)
- [RFC 8461: SMTP MTA Strict Transport Security (MTA-STS)](https://www.rfc-editor.org/rfc/rfc8461.html)
