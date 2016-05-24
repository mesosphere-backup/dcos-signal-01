package signal

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
)

type HealthReport struct {
	Units map[string]*Unit
	Nodes map[string]*Node
}
