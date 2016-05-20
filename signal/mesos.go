package signal

import (
	"encoding/json"

	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
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
}

func (d *Mesos) SetReport(body []byte) error {
	if err := json.Unmarshal(body, &d.Report); err != nil {
		return err
	}
	return nil
}

func (d *Mesos) GetReport() interface{} {
	return d.Report
}

func (d *Mesos) SetHeaders(headers map[string]string) {
	d.Headers = headers
}

func (d *Mesos) GetHeaders() map[string]string {
	return d.Headers
}

func (d *Mesos) SetEndpoints(url []string) {
	d.Endpoints = url
}

func (d *Mesos) GetEndpoints() []string {
	return d.Endpoints
}

func (d *Mesos) SetMethod(method string) {
	d.Method = method
}

func (d *Mesos) GetMethod() string {
	return d.Method
}

func (d *Mesos) SetTrack(c config.Config) error {
	properties := map[string]interface{}{
		"source":             "cluster",
		"customerKey":        c.CustomerKey,
		"environmentVersion": c.DCOSVersion,
		"clusterId":          c.ClusterID,
		"variant":            c.DCOSVariant,
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

func (d *Mesos) GetTrack() *analytics.Track {
	return d.Track
}

func (d *Mesos) SendTrack(c config.Config) error {
	ac := CreateSegmentClient(c.SegmentKey, c.FlagVerbose)
	defer ac.Close()
	err := ac.Track(d.Track)
	return err
}
