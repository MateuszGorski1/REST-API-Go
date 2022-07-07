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

func SumHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	sum := strconv.Itoa(int(a + b))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sum))
}

func DiffHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	diff := strconv.Itoa(int(a - b))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(diff))
}
func MulHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseInt(variables["a"], 10, 64)
	b, _ := strconv.ParseInt(variables["b"], 10, 64)
	mul := strconv.Itoa(int(a * b))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mul))
}
func DivHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	a, _ := strconv.ParseFloat(variables["a"], 64)
	b, _ := strconv.ParseFloat(variables["b"], 64)
	div := fmt.Sprintf("%f", a/b)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(div))
}

func FactorialHandler(w http.ResponseWriter, r *http.Request) {
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

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad Request"))
}

func StartServer() {
	logs.ServeLogs()
	wg := new(sync.WaitGroup)
	wg.Add(2)
	calculationEndpoints := mux.NewRouter()
	calculationEndpoints.HandleFunc("/{type:[a-z]+}", BadRequestHandler)
	calculationEndpoints.HandleFunc("/sum/{a:[0-9]+}/{b:[0-9]+}", SumHandler)
	calculationEndpoints.HandleFunc("/diff/{a:[0-9]+}/{b:[0-9]+}", DiffHandler)
	calculationEndpoints.HandleFunc("/mul/{a:[0-9]+}/{b:[0-9]+}", MulHandler)
	calculationEndpoints.HandleFunc("/div/{a:[0-9]+}/{b:[0-9]+}", DivHandler)
	calculationEndpoints.HandleFunc("/factorial/{a:[0-9]+}", FactorialHandler)
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(":8080", logs.C.Then(calculationEndpoints)); err != nil {
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
