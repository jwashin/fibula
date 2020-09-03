package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

// App is the center of it all. It connects http routing with a database
// in a REST configuration
type App struct {
	Router *mux.Router
	DB     *datastore.Client
}

// Initialize sets up the app
func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	c := context.Background()
	client, err := datastore.NewClient(c, "")
	if err != nil {
		fmt.Print("problem connecting to database")
	}
	a.DB = client
	a.initializeRoutes()
}

// Run runs the app
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getContest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idNum := vars["id"]
	var contest Contest

	if err := contest.getContest(a.DB, idNum); err != nil {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}
	respondWithJSON(w, http.StatusOK, contest)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getContests(w http.ResponseWriter, r *http.Request) {

	contests, err := getContests(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, contests)
}

func (a *App) createContest(w http.ResponseWriter, r *http.Request) {
	var e Contest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := e.createContest(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, e)
}

func (a *App) updateContest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idNum := vars["id"]
	var contest Contest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&contest); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := contest.updateContest(a.DB, idNum); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, contest)
}

func (a *App) deleteContest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idNum := vars["id"]
	var contest Contest
	if err := contest.deleteContest(a.DB, idNum); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/contests", a.getContests).Methods("GET")
	// a.Router.HandleFunc("/event/{id}", a.getEvent).Methods("GET")
	a.Router.HandleFunc("/contest", a.createContest).Methods("POST")
	a.Router.HandleFunc("/contest/{id:[0-9]+}", a.getContest).Methods("GET")
	a.Router.HandleFunc("/contest/{id:[0-9]+}", a.updateContest).Methods("PUT")
	a.Router.HandleFunc("/contest/{id:[0-9]+}", a.deleteContest).Methods("DELETE")
}
