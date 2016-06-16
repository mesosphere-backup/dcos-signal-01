# DCOS Signal Service [![velocity](http://velocity.mesosphere.com/service/velocity/buildStatus/icon?job=public-dcos-signal-service-master)](http://velocity.mesosphere.com/service/velocity/job/public-dcos-signal-service-master/)
The signal service is a passive data forwarding service for telemetry and analytics gathering. The signal service acts as a middleware which runs GET requests to 3DT, Mesos, and Cosmos on masters on a systemd timer.   

## SegmentIO Library Used 
[SegmentIO](https://segment.com/docs/libraries/go/)

## Build
```
make build
```

## Run
You can run an example with the ```run``` script:

```
dcos-signal
```

This will query a running 3DT environment and post the results to segmentIO.

## Test
Unittest:

```
make test
```

Integration test locally:

```
netcat -lk 4444 &
go run dcos_signal.go -v -test-url http://localhost:4444
```

Integration test in an actually integrated scenario:

```
ssh myuser@mymaster.com
netcat -lk 4444 &
./dcos_signal -v -test-url http://localhost:4444
```

## CLI Arguments
<pre>
Usage:
  -c                string | Path to dcos-signal-service.conf. (default "/opt/mesosphere/etc/dcos-signal-config.json")
  
  -cluster-id-path  string | Override path to DCOS anonymous ID. (default "/var/lib/dcos/cluster-id")
  
  -segment-key      string | Key for segmentIO.

  -test-url         string | URL to send would-be SegmentIO data to as JSON blob.
  
  -v                  bool | Verbose logging mode.
  
  -version            bool | Print version and exit.
</pre>
