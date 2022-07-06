package servers

import (
	"net/http"

	"github.com/Icikowski/kubeprobes"
)

func StartHealthcheck() {
	live := kubeprobes.NewStatefulProbe()
	ready := kubeprobes.NewStatefulProbe()
	_, err := http.Get("http://localhost:8080/sum/0/0")
	if err != nil {
		live.MarkAsDown()
		ready.MarkAsDown()
	} else {
		live.MarkAsUp()
		ready.MarkAsUp()
	}

	kp := kubeprobes.New(
		kubeprobes.WithLivenessProbes(live.GetProbeFunction()),
		kubeprobes.WithReadinessProbes(ready.GetProbeFunction()),
	)
	probes := &http.Server{
		Addr:    ":8081",
		Handler: kp,
	}
	probes.ListenAndServe()
}
