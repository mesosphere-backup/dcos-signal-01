
# DCOS Signal Service
The signal service is a passive data forwarding service for the system health API. The signal service acts as a middleware which runs GET requests to 3DT on masters every hour, forming a POST to send to SegementIO for our support team. 

## Go Requirements 
[SegmentIO](https://segment.com/docs/libraries/go/)

## Build
```
make build
```

Or if you're running in an EE install

```
make build VARIANT=enterprise
``` 

## Run 
You can run an example with the ```run``` script:

```
dcos-signal
```

This will query a running 3DT environment and post the results to segmentIO. 

## Test

```
make test
```

## CLI Arguments 
<pre>
Usage:
  -anonymous-id-path string
        Override path to DCOS anonymous ID. (default "/var/lib/dcos_anonymous_uuid.json")
  -c string
        Path to dcos-signal-service.conf. (default "/etc/dcos-signal-config.json")
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

