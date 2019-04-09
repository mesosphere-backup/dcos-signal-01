// +build unit

package signal

import (
	"testing"
	"time"
)

func TestCreateUnitTotalKey(t *testing.T) {
	testTotalKey := CreateUnitTotalKey("foo")
	if testTotalKey != "health-unit-foo-total" {
		t.Error("Expected \"health-unit-foo-total\", got ", testTotalKey)
	}
}

func TestCreateUnitUnhealthyKey(t *testing.T) {
	testTotalKey := CreateUnitUnhealthyKey("foo")
	if testTotalKey != "health-unit-foo-unhealthy" {
		t.Error("Expected \"health-unit-foo-unhealthy\", got ", testTotalKey)
	}
}

func TestCreateSegmentClient(t *testing.T) {
	tc := CreateSegmentClient("12345", false)
	if tc.Size != 100 {
		t.Error("Expected 100, got ", tc.Size)
	}
	if tc.Interval != 30*time.Second {
		t.Error("Expected 30 seconds, got ", tc.Interval)
	}
	if tc.Verbose != false {
		t.Error("Expected false, got ", tc.Verbose)
	}
}
