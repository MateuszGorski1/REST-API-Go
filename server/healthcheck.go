package servers

import (
	"fmt"
	"net/http"

	"github.com/Icikowski/kubeprobes"
	"gorski.mateusz/webcalc/logs"
)

var (
	LiveStatus  = kubeprobes.NewStatefulProbe()
	ReadyStatus = kubeprobes.NewStatefulProbe()
)

func PrepareHealthcheck() *http.Server {
	kp := kubeprobes.New(
		kubeprobes.WithLivenessProbes(LiveStatus.GetProbeFunction()),
		kubeprobes.WithReadinessProbes(ReadyStatus.GetProbeFunction()),
	)
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: logs.C.Then(kp),
	}
}
