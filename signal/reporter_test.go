package signal

import (
	"testing"

	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
)

type testReportType struct {
	Endpoints []string
	Headers   map[string]string
	Method    string
	Report    string
}

func (t *testReportType) SetEndpoints(url []string) {
	t.Endpoints = url
}

func (t *testReportType) GetEndpoints() []string {
	return t.Endpoints
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
		tr = testReportType{
			Endpoints: []string{
				"/package/list",
			},
			Method: "GET",
		}
		tc = config.Config{
			MasterURL: server.URL,
		}
	)

	goodReportErr := PullReport(&tr, tc)
	if goodReportErr != nil {
		t.Error("Expected nil error, got ", goodReportErr.Error)
	}
}
