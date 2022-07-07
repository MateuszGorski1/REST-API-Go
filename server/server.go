package servers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/Icikowski/kubeprobes"
	"github.com/gorilla/mux"
)

var live = kubeprobes.NewStatefulProbe()
var ready = kubeprobes.NewStatefulProbe()
var kp = kubeprobes.New(
	kubeprobes.WithLivenessProbes(live.GetProbeFunction()),
	kubeprobes.WithReadinessProbes(ready.GetProbeFunction()),
)

func SumHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	sum := strconv.Itoa(int(a + b))
	w.Write([]byte(sum))
}
func DiffHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	diff := strconv.Itoa(int(a - b))
	w.Write([]byte(diff))
}
func MulHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	mul := strconv.Itoa(int(a * b))
	w.Write([]byte(mul))
}
func DivHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseFloat(variables["a"], 64)
	b, _ := strconv.ParseFloat(variables["b"], 64)
	div := fmt.Sprintf("%f", a/b)
	w.Write([]byte(div))
}
func FactorialHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseFloat(variables["a"], 64)
	b, _ := strconv.ParseFloat(variables["b"], 64)
	div := fmt.Sprintf("%f", a/b)
	w.Write([]byte(div))
}

func StartServer() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	calculationEndpoints := mux.NewRouter()
	calculationEndpoints.HandleFunc("/sum/{a:[0-9]+}/{b:[0-9]+}", SumHandler)
	calculationEndpoints.HandleFunc("/diff/{a:[0-9]+}/{b:[0-9]+}", DiffHandler)
	calculationEndpoints.HandleFunc("/mul/{a:[0-9]+}/{b:[0-9]+}", MulHandler)
	calculationEndpoints.HandleFunc("/div/{a:[0-9]+}/{b:[0-9]+}", DivHandler)
	calculationEndpoints.HandleFunc("/factorial/{a:[0-9]+}", FactorialHandler)
	go func() {
		defer wg.Done()
		http.ListenAndServe(":8080", calculationEndpoints)
	}()
	go func() {
		defer wg.Done()
		http.ListenAndServe(":8081", kp)
	}()
	_, err := http.Get("http://localhost:8080/sum/0/0")
	if err != nil {
		live.MarkAsDown()
		ready.MarkAsDown()
	} else {
		live.MarkAsUp()
		ready.MarkAsUp()
	}
	wg.Wait()

}
