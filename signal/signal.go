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
				log.Info("====> DONE")
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
			Endpoints: []string{
				":1050/system/health/v1/report",
			},
			Method: "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		}

		cosmos = Cosmos{
			Endpoints: []string{
				":7070/package/list",
			},
			Method: "POST",
			Headers: map[string]string{
				"content-type": "application/vnd.dcos.package.list-request",
			},
		}

		mesos = Mesos{
			Endpoints: []string{
				":5050/frameworks",
				":5050/metrics/snapshot",
			},
			Method: "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		}

		errored = []error{}
	)

	// Add extra headers if any
	for k, v := range c.ExtraHeaders {
		diagnostics.Headers[k] = v
		cosmos.Headers[k] = v
		mesos.Headers[k] = v
	}

	// Might want to declare a []Reporter and do this async once we get a few more.
	if err := runner(&diagnostics, c); err != nil {
		errored = append(errored, err)
	}

	if err := runner(&cosmos, c); err != nil {
		errored = append(errored, err)
	}

	if err := runner(&mesos, c); err != nil {
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
		fmt.Printf("DCOS Signal Service\n Version: %s\n Revision: %s\n DC/OS Variant: %s\n", VERSION, REVISION, config.DCOSVariant)
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
