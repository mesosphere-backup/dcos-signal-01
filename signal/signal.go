package signal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dcos/dcos-signal/config"
)

var (
	VERSION  = "UNSET"
	REVISION = "UNSET"
)

type test struct {
	Event      string
	UserId     string
	ClusterId  string
	Properties map[string]interface{}
}

// StartSignalRunner accepts Config and runs the signal service once. It returns
// an error if encountered.
func executeRunner(c config.Config) error {
	log.Info("==> STARTING SIGNAL RUNNER")

	var (
		diagnostics = Diagnostics{
			URL:    c.DiagnosticsURL,
			Method: "GET",
		}
	)

	healthReport, err := pullHealthReport(healthURL, c.HealthEndpoint)
	if err != nil {
		log.Error("==> ERROR GETTING REPORT.")
		log.Error("Are you sure the URL, endport and port are correct?")
		return err
	}

	log.Info("Retrieved health report from ", c.HealthHost, ":", c.HealthAPIPort, c.HealthEndpoint)

	ac := CreateSegmentClient(c.SegmentKey, c.FlagVerbose)
	track, test := CreateSegmentTrack(healthReport, c)
	if c.TestFlag {
		pretty, _ := json.MarshalIndent(test, "", "    ")
		fmt.Printf(string(pretty))
		return nil
	}
	if err := ac.Track(track); err != nil {
		log.Error(err)
		return err
	}

	ac.Close()
	log.Info("==> SUCCESS")
	return nil
}

func Start() {
	config, configErr := config.ParseArgsReturnConfig(os.Args[1:])
	switch {
	case config.FlagVersion:
		fmt.Println("DCOS Signal Service: version", VERSION, "on revision", REVISION)
		os.Exit(0)
	default:
		if config.Enabled == "false" {
			os.Exit(0)
		}
		if config.FlagVerbose {
			log.SetLevel(log.DebugLevel)
		}
		if config.TestFlag {
			log.SetLevel(log.ErrorLevel)
		}
	}
	if configErr != nil {
		// There can be a number of errors during the config parsing. Several files,
		// as well as other factors and we should definitly at least attempt to send
		// data to segment even if we can't find things like the anon uuid file or
		// signal service config json since those would indicate that something is
		// no right, and signal service is all about surfacing that kind of data.
		log.Error(configErr)
	}
	if err := executeRunner(config); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}
