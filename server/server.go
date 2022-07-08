package servers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"gorski.mateusz/webcalc/logs"
)

func factorial(num int) int {
	if num == 1 || num == 0 {
		return num
	}
	return num * factorial(num-1)
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	sum := strconv.Itoa(int(a + b))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sum))
}

func diffHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	diff := strconv.Itoa(int(a - b))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(diff))
}
func mulHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	mul := strconv.Itoa(int(a * b))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mul))
}
func divHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseFloat(variables["a"], 64)
	b, _ := strconv.ParseFloat(variables["b"], 64)
	div := fmt.Sprintf("%f", a/b)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(div))
}

func factorialHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	var fact int
	a, _ := strconv.ParseInt(variables["a"], 10, 10)
	if a < 0 {
		fmt.Print("Factorial of negative number doesn't exist.")
	} else {
		fact = factorial(int(a))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(fact)))
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad Request"))
}

func badUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("503 Bad Url"))
}

func PrepareServer() *http.Server {
	calculationEndpoints := mux.NewRouter()
	calculationEndpoints.HandleFunc("/{type:[a-z]+}", badRequestHandler)
	calculationEndpoints.HandleFunc("/sum/{a:[0-9]+}/{b:[0-9]+}", sumHandler)
	calculationEndpoints.HandleFunc("/diff/{a:[0-9]+}/{b:[0-9]+}", diffHandler)
	calculationEndpoints.HandleFunc("/mul/{a:[0-9]+}/{b:[0-9]+}", mulHandler)
	calculationEndpoints.HandleFunc("/div/{a:[0-9]+}/{b:[0-9]+}", divHandler)
	calculationEndpoints.HandleFunc("/fact/{a:[0-9]+}", factorialHandler)
	calculationEndpoints.HandleFunc("/", badUrlHandler)
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: logs.C.Then(calculationEndpoints),
	}
}

func StartServer() {
	logs.ServeLogs()
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		server := PrepareServer()
		if err := server.ListenAndServe(); err != nil {
			LiveStatus.MarkAsDown()
			ReadyStatus.MarkAsDown()
		}
	}()
	go func() {
		defer wg.Done()
		healthServer := PrepareHealthcheck()
		if err := healthServer.ListenAndServe(); err != nil {
			LiveStatus.MarkAsDown()
			ReadyStatus.MarkAsDown()
		}
	}()

	LiveStatus.MarkAsUp()
	ReadyStatus.MarkAsUp()

	wg.Wait()

}
