package healthchecks

import (
	"fmt"
	"net/http"

	"github.com/Icikowski/kubeprobes"
)

var (
	LiveStatus  = kubeprobes.NewStatefulProbe()
	ReadyStatus = kubeprobes.NewStatefulProbe()
)

// PrepareHealthcheck returns server that check for app health
func PrepareHealthcheck() *http.Server {
	kp := kubeprobes.New(
		kubeprobes.WithLivenessProbes(LiveStatus.GetProbeFunction()),
		kubeprobes.WithReadinessProbes(ReadyStatus.GetProbeFunction()),
	)
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: kp,
	}
}
