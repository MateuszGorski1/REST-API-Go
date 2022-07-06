package servers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	calculationEndpoints := mux.NewRouter()
	calculationEndpoints.HandleFunc("/sum/{a:[0-9]+}/{b:[0-9]+}", SumHandler)
	calculationEndpoints.HandleFunc("/diff/{a:[0-9]+}/{b:[0-9]+}", DiffHandler)
	calculationEndpoints.HandleFunc("/mul/{a:[0-9]+}/{b:[0-9]+}", MulHandler)
	calculationEndpoints.HandleFunc("/div/{a:[0-9]+}/{b:[0-9]+}", DivHandler)
	calculationEndpoints.HandleFunc("/factorial/{a:[0-9]+}", FactorialHandler)
	http.ListenAndServe(":8080", calculationEndpoints)

}
