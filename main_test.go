package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturnsMethodNotAllowedForGetRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	config := Config{
		Reports: Reports{
			Save:        false,
			MaxBodySize: 5000,
			MaxJsonSize: 5000,
		},
	}

	handleReport(config)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected StatusMethodNotAllowed got %v", res.Status)
	}
}

func TestReturnsBadRequestForNonGzip(t *testing.T) {
	body := make([]byte, 50)
	reader := bytes.NewReader(body)

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	recorder := httptest.NewRecorder()

	config := Config{
		Reports: Reports{
			Save:        false,
			MaxBodySize: 5000,
			MaxJsonSize: 5000,
		},
	}

	handleReport(config)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest got %v", res.Status)
	}
}

func TestReturnsRequestEntityTooLargeForBody(t *testing.T) {
	reader := reportExampleGzip("rfc")
	defer reader.Close()

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	recorder := httptest.NewRecorder()

	config := Config{
		Reports: Reports{
			Save:        false,
			MaxBodySize: 25,
			MaxJsonSize: 5_000,
		},
	}

	handleReport(config)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("expected StatusRequestEntityTooLarge got %v", res.Status)
	}
}

func TestReturnsRequestEntityTooLargeForJson(t *testing.T) {
	reader := reportExampleGzip("rfc")
	defer reader.Close()

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	recorder := httptest.NewRecorder()

	config := Config{
		Reports: Reports{
			Save:        false,
			MaxBodySize: 5_000,
			MaxJsonSize: 25,
		},
	}

	handleReport(config)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("expected StatusRequestEntityTooLarge got %v", res.Status)
	}
}

func TestReturnsOk(t *testing.T) {
	reader := reportExampleGzip("rfc")
	defer reader.Close()

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	recorder := httptest.NewRecorder()

	config := Config{
		Reports: Reports{
			Save:        false,
			MaxBodySize: 5_000,
			MaxJsonSize: 5_000,
		},
	}

	handleReport(config)(recorder, req)
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected StatusOK got %v", res.Status)
	}
}
