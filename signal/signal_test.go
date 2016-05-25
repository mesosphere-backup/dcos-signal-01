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

var (
	mockHealthReport = &HealthReport{
		Units: mockUnits,
		Nodes: nil,
	}

	cosmosPkgs = CosmosPackages{
		AppID: "fooPackage",
	}

	mockCosmosReport = &CosmosReport{
		Packages: []CosmosPackages{
			cosmosPkgs,
		},
	}
)

func mockHealthReportHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(mockHealthReport)
}

func mockCosmosReportHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(mockCosmosReport)
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
	health := "/system/health/report/test"
	cosmos := "/package/list"
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(health, mockHealthReportHandler).Methods("GET")
	router.HandleFunc(cosmos, mockCosmosReportHandler).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/badjson", health), mockBadJson).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/500", health), mockFive).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/400", health), mockFour).Methods("GET")
	return router
}

func TestSignalRunner(t *testing.T) {
	var (
		healthServer = httptest.NewServer(mockRouter())
		port, _      = strconv.Atoi(strings.Split(healthServer.URL, ":")[2])
		ip           = strings.Split(strings.Split(healthServer.URL, ":")[1], "/")[1]
		endpoint     = "/system/health/report/test"

		cOk       = config.DefaultConfig()
		badJson   = config.DefaultConfig()
		badUserId = config.DefaultConfig()
		badHost   = config.DefaultConfig()
		version   = config.DefaultConfig()
		verbose   = config.DefaultConfig()
	)

	cOk.CosmosURL = fmt.Sprintf("http://%s:%d/package/list", ip, port)
	cOk.DiagnosticsURL = fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
	cOk.CustomerKey = "12345"

	verbose.FlagVerbose = true
	verbose.DiagnosticsURL = fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)

	verbose.CustomerKey = "12345"

	badUserId.DiagnosticsURL = fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)

	badJson.CustomerKey = "12345"
	badJson.DiagnosticsURL = fmt.Sprintf("http://%s:%d/%s/badjson", ip, port, endpoint)

	badHost.CustomerKey = "12345"
	badHost.DiagnosticsURL = "http://foo"

	version.FlagVersion = true

	var (
		errOk   = executeRunner(cOk)
		errJson = executeRunner(badJson)
		errHost = executeRunner(badHost)
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
}
