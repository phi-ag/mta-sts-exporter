package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RFC3339OptionalTimezone time.Time

func (t *RFC3339OptionalTimezone) UnmarshalJSON(input []byte) error {
	trimmed := strings.Trim(string(input), `"`)

	parsedTime, err := time.Parse(time.RFC3339, trimmed)
	if err == nil {
		*t = RFC3339OptionalTimezone(parsedTime)
		return nil
	}

	RFC3339NoTimeZone := "2006-01-02T15:04:05"
	parsedTime, err = time.Parse(RFC3339NoTimeZone, trimmed)
	if err == nil {
		*t = RFC3339OptionalTimezone(parsedTime)
		return nil
	}

	return fmt.Errorf("invalid time '%s'", trimmed)
}

func (t RFC3339OptionalTimezone) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t))
}

type DateRange struct {
	StartDatetime RFC3339OptionalTimezone `json:"start-datetime"`
	EndDatetime   RFC3339OptionalTimezone `json:"end-datetime"`
}

type MxHost []string

// NOTE: see https://www.rfc-editor.org/errata/eid6241
func (l *MxHost) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return nil
	}

	switch input[0] {
	case '"':
		*l = MxHost{strings.Trim(string(input), `"`)}
		return nil

	case '[':
		var tmp []string
		err := json.Unmarshal(input, &tmp)
		if err != nil {
			return err
		}
		*l = MxHost(tmp)
		return nil

	default:
		return fmt.Errorf("invalid mx-host '%s'", string(input))
	}
}

type ReportPolicy struct {
	PolicyType   string   `json:"policy-type"`
	PolicyString []string `json:"policy-string"`
	PolicyDomain string   `json:"policy-domain"`
	MxHost       MxHost   `json:"mx-host"`
}

type Summary struct {
	TotalSuccessfulSessionCount int64 `json:"total-successful-session-count"`
	TotalFailureSessionCount    int64 `json:"total-failure-session-count"`
}

type FailureDetail struct {
	ResultType            string `json:"result-type"`
	SendingMtaIp          string `json:"sending-mta-ip"`
	ReceivingIp           string `json:"receiving-ip"`
	ReceivingMxHostname   string `json:"receiving-mx-hostname"`
	FailedSessionCount    int64  `json:"failed-session-count"`
	FailureReasonCode     string `json:"failure-reason-code"`
	AdditionalInformation string `json:"additional-information"`
}

type PolicyItem struct {
	Policy         ReportPolicy    `json:"policy"`
	Summary        Summary         `json:"summary"`
	FailureDetails []FailureDetail `json:"failure-details"`
}

type Report struct {
	OrganizationName string       `json:"organization-name"`
	DateRange        DateRange    `json:"date-range"`
	ContactInfo      string       `json:"contact-info"`
	ReportId         string       `json:"report-id"`
	Policies         []PolicyItem `json:"policies"`
}

func parseReport(reader io.Reader) (Report, error) {
	report := &Report{}
	err := json.NewDecoder(reader).Decode(report)
	return *report, err
}

func saveReport(config Config, reader io.Reader) (io.Reader, *os.File, error) {
	err := os.MkdirAll(config.Reports.Save.Path, os.ModePerm)
	if err != nil {
		slog.Error("Failed to create directory", "path", config.Reports.Save.Path, "error", err)
		return reader, nil, err
	}

	filename := time.Now().Format(time.RFC3339Nano) + ".json"
	target := filepath.Join(config.Reports.Save.Path, filename)
	slog.Info("Saving report", "target", target)

	file, err := os.Create(target)
	if err != nil {
		slog.Error("Failed to create file", "target", target, "error", err)
		return reader, nil, err
	}

	/// NOTE: It seems `TeeReader` writes the complete stream even when parsing fails later.
	/// This is probably only true for small payloads.
	return io.TeeReader(reader, file), file, nil
}
