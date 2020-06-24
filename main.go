package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {}

/*Init only initializes the http server for testing*/
func Init() {
	r := makeRouter()
	http.Handle("/", r)
}

func makeRouter() *mux.Router {
	r := mux.NewRouter()
	app := r.Headers("X-Requested-With", "XMLHttpRequest").Subrouter()
	app.HandleFunc("/user/{email}/{id}", restHandler).Methods(http.MethodGet)

	return r
}

func restHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]
	fmt.Fprintf(w, email)
}
