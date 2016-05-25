package signal

import (
	"encoding/json"

	//	log "github.com/Sirupsen/logrus"
	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
)

type CosmosPackages struct {
	AppID string `json:"appId"`
}

type CosmosReport struct {
	Packages []CosmosPackages `json:"packages"`
}

type Cosmos struct {
	Report  *CosmosReport
	URL     string
	Method  string
	Headers map[string]string
	Track   *analytics.Track
}

func (c *Cosmos) SetReport(body []byte) error {
	var report *CosmosReport
	if err := json.Unmarshal(body, &report); err != nil {
		return err
	}
	c.Report = report
	return nil
}

func (c *Cosmos) GetReport() interface{} {
	return c.Report
}

func (c *Cosmos) SetHeaders(headers map[string]string) {
	c.Headers = headers
}

func (c *Cosmos) GetHeaders() map[string]string {
	return c.Headers
}

func (c *Cosmos) SetURL(url string) {
	c.URL = url
}

func (c *Cosmos) GetURL() string {
	return c.URL
}

func (c *Cosmos) SetMethod(method string) {
	c.Method = method
}

func (c *Cosmos) GetMethod() string {
	return c.Method
}

func (c *Cosmos) SetTrack(config config.Config) error {
	properties := make(map[string]interface{})
	properties["package_list"] = c.Report.Packages
	properties["source"] = "cluster"
	properties["customerKey"] = config.CustomerKey
	properties["environmentVersion"] = config.DCOSVersion
	properties["clusterId"] = config.ClusterID
	properties["variant"] = config.DCOSVariant
	properties["provider"] = config.GenProvider

	c.Track = &analytics.Track{
		Event:       "package_list",
		UserId:      config.CustomerKey,
		AnonymousId: config.ClusterID,
		Properties:  properties,
	}
	return nil
}

func (c *Cosmos) GetTrack() *analytics.Track {
	return c.Track
}

func (c *Cosmos) SendTrack(config config.Config) error {
	return nil
}
