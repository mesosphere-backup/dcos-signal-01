package signal

import (
	"github.com/dcos/dcos-signal/config"
)

func makeReporters(c config.Config) (chan Reporter, error) {

	var reporters = []Reporter{
		&Diagnostics{
			Name: "diagnostics",
			Endpoints: []string{
				":1050/system/health/v1/report",
			},
			Method: "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		},
		&Cosmos{
			Name: "cosmos",
			Endpoints: []string{
				":7070/package/list",
			},
			Method: "POST",
			Headers: map[string]string{
				"content-type": "application/vnd.dcos.package.list-request",
			},
		},
		&Mesos{
			Name: "mesos",
			Endpoints: []string{
				":5050/frameworks",
				":5050/metrics/snapshot",
			},
			Method: "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		},
	}

	reportChan := make(chan Reporter, len(reporters))
	for _, r := range reporters {
		r.addHeaders(c.ExtraHeaders)
		reportChan <- r
	}

	return reportChan, nil
}
