package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GatherCounters(config Config, metrics Metrics) map[string]float64 {
	registry := createRegistry(config, metrics)
	metricFamilies, err := registry.Gather()
	if err != nil {
		panic(err)
	}

	counters := make(map[string]float64)
	for _, family := range metricFamilies {
		counters[*family.Name] = *family.Metric[0].Counter.Value
	}

	return counters
}

func TestReportReturnsMethodNotAllowedForGetRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	config := Config{}

	metrics := createMetrics()
	handleReport(config, metrics)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected StatusMethodNotAllowed got %v", res.Status)
	}

	counters := GatherCounters(config, metrics)
	if counters["mta_sts_report_requests_total"] != 0 {
		t.Errorf("expected 0 report requests got %v", counters["mta_sts_report_requests_total"])
	}
}

func TestReportReturnsBadRequestForNonGzip(t *testing.T) {
	body := make([]byte, 50)
	reader := bytes.NewReader(body)

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	recorder := httptest.NewRecorder()

	config := Config{}

	metrics := createMetrics()
	handleReport(config, metrics)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest got %v", res.Status)
	}

	counters := GatherCounters(config, metrics)
	if counters["mta_sts_report_requests_total"] != 1 {
		t.Errorf("expected 1 report requests got %v", counters["mta_sts_report_requests_total"])
	}

	if counters["mta_sts_report_errors_total"] != 1 {
		t.Errorf("expected 1 report errors got %v", counters["mta_sts_report_errors_total"])
	}
}

func TestReportReturnsRequestEntityTooLargeForBody(t *testing.T) {
	reader := reportExampleGzip("rfc")
	defer reader.Close()

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	recorder := httptest.NewRecorder()

	config := Config{
		Reports: Reports{
			Max: ReportsMax{
				Body: 25,
				Json: 5_000,
			},
		},
	}

	metrics := createMetrics()
	handleReport(config, metrics)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("expected StatusRequestEntityTooLarge got %v", res.Status)
	}

	counters := GatherCounters(config, metrics)
	if counters["mta_sts_report_requests_total"] != 1 {
		t.Errorf("expected 1 report requests got %v", counters["mta_sts_report_requests_total"])
	}

	if counters["mta_sts_report_errors_total"] != 1 {
		t.Errorf("expected 1 report errors got %v", counters["mta_sts_report_errors_total"])
	}
}

func TestReportReturnsRequestEntityTooLargeForJson(t *testing.T) {
	reader := reportExampleGzip("rfc")
	defer reader.Close()

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	recorder := httptest.NewRecorder()

	config := Config{
		Reports: Reports{
			Max: ReportsMax{
				Body: 5_000,
				Json: 25,
			},
		},
	}

	metrics := createMetrics()
	handleReport(config, metrics)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("expected StatusRequestEntityTooLarge got %v", res.Status)
	}

	counters := GatherCounters(config, metrics)
	if counters["mta_sts_report_requests_total"] != 1 {
		t.Errorf("expected 1 report requests got %v", counters["mta_sts_report_requests_total"])
	}

	if counters["mta_sts_report_errors_total"] != 1 {
		t.Errorf("expected 1 report errors got %v", counters["mta_sts_report_errors_total"])
	}
}

func TestReportReturnsOk(t *testing.T) {
	reader := reportExampleGzip("rfc")
	defer reader.Close()

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	recorder := httptest.NewRecorder()

	config := Config{
		Reports: Reports{
			Max: ReportsMax{
				Body: 5_000,
				Json: 5_000,
			},
		},
	}

	metrics := createMetrics()
	handleReport(config, metrics)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected StatusOK got %v", res.Status)
	}

	counters := GatherCounters(config, metrics)

	if counters["mta_sts_policy_requests_total"] != 0 {
		t.Errorf("expected 0 policy requests got %v", counters["mta_sts_policy_requests_total"])
	}

	if counters["mta_sts_report_requests_total"] != 1 {
		t.Errorf("expected 1 report requests got %v", counters["mta_sts_report_requests_total"])
	}

	if counters["mta_sts_successful_sessions_total"] != 5326 {
		t.Errorf("expected 5326 successful sessions got %v", counters["mta_sts_successful_sessions_total"])
	}

	if counters["mta_sts_failure_sessions_total"] != 303 {
		t.Errorf("expected 303 failure sessions got %v", counters["mta_sts_failure_sessions_total"])
	}

	if counters["mta_sts_report_errors_total"] != 0 {
		t.Errorf("expected 0 report errors got %v", counters["mta_sts_report_errors_total"])
	}
}

func TestPolicyReturnsMethodNotAllowedForPostRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/.well-known/mta-sts.txt", nil)
	recorder := httptest.NewRecorder()

	config := Config{}

	metrics := createMetrics()
	handlePolicy(config, metrics)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected StatusMethodNotAllowed got %v", res.Status)
	}

	counters := GatherCounters(config, metrics)
	if counters["mta_sts_policy_requests_total"] != 0 {
		t.Errorf("expected 0 policy requests got %v", counters["mta_sts_policy_requests_total"])
	}
}

func TestPolicyReturnsOk(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/.well-known/mta-sts.txt", nil)
	recorder := httptest.NewRecorder()

	config := Config{
		Policy: Policy{
			Version: "STSv1",
			Mode:    "testing",
			Mx:      []string{"mx1.example.com", "mx2.example.com"},
			MaxAge:  600,
		},
	}

	expected := "version: STSv1\nmode: testing\nmx: mx1.example.com\nmx: mx2.example.com\nmax_age: 600\n"

	metrics := createMetrics()
	handlePolicy(config, metrics)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected StatusOK got %v", res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	body := string(bodyBytes)
	if body != expected {
		t.Errorf("expected policy got %v", body)
	}

	counters := GatherCounters(config, metrics)
	if counters["mta_sts_policy_requests_total"] != 1 {
		t.Errorf("expected 1 policy requests got %v", counters["mta_sts_policy_requests_total"])
	}
}
