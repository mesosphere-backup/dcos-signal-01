package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"
)

// Example response from licensing service's /licenses endpoint
// [
//   {
//     "id": "test_license_id",
//     "version": "1.11",
//     "decryption_key": "...",
//     "license_terms": {
//       "node_capacity": 10,
//       "start_timestamp": "2006-01-02T15:04:05Z",
//       "end_timestamp": "2026-01-02T15:04:05Z"
//     }
//   }
// ]
//

type Licenses []License

type License struct {
	ID            string       `json:"id"`
	Version       string       `json:"version"`
	DecryptionKey string       `json:"decryption_key"`
	LicenseTerms  LicenseTerms `json:"license_terms"`
}

type LicenseTerms struct {
	NodeCapacity   int       `json:"node_capacity"`
	StartTimestamp time.Time `json:"start_timestamp"`
	EndTimestamp   time.Time `json:"end_timestamp"`
}

func (l Licenses) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

const SOCK = "/tmp/dcos-licensing.socket"

func main() {

	expired := License{
		ID:            "expired_license",
		Version:       "1.11",
		DecryptionKey: "",
		LicenseTerms: LicenseTerms{
			NodeCapacity:   10,
			StartTimestamp: time.Now().AddDate(-2, 0, 0),
			EndTimestamp:   time.Now().AddDate(-1, 0, 0),
		},
	}

	current := License{
		ID:            "current_license",
		Version:       "1.12",
		DecryptionKey: "",
		LicenseTerms: LicenseTerms{
			NodeCapacity:   10,
			StartTimestamp: time.Now().AddDate(-1, 0, 0),
			EndTimestamp:   time.Now().AddDate(0, 1, 0),
		},
	}

	upgrade := License{
		ID:            "upgrade_license",
		Version:       "1.13",
		DecryptionKey: "",
		LicenseTerms: LicenseTerms{
			NodeCapacity:   10,
			StartTimestamp: time.Now().AddDate(0, -1, 0),
			EndTimestamp:   time.Now().AddDate(1, -1, 0),
		},
	}

	server := http.Server{
		Handler: Licenses{expired, current, upgrade},
	}

	os.Remove(SOCK)
	unixListener, err := net.Listen("unix", SOCK)
	if err != nil {
		panic(err)
	}
	server.Serve(unixListener)
}
