package config

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Config defines dcos-signal configuration
type Config struct {
	// MasterURL is gained from execution of ip-detect
	MasterURL string
	// Make requests using HTTP or HTTPS
	TLSEnabled bool `json:"tls_enabled"`

	// CA Configuration for TLS requests
	CACertPath string `json:"ca_cert_path"`
	CAPool     *x509.CertPool

	// Segment IO Settings
	SegmentKey   string
	SegmentEvent string
	CustomerKey  string `json:"customer_key"`
	ClusterID    string `json:"cluster_id"`

	// DCOS-Specific Data
	DCOSVersion       string
	DCOSVariant       string
	GenProvider       string `json:"gen_provider"`
	DCOSClusterIDPath string

	// External Config Path Generated at Install Time
	SignalServiceConfigPath string
	ExtraJSONConfigPath     string

	// Optional CLI Flags
	FlagVersion  bool
	FlagVerbose  bool
	TestFlag     bool
	TestURL      string
	TestHost     string
	TestPort     int
	TestEndpoint string
	Enabled      string `json:"enabled"`

	// Extra headers for all reporter{}'s
	ExtraHeaders map[string]string

	// DC/OS Variant: enterprise or open
	Variant string
}

var (
	defaultConfig = Config{
		SegmentEvent:            "health",
		DCOSVersion:             os.Getenv("DCOS_VERSION"),
		DCOSClusterIDPath:       "/var/lib/dcos/cluster-id",
		DCOSVariant:             "open",
		SignalServiceConfigPath: "/opt/mesosphere/etc/dcos-signal-config.json",
		ExtraJSONConfigPath:     "/opt/mesosphere/etc/dcos-signal-extra.json",
		TestFlag:                false,
		TestURL:                 "http://localhost:444/tester",
		TestHost:                "localhost",
		TestPort:                4444,
		TestEndpoint:            "/tester",
		ExtraHeaders:            make(map[string]string),
	}
)

// DefaultConfig returns default Config{}
func DefaultConfig() Config {
	return defaultConfig
}

func (c *Config) setFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.FlagVerbose, "v", c.FlagVerbose, "Verbose logging mode.")
	fs.BoolVar(&c.FlagVersion, "version", c.FlagVersion, "Print version and exit.")
	fs.StringVar(&c.DCOSClusterIDPath, "cluster-id-path", c.DCOSClusterIDPath, "Override path to DCOS anonymous ID.")
	fs.StringVar(&c.SignalServiceConfigPath, "c", c.SignalServiceConfigPath, "Path to dcos-signal-service.conf.")
	fs.BoolVar(&c.TestFlag, "test", c.TestFlag, "Test mode dumps a JSON object of the data that would be sent to Segment to STDOUT.")
	fs.StringVar(&c.SegmentKey, "segment-key", c.SegmentKey, "Key for segmentIO.")
	fs.StringVar(&c.TestHost, "test-host", c.TestHost, "Host (IP or domain name) to POST test JSON")
	fs.IntVar(&c.TestPort, "test-port", c.TestPort, "Host port to POST test JSON")
	fs.StringVar(&c.TestEndpoint, "test-endpoint", c.TestEndpoint, "Host endpoint to POST to")
	fs.StringVar(&c.MasterURL, "master-url", c.MasterURL, "Override the IP to the master host")
	fs.StringVar(&c.TestURL, "test-url", c.TestURL, "Override default test URL")
}

func (c *Config) setMasterURL() error {
	log.Debug("Calculating Master URL")
	if c.MasterURL != "" {
		log.Debugf("MasterURL already set in memory, not regenerating: %s", c.MasterURL)
		return nil
	}

	detectIPCommand := os.Getenv("MESOS_IP_DISCOVERY_COMMAND")
	if detectIPCommand == "" {
		detectIPCommand = "/opt/mesosphere/bin/detect_ip"
		log.Debugf("Environment variable MESOS_IP_DISCOVERY_COMMAND is not set, using default location: %s", detectIPCommand)
	}

	out, err := exec.Command(detectIPCommand).Output()
	if err != nil {
		log.Errorf("Unable to execute %s: %s", detectIPCommand, err.Error())
		// Best effort, some services are still available from this URI
		out = []byte("localhost")
	}

	masterIP := strings.TrimRight(string(out), "\n")
	c.MasterURL = fmt.Sprintf("http://%s", masterIP)
	if !c.TLSEnabled {
		log.Debug("TLS disabled, protocol set to HTTP")
	} else {
		log.Debug("TLS enabled, protocol set to HTTPS")
		c.MasterURL = fmt.Sprintf("https://%s", masterIP)
	}

	log.Debugf("Master URL Set: %s", c.MasterURL)

	return nil
}

func (c *Config) getClusterID() error {
	fileByte, err := ioutil.ReadFile(c.DCOSClusterIDPath)
	if err != nil {
		return err
	}
	c.ClusterID = string(fileByte)
	return nil
}

func (c *Config) getExternalConfig() error {
	fileByte, err := ioutil.ReadFile(c.SignalServiceConfigPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(fileByte, &c); err != nil {
		return err
	}
	// Check for extra config and load if available
	if extraJSON, err := ioutil.ReadFile(c.ExtraJSONConfigPath); err == nil {
		if jsonErr := json.Unmarshal(extraJSON, &c); jsonErr != nil {
			return jsonErr
		}
	}
	return nil
}

func (c *Config) tryLoadingCert() error {
	// If no ca found, return nil.
	if c.CACertPath == "" {
		return nil
	}

	caPool := x509.NewCertPool()
	f, err := os.Open(c.CACertPath)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if !caPool.AppendCertsFromPEM(b) {
		return errors.New("CACertFile parsing failed")
	}
	c.CAPool = caPool
	return nil
}

func ParseArgsReturnConfig(args []string) (Config, []error) {
	errAry := []error{}
	c := DefaultConfig()
	signalFlag := flag.NewFlagSet("", flag.ContinueOnError)
	c.setFlags(signalFlag)

	// Parse CLI flags to override default config
	if err := signalFlag.Parse(args); err != nil {
		errAry = append(errAry, err)
	}

	// Get the cluster-id generate by ZK consensus
	if err := c.getClusterID(); err != nil {
		c.ClusterID = err.Error()
		errAry = append(errAry, err)
	}

	// Get standard and extra JSON config off disk
	if err := c.getExternalConfig(); err != nil {
		c.GenProvider = err.Error()
		c.CustomerKey = err.Error()
		errAry = append(errAry, err)
	}

	// Set the MasterURL from default or ip-detect
	if err := c.setMasterURL(); err != nil {
		errAry = append(errAry, err)
	}

	// Once all the config has been loaded, we can attempted to make a CAPool from the
	// path passed in config
	if err := c.tryLoadingCert(); err != nil {
		errAry = append(errAry, err)
	}

	if len(errAry) > 0 {
		return c, errAry
	}
	return c, nil
}
