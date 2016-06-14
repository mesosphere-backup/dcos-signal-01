package signal

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/dcos/dcos-signal/config"
)

var (
	healthServer = httptest.NewServer(mockRouter())
	testDiag     = Diagnostics{
		Endpoints: []string{
			fmt.Sprintf("%s/system/health/v1/report", server.URL)},
		Method: "GET",
	}
)

func TestDiagnosticsTrack(t *testing.T) {
	c := config.DefaultConfig()
	c.CustomerKey = "12345"
	c.ClusterID = "anon"
	c.DCOSVersion = "test_version"
	c.GenProvider = "test_provider"
	c.DCOSVariant = "test_variant"

	pullErr := PullReport(&testDiag, c)
	if pullErr != nil {
		t.Error("Got error pulling from test server, ", pullErr)
	}

	setupErr := testDiag.setTrack(c)
	actualSegmentTrack := testDiag.getTrack()

	if setupErr != nil {
		t.Error("Expected no errors running diagnostics.SetTrack(), got ", setupErr)
	}

	if len(actualSegmentTrack.Properties) != 10 {
		t.Error("Expected 10 properties, got ", len(actualSegmentTrack.Properties))
	}

	if actualSegmentTrack.Event != "health" {
		t.Error("Expected actualSegmentTrack.Event to be 'health', got ", actualSegmentTrack.Event)
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

	if actualSegmentTrack.Properties["provider"] != "test_provider" {
		t.Error("Expected provder 'test_provider', got ", actualSegmentTrack.Properties["provider"])
	}

	if actualSegmentTrack.Properties["variant"] != "test_variant" {
		t.Error("Expected variant 'test_variant', got ", actualSegmentTrack.Properties["variant"])
	}

	if actualSegmentTrack.Properties["environmentVersion"] != "test_version" {
		t.Error("Expected environmenetVersion 'test_varsion', got ", actualSegmentTrack.Properties["environmentVersion"])
	}

	if _, ok := actualSegmentTrack.Properties["health-unit-foo-unit-2-total"]; !ok {
		t.Error("Expected key health-unit-foo-unit-2-total to exist, got ", ok)
	}

	if val, _ := actualSegmentTrack.Properties["health-unit-foo-unit-2-total"]; val != 2 {
		t.Error("Expected key health-unit-foo-unit-2-total to be 2, got ", val)
	}

	if _, ok := actualSegmentTrack.Properties["health-unit-foo-unit-2-unhealthy"]; !ok {
		t.Error("Expected key health-unit-foo-unit-2-unhealthy to exist, got ", ok)
	}

	if val, _ := actualSegmentTrack.Properties["health-unit-foo-unit-2-unhealthy"]; val != 1 {
		t.Error("Expected key health-unit-foo-unit-2-unhealthy to be 1, got ", val)
	}

	if _, ok := actualSegmentTrack.Properties["health-unit-foo-unit-1-total"]; !ok {
		t.Error("Expected key health-unit-foo-unit-1-total to exist, got ", ok)
	}

	if val, _ := actualSegmentTrack.Properties["health-unit-foo-unit-1-total"]; val != 2 {
		t.Error("Expected key health-unit-foo-unit-1-total to be 1, got ", val)
	}
	if _, ok := actualSegmentTrack.Properties["health-unit-foo-unit-1-unhealthy"]; !ok {
		t.Error("Expected key health-unit-foo-unit-1-unhealthy to exist, got ", ok)
	}

	if val, _ := actualSegmentTrack.Properties["health-unit-foo-unit-1-unhealthy"]; val != 0 {
		t.Error("Expected key health-unit-foo-unit-1-unhealthy to be 0, got ", val)
	}
}
