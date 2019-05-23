package signal

import (
	"encoding/json"
	"errors"
	"fmt"
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
		if len(r.getEndpoints()) == 0 {
			return fmt.Errorf("Reporter %s has no endpoints", r.getName())
		}
		for _, endpoint := range r.getEndpoints() {
			log.Debugf("Worker %d: Processing %s endpoint %s", w, r.getName(), endpoint)
			if err := PullReport(endpoint, r, c); err != nil {
				log.Errorf("Error setting track for %s: %s", r.getName(), err.Error())
				r.appendError(err.Error())
			}

			if err := r.setTrack(c); err != nil {
				log.Errorf("Error setting track for %s: %s", r.getName(), err.Error())
				r.appendError(err.Error())
			}
		}
		done <- r
	}
	return nil
}

func executeTester(data map[string]*analytics.Track, c config.Config) error {
	jsonStr, err := json.MarshalIndent(data, "", "    ")
	fmt.Print(string(jsonStr))
	return err
}

func executeRunner(c config.Config) error {
	log.Info("==> STARTING SIGNAL RUNNER")

	// Get our channel of jobs (reporters)
	reporters, err := makeReporters(c)
	if err != nil {
		return errors.New("unable to get reporters")
	}
	// Make a channel to dump the built tracks to
	done := make(chan Reporter)

	workers := len(reporters)
	for w := 1; w <= workers; w++ {
		log.Debugf("Deploying Worker %d", w)
		// runner probably shouldn't be returning an error but should send that to the reporters
		// channel. However that's a large change so we ignore this check here
		// nolint: errcheck
		go runner(done, reporters, c, w)
	}

	tester := make(map[string]*analytics.Track)

	processed := 1
	for r := range done {
		for _, err := range r.getError() {
			log.Errorf("%s: %s", r.getName(), err)
		}
		if c.FlagTest {
			log.Debugf("Adding test data for %s: %+v", r.getName(), r.getTrack())
			tester[r.getName()] = r.getTrack()
		} else if len(r.getError()) > 0 {
			for _, err := range r.getError() {
				log.Errorf("%s: %s", r.getName(), err)
			}
		} else {
			_ = r.sendTrack(c)
		}
		log.Warnf("processed %d, workers %d", processed, workers)
		processed++
	}

	if c.FlagTest {
		if err := executeTester(tester, c); err != nil {
			return err
		}
	}

	return nil
}

// Start starts the signal service
func Start() {
	config, configErr := config.ParseArgsReturnConfig(os.Args[1:])
	if configErr != nil {
		for _, err := range configErr {
			log.Error(err)
		}
		os.Exit(1)
	}
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
		if config.FlagTest {
			log.SetLevel(log.ErrorLevel)
		}
	}
	if err := executeRunner(config); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}
