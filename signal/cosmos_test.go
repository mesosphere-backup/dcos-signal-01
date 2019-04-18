// +build unit

package signal

import (
	"fmt"
	"testing"

	"encoding/json"

	"github.com/dcos/dcos-signal/config"
)

var (
	testCosmos = Cosmos{
		Endpoints: []string{
			fmt.Sprintf("%s/package/list", server.URL)},
		Method: "POST",
	}
)

func TestCosmosTrack(t *testing.T) {
	c := config.DefaultConfig()
	c.CustomerKey = "12345"
	c.ClusterID = "anon"
	c.LicenseID = "test_license"
	c.DCOSVersion = "test_version"
	c.GenPlatform = "test_platform"
	c.GenProvider = "test_provider"
	c.DCOSVariant = "test_variant"

	for _, endpoint := range testCosmos.Endpoints {
		pullErr := PullReport(endpoint, &testCosmos, c)
		if pullErr != nil {
			t.Error("Expected no errors pulling report from test server, got", pullErr)
		}
	}

	setupErr := testCosmos.setTrack(c)
	if setupErr != nil {
		t.Error("Expected no errors setting up track, got", setupErr)
	}

	actualSegmentTrack := testCosmos.getTrack()
	if actualSegmentTrack.Event != "package_list" {
		t.Error("Expected actualSegmentTrack.Event to be 'package_list', got ", actualSegmentTrack.Event)
	}

	if actualSegmentTrack.UserId != "12345" {
		t.Error("Expected actual segment track user ID to be 12345, got ", actualSegmentTrack.UserId)
	}

	if actualSegmentTrack.AnonymousId != "anon" {
		t.Error("Expected anon ID to be 'anon', got ", actualSegmentTrack.AnonymousId)
	}

	if actualSegmentTrack.Properties["clusterId"] != "anon" {
		t.Error("Expected clusterId to be anon, got ", actualSegmentTrack.Properties["clusterId"])
	}

	if actualSegmentTrack.Properties["source"] != "cluster" {
		t.Error("Expected source to be cluster, got ", actualSegmentTrack.Properties["source"])
	}

	if actualSegmentTrack.Properties["customerKey"] != "12345" {
		t.Error("Expected customerKey to be 12345, got ", actualSegmentTrack.Properties["customerKey"])
	}

	if actualSegmentTrack.Properties["licenseId"] != "test_license" {
		t.Error("Expected licenseId to be 'test_license', got ", actualSegmentTrack.Properties["licenseId"])
	}

	if actualSegmentTrack.Properties["platform"] != "test_platform" {
		t.Error("Expected platform 'test_platform', got ", actualSegmentTrack.Properties["platform"])
	}

	if actualSegmentTrack.Properties["provider"] != "test_provider" {
		t.Error("Expected provider 'test_provider', got ", actualSegmentTrack.Properties["provider"])
	}

	if actualSegmentTrack.Properties["variant"] != "test_variant" {
		t.Error("Expected variant 'test_variant', got ", actualSegmentTrack.Properties["variant"])
	}

	if actualSegmentTrack.Properties["environmentVersion"] != "test_version" {
		t.Error("Expected environmenetVersion 'test_varsion', got ", actualSegmentTrack.Properties["environmentVersion"])
	}

	if len(actualSegmentTrack.Properties["package_list"].([]CosmosPackages)) != 1 {
		t.Error("Expected 1 package in list, got", len(actualSegmentTrack.Properties["package_list"].([]CosmosPackages)))
	}

	if actualSegmentTrack.Properties["package_list"].([]CosmosPackages)[0].AppID != "fooPackage" {
		t.Error("Expected 'fooPackage', got", actualSegmentTrack.Properties["package_list"].([]CosmosPackages)[0].AppID)
	}
}

func TestCosmosStringer(t *testing.T) {
	report := CosmosReport{}
	response := `
	{
	  "packages": [
	    {
	      "packageInformation": {
	        "packageDefinition": {
	          "name": "hello-world",
	          "version": "0.0.1"
	        }
	      }
	    },
	    {
	      "packageInformation": {
	        "packageDefinition": {
	          "name": "test-pkg",
	          "version": "0.0.2"
	        }
	      }
	    }
	  ]
	}`

	if err := json.Unmarshal([]byte(response), &report); err != nil {
		t.Fatal(err)
	}

	expectedLine := "hello-world 0.0.1, test-pkg 0.0.2"
	if report.String() != expectedLine {
		t.Fatalf("Expect %s. Got %s", expectedLine, report)
	}
}
