# DCOS Signal Service [![Build Status](https://jenkins.mesosphere.com/service/jenkins/buildStatus/icon?job=public-dcos-cluster-ops/dcos-signal/add-standard-jenkinsfile)](https://jenkins.mesosphere.com/service/jenkins/job/public-dcos-cluster-ops/job/dcos-signal/job/add-standard-jenkinsfile/)
The signal service is a passive data forwarding service for telemetry and analytics gathering. The signal service acts as a middleware which runs GET requests to 3DT, Mesos, and Cosmos on masters on a systemd timer.   

## SegmentIO Library Used
[SegmentIO](https://segment.com/docs/libraries/go/)

## Build
```
make build
```

## Test

#### Unit Tests:

```
make unit
```

#### Integration Tests:

NOTE: Some integration tests require the `SEGMENT_WRITE_KEY` to be set in order to integrate with Segment. The `dcos-dev` segment write key can be found [here](https://mesosphere.onelogin.com/notes/71322). With the key set, you should be able to view new "mesos_integration_test" events [here](https://app.segment.com/mesosphere/sources/dcos-dev/debugger) in Segment.

For access to segment.com, please contact Christopher Gutierrez on the analytics team.

```
make integration
```

#### All Tests

```
make test
```

#### Local End-to-End Tests:

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
