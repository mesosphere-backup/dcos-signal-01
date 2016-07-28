package signal

import (
	"encoding/json"
	"fmt"
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
	IP     string
	Host   string
	Health int
	Output map[string]string
	Units  []Unit
}

type Diagnostics struct {
	Report    *HealthReport
	Name      string
	Endpoints []string
	Method    string
	Headers   map[string]string
	Track     *analytics.Track
	Error     []string
}

func (d *Diagnostics) getName() string {
	return d.Name
}

func (d *Diagnostics) setReport(body []byte) error {
	if err := json.Unmarshal(body, &d.Report); err != nil {
		return err
	}
	return nil
}

func (d *Diagnostics) getReport() interface{} {
	return d.Report
}

func (d *Diagnostics) addHeaders(head map[string]string) {
	for k, v := range head {
		d.Headers[k] = v
	}
}

func (d *Diagnostics) getHeaders() map[string]string {
	return d.Headers
}

func (d *Diagnostics) getEndpoints() []string {
	if len(d.Endpoints) != 1 {
		log.Errorf("Diagnostics needs 1 endpoint, got %d", len(d.Endpoints))
	}
	return d.Endpoints
}

func (d *Diagnostics) getMethod() string {
	return d.Method
}

func (d *Diagnostics) getError() []string {
	return d.Error
}

func (d *Diagnostics) appendError(err string) {
	d.Error = append(d.Error, err)
}

func (d *Diagnostics) setTrack(c config.Config) error {
	properties := map[string]interface{}{
		"source":             "cluster",
		"customerKey":        c.CustomerKey,
		"environmentVersion": c.DCOSVersion,
		"clusterId":          c.ClusterID,
		"variant":            c.DCOSVariant,
		"provider":           c.GenProvider,
	}

	if d.Report == nil {
		return fmt.Errorf("%s is report is nil, bailing out.", d.Name)
	}

	for _, unit := range d.Report.Units {
		totalUnits := len(unit.Nodes)
		totalUnhealthyUnits := 0
		for _, node := range unit.Nodes {
			if node.Health != 0 {
				log.Debugf("UNHEALTHY NODE: %s", node.IP)
				totalUnhealthyUnits++
			} else {
				for _, nodeUnit := range node.Units {
					if unit.UnitName == nodeUnit.UnitName {
						if nodeUnit.Health != 0 {
							log.Debugf("UNHEALTHY UNIT: %s", node.Output[unit.UnitName])
							totalUnhealthyUnits++
						}
					}
				}
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

func (d *Diagnostics) getTrack() *analytics.Track {
	return d.Track
}

func (d *Diagnostics) sendTrack(c config.Config) error {
	ac := CreateSegmentClient(c.SegmentKey, c.FlagVerbose)
	defer ac.Close()

	err := ac.Track(d.Track)
	return err
}
