package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

// App is the center of it all
type App struct {
	Router *mux.Router
	DB     *datastore.Client
}

// Initialize sets up the app
func (a *App) Initialize() {
	// connectionString :=
	// 	// fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	// 	fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbname)
	// var err error
	// a.DB, err = sql.Open("postgres", connectionString)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
func (a *App) Run(addr string) {}

func (a *App) getEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idNum, err := strconv.Atoi(vars["id"])
	// key, err := datastore.DecodeKey(id)
	// key := datastore.IDKey("Event", int64(idNum), nil)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Bad Key")
	}
	// event := Event{ID: key}
	var event Event

	if err = event.getEvent(a.DB, idNum); err != nil {
		respondWithError(w, http.StatusNotFound, "Event not found")
		return
	}
	respondWithJSON(w, http.StatusOK, event)
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

func (a *App) getEvents(w http.ResponseWriter, r *http.Request) {

	events, err := getEvents(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}

func (a *App) createEvent(w http.ResponseWriter, r *http.Request) {
	var e Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := e.createEvent(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, e)
}

func (a *App) updateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idNum, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}
	var event Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := event.updateEvent(a.DB, idNum); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, event)
}

func (a *App) deleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idNum, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Event ID")
		return
	}
	// key := datastore.IDKey("Event", int64(idNum), nil)
	// e := Event{ID: key}
	var event Event
	if err := event.deleteEvent(a.DB, idNum); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/events", a.getEvents).Methods("GET")
	// a.Router.HandleFunc("/event/{id}", a.getEvent).Methods("GET")
	a.Router.HandleFunc("/event", a.createEvent).Methods("POST")
	a.Router.HandleFunc("/event/{id:[0-9]+}", a.getEvent).Methods("GET")
	a.Router.HandleFunc("/event/{id:[0-9]+}", a.updateEvent).Methods("PUT")
	a.Router.HandleFunc("/event/{id:[0-9]+}", a.deleteEvent).Methods("DELETE")
}
