# DCOS Signal Service
The signal service is a passive data forwarding service for the system health API. The signal service acts as a middleware which runs GET requests to 3DT on masters every hour, forming a POST to send to SegementIO for our support team. 

## Go Requirements 
[SegmentIO](https://segment.com/docs/libraries/go/)
[Ginkgo](https://github.com/onsi/ginkgo)
[Logrus](https://github.com/Sirupsen/logrus)

## DCOS CLI Integration
TBD

## Build
```
go build dcos_signal_service.go segmentizer.go signaler.go config.go
```

## Run 
You can run an example with the ```run``` script:

```
./run
```

This will query a running 3DT environment and post the results to segmentIO. You can override the ENV parameters in this scirpt to test out different configurations and settings. 

## Test
Test suite is Ginkgo

1. Get Ginkgo 
```
go get github.com/onsi/ginkgo/ginkgo  # installs the ginkgo CLI
go get github.com/onsi/gomega         # fetches the matcher library
```

2. Install Ginkgo
```
cd /path/to/ginkgo/
go install
```

3. Test Signal Service
```
cd /path/to/dcos-signal/
ginkgo 
```

## CLI Arguments 
<pre>
Usage:
  -anonymous-id-path string
        Override path to DCOS anonymous ID. (default "/var/lib/dcos_anonymous_uuid.json")
  -c string
        Path to dcos-signal-service.conf. (default "/etc/dcos-signal-config.json")
  -ee
        Set the EE flag.
  -report-endpoint string
        Override default health endpoint. (default "/system/health/v1/report")
  -report-host string
        Override the default host to query the health report from. (default "localhost")
  -report-port int
        Override the default health API port. (default 1050)
  -v    Verbose logging mode.
  -version
        Print version and exit.
</pre>

## Local Testing
You can modify the run script and point it to any master where :1050 is exposed and available. 
