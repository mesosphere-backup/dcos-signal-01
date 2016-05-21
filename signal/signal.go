package signal

import (
	"encoding/json"
	"fmt"
	"os"

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
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)

	if err := PullReport(&diagnostics, c); err != nil {
		log.Error("Error getting diagnostics report")
		return err
	}

	if err := diagnostics.SetTrack(c); err != nil {
		log.Error("Unable to set diagnostics .track, ", err)
		return err
	}

	if c.TestFlag {
		pretty, _ := json.MarshalIndent(diagnostics.GetTrack(), "", "    ")
		fmt.Printf(string(pretty))
		return nil
	} else {
		if err := diagnostics.SendTrack(c); err != nil {
			log.Error("Error sending diagnostics track data")
			return err
		}
	}

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
