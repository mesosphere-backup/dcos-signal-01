// +build integration

package signal

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"gopkg.in/segmentio/analytics-go.v2"
)

// TestSignalIntegration requires the SEGMENT_WRITE_KEY environment
// variable.
func TestSignalIntegration(t *testing.T) {

	m := Mesos{
		Track: &analytics.Track{
			Event:  "mesos_integration_test",
			UserId: uuid.New().String(),
			Properties: map[string]interface{}{
				"license_id": "local_integration_test",
			},
		},
	}

	client := analytics.New(os.Getenv("SEGMENT_WRITE_KEY"))
	client.Size = 1
	defer client.Close()

	err := client.Track(m.Track)
	if err != nil {
		t.Error("Failed to send track. Got error: ", err)
	}
}
