package signal

import (
	//	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
)

type testReportType struct {
	URL     string
	Headers map[string]string
	Method  string
	Report  string
}

func (t *testReportType) SetURL(url string) {
	t.URL = url
}

func (t *testReportType) GetURL() string {
	return t.URL
}

func (t *testReportType) SetMethod(meth string) {
	t.Method = meth
}

func (t *testReportType) GetMethod() string {
	return t.Method
}

func (t *testReportType) SetHeaders(head map[string]string) {
	t.Headers = head
}

func (t *testReportType) GetHeaders() map[string]string {
	return t.Headers
}

func (t *testReportType) SetReport(report []byte) error {
	t.Report = string(report)
	return nil
}

func (t *testReportType) GetReport() interface{} {
	return t.Report
}

func (t *testReportType) SetTrack(config.Config) error {
	return nil
}

func (t *testReportType) GetTrack() (a *analytics.Track) {
	return a
}

func (t *testReportType) SendTrack(config.Config) error {
	return nil
}

func TestPullHealthReport(t *testing.T) {
	var (
		server = httptest.NewServer(mockRouter())
		//baseUrl = "/system/health/report/test"
		tr = testReportType{
			URL:    server.URL, //fmt.Sprintf("%s%s", server.URL, baseUrl),
			Method: "GET",
		}
		tc            = config.Config{}
		goodReportErr = PullReport(&tr, tc)
	)
	//		badJsonReport, badJsonErr   = pullHealthReport(server.URL, fmt.Sprintf("%s/badjson", baseUrl))
	//		badProtoReport, badProtoErr = pullHealthReport("foo.com", baseUrl)
	//		badHostReport, badHostErr   = pullHealthReport("http://foo:80", baseUrl)
	//		badUrl, badUrlErr           = pullHealthReport("", "")
	//	)
	//
	if goodReportErr != nil {
		t.Error("Expected nil error, got ", goodReportErr.Error)
	}

	//	if _, ok := goodReport.Units["10.0.0.1"]; !ok {
	//		t.Error("Expected key '10.0.0.1', got ", ok)
	//	}
	//
	//	if _, ok := goodReport.Units["10.0.0.2"]; !ok {
	//		t.Error("Expected key '10.0.0.2', got ", ok)
	//	}
	//
	//	if badJsonReport != nil {
	//		t.Error("Expected nil report from bad JSON, got ", badJsonReport)
	//	}
	//	if badJsonErr == nil {
	//		t.Error("Expected from bad JSON, got ", badJsonErr.Error())
	//	}
	//
	//	if badJsonErr.Error() != "json: cannot unmarshal string into Go value of type signal.HealthReport" {
	//		t.Error("Expected \"json: cannot unmarshal string into Go value of type signal.HealthReport\", got ", badJsonErr.Error())
	//	}
	//
	//	if badProtoReport != nil {
	//		t.Error("Expected bad protocol to throw err, got ", badProtoReport)
	//	}
	//
	//	if badProtoErr == nil {
	//		t.Error("Expected thrown error, got ", badProtoErr)
	//	}
	//
	//	if badProtoErr.Error() != "Get foo.com/system/healt/report/test: unsupported protocol scheme \"\"" {
	//		t.Error(`Expected "Get foo.com/system/healt/report/test: unsupported protocol scheme", got `, badProtoErr.Error())
	//	}
	//
	//	if badHostReport != nil {
	//		t.Error("Expected error, got ", badHostReport)
	//	}
	//
	//	if badHostErr == nil {
	//		t.Error("Expected error, got ", badHostErr)
	//	}
	//
	//	if badUrl != nil {
	//		t.Error("Expected bad url to return nil, got ", badUrl)
	//	}
	//
	//	if badUrlErr == nil {
	//		t.Error("Expected bad url to return err, got ", badUrlErr)
	//	}
}
