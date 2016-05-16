package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	expectedDefault := Config{
		HealthAPIPort:           1050,
		HealthEndpoint:          "/system/health/v1/report",
		HealthHost:              "localhost",
		SegmentEvent:            "health",
		SegmentKey:              "",
		CustomerKey:             "",
		ClusterID:               "",
		DCOSVersion:             "",
		DCOSClusterIDPath:       "/var/lib/dcos/cluster-id",
		SignalServiceConfigPath: "/opt/mesosphere/etc/dcos-signal-config.json",
		ExtraJSONConfigPath:     "/opt/mesosphere/etc/dcos-signal-extra.json",
		FlagEE:                  false,
		DCOSVariant:             "UNSET",
		GenProvider:             "",
	}
	gotDefault := DefaultConfig()
	if gotDefault != expectedDefault {
		t.Error("Expected ", expectedDefault, ", got ", gotDefault)
	}
}

func TestFlagParsing(t *testing.T) {
	json := []byte(`{"cluster_uuid": "12345"}`)
	tempAnonJson, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(tempAnonJson.Name())
	tempAnonJson.Write(json)

	config := []byte(`{"customer_key": "someuser-enterprise-key"}`)
	tempConfig, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(tempConfig.Name())
	tempConfig.Write(config)

	var (
		verboseConfig, verboseErr = ParseArgsReturnConfig([]string{
			"-v",
			"-cluster-id-path", tempAnonJson.Name(),
			"-c", tempConfig.Name()})

		versionConfig, versionErr = ParseArgsReturnConfig([]string{
			"-version",
			"-cluster-id-path", tempAnonJson.Name(),
			"-c", tempConfig.Name()})

		reportHostConfig, hostErr = ParseArgsReturnConfig([]string{
			"-report-host", "10.0.0.1",
			"-cluster-id-path", tempAnonJson.Name(),
			"-c", tempConfig.Name()})

		reportPortConfig, portErr = ParseArgsReturnConfig([]string{
			"-report-port", "8080",
			"-cluster-id-path", tempAnonJson.Name(),
			"-c", tempConfig.Name()})

		testConfig, testConfigErr = ParseArgsReturnConfig([]string{
			"-test",
			"-cluster-id-path", tempAnonJson.Name(),
			"-c", tempConfig.Name()})

		testNoFile, noFileErr = ParseArgsReturnConfig([]string{})
	)
	// Test No Config Files (anon ID or config.json)
	if testNoFile.ClusterID != "open /var/lib/dcos/cluster-id: no such file or directory" {
		t.Error("Expected 'open /var/lib/dcos/cluster-id: no such file or directory', got ", testNoFile.ClusterID)
	}
	if noFileErr == nil {
		t.Error("Expected error with no config, got ", noFileErr)
	}

	// -test
	if testConfig.TestFlag != true {
		t.Error("Expected test flag to be true, got ", testConfig.TestFlag)
	}
	if testConfigErr != nil {
		t.Error("Expected test error to be nil, got ", testConfigErr)
	}

	// -v
	if verboseConfig.FlagVerbose != true {
		t.Error("Expected true, got ", verboseConfig.FlagVerbose)
	}
	if verboseErr != nil {
		t.Error("Expected nil, got ", verboseErr)
	}

	// -version
	if versionConfig.FlagVersion != true {
		t.Error("Expected true, got ", versionConfig.FlagVersion)
	}
	if versionErr != nil {
		t.Error("Expected nil, got ", versionErr)
	}

	// -report-host
	if reportHostConfig.HealthHost != "10.0.0.1" {
		t.Error("Expected 10.0.0.1, got ", reportHostConfig.HealthHost)
	}
	if hostErr != nil {
		t.Error("Expected nil, got ", hostErr)
	}

	// -report-port
	if reportPortConfig.HealthAPIPort != 8080 {
		t.Error("Expected 8080, got ", reportPortConfig.HealthAPIPort)
	}
	if portErr != nil {
		t.Error("Expected nil, got ", portErr)
	}
}

func TestGetClusterID(t *testing.T) {
	file := []byte(`12345`)
	tempAnon, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(tempAnon.Name())

	tempAnon.Write(file)
	c := DefaultConfig()
	c.DCOSClusterIDPath = tempAnon.Name()

	if err := c.getClusterID(); err != nil {
		t.Error("Expected no errors from getClusterID(), got ", err)
	}

	if c.ClusterID != "12345" {
		t.Error("Expected cluster ID to be 12345, got ", c.ClusterID)
	}
}

func TestGetExternalConfig(t *testing.T) {
	// Test basic config
	config := []byte(`
		{
			"customer_key": "foo-123-bar",
			"gen_provider": "onprem"	
		}`)
	tempConfig, _ := ioutil.TempFile(os.TempDir(), "")

	defer os.Remove(tempConfig.Name())

	tempConfig.Write(config)
	c := DefaultConfig()
	c.SignalServiceConfigPath = tempConfig.Name()

	if err := c.getExternalConfig(); err != nil {
		t.Error("Expected config, got ", err)
	}

	if c.CustomerKey != "foo-123-bar" {
		t.Error("Expected customer ID to be foo-123-bar, got ", c.CustomerKey)
	}
	if c.GenProvider != "onprem" {
		t.Error("Expected onprem, got ", c.GenProvider)
	}

	// Test no enterprise
	noEntConfig := []byte(`
		{
			"customer_key": "",
			"gen_provider": "onprem"	
		}`)
	noEntFile, _ := ioutil.TempFile(os.TempDir(), "")

	defer os.Remove(noEntFile.Name())

	noEntFile.Write(noEntConfig)
	noEntC := DefaultConfig()
	noEntC.SignalServiceConfigPath = noEntFile.Name()

	if err := noEntC.getExternalConfig(); err != nil {
		t.Error("Expected no errors with no enterprise in config, got ", err)
	}

	if noEntC.CustomerKey != "" {
		t.Error("Expected customer ID to be empty, got ", noEntC.CustomerKey)
	}

	if noEntC.GenProvider != "onprem" {
		t.Error("Expected onprem, got ", noEntC.GenProvider)
	}

}
