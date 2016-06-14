package signal

import (
	"github.com/dcos/dcos-signal/config"
)

func makeReporters(c config.Config) (chan Reporter, error) {

	var reporters = []Reporter{
		&Diagnostics{
			Name: "diagnostics",
			// Open Endpoints:
			//   :1050/system/health/v1/report
			// EE Endpoints:
			//   /system/health/v1/report
			Endpoints: c.DiagnosticsURLs,
			Method:    "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		},
		&Cosmos{
			Name: "cosmos",
			// Open Endpoints:
			//   :7070/package/list
			// EE Endpoints:
			//   /cosmos/package/list
			Endpoints: c.CosmosURLs,
			Method:    "POST",
			Headers: map[string]string{
				"content-type": "application/vnd.dcos.package.list-request",
			},
		},
		&Mesos{
			Name: "mesos",
			// Open Endpoints:
			//   :5050/frameworks, :5050/metrics/snapshot
			// EE Endpoints:
			//   /mesos/frameworks, /mesos/metrics/snapshot
			Endpoints: c.MesosURLs,
			Method:    "GET",
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
