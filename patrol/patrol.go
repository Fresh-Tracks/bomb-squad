package patrol

import (
	"log"
	"net/http"
	"time"

	configmap "github.com/Fresh-Tracks/bomb-squad/k8s/configmap"
	"github.com/Fresh-Tracks/bomb-squad/prom"
	promcfg "github.com/Fresh-Tracks/bomb-squad/prom/config"
)

var (
	iq prom.InstantQuery
)

type Patrol struct {
	PromURL           string
	Interval          time.Duration
	HighCardN         int
	HighCardThreshold float64
	Client            *http.Client
	ConfigMap         *configmap.ConfigMap
	PromConfig        *promcfg.Config
}

func (p *Patrol) Run() {
	//p.Bootstrap()
	ticker := time.NewTicker(time.Duration(p.Interval) * time.Second)
	for _ = range ticker.C {
		err := p.getTopCardinalities()
		if err != nil {
			log.Fatal(err)
		}
	}
}
