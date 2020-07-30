// event_test.go

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"cloud.google.com/go/datastore"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize()
	code := m.Run()
	// clearTable()
	os.Exit(code)
}

func clearTempEvents() {
	// get all the keys for Products
	ctx := context.Background()
	// add one first just so all this doesn't fail
	// addProducts(1)
	client, _ := datastore.NewClient(ctx, "")
	qry := datastore.NewQuery("Event").KeysOnly().Filter("InfoURL=", "http://hello.org").Filter("Title <", "Event9")
	// var keylist []*datastore.Key
	keys, err := client.GetAll(ctx, qry, nil)

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	// if any, delete them
	if len(keys) > 0 {
		client.DeleteMulti(ctx, keys)
	}

}

func TestNoEvents(t *testing.T) {
	clearTempEvents()

	req, _ := http.NewRequest("GET", "/events", nil)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	clearTempEvents()

	req, _ := http.NewRequest("GET", "/event/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Event not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Event not found'. Got '%s'", m["error"])
	}
}

func TestCreateEvent(t *testing.T) {

	clearTempEvents()

	var jsonStr = []byte(`{"title":"Event702", "active": true}`)
	req, _ := http.NewRequest("POST", "/event", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "Event702" {
		t.Errorf("Expected event name to be 'Event702'. Got '%v'", m["name"])
	}

	if m["active"] != true {
		t.Errorf("Expected product Active to be true. Got '%v'", m["active"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	// if m["id"] != "" {
	// 	t.Errorf("Expected product ID to be positive number. Got '%v'", m["id"])
	// }
}

func TestGetEvent(t *testing.T) {
	clearTempEvents()
	addTempEvents(1)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addTempEvents(count int) {
	if count < 1 {
		count = 1
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		fmt.Printf("Could not connect to database: %v", err)
	}
	for i := 0; i < count; i++ {
		idNum := i + 1
		s := Event{Title: "Event" + strconv.Itoa(idNum), Active: true,
			InfoURL: "http://hello.org"}
		k := datastore.IDKey("Event", int64(idNum), nil)
		_, err := client.Put(ctx, k, &s)
		if err != nil {
			fmt.Printf("Could not add Event to database: %v ", err)
		}
	}
}

func TestUpdateEvent(t *testing.T) {

	clearTempEvents()
	addTempEvents(1)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	response := executeRequest(req)
	var originalEvent map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalEvent)

	var jsonStr = []byte(`{"organization":"happy nonprofit", "active":false, "regions":["A","B","C"]}`)
	req, _ = http.NewRequest("PUT", "/event/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	// fmt.Printf("Original: %v, New: %v", originalEvent, m)

	if m["id"] != originalEvent["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalEvent["id"], m["id"])
	}

	if m["organization"] == originalEvent["organization"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalEvent["organization"], m["organization"], m["organization"])
	}

	if m["active"] == originalEvent["active"] {
		t.Errorf("Expected activity to change from '%v' to '%v'. Got '%v'", originalEvent["price"], m["price"], m["price"])
	}
}

func TestDeleteEvent(t *testing.T) {
	clearTempEvents()
	addTempEvents(1)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/event/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/event/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
