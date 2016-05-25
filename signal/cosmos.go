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
	Report    *CosmosReport
	Endpoints []string
	Method    string
	Headers   map[string]string
	Track     *analytics.Track
}

func (c *Cosmos) SetReport(body []byte) error {
	if err := json.Unmarshal(body, &c.Report); err != nil {
		return err
	}
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

func (c *Cosmos) SetEndpoints(url []string) {
	c.Endpoints = url
}

func (c *Cosmos) GetEndpoints() []string {
	return c.Endpoints
}

func (c *Cosmos) SetMethod(method string) {
	c.Method = method
}

func (c *Cosmos) GetMethod() string {
	return c.Method
}

func (c *Cosmos) SetTrack(config config.Config) error {
	properties := map[string]interface{}{
		"package_list":       c.Report.Packages,
		"source":             "cluster",
		"customerKey":        config.CustomerKey,
		"environmentVersion": config.DCOSVersion,
		"clusterId":          config.ClusterID,
		"variant":            config.DCOSVariant,
		"provider":           config.GenProvider,
	}

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
	ac := CreateSegmentClient(config.SegmentKey, config.FlagVerbose)
	defer ac.Close()
	err := ac.Track(c.Track)
	return err
}
