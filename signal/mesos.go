package signal

import (
	"encoding/json"
	"fmt"

	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"

	log "github.com/Sirupsen/logrus"
)

// Complete report used by signal service, composed of all requests
type MesosReport struct {
	Frameworks      []Framework `json:"frameworks"`
	CPUTotal        float64     `json:"master/cpus_total"`
	CPUUsed         float64     `json:"master/cpus_used"`
	DiskTotal       float64     `json:"master/disk_total"`
	DiskUsed        float64     `json:"master/disk_used"`
	MemTotal        float64     `json:"master/mem_total"`
	MemUsed         float64     `json:"master/mem_used"`
	TaskCount       float64     `json:"master/tasks_running"`
	FrameworkCount  float64     `json:"master/frameworks_active"`
	AgentsConnected float64     `json:"master/slaves_connected"`
	AgentsActive    float64     `json:"master/slaves_active"`
}

type Framework struct {
	Name string `json:"name"`
}

type Mesos struct {
	Report    *MesosReport
	Endpoints []string
	Method    string
	Headers   map[string]string
	Track     *analytics.Track
	Error     []string
	Name      string
}

func (d *Mesos) getName() string {
	return d.Name
}

func (d *Mesos) setReport(body []byte) error {
	if err := json.Unmarshal(body, &d.Report); err != nil {
		return err
	}
	return nil
}

func (d *Mesos) getReport() interface{} {
	return d.Report
}

func (d *Mesos) addHeaders(head map[string]string) {
	for k, v := range head {
		d.Headers[k] = v
	}
}
func (d *Mesos) getHeaders() map[string]string {
	return d.Headers
}

func (d *Mesos) getEndpoints() []string {
	if len(d.Endpoints) != 2 {
		log.Errorf("Mesos needs 2 endpoints, got %d", len(d.Endpoints))
	}
	return d.Endpoints
}

func (d *Mesos) getMethod() string {
	return d.Method
}

func (d *Mesos) getError() []string {
	return d.Error
}

func (d *Mesos) appendError(err string) {
	d.Error = append(d.Error, err)
}

func (d *Mesos) setTrack(c config.Config) error {
	if d.Report == nil {
		return fmt.Errorf("%s report is nil, bailing out.", d.Name)
	}

	properties := map[string]interface{}{
		"source":             "cluster",
		"customerKey":        c.CustomerKey,
		"environmentVersion": c.DCOSVersion,
		"clusterId":          c.ClusterID,
		"licenseId":          c.LicenseID,
		"variant":            c.DCOSVariant,
		"platform":           c.GenPlatform,
		"provider":           c.GenProvider,
		"frameworks":         d.Report.Frameworks,
		"cpu_total":          d.Report.CPUTotal,
		"cpu_used":           d.Report.CPUUsed,
		"mem_total":          d.Report.MemTotal,
		"mem_used":           d.Report.MemUsed,
		"disk_total":         d.Report.DiskTotal,
		"disk_used":          d.Report.DiskUsed,
		"task_count":         d.Report.TaskCount,
		"framework_count":    d.Report.FrameworkCount,
		"agents_connected":   d.Report.AgentsConnected,
		"agents_active":      d.Report.AgentsActive,
	}

	d.Track = &analytics.Track{
		Event:       "mesos_track",
		UserId:      c.CustomerKey,
		AnonymousId: c.ClusterID,
		Properties:  properties,
	}
	return nil
}

func (d *Mesos) getTrack() *analytics.Track {
	return d.Track
}

func (d *Mesos) sendTrack(c config.Config) error {
	ac := CreateSegmentClient(c.SegmentKey, c.FlagVerbose)
	defer ac.Close()
	err := ac.Track(d.Track)
	return err
}
