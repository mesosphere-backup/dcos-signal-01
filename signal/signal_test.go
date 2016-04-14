package signal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"testing"

	"github.com/dcos/dcos-signal/config"
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

func TestSignalRunner(t *testing.T) {
	var (
		healthServer = httptest.NewServer(mockRouter())
		port, _      = strconv.Atoi(strings.Split(healthServer.URL, ":")[2])
		ip           = strings.Split(strings.Split(healthServer.URL, ":")[1], "/")[1]

		cOk       = config.DefaultConfig()
		badJson   = config.DefaultConfig()
		badUserId = config.DefaultConfig()
		badHost   = config.DefaultConfig()
		version   = config.DefaultConfig()
		verbose   = config.DefaultConfig()
	)

	cOk.HealthAPIPort = port
	cOk.HealthHost = ip
	cOk.HealthEndpoint = "/system/healt/report/test"
	cOk.CustomerKey = "12345"

	verbose.FlagVerbose = true
	verbose.HealthAPIPort = port
	verbose.HealthHost = ip
	verbose.HealthEndpoint = "/system/healt/report/test"
	verbose.CustomerKey = "12345"

	badUserId.HealthAPIPort = port
	badUserId.HealthHost = ip
	badUserId.HealthEndpoint = "/system/healt/report/test"

	badJson.HealthAPIPort = port
	badJson.HealthHost = ip
	badJson.CustomerKey = "12345"
	badJson.HealthEndpoint = "/system/health/report/test/badjson"

	badHost.HealthEndpoint = "/system/healt/report/test"
	badHost.CustomerKey = "12345"
	badHost.HealthAPIPort = 80
	badHost.HealthHost = "localhost"

	version.FlagVersion = true

	var (
		errOk     = executeRunner(cOk)
		errJson   = executeRunner(badJson)
		errUserId = executeRunner(badUserId)
		errHost   = executeRunner(badHost)
	)

	if errOk != nil {
		t.Error("Expected nil error with good config, got ", errOk)
	}
	if errJson == nil {
		t.Error("Expected bad JSON to throw err, got ", errJson)
	}
	if errHost == nil {
		t.Error("Expected bad route to host error, got ", errHost)
	}
	if errUserId == nil {
		t.Error("Expected bad segment user ID to throw err, got ", errUserId)
	}
}
