package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func reportExample(name string) *os.File {
	reader, err := os.Open(fmt.Sprintf("examples/%s.json", name))
	if err != nil {
		panic(err)
	}
	return reader
}

func reportExampleGzip(name string) *io.PipeReader {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()

		gzipWriter := gzip.NewWriter(writer)
		defer gzipWriter.Close()

		reportReader := reportExample(name)
		defer reportReader.Close()

		_, err := io.Copy(gzipWriter, reportReader)
		if err != nil {
			panic(err)
		}
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
	reportReader := reportExample("rfc")
	defer reportReader.Close()

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
