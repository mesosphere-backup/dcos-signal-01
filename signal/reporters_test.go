package signal

import (
	"testing"

	"github.com/dcos/dcos-signal/config"
)

func TestMakeReporters(t *testing.T) {
	var (
		c       = config.DefaultConfig()
		r, rErr = makeReporters(c)
	)

	if rErr != nil {
		t.Error("Expected nil errors getting reporters, got", rErr)
	}

	if len(r) != 3 {
		t.Error("Expected 3 reporters, got", len(r))
	}
}
