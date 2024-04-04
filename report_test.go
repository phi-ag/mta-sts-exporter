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
	_start, _ := time.Parse(time.RFC3339, "2016-04-01T00:00:00Z")
	_end, _ := time.Parse(time.RFC3339, "2016-04-01T23:59:59Z")

	start := RFC3339OptionalTimezone(_start)
	end := RFC3339OptionalTimezone(_end)

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
	examples := []string{"rfc", "rfc-errata-6241"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			reportReader := reportExample(example)
			defer reportReader.Close()

			parsed, err := parseReport(reportReader)
			if err != nil {
				t.Error("failed to parse report example", err)
			}

			if parsed.OrganizationName != "Company-X" {
				t.Errorf("expected OrganizationName Company-X got %v", parsed.OrganizationName)
			}

			if len(parsed.Policies) != 1 {
				t.Errorf("expected 1 Policy got %v", len(parsed.Policies))
			}

			if parsed.Policies[0].Policy.MxHost[0] != "*.mail.company-y.example" {
				t.Errorf("expected MxHost *.mail.company-y.example got %v", parsed.Policies[0].Policy.MxHost[0])
			}

			if len(parsed.Policies[0].FailureDetails) != 3 {
				t.Errorf("expected 3 FailureDetails got %v", len(parsed.Policies[0].FailureDetails))
			}
		})
	}
}

func TestParseReportGoogleExample(t *testing.T) {
	reportReader := reportExample("google")
	defer reportReader.Close()

	parsed, err := parseReport(reportReader)
	if err != nil {
		t.Error("failed to parse report example", err)
	}

	if parsed.OrganizationName != "Google Inc." {
		t.Errorf("expected OrganizationName Google Inc. got %v", parsed.OrganizationName)
	}

	if len(parsed.Policies) != 1 {
		t.Errorf("expected 1 Policy got %v", len(parsed.Policies))
	}

	if len(parsed.Policies[0].Policy.MxHost) != 1 {
		t.Errorf("expected 1 MxHost got %v", len(parsed.Policies[0].Policy.MxHost))
	}

	if parsed.Policies[0].Policy.MxHost[0] != "example.com" {
		t.Errorf("expected MxHost example.com got %v", parsed.Policies[0].Policy.MxHost[0])
	}
}

func TestParseReportMicrosoftExample(t *testing.T) {
	reportReader := reportExample("microsoft")
	defer reportReader.Close()

	parsed, err := parseReport(reportReader)
	if err != nil {
		t.Error("failed to parse report example", err)
	}

	if parsed.OrganizationName != "Microsoft Corporation" {
		t.Errorf("expected OrganizationName Microsoft Corporation got %v", parsed.OrganizationName)
	}

	if len(parsed.Policies) != 1 {
		t.Errorf("expected 1 Policy got %v", len(parsed.Policies))
	}

	if len(parsed.Policies[0].Policy.MxHost) != 0 {
		t.Errorf("expected no MxHost got %v", len(parsed.Policies[0].Policy.MxHost))
	}
}

func TestParseReportMicrosoftExampleWithInvalidTime(t *testing.T) {
	reportReader := reportExample("microsoft-2")
	defer reportReader.Close()

	parsed, err := parseReport(reportReader)
	if err != nil {
		t.Error("failed to parse report example", err)
	}

	if parsed.OrganizationName != "Microsoft Corporation" {
		t.Errorf("expected OrganizationName Microsoft Corporation got %v", parsed.OrganizationName)
	}

	if len(parsed.Policies) != 1 {
		t.Errorf("expected 1 Policy got %v", len(parsed.Policies))
	}

	if len(parsed.Policies[0].Policy.MxHost) != 0 {
		t.Errorf("expected no MxHost got %v", len(parsed.Policies[0].Policy.MxHost))
	}

	expectedStart, _ := time.Parse(time.RFC3339, "2024-03-26T00:00:00Z")

	if parsed.DateRange.StartDatetime != RFC3339OptionalTimezone(expectedStart) {
		t.Errorf("expected start to equal %v got %v", expectedStart, parsed.DateRange.StartDatetime)
	}
}
