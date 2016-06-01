package signal

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
)

// HealthReport defines the JSON received from the /system/health/report endpoint
// The health report returns keys that are not formatted for JSON specifically, so
// we do not modify them and instead pass the param as the key, unmodified.
type HealthReport struct {
	Units map[string]*Unit
	Nodes map[string]*Node
}

// Unit defines the JSON for the unit field in HealthReport
type Unit struct {
	UnitName  string
	Nodes     []*Node
	Health    int
	Title     string
	Timestamp time.Time
}

// Node defines the JSON for the node field in the HealthReport
type Node struct {
	Role   string
	Ip     string
	Host   string
	Health int
	Output map[string]string
	Units  []Unit
}

type Diagnostics struct {
	Report    *HealthReport
	Endpoints []string
	Method    string
	Headers   map[string]string
	Track     *analytics.Track
}

func (d *Diagnostics) SetReport(body []byte) error {
	if err := json.Unmarshal(body, &d.Report); err != nil {
		return err
	}
	return nil
}

func (d *Diagnostics) GetReport() interface{} {
	return d.Report
}

func (d *Diagnostics) SetHeaders(headers map[string]string) {
	d.Headers = headers
}

func (d *Diagnostics) GetHeaders() map[string]string {
	return d.Headers
}

func (d *Diagnostics) SetEndpoints(url []string) {
	d.Endpoints = url
}

func (d *Diagnostics) GetEndpoints() []string {
	return d.Endpoints
}

func (d *Diagnostics) SetMethod(method string) {
	d.Method = method
}

func (d *Diagnostics) GetMethod() string {
	return d.Method
}

func (d *Diagnostics) SetTrack(c config.Config) error {
	properties := map[string]interface{}{
		"source":             "cluster",
		"customerKey":        c.CustomerKey,
		"environmentVersion": c.DCOSVersion,
		"clusterId":          c.ClusterID,
		"variant":            c.DCOSVariant,
		"provider":           c.GenProvider,
	}

	for _, unit := range d.Report.Units {
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
			segmentUnitTotalKey := CreateUnitTotalKey(unit.UnitName)
			segmentUnitUnhealthyKey := CreateUnitUnhealthyKey(unit.UnitName)
			properties[segmentUnitTotalKey] = totalUnits
			properties[segmentUnitUnhealthyKey] = totalUnhealthyUnits
		}
	}
	d.Track = &analytics.Track{
		Event:       c.SegmentEvent,
		UserId:      c.CustomerKey,
		AnonymousId: c.ClusterID,
		Properties:  properties,
	}
	return nil
}

func (d *Diagnostics) GetTrack() *analytics.Track {
	return d.Track
}

func (d *Diagnostics) SendTrack(c config.Config) error {
	ac := CreateSegmentClient(c.SegmentKey, c.FlagVerbose)
	defer ac.Close()

	err := ac.Track(d.Track)
	return err
}
