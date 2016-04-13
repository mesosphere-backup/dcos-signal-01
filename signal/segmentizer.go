package signal

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	config "github.com/mesosphere/dcos-signal/config"
	"github.com/segmentio/analytics-go"
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

// CreateSegmentTrack accepts a health report and transforms it into an analytics.Track
// to run our data to SegmentIO.
// We're not returning an error here since strict typeing in pulling the health and unmarshaling to
// the struct in PullHealthReport should be safe enough.
func CreateSegmentTrack(hr *HealthReport, c config.Config) (*analytics.Track, test) {
	properties := make(map[string]interface{})
	properties["source"] = "cluster"
	properties["customerKey"] = c.CustomerKey
	properties["environmentVersion"] = c.DCOSVersion
	properties["clusterId"] = c.ClusterID
	properties["variant"] = c.DCOSVariant
	properties["provider"] = c.GenProvider

	for _, unit := range hr.Units {
		totalUnits := len(unit.Nodes)
		totalUnhealthyUnits := 0
		for _, node := range unit.Nodes {
			// If the length of the output is greater than 0, then the unit can be considered unhealthy on that
			// specific node. As of writing this, we had no other way to determine by node how many unhealthy
			// units exist. This is because if any unit is unhealthy, units are unhealthy.
			if len(node.Output[unit.UnitName]) > 0 {
				log.Debug("==> UNHEALTHY HOST:")
				log.Debug(node.Output[unit.UnitName])
				totalUnhealthyUnits++
			}
		}
		segmentUnitTotalKey := CreateUnitTotalKey(unit.UnitName)
		segmentUnitUnhealthyKey := CreateUnitUnhealthyKey(unit.UnitName)
		properties[segmentUnitTotalKey] = totalUnits
		properties[segmentUnitUnhealthyKey] = totalUnhealthyUnits
	}

	var t = test{
		Event:      c.SegmentEvent,
		UserId:     c.CustomerKey,
		ClusterId:  c.ClusterID,
		Properties: properties,
	}
	pretty, _ := json.MarshalIndent(t, "", "    ")
	log.Debug("Data:\n", string(pretty))

	return &analytics.Track{
		Event:       c.SegmentEvent,
		UserId:      c.CustomerKey,
		AnonymousId: c.ClusterID,
		Properties:  properties,
	}, t
}
