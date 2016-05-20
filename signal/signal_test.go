package signal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

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

var (
	mockHealthReport = &HealthReport{
		Units: mockUnits,
		Nodes: nil,
	}

	cosmosPkgs = CosmosPackages{
		AppID: "fooPackage",
	}

	mesosFrameworks = map[string][]string{
		"frameworks": []string{
			"fooFramework1",
			"fooFramework2",
		},
	}

	mesosMetricsSnapshot = map[string]int{
		"master/cpus_total":        10,
		"master/cpus_used":         2,
		"master/disk_total":        1000,
		"master/disk_used":         20,
		"master/mem_total":         2000,
		"master/mem_used":          200,
		"master/tasks_running":     4,
		"master/frameworks_active": 2,
		"master/slaves_connected":  3,
		"master/slaves_active":     1,
	}

	mockCosmosReport = &CosmosReport{
		Packages: []CosmosPackages{
			cosmosPkgs,
		},
	}
	server = httptest.NewServer(mockRouter())
)

func mockHealthReportHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(mockHealthReport)
}

func mockFrameworksHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(mesosFrameworks)
}

func mockMesosStatsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(mesosMetricsSnapshot)
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
	var (
		health     = "/system/health/v1/report"
		cosmos     = "/package/list"
		frameworks = "/frameworks"
		mesosStats = "/metrics/snapshot"
	)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(health, mockHealthReportHandler).Methods("GET")
	router.HandleFunc(frameworks, mockFrameworksHandler).Methods("GET")
	router.HandleFunc(mesosStats, mockMesosStatsHandler).Methods("GET")
	router.HandleFunc(cosmos, mockCosmosReportHandler).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/badjson", health), mockBadJson).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/500", health), mockFive).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/400", health), mockFour).Methods("GET")
	return router
}
