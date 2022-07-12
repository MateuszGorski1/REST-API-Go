package servers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"gorski.mateusz/webcalc/healthchecks"
	"gorski.mateusz/webcalc/logs"
)

func factorial(num int) int {
	if num == 1 || num == 0 {
		return num
	}
	return num * factorial(num-1)
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	var sum string
	variables := mux.Vars(r)
	a, err1 := strconv.ParseFloat(variables["a"], 64)
	b, err2 := strconv.ParseFloat(variables["b"], 64)
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		sum = fmt.Sprintf("%f", a+b)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sum))
	}
}

func diffHandler(w http.ResponseWriter, r *http.Request) {
	var diff string
	variables := mux.Vars(r)
	a, err1 := strconv.ParseFloat(variables["a"], 64)
	b, err2 := strconv.ParseFloat(variables["b"], 64)
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		diff = fmt.Sprintf("%f", a-b)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(diff))
	}
}
func mulHandler(w http.ResponseWriter, r *http.Request) {
	var mul string
	variables := mux.Vars(r)
	a, err1 := strconv.ParseFloat(variables["a"], 64)
	b, err2 := strconv.ParseFloat(variables["b"], 64)
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		mul = fmt.Sprintf("%f", a*b)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mul))
	}
}
func divHandler(w http.ResponseWriter, r *http.Request) {
	var div string
	variables := mux.Vars(r)
	a, err1 := strconv.ParseFloat(variables["a"], 64)
	b, err2 := strconv.ParseFloat(variables["b"], 64)
	if err1 != nil || err2 != nil || b == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		div = fmt.Sprintf("%f", a/b)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(div))
	}
}

func factorialHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	var fact int
	a, err1 := strconv.ParseInt(variables["a"], 10, 10)
	if err1 != nil || a < 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		fact = factorial(int(a))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(fact)))
	}
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

// PrepareServer returns server listening on port 8080 which handles calculator
func PrepareServer() *http.Server {
	calculationEndpoints := mux.NewRouter()
	calculationEndpoints.HandleFunc("/sum/{a}/{b}", sumHandler)
	calculationEndpoints.HandleFunc("/diff/{a}/{b}", diffHandler)
	calculationEndpoints.HandleFunc("/mul/{a}/{b}", mulHandler)
	calculationEndpoints.HandleFunc("/div/{a}/{b}", divHandler)
	calculationEndpoints.HandleFunc("/fact/{a}", factorialHandler)
	calculationEndpoints.NotFoundHandler = http.HandlerFunc(badRequestHandler)
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: logs.C.Then(calculationEndpoints),
	}
}

// StartServer runs server and its healthcheck in separate go routines
func StartServer() {
	logs.ServeLogs()
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		server := PrepareServer()
		if err := server.ListenAndServe(); err != nil {
			healthchecks.LiveStatus.MarkAsDown()
			healthchecks.ReadyStatus.MarkAsDown()
		}
	}()
	go func() {
		defer wg.Done()
		healthServer := healthchecks.PrepareHealthcheck()
		if err := healthServer.ListenAndServe(); err != nil {
			healthchecks.LiveStatus.MarkAsDown()
			healthchecks.ReadyStatus.MarkAsDown()
		}
	}()
	healthchecks.LiveStatus.MarkAsUp()
	healthchecks.ReadyStatus.MarkAsUp()

	wg.Wait()
}
