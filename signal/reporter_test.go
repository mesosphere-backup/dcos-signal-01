package signal

import (
	"fmt"
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

func (t *testReportType) getReport() interface{} { return t.Report }

func (t *testReportType) setTrack(config.Config) error { return nil }

func (t *testReportType) getTrack() (a *analytics.Track) { return a }

func (t *testReportType) sendTrack(config.Config) error { return nil }

func (t *testReportType) getName() string { return "" }

func (t *testReportType) setError(string) {}

func (t *testReportType) getError() string { return "" }

func (t *testReportType) setEndpoints(url []string) { t.Endpoints = url }

func (t *testReportType) getEndpoints() []string { return t.Endpoints }

func (t *testReportType) setMethod(meth string) { t.Method = meth }

func (t *testReportType) getMethod() string { return t.Method }

func (t *testReportType) getHeaders() map[string]string { return t.Headers }

func (t *testReportType) addHeaders(head map[string]string) {
	for k, v := range head {
		t.Headers[k] = v
	}
}

func (t *testReportType) setReport(report []byte) error {
	t.Report = string(report)
	return nil
}

func TestPullHealthReport(t *testing.T) {
	var (
		tr = testReportType{
			Endpoints: []string{
				fmt.Sprintf("%s/package/list", server.URL),
			},
			Method: "GET",
		}
		tc = config.Config{}
	)

	goodReportErr := PullReport(&tr, tc)
	if goodReportErr != nil {
		t.Error("Expected nil error, got ", goodReportErr.Error)
	}
}
