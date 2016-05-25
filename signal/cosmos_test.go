package signal

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/dcos/dcos-signal/config"
)

var (
	server     = httptest.NewServer(mockRouter())
	testCosmos = Cosmos{
		URL:    fmt.Sprintf("%s/package/list", server.URL),
		Method: "POST",
	}
)

func TestCosmosTrack(t *testing.T) {
	c := config.DefaultConfig()
	c.CustomerKey = "12345"
	c.ClusterID = "anon"
	c.DCOSVersion = "test_version"
	c.GenProvider = "test_provider"
	c.DCOSVariant = "test_variant"

	pullErr := PullReport(&testCosmos, c)
	if pullErr != nil {
		t.Error("Expected no errors pulling report from test server, got", pullErr)
	}

	setupErr := testCosmos.SetTrack(c)
	if setupErr != nil {
		t.Error("Expected no errors setting up track, got", setupErr)
	}

	actualTrack := testCosmos.GetTrack()

}
