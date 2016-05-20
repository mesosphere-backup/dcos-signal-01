package signal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dcos/dcos-signal/config"
	"github.com/gorilla/mux"
)

var mockNodes = []*Node{
	&Node{
		Role:   "master",
		Ip:     "10.0.0.1",
		Host:   "foo.master",
		Health: 0,
		Output: map[string]string{
			"foo-unit.2": "",
			"foo-unit.1": "",
		},
		Units: nil,
	},
	&Node{
		Role:   "slave",
		Ip:     "10.0.0.2",
		Host:   "foo.slave",
		Health: 1,
		Output: map[string]string{
			"foo-unit.2": "Something is broken!!",
			"foo-unit.1": "",
		},
		Units: nil,
	},
}

var mockUnits = map[string]*Unit{
	"10.0.0.1": {
		UnitName: "foo-unit.1",
		Nodes:    mockNodes,
		Health:   1,
		Title:    "Foo Test 1",
	},
	"10.0.0.2": {
		UnitName: "foo-unit.2",
		Nodes:    mockNodes,
		Health:   1,
		Title:    "Foo Test 2",
	},
}

var mockHealthReport = &HealthReport{
	Units: mockUnits,
	Nodes: nil,
}

// Mock a /system/health/report endpoint
func mockReportHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(mockHealthReport)
}

func mockBadJson(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("foo")
}

func mockFive(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(500), 500)
}

func mockFour(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(400), 400)
}

func mockRouter() *mux.Router {
	baseUrl := "/system/healt/report/test"
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(baseUrl, mockReportHandler).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/badjson", baseUrl), mockBadJson).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/500", baseUrl), mockFive).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/400", baseUrl), mockFour).Methods("GET")
	return router
}

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

func (t *testReportType) Track(config.Config) error {
	return nil
}

func TestPullHealthReport(t *testing.T) {
	var (
		server  = httptest.NewServer(mockRouter())
		baseUrl = "/system/health/report/test"
		tr      = testReportType{
			URL:    "server.URL",
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
