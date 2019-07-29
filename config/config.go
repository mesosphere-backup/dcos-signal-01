package config

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type DCOSVariant struct {
	Name string
}

func (v DCOSVariant) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", v.Name)), nil
}

func (v DCOSVariant) String() string {
	return v.Name
}

func (v *DCOSVariant) Set(variant string) error {
	if variant == "enterprise" {
		initEnterprise()
	} else if variant != "open" {
		return fmt.Errorf("unknown variant '%s'. Only 'open' or 'enterprise' are allowed", variant)
	}
	v.Name = variant
	return nil
}

// Config defines dcos-signal configuration
type Config struct {
	// URL Configuration for Reports
	DiagnosticsURLs []string `json:"diagnostics_urls"`
	CosmosURLs      []string `json:"cosmos_urls"`
	MesosURLs       []string `json:"mesos_urls"`

	// CA Configuration for TLS requests
	CACertPath string `json:"ca_cert_path"`
	CAPool     *x509.CertPool

	// Segment IO Settings
	SegmentKey   string
	SegmentEvent string
	CustomerKey  string `json:"customer_key"`
	ClusterID    string `json:"cluster_id"`
	LicenseID    string `json:"license_id"`

	// DCOS-Specific Data
	DCOSVersion       string
	DCOSVariant       DCOSVariant
	GenPlatform       string `json:"gen_platform"`
	GenProvider       string `json:"gen_provider"`
	DCOSClusterIDPath string

	// External Config Path Generated at Install Time
	LicensingSocket         string
	SignalServiceConfigPath string
	ExtraJSONConfigPath     string

	// Optional CLI Flags
	FlagVersion bool
	FlagVerbose bool
	FlagTest    bool
	Enabled     string `json:"enabled"`

	// Extra headers for all reporter{}'s
	ExtraHeaders map[string]string
}

var (
	defaultConfig = Config{
		SegmentEvent:            "health",
		DCOSVersion:             os.Getenv("DCOS_VERSION"),
		DCOSClusterIDPath:       "/var/lib/dcos/cluster-id",
		DCOSVariant:             DCOSVariant{"open"},
		LicensingSocket:         "/tmp/dcos-licensing.socket",
		SignalServiceConfigPath: "/opt/mesosphere/etc/dcos-signal-config.json",
		ExtraJSONConfigPath:     "/opt/mesosphere/etc/dcos-signal-extra.json",
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
	fs.StringVar(&c.LicensingSocket, "licensing-socket", c.LicensingSocket, "Path to licensing socket.")
	fs.StringVar(&c.SignalServiceConfigPath, "c", c.SignalServiceConfigPath, "Path to dcos-signal-service.conf.")
	fs.StringVar(&c.SegmentKey, "segment-key", c.SegmentKey, "Key for segmentIO.")
	fs.BoolVar(&c.FlagTest, "test", c.FlagTest, "Dump the data sent to segment to stdout.")
	fs.Var(&c.DCOSVariant, "dcos-variant", "Variant of DC/OS ('open' or 'enterprise')")
}

func (c *Config) getLicenseID() error {
	// Build an http client that connects via unix domain socket
	httpc := http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, "unix", c.LicensingSocket)
			},
		},
	}

	// Call the /licenses endpoint on the dcos-licensing service
	resp, err := httpc.Get("http://unix/licenses")
	if err != nil {
		return err
	}

	// Response from /licenses endpoint of dcos-licensing service
	licenses := []struct {
		ID            string `json:"id"`
		Version       string `json:"version"`
		DecryptionKey string `json:"decryption_key"`
		LicenseTerms  struct {
			NodeCapacity   int       `json:"node_capacity"`
			StartTimestamp time.Time `json:"start_timestamp"`
			EndTimestamp   time.Time `json:"end_timestamp"`
		} `json:"license_terms"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&licenses); err != nil {
		return err
	}

	// Loop through licenses and find the one with the latest expirary
	// date. That is probably the valid license. All we care about is
	// the license_id, which is probably the same for all licenses on
	// the cluster. But this is an extra precaution to try to determine
	// which of the licenses is valid.
	if len(licenses) >= 1 {
		id := licenses[0].ID
		end := licenses[0].LicenseTerms.EndTimestamp
		for _, l := range licenses {
			if l.LicenseTerms.EndTimestamp.After(end) {
				id = l.ID
				end = l.LicenseTerms.EndTimestamp
			}
		}
		c.LicenseID = id
	}

	return nil
}

func (c *Config) getClusterID() error {
	fileByte, err := ioutil.ReadFile(c.DCOSClusterIDPath)
	if err != nil {
		return err
	}
	c.ClusterID = strings.TrimSpace(string(fileByte))
	log.Debugf("Detected Cluster ID: %s", c.ClusterID)
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

// ParseArgsReturnConfig does exactly that
func ParseArgsReturnConfig(args []string) (Config, []error) {
	errAry := []error{}
	c := DefaultConfig()
	signalFlag := flag.NewFlagSet("", flag.ContinueOnError)
	c.setFlags(signalFlag)

	// Parse CLI flags to override default config
	if err := signalFlag.Parse(args); err != nil {
		errAry = append(errAry, err)
	}

	// Not all clusters will have a license, including open source clusters.
	if err := c.getLicenseID(); err != nil {
		log.Errorf("error getting LicenseID. Got error: %v", err)
	}

	// Get the cluster-id generate by ZK consensus
	if err := c.getClusterID(); err != nil {
		// If cluster-id is not found signal service should fail.
		errAry = append(errAry, err)
		return c, errAry
	}

	// Get standard and extra JSON config off disk
	if err := c.getExternalConfig(); err != nil {
		c.GenPlatform = err.Error()
		c.GenProvider = err.Error()
		c.CustomerKey = err.Error()
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
