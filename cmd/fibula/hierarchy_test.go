// hierarchy_test.go

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"cloud.google.com/go/datastore"
)

func createContest() {
	// newID := makeXID()
	clearTempContests()
	// ID := "dw91aj70xxxa"
	var jsonStr = []byte(`{"title":"TestingTesting Contest", "id": "dw91aj70xxxa", "active": true}`)
	req, _ := http.NewRequest("POST", "/contest", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	executeRequest(req)

}

func clearAll() {
	req, _ := http.NewRequest("DELETE", "/contest/dw91aj70xxxa", nil)
	executeRequest(req)
}

// func clearTempContests() {
// 	// get all the keys for Products
// 	ctx := context.Background()
// 	// add one first just so all this doesn't fail
// 	// addProducts(1)
// 	client, _ := datastore.NewClient(ctx, "")
// 	qry := datastore.NewQuery("Contest").KeysOnly()
// 	// var keylist []*datastore.Key
// 	keys, err := client.GetAll(ctx, qry, nil)

// 	if err != nil {
// 		fmt.Printf("%v", err)
// 		return
// 	}

// 	// if any, delete them
// 	if len(keys) > 0 {
// 		client.DeleteMulti(ctx, keys)
// 	}

// }

func TestNoContests(t *testing.T) {
	clearTempContests()

	req, _ := http.NewRequest("GET", "/contests", nil)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "null" {
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

func TestGetNonExistentContest(t *testing.T) {
	clearTempContests()

	req, _ := http.NewRequest("GET", "/contest/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Not found'. Got '%s'", m["error"])
	}
}

func TestCreateContest(t *testing.T) {
	// newID := makeXID()
	clearTempContests()
	newID := "dw91aj70"
	var jsonStr = []byte(`{"title":"Contest702", "id": "dw91aj70", "active": true}`)
	req, _ := http.NewRequest("POST", "/contest", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "Contest702" {
		t.Errorf("Expected contest name to be 'Contest702'. Got '%v'", m["name"])
	}

	if m["active"] != true {
		t.Errorf("Expected Active to be true. Got '%v'", m["active"])
	}

	if m["id"] != newID {
		t.Errorf("Expected ID to be %v. Got '%v'", newID, m["id"])

	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	// if m["id"] != "" {
	// 	t.Errorf("Expected product ID to be positive number. Got '%v'", m["id"])
	// }
}

func TestGetContest(t *testing.T) {
	clearTempContests()
	addTempContests(1)

	req, _ := http.NewRequest("GET", "/contest/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetContests(t *testing.T) {
	clearTempContests()
	addTempContests(4)

	req, _ := http.NewRequest("GET", "/contests", nil)

	response := executeRequest(req)
	// fmt.Printf("%v", response.Body)
	var m []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if len(m) != 4 {
		t.Errorf("Expected %v contests. Got %v", 4, len(m))
	}

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTempContests(count int) {
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
		s := Contest{Title: "Contest" + strconv.Itoa(idNum), Active: true,
			ID:      strconv.Itoa(idNum),
			InfoURL: "http://hello.org"}

		k := datastore.NameKey("Contest", strconv.Itoa(idNum), nil)
		_, err := client.Put(ctx, k, &s)
		if err != nil {
			fmt.Printf("Could not add Event to database: %v ", err)
		}
	}
}

func TestUpdateContest(t *testing.T) {

	clearTempContests()
	addTempContests(1)

	req, _ := http.NewRequest("GET", "/contest/1", nil)
	response := executeRequest(req)
	var originalEvent map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalEvent)

	var jsonStr = []byte(`{"id":"1", "organization":"happy nonprofit", "active":false, "regions":["A","B","C"]}`)
	req, _ = http.NewRequest("PUT", "/contest/1", bytes.NewBuffer(jsonStr))
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
		t.Errorf("Expected activity to change from '%v' to '%v'. Got '%v'", originalEvent["activity"], m["price"], m["price"])
	}
}

func TestDeleteContest(t *testing.T) {
	clearTempContests()
	addTempContests(1)

	req, _ := http.NewRequest("GET", "/contest/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/contest/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/contest/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
