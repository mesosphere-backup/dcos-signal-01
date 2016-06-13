package signal

import (
	"github.com/dcos/dcos-signal/config"
)

func makeReporters(c config.Config) (chan Reporter, error) {

	var reporters = []Reporter{
		&Diagnostics{
			Name: "diagnostics",
			Endpoints: []string{
				"/system/health/v1/report",
			},
			Method: "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		},
		&Cosmos{
			Name: "cosmos",
			Endpoints: []string{
				"/cosmos/package/list",
			},
			Method: "POST",
			Headers: map[string]string{
				"content-type": "application/vnd.dcos.package.list-request",
			},
		},
		&Mesos{
			Name: "mesos",
			Endpoints: []string{
				"/mesos/frameworks",
				"/mesos/metrics/snapshot",
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
