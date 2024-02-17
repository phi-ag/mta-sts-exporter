package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

var reportExample = `{
	"organization-name": "Company-X",
	"date-range": {
		"start-datetime": "2016-04-01T00:00:00Z",
		"end-datetime": "2016-04-01T23:59:59Z"
	},
	"contact-info": "sts-reporting@company-x.example",
	"report-id": "5065427c-23d3-47ca-b6e0-946ea0e8c4be",
	"policies": [{
		"policy": {
			"policy-type": "sts",
			"policy-string": ["version: STSv1","mode: testing","mx: *.mail.company-y.example","max_age: 86400"],
			"policy-domain": "company-y.example",
			"mx-host": "*.mail.company-y.example"
		},
		"summary": {
			"total-successful-session-count": 5326,
			"total-failure-session-count": 303
		},
		"failure-details": [{
			"result-type": "certificate-expired",
			"sending-mta-ip": "2001:db8:abcd:0012::1",
			"receiving-mx-hostname": "mx1.mail.company-y.example",
			"failed-session-count": 100
		}, {
			"result-type": "starttls-not-supported",
			"sending-mta-ip": "2001:db8:abcd:0013::1",
			"receiving-mx-hostname": "mx2.mail.company-y.example",
			"receiving-ip": "203.0.113.56",
			"failed-session-count": 200,
			"additional-information": "https://reports.company-x.example/report_info ? id = 5065427 c - 23 d3# StarttlsNotSupported "
		}, {
			"result-type": "validation-failure",
			"sending-mta-ip": "198.51.100.62",
			"receiving-ip": "203.0.113.58",
			"receiving-mx-hostname": "mx-backup.mail.company-y.example",
			"failed-session-count": 3,
			"failure-reason-code": "X509_V_ERR_PROXY_PATH_LENGTH_EXCEEDED"
		}]
	}]
}
`

func reportExampleGzip() *io.PipeReader {
	reader, writer := io.Pipe()
	gzipWriter := gzip.NewWriter(writer)

	go func() {
		dataReader := strings.NewReader(reportExample)
		_, err := io.Copy(gzipWriter, dataReader)
		if err != nil {
			panic(err)
		}

		gzipWriter.Close()
		writer.Close()
	}()

	return reader
}

func TestParseReportTime(t *testing.T) {
	start, _ := time.Parse(time.RFC3339, "2016-04-01T00:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2016-04-01T23:59:59Z")

	report := Report{
		OrganizationName: "Foo Ltd.",
		DateRange: DateRange{
			StartDatetime: start,
			EndDatetime:   end,
		},
	}

	json, _ := json.Marshal(report)
	reportReader := bytes.NewReader(json)

	parsed, _ := parseReport(reportReader)

	if parsed.OrganizationName != "Foo Ltd." {
		t.Errorf("expected Foo Ltd. got %v", parsed.OrganizationName)
	}

	if parsed.DateRange.StartDatetime != start {
		t.Errorf("expected start to equal initial value got %v", parsed.DateRange.StartDatetime)
	}
}

func TestParseReportRfcExample(t *testing.T) {
	reportReader := strings.NewReader(reportExample)
	parsed, err := parseReport(reportReader)
	if err != nil {
		t.Error("failed to parse report example", err)
	}

	if parsed.OrganizationName != "Company-X" {
		t.Errorf("expected OrganizationName Company-X got %v", parsed.OrganizationName)
	}

	if len(parsed.Policies) != 1 {
		t.Errorf("expected single Policy got %v", len(parsed.Policies))
	}

	if len(parsed.Policies[0].FailureDetails) != 3 {
		t.Errorf("expected three FailureDetails got %v", len(parsed.Policies[0].FailureDetails))
	}
}
