package signal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
)

var (
	VERSION  = "UNSET"
	REVISION = "UNSET"
)

func runner(done chan Reporter, reporters chan Reporter, c config.Config, w int) error {
	for r := range reporters {
		log.Debugf("Worker %d: Processing job for %+v", w, r)
		err := PullReport(r, c)
		if err != nil {
			log.Errorf("Worker %d: %s", w, err)
			r.setError(err.Error())
			done <- r
			return err
		}

		err = r.setTrack(c)
		if err != nil {
			log.Errorf("Worker %d: %s", w, err)
			r.setError(err.Error())
			done <- r
			return err
		}
		done <- r
	}
	return nil
}

func executeTester(data map[string]*analytics.Track, c config.Config) error {
	log.Info("Executing POST to test server")
	jsonStr, _ := json.MarshalIndent(data, "", "    ")
	url := fmt.Sprintf("http://%s:%d", c.TestHost, c.TestPort)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Infof("Test server response: %s", resp.Status)
	return nil
}

func executeRunner(c config.Config) error {
	log.Info("==> STARTING SIGNAL RUNNER")

	// Get our channel of jobs (reporters)
	reporters, err := makeReporters(c)
	if err != nil {
		return errors.New("Unable to get reporters.")
	}
	// Make a channel to dump the built tracks to
	done := make(chan Reporter)

	workers := len(reporters)
	for w := 1; w <= workers; w++ {
		log.Debugf("Deploying Worker %d", w)
		go runner(done, reporters, c, w)
	}

	tester := make(map[string]*analytics.Track)

	processed := 1
	for processed <= workers {
		select {
		case r := <-done:
			if c.TestFlag {
				log.Info("Adding test data")
				tester[r.getName()] = r.getTrack()
			} else if len(r.getError()) > 0 {
				log.Errorf("%s: %s", r.getName(), r.getError())
			} else {
				r.sendTrack(c)
			}
			processed += 1
		}
	}

	if c.TestFlag {
		if err := executeTester(tester, c); err != nil {
			return err
		}
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
			log.SetLevel(log.DebugLevel)
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
