package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	expectedDefault := Config{
		DiagnosticsURL:          "http://localhost:1050/system/health/v1/report",
		CosmosURL:               "http://localhost:7070/package/list",
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
			"gen_provider": "onprem",
			"uid": "signal",
			"secret": "-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQCnNDnC7zjA5+3KgeiEjq6TbEYu3YKBLEWN60PF+KQQOIH10Wex0jS8O+sIFKfvKuALAOjA6GUwRFq8ZD9RfDz72MpRHtXsgYWQOWWnVa9OG9mgjJHzLtdKPYKjUVcjPf6QUNzJCZw+OIAKk16+5bJvEJJO5rt46zmkd4Gtaql8hQIDAQABAoGAYRY8K+qIA8soEhxYjQ/kYonOPsw0SRkR0hQ3qC515U1KeRf8pA4wvNP15x1HXeKBcSI4BDts9hfar+VttrzzE0E4pTsYgmIlloU6z1edjYOB6AIle7/1rPA/d5xt46XDEzCR55MHsIoQOYbEn5m91ME7ud9T0IZ4xAYvCg7Aq0kCQQDUJYK1+qtXeBUvXDlDioXdU83pQ0c6vcs/n4a6BNbsL+c/Biu3gWfYqrI3f8LA6YyUWAYriG2Wm3ZcSm4TIj7rAkEAycRr5aLGmitz8XJHDL4Mi8F/PwU5UwnzMwerLurb4dHv20XQxzRxv9OZTCDQzdhap36LjZBJDuM9MimCOde2TwJAZtqg2tXjiI7hxopyAPsCF+JvrK4/tI0cI4aWbU23Xd+DwByfyWJmFLf9m8bHh3wz+iALLcQBTcmlwu0bHQ+3bQJBAIXRWmZhQSs7Kpi2XF0dJyEB4q0ff9eNP9lWeriRV+g73sMlWMTmCZNaec+96/66QdXY3iGz0mCnYg0E7rQCV40CQQCLNT2vDBFgSH7dCUq+FRoDvASsBiqWGSE/njGABnl61BJNKH8jqmMcrmxMoNc9kQiYs5vQqiPtZxJQenUi7zL7\n-----END RSA PRIVATE KEY-----" 
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
	if c.ID != "signal" {
		t.Error("Expected ID, got ", c.ID)
	}
	if c.Secret != "-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQCnNDnC7zjA5+3KgeiEjq6TbEYu3YKBLEWN60PF+KQQOIH10Wex0jS8O+sIFKfvKuALAOjA6GUwRFq8ZD9RfDz72MpRHtXsgYWQOWWnVa9OG9mgjJHzLtdKPYKjUVcjPf6QUNzJCZw+OIAKk16+5bJvEJJO5rt46zmkd4Gtaql8hQIDAQABAoGAYRY8K+qIA8soEhxYjQ/kYonOPsw0SRkR0hQ3qC515U1KeRf8pA4wvNP15x1HXeKBcSI4BDts9hfar+VttrzzE0E4pTsYgmIlloU6z1edjYOB6AIle7/1rPA/d5xt46XDEzCR55MHsIoQOYbEn5m91ME7ud9T0IZ4xAYvCg7Aq0kCQQDUJYK1+qtXeBUvXDlDioXdU83pQ0c6vcs/n4a6BNbsL+c/Biu3gWfYqrI3f8LA6YyUWAYriG2Wm3ZcSm4TIj7rAkEAycRr5aLGmitz8XJHDL4Mi8F/PwU5UwnzMwerLurb4dHv20XQxzRxv9OZTCDQzdhap36LjZBJDuM9MimCOde2TwJAZtqg2tXjiI7hxopyAPsCF+JvrK4/tI0cI4aWbU23Xd+DwByfyWJmFLf9m8bHh3wz+iALLcQBTcmlwu0bHQ+3bQJBAIXRWmZhQSs7Kpi2XF0dJyEB4q0ff9eNP9lWeriRV+g73sMlWMTmCZNaec+96/66QdXY3iGz0mCnYg0E7rQCV40CQQCLNT2vDBFgSH7dCUq+FRoDvASsBiqWGSE/njGABnl61BJNKH8jqmMcrmxMoNc9kQiYs5vQqiPtZxJQenUi7zL7\n-----END RSA PRIVATE KEY-----" {
		t.Error("Expected secret, got", c.Secret)
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

func TestGenerateJWTToken(t *testing.T) {
	config := []byte(`
		{
			"customer_key": "foo-123-bar",
			"gen_provider": "onprem",
			"uid": "signal",
			"secret": "-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQCnNDnC7zjA5+3KgeiEjq6TbEYu3YKBLEWN60PF+KQQOIH10Wex0jS8O+sIFKfvKuALAOjA6GUwRFq8ZD9RfDz72MpRHtXsgYWQOWWnVa9OG9mgjJHzLtdKPYKjUVcjPf6QUNzJCZw+OIAKk16+5bJvEJJO5rt46zmkd4Gtaql8hQIDAQABAoGAYRY8K+qIA8soEhxYjQ/kYonOPsw0SRkR0hQ3qC515U1KeRf8pA4wvNP15x1HXeKBcSI4BDts9hfar+VttrzzE0E4pTsYgmIlloU6z1edjYOB6AIle7/1rPA/d5xt46XDEzCR55MHsIoQOYbEn5m91ME7ud9T0IZ4xAYvCg7Aq0kCQQDUJYK1+qtXeBUvXDlDioXdU83pQ0c6vcs/n4a6BNbsL+c/Biu3gWfYqrI3f8LA6YyUWAYriG2Wm3ZcSm4TIj7rAkEAycRr5aLGmitz8XJHDL4Mi8F/PwU5UwnzMwerLurb4dHv20XQxzRxv9OZTCDQzdhap36LjZBJDuM9MimCOde2TwJAZtqg2tXjiI7hxopyAPsCF+JvrK4/tI0cI4aWbU23Xd+DwByfyWJmFLf9m8bHh3wz+iALLcQBTcmlwu0bHQ+3bQJBAIXRWmZhQSs7Kpi2XF0dJyEB4q0ff9eNP9lWeriRV+g73sMlWMTmCZNaec+96/66QdXY3iGz0mCnYg0E7rQCV40CQQCLNT2vDBFgSH7dCUq+FRoDvASsBiqWGSE/njGABnl61BJNKH8jqmMcrmxMoNc9kQiYs5vQqiPtZxJQenUi7zL7\n-----END RSA PRIVATE KEY-----" 
		}`)
	tempConfig, _ := ioutil.TempFile(os.TempDir(), "")

	defer os.Remove(tempConfig.Name())

	tempConfig.Write(config)
	c := DefaultConfig()
	extErr := c.getExternalConfig()
	noErr := c.generateJWTToken()

	if extErr != nil {
		t.Error("Expected no errors loading external config, got", extErr)
	}
	if noErr != nil {
		t.Error("Expected no errors, got", noErr)
	}
	if len(c.JWTToken) == 0 {
		t.Error("Expected generated token, got", c.JWTToken)
	}

	c.Secret = "foobar"
	expectErr := c.generateJWTToken()

	if expectErr == nil {
		t.Error("Expected error loading bad PEM config, got", expectErr)
	}
}
