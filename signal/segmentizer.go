package signal

import (
	"strings"
	"time"

	"gopkg.in/segmentio/analytics-go.v2"
)

// CreateUnitTotalKey creates the key for segmentIO properties for total hosts. This key
// has the format: health-unit-$UNIT_ID-total
func CreateUnitTotalKey(name string) string {
	return "health-unit-" + strings.Replace(name, ".", "-", -1) + "-total"
}

// CreateUnitUnhealthyKey creates the key for segemntIO properties for unhealthy hosts. This
// key has the format: health-units-$UNIT_ID-unhealthy
func CreateUnitUnhealthyKey(name string) string {
	return "health-unit-" + strings.Replace(name, ".", "-", -1) + "-unhealthy"
}

// CreateSegmentClient returns our specific client implementation
func CreateSegmentClient(segmentKey string, verbose bool) *analytics.Client {
	client := analytics.New(segmentKey)
	client.Interval = 30 * time.Second
	client.Size = 100
	client.Verbose = verbose
	return client
}
