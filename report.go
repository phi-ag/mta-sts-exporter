package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type DateRange struct {
	StartDatetime time.Time `json:"start-datetime"`
	EndDatetime   time.Time `json:"end-datetime"`
}

type Policy struct {
	PolicyType   string   `json:"policy-type"`
	PolicyString []string `json:"policy-string"`
	PolicyDomain string   `json:"policy-domain"`
	MxHost       string   `json:"mx-host"`
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
	Policy         Policy          `json:"policy"`
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
	err := os.MkdirAll(config.Reports.SavePath, os.ModePerm)
	if err != nil {
		slog.Error("Failed to create directory", "path", config.Reports.SavePath, "error", err)
		return reader, nil, err
	}

	filename := time.Now().Format(time.RFC3339Nano) + ".json"
	target := filepath.Join(config.Reports.SavePath, filename)
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
