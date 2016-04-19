package signal

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestPullHealthReport(t *testing.T) {
	var (
		server  = httptest.NewServer(mockRouter())
		baseUrl = "/system/healt/report/test"

		goodReport, goodReportErr   = pullHealthReport(server.URL, baseUrl)
		badJsonReport, badJsonErr   = pullHealthReport(server.URL, fmt.Sprintf("%s/badjson", baseUrl))
		badProtoReport, badProtoErr = pullHealthReport("foo.com", baseUrl)
		badHostReport, badHostErr   = pullHealthReport("http://foo:80", baseUrl)
		badUrl, badUrlErr           = pullHealthReport("", "")
	)

	if goodReportErr != nil {
		t.Error("Expected nil error, got ", goodReportErr.Error)
	}

	if _, ok := goodReport.Units["10.0.0.1"]; !ok {
		t.Error("Expected key '10.0.0.1', got ", ok)
	}

	if _, ok := goodReport.Units["10.0.0.2"]; !ok {
		t.Error("Expected key '10.0.0.2', got ", ok)
	}

	if badJsonReport != nil {
		t.Error("Expected nil report from bad JSON, got ", badJsonReport)
	}
	if badJsonErr == nil {
		t.Error("Expected from bad JSON, got ", badJsonErr.Error())
	}

	if badJsonErr.Error() != "json: cannot unmarshal string into Go value of type signal.HealthReport" {
		t.Error("Expected \"json: cannot unmarshal string into Go value of type signal.HealthReport\", got ", badJsonErr.Error())
	}

	if badProtoReport != nil {
		t.Error("Expected bad protocol to throw err, got ", badProtoReport)
	}

	if badProtoErr == nil {
		t.Error("Expected thrown error, got ", badProtoErr)
	}

	if badProtoErr.Error() != "Get foo.com/system/healt/report/test: unsupported protocol scheme \"\"" {
		t.Error(`Expected "Get foo.com/system/healt/report/test: unsupported protocol scheme", got `, badProtoErr.Error())
	}

	if badHostReport != nil {
		t.Error("Expected error, got ", badHostReport)
	}

	if badHostErr == nil {
		t.Error("Expected error, got ", badHostErr)
	}

	if badUrl != nil {
		t.Error("Expected bad url to return nil, got ", badUrl)
	}

	if badUrlErr == nil {
		t.Error("Expected bad url to return err, got ", badUrlErr)
	}
}
