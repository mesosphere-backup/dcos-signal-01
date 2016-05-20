package signal

import (
	"encoding/json"
	"time"
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
	Report *HealthReport
	URL    string
	Method string
}

func (d *Diagnostics) SetReport(body []byte) error {
	var hr *HealthReport
	if err := json.Unmarshal(body, &hr); err != nil {
		return err
	}
	d.Report = hr
	return nil
}

func (d *Diagnostics) GetReport() interface{} {
	return d.Report
}

func (d *Diagnostics) SetURL(url string) {
	d.URL = url
}

func (d *Diagnostics) GetURL(url string) string {
	return d.URL
}

func (d *Diagnostics) SetMethod(method string) {
	d.Method = method
}

func (d *Diagnostics) GetMethod() string {
	return d.Method
}
