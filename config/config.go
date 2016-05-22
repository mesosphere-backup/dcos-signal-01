package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

var VARIANT = "UNSET"

// Config defines dcos-signal configuration
type Config struct {
	// Service URLs
	DiagnosticsURL string `json:"diagnostics_url"`
	CosmosURL      string `json:"cosmos_url"`

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
	FlagVersion bool
	FlagVerbose bool
	FlagEE      bool
	TestFlag    bool
	Enabled     string `json:"enabled"`

	// Service account configuration
	ID         string `json:"uid"`
	SecretPath string `json:"secret_path"`
	Secret     string
	JWTToken   string
}

// DefaultConfig returns default Config{}
func DefaultConfig() Config {
	return Config{
		DiagnosticsURL: "http://localhost:1050/system/health/v1/report",
		CosmosURL:      "http://localhost:7070/package/list",

		SegmentEvent:            "health",
		SegmentKey:              "",
		CustomerKey:             "",
		ClusterID:               "",
		DCOSVersion:             os.Getenv("DCOS_VERSION"),
		DCOSClusterIDPath:       "/var/lib/dcos/cluster-id",
		FlagEE:                  false,
		DCOSVariant:             VARIANT,
		GenProvider:             "",
		SignalServiceConfigPath: "/opt/mesosphere/etc/dcos-signal-config.json",
		ExtraJSONConfigPath:     "/opt/mesosphere/etc/dcos-signal-extra.json",
		TestFlag:                false,
	}
}

func (c *Config) setFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.FlagVerbose, "v", c.FlagVerbose, "Verbose logging mode.")
	fs.BoolVar(&c.FlagVersion, "version", c.FlagVersion, "Print version and exit.")
	fs.StringVar(&c.DCOSClusterIDPath, "cluster-id-path", c.DCOSClusterIDPath, "Override path to DCOS anonymous ID.")
	fs.BoolVar(&c.FlagEE, "ee", c.FlagEE, "Set the EE flag.")
	fs.StringVar(&c.SignalServiceConfigPath, "c", c.SignalServiceConfigPath, "Path to dcos-signal-service.conf.")
	fs.BoolVar(&c.TestFlag, "test", c.TestFlag, "Test mode dumps a JSON object of the data that would be sent to Segment to STDOUT.")
	fs.StringVar(&c.SegmentKey, "segment-key", c.SegmentKey, "Key for segmentIO.")
}

func (c *Config) generateJWTToken() error {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["uid"] = c.ID
	token.Claims["exp"] = time.Now().Add(time.Hour).Unix()
	tokenStr, err := token.SignedString([]byte(c.Secret))
	if err != nil {
		return err
	}
	c.JWTToken = tokenStr
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
	// Attempt the load the secret file
	if len(c.SecretPath) > 0 {
		log.Warnf("Attempting to load secret file %s", c.SecretPath)
		if secretFile, err := ioutil.ReadFile(c.SecretPath); err != nil {
			return err
		} else {
			c.Secret = string(secretFile)
		}
	}
	return nil
}

func ParseArgsReturnConfig(args []string) (Config, []error) {
	errAry := []error{}
	c := DefaultConfig()
	signalFlag := flag.NewFlagSet("", flag.ContinueOnError)
	c.setFlags(signalFlag)
	if err := signalFlag.Parse(args); err != nil {
		errAry = append(errAry, err)
	}

	if err := c.getClusterID(); err != nil {
		c.ClusterID = err.Error()
		errAry = append(errAry, err)
	}

	if err := c.getExternalConfig(); err != nil {
		c.GenProvider = err.Error()
		c.CustomerKey = err.Error()
		errAry = append(errAry, err)
	}

	if c.FlagEE {
		c.DCOSVariant = "enterprise"
	}

	if len(c.ID) > 0 || len(c.Secret) > 0 {
		if err := c.generateJWTToken(); err != nil {
			errAry = append(errAry, err)
		}
	}

	if len(errAry) > 0 {
		return c, errAry
	}
	return c, nil
}
