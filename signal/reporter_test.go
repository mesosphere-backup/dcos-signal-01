package signal

import (
	//	"fmt"
	"fmt"
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
		tr     = testReportType{
			URL:    server.URL,
			Method: "GET",
		}
		tc = config.Config{}
	)

	// Get the good report first
	goodReportErr := PullReport(&tr, tc)

	// Break it with bad JSON
	tr.URL = fmt.Sprintf("%s/500", server.URL)

	if goodReportErr != nil {
		t.Error("Expected nil error, got ", goodReportErr.Error)
	}
}
