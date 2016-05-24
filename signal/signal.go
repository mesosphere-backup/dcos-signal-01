package signal

import (
	"encoding/json"
	"errors"
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

func runner(r Reporter, c config.Config) error {
	if err := PullReport(r, c); err == nil {
		if err := r.SetTrack(c); err == nil {
			if c.TestFlag {
				pretty, _ := json.MarshalIndent(r.GetTrack(), "", "    ")
				fmt.Printf(string(pretty))
				return nil
			} else {
				if err := r.SendTrack(c); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
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
				"content-type": "application/json",
			},
		}

		cosmos = Cosmos{
			URL:    c.CosmosURL,
			Method: "POST",
			Headers: map[string]string{
				"content-type": "application/vnd.dcos.package.list-request",
			},
		}

		errored = []error{}
	)

	if err := runner(&diagnostics, c); err != nil {
		errored = append(errored, err)
	}

	if err := runner(&cosmos, c); err != nil {
		errored = append(errored, err)
	}

	if len(errored) > 0 {
		for _, err := range errored {
			log.Error(err)
		}
		return errors.New("Errors encountered executing report runners")
	}

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
		for _, err := range configErr {
			log.Error(err)
		}
	}
	if err := executeRunner(config); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}
